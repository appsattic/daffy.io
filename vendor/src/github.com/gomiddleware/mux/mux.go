package mux

import (
	"context"
	"errors"
	"log"
	"net/http"
	"path"
	"strings"
)

type key int

const valsIdKey key = 999

// Errors that can be returned from this package.
var (
	// ErrMultipleHandlers is returned when you create a route with multiple handlers.
	ErrMultipleHandlers = errors.New("mux: route has been given two handlers but only one can be provider")

	// ErrMiddlewareAfterHandler is returned when you create a route which has some middleware defined after the
	// handler.
	ErrMiddlewareAfterHandler = errors.New("mux: route can't have middleware defined after the handler")

	// ErrUnknownTypeInRoute is returned when something unexpected is passed to a route function.
	ErrUnknownTypeInRoute = errors.New("mux: unexpected type passed to route")
)

// Route is an internal method/path/middlewares/handler type created when each route is added.
type Route struct {
	Method      string
	Path        string
	Segments    []string
	Length      int
	Middlewares []func(http.Handler) http.Handler
	Handler     http.Handler
}

// Mux is just an array of Route.
type Mux struct {
	routes []Route
}

// Make sure the Mux conforms with the http.Handler interface.
var _ http.Handler = New()

// New returns a new initialized Mux.  Nothing is automatic. You must do slash/non-slash redirection yourself.
func New() *Mux {
	return &Mux{}
}

// Get is a shortcut for mux.add("GET", path, things...)
func (m *Mux) Get(path string, things ...interface{}) error {
	log.Printf("NewGet()\n")
	return m.add("GET", path, things...)
}

// Post is a shortcut for mux.add("POST", path, things...)
func (m *Mux) Post(path string, things ...interface{}) {
	m.add("POST", path, things...)
}

// Put is a shortcut for mux.add("PUT", path, things...)
func (m *Mux) Put(path string, things ...interface{}) {
	m.add("PUT", path, things...)
}

// Patch is a shortcut for mux.add("PATCH", path, things...)
func (m *Mux) Patch(path string, things ...interface{}) {
	m.add("PATCH", path, things...)
}

// Delete is a shortcut for mux.add("DELETE", path, things...)
func (m *Mux) Delete(path string, things ...interface{}) {
	m.add("DELETE", path, things...)
}

// Use adds some middleware to a path prefix. Unlike other methods such as Get, Post, Put, Patch, and Delete, Use
// matches for the prefix only and not the entire path. (Though of course, the entire exact path also matches.)
//
// e.g. m.Use("/profile/", ...) matches the requests "/profile/", "/profile/settings", and "/profile/a/path/to/".
//
// Note however, m.Use("/profile/", ...) doesn't match "/profile" since it contains too many slashes. But
// m.Use("/profile", ...) does match "/profile/" and "/profile/..." (but check that's actually what you want here).
//
// Also note that if you Use("/s/"), then this can be used to handle all static files inside /s/.
func (m *Mux) Use(path string, things ...interface{}) error {
	return m.add("USE", path, things...)
}

// add registers a new request handle with the given path and method.
//
// The respective shortcuts (for GET, POST, PUT, PATCH and DELETE) can also be used.
func (m *Mux) add(method, path string, things ...interface{}) error {
	log.Printf("add()\n")

	if path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	if m.routes == nil {
		m.routes = make([]Route, 0)
	}

	// collect up some things like the middlewares and the handler
	var handler http.Handler
	var middlewares []func(http.Handler) http.Handler

	segments := strings.Split(path, "/")[1:]

	log.Printf("Things = %#v\n", things)

	for i, thing := range things {
		log.Printf("Loop %d %#v\n", i, thing)
		switch val := thing.(type) {
		case func(http.Handler) http.Handler:
			log.Printf("got func(http.Handler) http.Handler\n")
			// if we already have a handler, then we should bork
			if len(middlewares) > 0 {
				return ErrMiddlewareAfterHandler
			}
			// all good, so add the middleware
			log.Printf("adding to middlewares")
			middlewares = append(middlewares, val)
		case http.Handler:
			log.Printf("got http.Handler\n")
			if handler != nil {
				log.Printf("already got a handler")
				return ErrMultipleHandlers
			}
			// all good, so remember the handler
			log.Printf("adding a handler")
			handler = val
		case func(http.ResponseWriter, *http.Request):
			log.Printf("got func(http.ResponseWriter, *http.Request)\n")
			if handler != nil {
				log.Printf("already got a handler")
				return ErrMultipleHandlers
			}
			// all good, so remember the handler
			log.Printf("adding a HandlerFunc")
			handler = http.HandlerFunc(val)
		default:
			return ErrUnknownTypeInRoute
		}
	}

	log.Printf("add(): now adding to the handlers\n")

	// create our handler which contains everything we need
	route := Route{
		Method:      method,
		Path:        path,
		Segments:    segments,
		Length:      len(segments),
		Middlewares: middlewares,
		Handler:     handler,
	}

	// add it to the handlers
	m.routes = append(m.routes, route)

	// log.Printf("routes=%#v\n", m.routes)
	return nil
}

