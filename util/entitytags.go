package util

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
)

var entityTagRegexp = regexp.MustCompile(`^(W/)?"([\x21\x23-\x7E]*)"|\*`)

// EntityTag represents an RFC7230 entity-tag
type EntityTag struct {
	// True if the tag is "*"
	All bool

	// True for weak tags
	Weak bool

	// The value of the tag
	Tag string
}

// ParseEntityTags returns the list of entity-tag values
// from an RFC7232 If-(None-)Match header.
func ParseEntityTags(header string) (tags []EntityTag, err error) {
	rest := []byte(header)
	needComma := false
	for len(rest) > 0 {
		switch rest[0] {
		case ' ', '\t', '\r', '\n':
			rest = rest[1:]
			continue
		case ',':
			// RFC7230 s7 MUST accept empty lists elements
			rest = rest[1:]
			needComma = false
			continue
		}
		if needComma {
			err = errors.New("malformed entity-tag list: missing comma")
			return
		}
		var m [][]byte
		if m = entityTagRegexp.FindSubmatch(rest); m == nil {
			err = errors.New("malformed entity-tag list: does not match entity-tag")
			return
		}
		if rest[0] == '*' {
			tags = append(tags, EntityTag{
				All: true,
			})
		} else {
			tags = append(tags, EntityTag{
				Weak: rest[0] == 'W',
				Tag:  string(m[2]),
			})
		}
		rest = rest[len(m[0]):]
	}
	if len(tags) == 0 {
		err = errors.New("malformed entity-tag list: no tags")
		return
	}
	return
}

// CheckEntityTag sets the ETag response header,
// and then checks the If-None-Match request header.
// If it is set and has a match for the entity tag,
// it returns http.StatusNotModified; in any other
// case it returns http.StatusOK.
func CheckEntityTag(w http.ResponseWriter, r *http.Request, etag string) (status int) {
	status = http.StatusOK
	if etag == "" {
		return
	}
	w.Header().Set("ETag", fmt.Sprintf(`"%s"`, etag))
	// See if the client has seen it before
	for _, h := range r.Header["If-None-Match"] {
		if tags, err := ParseEntityTags(h); err != nil {
			continue // ignore malformed headers
		} else {
			for _, tag := range tags {
				if tag.All || etag == tag.Tag {
					status = http.StatusNotModified
					break
				}
			}
		}
	}
	return
}
