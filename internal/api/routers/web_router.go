package routers

import "net/http"

func webRouter(mux *http.ServeMux) {
	mux.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, "web/index.html")
	})
	mux.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
}
