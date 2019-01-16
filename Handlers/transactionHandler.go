package handlers

import (
	debuggin "kernel-concierge/Debuggin"
	"log"
	"net/http"
	"time"
)

// StreamHandler wraps requests for new transactions
func StreamHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "POST" {
			http.Error(w, "Invalid request method.", http.StatusMethodNotAllowed)
			return
		}

		log.Println(debuggin.Tracer(), "Logged connection from: ", r.RemoteAddr)
		log.Println(debuggin.Tracer(), "starting at: ", time.Now())

		next.ServeHTTP(w, r)

		log.Println(debuggin.Tracer(), "finished at: ", time.Now())
	}
}
