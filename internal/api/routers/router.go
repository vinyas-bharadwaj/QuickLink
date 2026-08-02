package routers

import "net/http"

func Router() *http.ServeMux {
	mux := http.NewServeMux()
	shortenRourter(mux)
	webRouter(mux)
	return mux
}
