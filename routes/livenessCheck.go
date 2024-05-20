package routes

import (
	"fmt"
	"net/http"
)

func handleLivenessCheck(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("I'm alive!"))

	fmt.Println("I'm alive, but from the console xD!")

}
