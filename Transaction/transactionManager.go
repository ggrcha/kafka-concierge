package transaction

import (
	debuggin "kernel-concierge/Debuggin"
	"log"
	"net/http"
)

// Manager starts the process of a new transaction
func Manager(w http.ResponseWriter, r *http.Request) {
	log.Println(debuggin.Tracer(), "postei a mensagem no kafka")
}
