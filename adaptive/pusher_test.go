package adaptive

import (
	"net/http"
	"testing"

	httpmiddlewarevet "github.com/fd/httpmiddlewarevet"
)

func ExampleHandler() {

	// Handler wraps around the static file HandlerFunc
	http.HandleFunc("/",
		Handler(http.FileServer(http.Dir("./cmd/static"))).ServeHTTP,
	)

	err := http.ListenAndServeTLS(":4430", "cmd/localhost.crt", "cmd/localhost.key", nil)
	if err != nil {
		panic(err)
	}
}

func Test(t *testing.T) {
	httpmiddlewarevet.Vet(t, func(h http.Handler) http.Handler {
		return Handler(h)
	})
}