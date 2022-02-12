package main

import (
	"crypto/rand"
	"crypto/sha256"
	"embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ewxrjk/bank"
	"github.com/ewxrjk/bank/util"
	"github.com/gorilla/handlers"
	"github.com/spf13/cobra"
)

var serverAddress, serverKey, serverCert string
var staticPageLifetime int

func init() {
	serverCmd.PersistentFlags().StringVarP(&serverAddress, "address", "a", "localhost:80", "listen address")
	serverCmd.PersistentFlags().StringVarP(&serverCert, "cert", "c", "", "server certificate")
	serverCmd.PersistentFlags().StringVarP(&serverKey, "key", "k", "", "server private key")
	serverCmd.PersistentFlags().IntVarP(&staticPageLifetime, "lifetime", "L", 60, "static page lifetime")
}

var secure bool

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the bank web service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) != 0 {
			return errors.New("usage: bank server [OPTIONS]")
		}
		if err = setup(); err != nil {
			return
		}
		if err = namespace.Initialize(); err != nil {
			return
		}
		http.Handle("/", handlers.LoggingHandler(os.Stderr, &namespace))
		secure = serverCert != "" && serverKey != ""
		if secure {
			err = http.ListenAndServeTLS(serverAddress, serverCert, serverKey, nil)
		} else {
			err = http.ListenAndServe(serverAddress, nil)
		}
		return
	},
}

var namespace = util.HTTPNamespace{
	Paths: []*util.HTTPPath{
		{"POST", "^/v1/login$", handlePostLogin},
		{"POST", "^/v1/logout$", handlePostLogout},
		{"GET", "^/v1/user/?$", handleGetUser},
		{"POST", "^/v1/user/?$", handlePostUser},
		{"PUT", "^/v1/user/([^/]+)/password$", handlePutUserPassword},
		{"DELETE", "^/v1/user/([^/]+)$", handleDeleteUser},
		{"GET", "^/v1/account/?$", handleGetAccount},
		{"POST", "^/v1/account/?$", handlePostAccount},
		{"DELETE", "^/v1/account/([^/]+)$", handleDeleteAccount},
		{"GET", "^/v1/transaction/?$", handleGetTransaction},
		{"POST", "^/v1/transaction/?$", handlePostTransaction},
		{"POST", "^/v1/distribute/?$", handlePostDistribute},
		{"GET", "^/v1/config/?$", handleGetConfig},
		{"GET", "^/v1/config/([^/]+)$", handleGetConfigKey},
		{"PUT", "^/v1/config/([^/]+)$", handlePutConfigKey},
		{"GET", "^/.*", handleGetRoot},
	},
}

// Session handling

// Session defines a login session.
type Session struct {
	// User who may be logged in
	user string

	// Non-cookie coken
	token string

	// Epxiry time
	expires time.Time

	// Tag for cache (in)validation
	tag string
}

var sessions = map[string]*Session{}
var sessionLock sync.Mutex

const cookieName = "bank"

// LoginRequest is the JSON request for a new login
type LoginRequest struct {
	User     string
	Password string
}

// LoginResponse is  the JSON respons for a succesful login
type LoginResponse struct {
	Token string
}

