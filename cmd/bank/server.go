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
	"reflect"
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
var debug bool

func init() {
	serverCmd.PersistentFlags().StringVarP(&serverAddress, "address", "a", "localhost:80", "listen address")
	serverCmd.PersistentFlags().StringVarP(&serverCert, "cert", "c", "", "server certificate")
	serverCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "log extra information for debugging (insecure)")
	// see use sites for why debug is 'isnecure'
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
		{
			Method:    "POST",
			Path:      "^/v1/login$",
			ServeHTTP: handlePostLogin,
		},
		{
			Method:    "POST",
			Path:      "^/v1/logout$",
			ServeHTTP: handlePostLogout,
		},
		{
			Method:    "GET",
			Path:      "^/v1/user/?$",
			ServeHTTP: handleGetUser,
		},
		{
			Method:    "POST",
			Path:      "^/v1/user/?$",
			ServeHTTP: handlePostUser,
		},
		{
			Method:    "PUT",
			Path:      "^/v1/user/([^/]+)/password$",
			ServeHTTP: handlePutUserPassword,
		},
		{
			Method:    "DELETE",
			Path:      "^/v1/user/([^/]+)$",
			ServeHTTP: handleDeleteUser,
		},
		{
			Method:    "GET",
			Path:      "^/v1/account/?$",
			ServeHTTP: handleGetAccount,
		},
		{
			Method:    "POST",
			Path:      "^/v1/account/?$",
			ServeHTTP: handlePostAccount,
		},
		{
			Method:    "DELETE",
			Path:      "^/v1/account/([^/]+)$",
			ServeHTTP: handleDeleteAccount,
		},
		{
			Method:    "GET",
			Path:      "^/v1/transaction/?$",
			ServeHTTP: handleGetTransaction,
		},
		{
			Method:    "POST",
			Path:      "^/v1/transaction/?$",
			ServeHTTP: handlePostTransaction,
		},
		{
			Method:    "POST",
			Path:      "^/v1/distribute/?$",
			ServeHTTP: handlePostDistribute,
		},
		{
			Method:    "GET",
			Path:      "^/v1/config/?$",
			ServeHTTP: handleGetConfig,
		},
		{
			Method:    "GET",
			Path:      "^/v1/config/([^/]+)$",
			ServeHTTP: handleGetConfigKey,
		},
		{
			Method:    "PUT",
			Path:      "^/v1/config/([^/]+)$",
			ServeHTTP: handlePutConfigKey,
		},
		{
			Method:    "GET",
			Path:      "^/.*",
			ServeHTTP: handleGetRoot,
		},
	},
}

// Session handling

