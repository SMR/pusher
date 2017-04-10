// +build go1.8

package link

import (
	"net/http"
	"testing"

	httpmiddlewarevet "github.com/fd/httpmiddlewarevet"
)

func TestMiddleware(t *testing.T) {
	httpmiddlewarevet.Vet(t, func(h http.Handler) http.Handler {
		return Handler(h)
	})
}