// POST /v1/login/
func handlePostLogin(w http.ResponseWriter, r *http.Request, matches []string) {
	var jreq LoginRequest
	var err error
	if err = json.NewDecoder(r.Body).Decode(&jreq); err != nil {
		log.Printf("decoding JSON: %v", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err = b.CheckPassword(jreq.User, jreq.Password); err != nil {
		log.Printf("CheckPassword for %v: %v", jreq.User, err)
		http.Error(w, "invalid credentials", http.StatusForbidden)
		return
	}
	b := make([]byte, 3*18)
	if _, err = rand.Read(b); err != nil {
		log.Printf("rand.Read: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	ident := base64.URLEncoding.EncodeToString(b[:18])
	token := base64.URLEncoding.EncodeToString(b[18:36])
	tag := base64.URLEncoding.EncodeToString(b[36:])
	sessionLock.Lock()
	defer sessionLock.Unlock()
	expires := time.Now().Add(time.Hour * 8)
	sessions[ident] = &Session{jreq.User, token, expires, tag}
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    ident,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		Secure:   secure,
	})
	util.HTTPRespond(w, &LoginResponse{token})
}

// POST /v1/logout/
func handlePostLogout(w http.ResponseWriter, r *http.Request, matches []string) {
	var c *http.Cookie
	var err error
	if c, err = r.Cookie(cookieName); err != nil {
		http.Error(w, "not logged in", http.StatusOK)
		return
	}
	sessionLock.Lock()
	defer sessionLock.Unlock()
	delete(sessions, c.Value)
	expires := time.Now().Add(-24 * 365 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "logged out",
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		Secure:   secure,
	})
	http.Error(w, "logged out", http.StatusOK)
}

func getSession(w http.ResponseWriter, r *http.Request) (session *Session) {
	var c *http.Cookie
	var err error
	if c, err = r.Cookie(cookieName); err != nil {
		return
	}
	sessionLock.Lock()
	defer sessionLock.Unlock()
	var ok bool
	if session, ok = sessions[c.Value]; !ok {
		log.Printf("session %v unrecognized", c.Value)
		return
	}
	if time.Now().After(session.expires) {
		log.Printf("session %v expired", c.Value)
		delete(sessions, c.Value)
		session = nil
		return
	}
	return
}

func mustSession(w http.ResponseWriter, r *http.Request) (session *Session) {
	if session = getSession(w, r); session == nil {
		http.Error(w, "not logged in", http.StatusForbidden)
	}
	return
}

func handleGetUser(w http.ResponseWriter, r *http.Request, matches []string) {
	var err error
	var session *Session
	if session = mustSession(w, r); session == nil {
		return
	}
	var users []string
	if users, err = b.GetUsers(); err != nil {
		util.HTTPErrorResponse(w, err, "cannot get users")
		return
	}
	util.HTTPRespond(w, &users)
}

// NewUserRequest is the JSON request to create a new user.
type NewUserRequest struct {
	User     string
	Password string
	Token    string
}

func handlePostUser(w http.ResponseWriter, r *http.Request, matches []string) {
	var err error
	var session *Session
	if session = mustSession(w, r); session == nil {
		return
	}
	var jreq NewUserRequest
	if err = json.NewDecoder(r.Body).Decode(&jreq); err != nil {
		log.Printf("decoding JSON: %v", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if jreq.Token != session.token {
		log.Printf("token mismatch")
		http.Error(w, "inconsistent session", http.StatusForbidden)
		return
	}
	if jreq.User == "" {
		http.Error(w, "empty user name", http.StatusBadRequest)
		return
	}
	if jreq.Password == "" {
		http.Error(w, "empty password", http.StatusBadRequest)
		return
	}
	if err = b.NewUser(jreq.User, jreq.Password); err != nil {
		util.HTTPErrorResponse(w, err, "cannot create user")
		return
	}
	http.Error(w, "created user", http.StatusOK)
}

// ChangePasswordRequest is the JSON request to change a password.
type ChangePasswordRequest struct {
	Password string
	Token    string
}

func handlePutUserPassword(w http.ResponseWriter, r *http.Request, matches []string) {
	var err error
	var session *Session
	if session = mustSession(w, r); session == nil {
		return
	}
	var jreq ChangePasswordRequest
	if err = json.NewDecoder(r.Body).Decode(&jreq); err != nil {
		log.Printf("decoding JSON: %v", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if jreq.Token != session.token {
		log.Printf("token mismatch")
		http.Error(w, "inconsistent session", http.StatusForbidden)
		return
	}
	if err = b.SetPassword(matches[1], jreq.Password); err != nil {
		util.HTTPErrorResponse(w, err, "cannot set password")
		return
	}
	http.Error(w, "changed password", http.StatusOK)
}

func handleDeleteUser(w http.ResponseWriter, r *http.Request, matches []string) {
	var err error
	var session *Session
	if session = mustSession(w, r); session == nil {
		return
	}
	if err = b.DeleteUser(matches[1]); err != nil {
		util.HTTPErrorResponse(w, err, "cannot delete user")
		return
	}
}

// NewAccountRequest is the JSON request to create new account.
type NewAccountRequest struct {
	Account string
	Token   string
}

func handlePostAccount(w http.ResponseWriter, r *http.Request, matches []string) {
	var session *Session
	if session = mustSession(w, r); session == nil {
		return
	}
	var jreq NewAccountRequest
	var err error
	if err = json.NewDecoder(r.Body).Decode(&jreq); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if jreq.Token != session.token {
		log.Printf("token mismatch")
		http.Error(w, "inconsistent session", http.StatusForbidden)
		return
	}
	if jreq.Account == "" {
		http.Error(w, "empty account name", http.StatusBadRequest)
		return
	}
	if err = b.NewAccount(jreq.Account); err != nil {
		util.HTTPErrorResponse(w, err, "cannot create account")
		return
	}
	http.Error(w, "created account", http.StatusOK)
}

func handleGetAccount(w http.ResponseWriter, r *http.Request, matches []string) {
	var session *Session
	if session = mustSession(w, r); session == nil {
		return
	}
	var accounts []string
	var err error
	if accounts, err = b.GetAccounts(); err != nil {
		util.HTTPErrorResponse(w, err, "getting accounts")
		return
	}
	util.HTTPRespond(w, &accounts)
}

func handleDeleteAccount(w http.ResponseWriter, r *http.Request, matches []string) {
	var err error
	var session *Session
	if session = mustSession(w, r); session == nil {
		return
	}
	if err = b.DeleteAccount(matches[1]); err != nil {
		util.HTTPErrorResponse(w, err, "cannot delete account")
		return
	}
}

// NewTransactionRequest is the JSON request to create a new transaction.
type NewTransactionRequest struct {
	Token       string
	Origin      string
	Destination string
	Description string
	Amount      int
}

func handlePostTransaction(w http.ResponseWriter, r *http.Request, matches []string) {
	var session *Session
	if session = mustSession(w, r); session == nil {
		return
	}
	var err error
	var jreq NewTransactionRequest
	if err = json.NewDecoder(r.Body).Decode(&jreq); err != nil {
		log.Printf("decoding JSON: %v", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if jreq.Token != session.token {
		log.Printf("token mismatch")
		http.Error(w, "inconsistent session", http.StatusForbidden)
		return
	}
	if err = b.NewTransaction(session.user, jreq.Origin, jreq.Destination, jreq.Description, jreq.Amount); err != nil {
		util.HTTPErrorResponse(w, err, "cannot create transaction")
		return
	}
	http.Error(w, "created transaction", http.StatusOK)
}

func handleGetTransaction(w http.ResponseWriter, r *http.Request, matches []string) {
	var session *Session
	if session = mustSession(w, r); session == nil {
		return
	}
	var err error
	if err = r.ParseForm(); err != nil {
		log.Printf("parsing query: %v", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	offset, _ := strconv.Atoi(r.FormValue("offset"))
	after, _ := strconv.Atoi(r.FormValue("after"))
	var transactions []bank.Transaction
	if transactions, err = b.GetTransactions(limit, offset, after); err != nil {
		util.HTTPErrorResponse(w, err, "cannot get transactions")
		return
	}
	util.HTTPRespond(w, transactions)
}

// DistributeRequest is the JSON request to create distribution transactions.
type DistributeRequest struct {
	Token        string
	Origin       string
	Destinations []string
	Description  string
}

// POST /v1/distribute/
func handlePostDistribute(w http.ResponseWriter, r *http.Request, matches []string) {
	var session *Session
	if session = mustSession(w, r); session == nil {
		return
	}
	var err error
	if r.Method == "POST" {
		var jreq DistributeRequest
		if err = json.NewDecoder(r.Body).Decode(&jreq); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		if jreq.Token != session.token {
			log.Printf("token mismatch")
			http.Error(w, "inconsistent session", http.StatusForbidden)
			return
		}
		if err = b.Distribute(session.user, jreq.Origin, jreq.Destinations, jreq.Description); err != nil {
			util.HTTPErrorResponse(w, err, "cannot distribute")
			return
		}
		http.Error(w, "distributed", http.StatusOK)
	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func handleGetConfig(w http.ResponseWriter, r *http.Request, matches []string) {
	var session *Session
	if session = mustSession(w, r); session == nil {
		return
	}
	var configs map[string]string
	var err error
	if configs, err = b.GetConfigs(); err != nil {
		util.HTTPErrorResponse(w, err, "cannot get configuration item")
		return
	}
	util.HTTPRespond(w, configs)
}

func handleGetConfigKey(w http.ResponseWriter, r *http.Request, matches []string) {
	var session *Session
	if session = mustSession(w, r); session == nil {
		return
	}
	var value string
	var err error
	if value, err = b.GetConfig(matches[1]); err != nil {
		if err == bank.ErrNoSuchConfig {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			log.Printf("getting configs: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(value))
}

// ConfigRequest is the JSON request to set a configuration item.
type ConfigRequest struct {
	Value string
	Token string
}

func handlePutConfigKey(w http.ResponseWriter, r *http.Request, matches []string) {
	var session *Session
	if session = mustSession(w, r); session == nil {
		return
	}
	var jreq ConfigRequest
	var err error
	if err = json.NewDecoder(r.Body).Decode(&jreq); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if jreq.Token != session.token {
		log.Printf("token mismatch")
		http.Error(w, "inconsistent session", http.StatusForbidden)
		return
	}
	if err = b.PutConfig(matches[1], jreq.Value); err != nil {
		util.HTTPErrorResponse(w, err, "cannot set configuration item")
		return
	}
	http.Error(w, "set configuration item", http.StatusOK)
}

// embedTemplate is the map of paths to parsed templates.
var embedTemplate = map[string]*template.Template{}

// embedHashes is the map of fixed path to hashes.
// For templates, it contains the hash of the template;
// it can't be directly used as an entity tag.
var embedHashes = map[string]string{}

func init() {
	initializeTags("ui")
}

func initializeTags(dir string) {
	var entries []fs.DirEntry
	var err error
	if entries, err = content.ReadDir(dir); err != nil {
		log.Fatalf("content.ReadDir: %v", err)
	}
	for _, entry := range entries {
		name := fmt.Sprintf("%s/%s", dir, entry.Name())
		if entry.IsDir() {
			initializeTags(name)
		} else {
			var data []byte
			if data, err = content.ReadFile(name); err != nil {
				log.Fatalf("content.ReadFile %s: %v", name, err)
			}
			hash := sha256.Sum256(data)
			tag := base64.RawURLEncoding.EncodeToString(hash[:18])
			embedHashes[name] = tag
			if path.Ext(name) == ".html" {
				embedTemplate[name] = template.Must(template.New(name).Parse(string(data)))
			}
		}
	}
}

// TemplateData is the data object type for template execution.
type TemplateData struct {
	Token string
	Title string
	User  string
}

//go:embed ui/*[^~]
var content embed.FS

var mimeTypes = map[string]string{
	".css":  "text/css",
	".html": "text/html;charset=utf-8",
	".js":   "text/javascript",
	".png":  "image/png",
}

// GET /
func handleGetRoot(w http.ResponseWriter, r *http.Request, matches []string) {
	var err error
	name := r.URL.Path[1:]
	if name == "" {
		name = "index.html"
	}
	// Everything lives under ui/
	name = "ui/" + name
	var weak, etag string
	var data []byte
	var ok bool
	// Prepare the content and compute the etag
	var template *template.Template
	var templateData TemplateData
	if template, ok = embedTemplate[name]; ok {
		templateData.Token = "not logged in"
		if templateData.Title, err = b.GetConfig("title"); err != nil {
			log.Printf("getting title: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		// We exclude the title from the calculation of the tag,
		// making it weak validator.
		weak = "W/"
		if session := getSession(w, r); session != nil {
			templateData.Token = session.token
			templateData.User = session.user
			etag = embedHashes[name] + session.tag
		} else {
			etag = embedHashes[name]
		}
	} else if data, err = content.ReadFile(name); err == nil {
		etag = embedHashes[name]
	} else {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	// Maybe the request is conditional
	status := util.CheckEntityTag(w, r, etag)
	// Generate the content (if not static)
	if status == http.StatusOK && template != nil {
		writer := strings.Builder{}
		if err = template.Execute(&writer, templateData); err != nil {
			log.Printf("cannot serve page: executing template %s: %v", name, err)
			http.Error(w, "cannot serve page", http.StatusInternalServerError)
			return
		}
		data = []byte(writer.String())
	}
	ct := mimeTypes[filepath.Ext(name)]
	if ct == "" {
		ct = "application/octet-string"
	}
	w.Header().Set("Content-Type", ct)
	if template == nil {
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", staticPageLifetime))
	} else {
		w.Header().Set("Cache-Control", "no-cache")
	}
	if etag != "" {
		w.Header().Set("ETag", fmt.Sprintf(`%s"%s"`, weak, etag))
	}
	w.WriteHeader(status)
	if status == http.StatusOK {
		w.Write(data)
	}
}
