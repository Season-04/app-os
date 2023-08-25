package auth

import (
	"log"
	"net/http"
	"strconv"

	"github.com/staugaard/app-os/core/internal/pb"
)

type server struct {
	users pb.UsersServiceServer
}

func (s *server) check(w http.ResponseWriter, r *http.Request) {
	log.Println("check")
	//log.Println(r.Header)

	resp, err := s.users.GetById(r.Context(), &pb.GetUserByIdRequest{Id: 1})
	if err != nil {
		log.Println("Failed to get user", err, "at path", r.Header.Get("X-Forwarded-Uri"))
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
	} else {
		log.Println("Got user", resp.User.Id, "at path", r.Header.Get("X-Forwarded-Uri"))
		w.Header().Set("X-AppOS-User-ID", strconv.FormatUint(uint64(resp.User.Id), 10))
		w.WriteHeader(http.StatusOK)
	}
}

func RunHTTPServer(usersServer pb.UsersServiceServer) {
	s := &server{
		users: usersServer,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/auth/check", s.check)

	log.Printf("Listening HTTP at %v", 3000)

	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
