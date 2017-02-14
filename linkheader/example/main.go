package main

import (
	"encoding/json"
	"net/http"

	"github.com/romainmenke/pusher/linkheader"
)

func main() {

	linkHeaderMux := linkheader.New()
	err := linkHeaderMux.Read("./linkheader/example/linkheaders.txt")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			w.Header().Set("Cache-Control", "public, max-age=86400")
			linkHeaderMux.SetLinkHeaders(w, r)

			http.FileServer(http.Dir("./example/static")).ServeHTTP(w, r)
		}),
	)

	http.HandleFunc("/call.json",
		apiCall,
	)

	err = http.ListenAndServe(":4430", nil)
	if err != nil {
		panic(err)
	}

}

func apiCall(w http.ResponseWriter, r *http.Request) {
	a := struct {
		Some string
	}{Some: "Remote Data"}
	json.NewEncoder(w).Encode(a)
}