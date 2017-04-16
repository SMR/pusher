// +build go1.8

package parser

import (
	"net/http"

	"github.com/romainmenke/pusher/common"
)

const (
	protoHTTP11    = "HTTP/1.1"
	protoHTTP11TLS = "HTTP/1.1+TLS"
	protoHTTP20    = "HTTP/2.0"
)

// Handler wraps an http.Handler reading the response body and setting Link Headers or generating Pushes
func Handler(handler http.Handler, options ...Option) http.Handler {

	s := &settings{}
	for _, o := range options {
		o(s)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if s.withCache {
			preloads := getFromCache(r.URL.RequestURI())
			if preloads != nil {
				defer handler.ServeHTTP(w, r)

				if pusher, ok := w.(http.Pusher); ok && r.Header.Get(common.XForwardedFor) == "" {
					for _, link := range preloads {
						pusher.Push(link.Path(), &http.PushOptions{
							Header: r.Header,
						})
					}
				} else {
					for _, link := range preloads {
						w.Header().Add(common.Link, link.LinkHeader())
					}
				}

				return
			}
		}

		// Get a responseWriter from the sync.Pool.
		var rw = getResponseWriter(s, w, r)
		// defer close the responseWriter.
		// This returns it to the sync.Pool and zeroes all values and pointers.
		defer rw.close()

		var protoW http.ResponseWriter
		switch r.Proto {
		case protoHTTP11:
			protoW = &responseWriterHTTP11{
				responseWriter: rw,
			}
		case protoHTTP11TLS:
			protoW = &responseWriterHTTP11TLS{
				responseWriter: rw,
			}
		case protoHTTP20:
			protoW = &responseWriterHTTP2{
				responseWriter: rw,
			}
		}

		// handle.
		handler.ServeHTTP(protoW, r)

	})
}