func isPrefixMatch(segments []string, route *Route) bool {
	log.Printf("isPrefixMatch: %v\n", segments)

	log.Printf("Checking against %#v\n", route)

	// if segments is just []string{''} (ie, from "/"), then this will match everything
	if route.Length == 1 && route.Segments[0] == "" {
		return true
	}

	// can't match if the route prefix length is longer than the URL
	if route.Length > len(segments) {
		return false
	}

	// check each segment is the same (for the length of the prefix)
	for i, segment := range route.Segments {
		log.Printf("isPrefixMatch: checking '%s' against '%s'\n", segments[i], segment)

		// if both segments are empty, then this matches
		if segment == "" && segments[i] == "" {
			log.Printf(" - both empty, fine\n")
			continue
		}

		// check if segment start with a ":"
		if segment[0:0] == ":" {
			log.Printf("Placeholder = %s\n", segment)
			continue
		}

		if segments[i] != segment {
			log.Printf(" - not the same, this prefix doesn't match\n")
			return false
		}
	}

	// nothing stopped us from matching, so it must be true
	return true
}

func isMatch(method string, segments []string, route *Route) (map[string]string, bool) {
	log.Printf("isMatch: %s %v\n", method, segments)

	// can't match if the methods are different
	if route.Method != method {
		log.Printf("isMatch: different method (got %s, this route is %s)\n", method, route.Method)
		return nil, false
	}

	// can't match if the url length is different from the route length
	if route.Length != len(segments) {
		log.Printf("isMatch: different path length (got %d, this route is %d long)\n", len(segments), route.Length)
		return nil, false
	}

	vals := make(map[string]string)

	// check each segment is the same (for the length of the prefix)
	for i, segment := range route.Segments {
		log.Printf("isMatch: checking '%s' against '%s'\n", segments[i], segment)

		// if both segments are empty, then this matches
		if segment == "" && segments[i] == "" {
			log.Printf(" - both empty, fine\n")
			continue
		}

		// check if segment start with a ":"
		if segment != "" && segment[0:1] == ":" {
			log.Printf("Placeholder = %s\n", segment)
			// ToDo: store/return this value somewhere
			vals[segment[1:]] = segments[i]
			continue
		}

		if segments[i] != segment {
			return nil, false
		}
	}

	// nothing stopped us from matching, so it must be true
	return vals, true
}

// ServeHTTP
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("--- NEW REQUEST %s %s ---\n", r.Method, r.URL.Path)

	method := r.Method
	normPath := path.Clean(r.URL.Path)
	log.Printf("request: method=%#v\n", method)
	log.Printf("request: path1=%#v\n", r.URL.Path)
	log.Printf("request: path2=%#v\n", normPath)

	// if the original path ends in a slash
	if normPath != "/" {
		if strings.HasSuffix(r.URL.Path, "/") {
			normPath = normPath + "/"
		}
	}

	// if these paths differ, then redirect to the real one
	if normPath != r.URL.Path {
		http.Redirect(w, r, normPath, http.StatusFound)
		return
	}

	log.Printf("request: split=%#v\n", strings.Split(normPath, "/"))
	segments := strings.Split(normPath, "/")[1:]

	log.Printf("request: segments=%#v\n", segments)

	for i, route := range m.routes {
		log.Printf("--- Route(%d): %s /%s\n", i, route.Method, strings.Join(route.Segments, "/"))

		var vals map[string]string
		var matched bool
		// check to see if this is just a USE (and therefore, a prefix match)
		if route.Method == "USE" {
			log.Printf("This route is a USE, therefore, we just check the prefix\n")
			matched = isPrefixMatch(segments, &route)
		} else {
			log.Printf("This route is a GET/POST/PUT/etc, therefore, we match the entire segment length\n")
			vals, matched = isMatch(method, segments, &route)
			log.Printf("vals1=%#v\n", vals)

			ctx := context.WithValue(r.Context(), valsIdKey, vals)
			r = r.WithContext(ctx)
		}

		// if matched, then we have a non-nil 'vals' too, even if it contains no values
		if matched {
			log.Printf("matched, calling all middlewares in this route:")
			// loop over all middlewares in this route
			for j, middleware := range route.Middlewares {
				log.Printf(" - route #%d\n", j)

				// presume the middleware deals with this request fully (and doesn't call `next`)
				finished := true

				// call the middleware ... passing this `next` function to see if we want the next handler
				next := func(w http.ResponseWriter, r *http.Request) {
					log.Printf("*** next has been called\n")
					finished = false
				}

				nextHandlerFunc := http.HandlerFunc(next)

				// now call the middleware with our next
				log.Printf("before middleware\n")
				// log.Printf("vals2=%#v\n", vals)
				// ctx := context.WithValue(r.Context(), valsIdKey, vals)
				middleware(nextHandlerFunc).ServeHTTP(w, r)
				log.Printf("after middleware\n")

				// if we're finished, then just return
				if finished {
					return
				}
			}

			// finally, check if we have a handler and if so, call it - presume it is the last in the chain
			if route.Handler != nil {
				route.Handler.ServeHTTP(w, r)
				return
			}
		} else {
			log.Printf("NO match")
		}
	}

	// if we got through to here, then either nothing matched, or any middlewares did match but still called `next` and
	// hence there is no final route to deal with the request
	http.NotFound(w, r)
}

func Vals(r *http.Request) map[string]string {
	return r.Context().Value(valsIdKey).(map[string]string)
}
