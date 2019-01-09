package transaction

import (
	"fmt"
	"net/http"
)

// Manager starts the process of a new transaction
func Manager(w http.ResponseWriter, r *http.Request) {
	fmt.Println("postei a mensagem no kafka")
}
