package auth

import (
	"log"
	"net/http"
)

type server struct {
}

func (s *server) check(w http.ResponseWriter, r *http.Request) {
	log.Println("check")
	w.WriteHeader(http.StatusOK)
}

func RunHTTPServer() {
	s := &server{}

	mux := http.NewServeMux()

	mux.HandleFunc("/auth/check", s.check)

	log.Printf("Listening HTTP at %v", 3000)

	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
