package util

import (
	"encoding/json"
	"git.anjou.terraraq.org.uk/bank"
	"log"
	"net/http"
	"regexp"
)

// HTTPPath represents a path and method.
type HTTPPath struct {
	// Only requests with this method are matched.
	Method string

	// Only requests matching this regexp are matched.
	Path string

	// Matching requests are dispatched to this function.
	//
	// matches consists of the matching string and the subexpression matches.
	ServeHTTP func(w http.ResponseWriter, r *http.Request, matches []string)
}

// HTTPNamespace represents a complete HTTP namespace.
//
// You must call Initialize() before handling any requests.
type HTTPNamespace struct {
	// Paths is processed in order until the request matches.
	Paths []*HTTPPath

	// regexps is the compiler regexps from Path fields.
	regexps []*regexp.Regexp
}

// Initialize prepares an HTTPNamespace for use.
func (ns *HTTPNamespace) Initialize() (err error) {
	for _, path := range ns.Paths {
		var re *regexp.Regexp
		if re, err = regexp.Compile(path.Path); err != nil {
			return err
		}
		ns.regexps = append(ns.regexps, re)
	}
	return
}

// ServeHTTP dispatches HTTP requests.
func (ns *HTTPNamespace) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	notFoundError := http.StatusNotFound
	for i, path := range ns.Paths {
		if matches := ns.regexps[i].FindStringSubmatch(r.URL.Path); matches != nil {
			if path.Method == r.Method {
				path.ServeHTTP(w, r, matches)
				return
			} else {
				notFoundError = http.StatusMethodNotAllowed
			}
		}
	}
	http.Error(w, "not found", notFoundError)
}

// HTTPRespond issue an HTTP response with JSON content.
func HTTPRespond(w http.ResponseWriter, jres interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&jres)
}

// HTTPErrorResponse issues an error response appropriate to err.
func HTTPErrorResponse(w http.ResponseWriter, err error, action string) {
	log.Printf("%s: %v", action, err)
	switch err {
	case bank.ErrUserExists,
		bank.ErrAccountExists,
		bank.ErrAccountHasBalance,
		bank.ErrInsufficientFunds,
		bank.ErrUnsuitableParties:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case bank.ErrNoSuchAccount,
		bank.ErrNoSuchConfig,
		bank.ErrNoSuchUser:
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		http.Error(w, action, http.StatusInternalServerError)
	}
}