// Session defines a login session.
type Session struct {
	// User who may be logged in
	user string

	// Non-cookie coken
	token string

	// Expiry time
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

// LoginResponse is  the JSON response for a succesful login
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
		// Sometimes passwords are accidentally entered into username fields,
		// so don't log usernames unless debug is enabled.
		if debug {
			log.Printf("CheckPassword for %v: %v", jreq.User, err)
		} else {
			log.Printf("CheckPassword: %v", err)
		}
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

func handleGetUser(w http.ResponseWriter, r *http.Request, matches []string) {
	var err error
	if _, ok := decodeRequest(w, r, nil, true); !ok {
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
	var jreq NewUserRequest
	if _, ok := decodeRequest(w, r, &jreq, true); !ok {
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
	var jreq ChangePasswordRequest
	if _, ok := decodeRequest(w, r, &jreq, true); !ok {
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
	if _, ok := decodeRequest(w, r, nil, true); !ok {
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
	var err error
	var jreq NewAccountRequest
	if _, ok := decodeRequest(w, r, &jreq, true); !ok {
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
	if _, ok := decodeRequest(w, r, nil, true); !ok {
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
	if _, ok := decodeRequest(w, r, nil, true); !ok {
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
	var ok bool
	var err error
	var jreq NewTransactionRequest
	if session, ok = decodeRequest(w, r, &jreq, true); !ok {
		return
	}
	if err = b.NewTransaction(session.user, jreq.Origin, jreq.Destination, jreq.Description, jreq.Amount); err != nil {
		util.HTTPErrorResponse(w, err, "cannot create transaction")
		return
	}
	http.Error(w, "created transaction", http.StatusOK)
}

func handleGetTransaction(w http.ResponseWriter, r *http.Request, matches []string) {
	if _, ok := decodeRequest(w, r, nil, true); !ok {
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
	var err error
	var jreq DistributeRequest
	var session *Session
	var ok bool
	if session, ok = decodeRequest(w, r, &jreq, true); !ok {
		return
	}
	if err = b.Distribute(session.user, jreq.Origin, jreq.Destinations, jreq.Description); err != nil {
		util.HTTPErrorResponse(w, err, "cannot distribute")
		return
	}
	http.Error(w, "distributed", http.StatusOK)
}

func handleGetConfig(w http.ResponseWriter, r *http.Request, matches []string) {
	if _, ok := decodeRequest(w, r, nil, true); !ok {
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
	if _, ok := decodeRequest(w, r, nil, true); !ok {
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
	var jreq ConfigRequest
	var err error
	if _, ok := decodeRequest(w, r, &jreq, true); !ok {
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
	Token   string
	Title   string
	User    string
	Version string
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
	templateData.Version = bank.Version
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

// decodeRequest decodes and authenticates a request.
//
// w and r are the normal HTTP request response/request pair.
//
// req is a pointer to the request object to be filled in, or nil if there isn't one.
// If it has a Token field then it must match the token from the session.
//
// If mustAuth is true then a cookie identifying a live session must be specified
// (and it will be returned).
func decodeRequest(w http.ResponseWriter, r *http.Request, req interface{}, mustAuth bool) (session *Session, ok bool) {
	var err error
	if mustAuth {
		// We need a live session
		if session = getSession(w, r); session == nil {
			http.Error(w, "not logged in", http.StatusForbidden)
			return
		}
	}
	if req != nil {
		// Decode the request
		if err = json.NewDecoder(r.Body).Decode(req); err != nil {
			log.Printf("decoding JSON: %v", err)
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		// If there the request expects a token then check it matches the session
		// This is an anti-CSRF measure.
		if f, found := reflect.TypeOf(req).Elem().FieldByName("Token"); found {
			fv := reflect.ValueOf(req).Elem().FieldByIndex(f.Index)
			if fv.String() != session.token {
				log.Printf("token mismatch")
				http.Error(w, "inconsistent session", http.StatusForbidden)
				return
			}
		}
	}
	ok = true
	return
}

// getSession get a valid session from an HTTP request.
// On error it writes an HTTP error and returns nil.
func getSession(w http.ResponseWriter, r *http.Request) (session *Session) {
	var c *http.Cookie
	var err error
	// Find the cookie from the HTTP request
	if c, err = r.Cookie(cookieName); err != nil {
		return
	}
	// Synchronize access to the session store
	sessionLock.Lock()
	defer sessionLock.Unlock()
	// Unrecognized and stale sessions should never become valid again,
	// so should be harmless to log - nevertheless we hide them unless
	// debug is enabled, in case something undermines our assumptions
	// (e.g. VM snapshot restoration).

	// Find the session
	var ok bool
	if session, ok = sessions[c.Value]; !ok {
		if debug {
			log.Printf("session %v unrecognized", c.Value)
		} else {
			log.Printf("session unrecognized")
		}
		return
	}
	// Check the session has not expired
	if time.Now().After(session.expires) {
		if debug {
			log.Printf("session %v expired", c.Value)
		} else {
			log.Printf("session expired")
		}
		// Garbage-collect expired sesions
		delete(sessions, c.Value)
		session = nil
		return
	}
	return
}
