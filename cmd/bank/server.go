package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"git.anjou.terraraq.org.uk/bank"
	"github.com/gorilla/handlers"
	"github.com/spf13/cobra"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var serverAddress, serverKey, serverCert, serviceName string

func init() {
	serverCmd.PersistentFlags().StringVarP(&serverAddress, "address", "a", "localhost:80", "listen address")
	serverCmd.PersistentFlags().StringVarP(&serverCert, "cert", "c", "", "server certificate")
	serverCmd.PersistentFlags().StringVarP(&serverKey, "key", "k", "", "server private key")
	serverCmd.PersistentFlags().StringVarP(&serviceName, "name", "n", "Bank", "service name")
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
		for _, handler := range namespace {
			http.Handle(handler.Path, handlers.LoggingHandler(os.Stderr, handler.Handler))
		}
		secure = serverCert != "" && serverKey != ""
		if secure {
			err = http.ListenAndServeTLS(serverAddress, serverCert, serverKey, nil)
		} else {
			err = http.ListenAndServe(serverAddress, nil)
		}
		return
	},
}

type handlerPath struct {
	Path    string
	Handler http.HandlerFunc
}

var namespace = []handlerPath{
	{"/v1/login", handleLogin},
	{"/v1/logout", handleLogout},
	{"/v1/user/", handleUser},
	{"/v1/user/password", handleUserPassword},
	{"/v1/account/", handleAccount},
	{"/v1/transaction/", handleTransaction},
	{"/v1/distribute/", handleDistribute},
	{"/", handleRoot},
}

// Session handling

// Session defines a login session.
type Session struct {
	user    string
	token   string
	expires time.Time
}

var sessions = map[string]*Session{}
var sessionLock sync.Mutex

const cookieName = "bank"

type LoginRequest struct {
	User     string
	Password string
}

type LoginResponse struct {
	Token string
}

// POST /v1/login/
func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "must use POST", http.StatusMethodNotAllowed)
		return
	}
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
	b := make([]byte, 36)
	if _, err = rand.Read(b); err != nil {
		log.Printf("rand.Read: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	var ident, token string
	ident = base64.URLEncoding.EncodeToString(b[:18])
	token = base64.URLEncoding.EncodeToString(b[18:])
	sessionLock.Lock()
	defer sessionLock.Unlock()
	expires := time.Now().Add(time.Hour * 8)
	sessions[ident] = &Session{jreq.User, token, expires}
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    ident,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		Secure:   secure,
	})
	respond(w, &LoginResponse{token})
}

// POST /v1/logout/
func handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "must use POST", http.StatusMethodNotAllowed)
		return
	}
	var c *http.Cookie
	var err error
	if c, err = r.Cookie(cookieName); err != nil {
		http.Error(w, "not logged in", http.StatusOK)
		return
	}
	sessionLock.Lock()
	defer sessionLock.Unlock()
	delete(sessions, c.Value)
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

func respond(w http.ResponseWriter, jres interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&jres)
}

type NewUserRequest struct {
	User     string
	Password string
	Token    string
}

// GET /v1/user/
// POST /v1/user/
func handleUser(w http.ResponseWriter, r *http.Request) {
	var err error
	var session *Session
	if session = mustSession(w, r); session == nil {
		return
	}
	if r.Method == "POST" {
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
			if err == bank.ErrUserExists {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				log.Printf("NewUser for %v: %v", jreq.User, err)
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
			return
		}
		http.Error(w, "created user", http.StatusOK)
	} else if r.Method == "GET" {
		var users []string
		if users, err = b.GetUsers(); err != nil {
			log.Printf("GetUsers: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		respond(w, &users)
	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

type ChangePasswordRequest struct {
	User     string
	Password string
	Token    string
}

// POST /v1/user/password
func handleUserPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "must use POST", http.StatusMethodNotAllowed)
		return
	}
	var session *Session
	if session = mustSession(w, r); session == nil {
		return
	}
	var jreq ChangePasswordRequest
	var err error
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
		jreq.User = session.user
	}
	if err = b.SetPassword(jreq.User, jreq.Password); err != nil {
		log.Printf("SetPassword for %v: %v", jreq.User, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.Error(w, "changed password", http.StatusOK)
}

type NewAccountRequest struct {
	Account string
	Token   string
}

// POST /v1/account/
// GET /v1/account/
func handleAccount(w http.ResponseWriter, r *http.Request) {
	var session *Session
	if session = mustSession(w, r); session == nil {
		return
	}
	var err error
	if r.Method == "POST" {
		var jreq NewAccountRequest
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
			if err == bank.ErrAccountExists {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				log.Printf("creating account: %v", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
			return
		}
		http.Error(w, "created account", http.StatusOK)
	} else if r.Method == "GET" {
		var accounts []string
		if accounts, err = b.GetAccounts(); err != nil {
			log.Printf("GetAccounts: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		respond(w, &accounts)
	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

type NewTransactionRequest struct {
	Token       string
	Origin      string
	Destination string
	Description string
	Amount      int
}

// POST /v1/transaction/
// GET /v1/transaction/?limit=LIMIT&offset=OFFSET&after=AFTER
func handleTransaction(w http.ResponseWriter, r *http.Request) {
	var session *Session
	if session = mustSession(w, r); session == nil {
		return
	}
	var err error
	if r.Method == "POST" {
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
			log.Printf("creating transaction: %v", err)
			// TODO invalid origin/destination are Bad Request, anything else is internal error
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "created transaction", http.StatusOK)
	} else if r.Method == "GET" {
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
			log.Printf("get transactions: %v", err)
			http.Error(w, "cannot get transactions", http.StatusInternalServerError)
			return
		}
		respond(w, transactions)
	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

type DistributeRequest struct {
	Token        string
	Origin       string
	Destinations []string
	Description  string
}

// POST /v1/distribute/
func handleDistribute(w http.ResponseWriter, r *http.Request) {
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
			log.Printf("creating transaction: %v", err)
			// TODO invalid origin/destination are Bad Request, anything else is internal error (NB not true, e.g. not enough funds)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "distributed", http.StatusOK)
	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

var embedTemplate = map[string]*template.Template{}

func init() {
	for name, content := range embedContent {
		if embedType[name] == "text/html" {
			embedTemplate[name] = template.Must(template.New(name).Parse(content))
		}
	}
}

type TemplateData struct {
	Token string
	Title string
}

// GET /
func handleRoot(w http.ResponseWriter, r *http.Request) {
	var err error
	path := r.URL.Path[1:]
	if path == "" {
		path = "index.html"
	}
	var content string
	var ok bool
	var template *template.Template
	if template, ok = embedTemplate[path]; ok {
		writer := strings.Builder{}
		data := TemplateData{
			Token: "not logged in",
			Title: serviceName,
		}
		if session := getSession(w, r); session != nil {
			data.Token = session.token
		} else {
			log.Printf("found no session")
		}
		if err = template.Execute(&writer, data); err != nil {
			log.Printf("executing template %s: %v", path, err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		content = writer.String()
	} else if content, ok = embedContent[path]; !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", embedType[path])
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(content))
}
