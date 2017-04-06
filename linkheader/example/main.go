package main

import (
	"encoding/json"
	"net/http"

	"github.com/romainmenke/pusher/linkheader"
)

func main() {

	http.Handle("/",
		linkheader.Handler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				w.Header().Set("Cache-Control", "public, max-age=86400")

				http.FileServer(http.Dir("./example/static")).ServeHTTP(w, r)
			}),
			linkheader.PathOption("./linkheader/example/linkheaders.txt"),
		),
	)

	http.HandleFunc("/call.json",
		apiCall,
	)

	err := http.ListenAndServe(":8000", nil)
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
