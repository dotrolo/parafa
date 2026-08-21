package api

import "net/http"

func Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", ping)
	return mux
}
