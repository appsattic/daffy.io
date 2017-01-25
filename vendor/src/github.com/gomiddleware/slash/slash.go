// Copyright (c) 2016, Apps Attic Ltd (https://appsattic.com/) <chilts@appsattic.com>.

// Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby
// granted, provided that the above copyright notice and this permission notice appear in all copies.

// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
// INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN
// AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
// PERFORMANCE OF THIS SOFTWARE.

package slash

import (
	"net/http"
	"strings"
)

//  Add returns middleware which checks the request's path and redirects to the same path with a trailing slash if it
//  doesn't already have one.
func Add(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			// already has a trailing slash
			next.ServeHTTP(w, r)
			return
		}

		// redirect to the new URL
		r.URL.Path = r.URL.Path + "/"
		http.Redirect(w, r, r.URL.String(), http.StatusFound)
	}

	return http.HandlerFunc(fn)
}

//  Remove returns middleware which checks the request's path and redirects to the same path but minus any trailing slash if
//  it currently has one. The "/" path is ignored by this middleware.
func Remove(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			// ignore plain "/"
			next.ServeHTTP(w, r)
			return
		}

		if !strings.HasSuffix(r.URL.Path, "/") {
			// no trailing slash already
			next.ServeHTTP(w, r)
			return
		}

		// redirect to the new URL
		r.URL.Path = strings.TrimRight(r.URL.Path, "/")
		http.Redirect(w, r, r.URL.String(), http.StatusFound)
	}

	return http.HandlerFunc(fn)
}
