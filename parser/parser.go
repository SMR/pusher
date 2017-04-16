package parser

import (
	"net/http"

	"github.com/romainmenke/pusher/common"
	"golang.org/x/net/html"
)

// Handler wraps an http.Handler reading the response body and setting Link Headers or generating Pushes
func Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get a responseWriter from the sync.Pool.
		var rw = getResponseWriter(w, r)
		// defer close the responseWriter.
		// This returns it to the sync.Pool and zeroes all values and pointers.
		defer rw.close()

		// handle.
		handler.ServeHTTP(rw, r)

	})
}

func (w *responseWriter) extractLinks() []common.Preloadable {
	links := make(map[common.Preloadable]struct{})
	preloads := make(map[string]struct{})

	contentType := http.DetectContentType(w.body.Bytes())
	if contentType != "text/html; charset=utf-8" {
		return nil
	}

	path := w.request.URL.RequestURI()

	z := html.NewTokenizer(w.body)

TOKENIZER:
	for {
		tt := z.Next()

		var asset common.Preloadable
		var preload string

		switch tt {
		case html.ErrorToken:
			// End of the document, we're done
			break TOKENIZER
		case html.SelfClosingTagToken:

			t := z.Token()
			asset, preload = parseToken(t, path)

			if asset != nil {
				if _, found := preloads[asset.Path()]; !found {
					links[asset] = struct{}{}
					asset = nil
				}
			} else if preload != "" {
				preloads[preload] = struct{}{}
				preload = ""
			}

		case html.StartTagToken:

			t := z.Token()
			asset, preload = parseToken(t, path)

			if asset != nil {
				if _, found := preloads[asset.Path()]; !found {
					links[asset] = struct{}{}
					asset = nil
				}
			} else if preload != "" {
				preloads[preload] = struct{}{}
				preload = ""
			}

		}
	}

	linkSlice := make([]common.Preloadable, len(links))

	index := 0
	for key := range links {
		if _, found := preloads[key.Path()]; found {
			continue
		}
		if index >= len(linkSlice) {
			break
		}
		linkSlice[index] = key
		index++
	}

	return linkSlice[:index]
}

const (
	hrefStr   = "href"
	imgStr    = "img"
	linkStr   = "link"
	relStr    = "rel"
	scriptStr = "script"
	srcStr    = "src"
)

func parseToken(t html.Token, path string) (common.Preloadable, string) {

	var (
		asset     common.Preloadable
		isPreload bool
	)

	switch t.Data {
	case linkStr:

		for _, attr := range t.Attr {
			switch attr.Key {
			case relStr:
				if attr.Val == common.Preload {
					isPreload = true
				}
			case common.NoPush:
				return nil, ""
			case hrefStr:
				if common.IsAbsolute(attr.Val) || attr.Val == path {
					return nil, ""
				}
				asset = common.CSS(attr.Val)
			}
		}

	case scriptStr:

		for _, attr := range t.Attr {
			switch attr.Key {
			case relStr:
				if attr.Val == common.Preload {
					return nil, ""
				}
			case common.NoPush:
				return nil, ""
			case srcStr:
				if common.IsAbsolute(attr.Val) || attr.Val == path {
					return nil, ""
				}
				asset = common.JS(attr.Val)
			}
		}

	case imgStr:

		for _, attr := range t.Attr {
			switch attr.Key {
			case relStr:
				if attr.Val == common.Preload {
					return nil, ""
				}
			case common.NoPush:
				return nil, ""
			case srcStr:
				if common.IsAbsolute(attr.Val) || attr.Val == path {
					return nil, ""
				}
				asset = common.Img(attr.Val)
			}
		}
	}

	if isPreload {
		return nil, asset.Path()
	}

	return asset, ""
}
