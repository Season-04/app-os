package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/Season-04/appos/core/internal/pb"
	"github.com/Season-04/appos/core/types"
	"github.com/golang-jwt/jwt/v5"
)

type Server struct {
	users pb.UsersServiceServer
}

func NewServer(users pb.UsersServiceServer) *Server {
	return &Server{users: users}
}

func (s *Server) check(w http.ResponseWriter, r *http.Request) {
	log.Println("check")
	//log.Println(r.Header)

	authCookie, err := r.Cookie("appos-auth")
	if err != nil {
		log.Println("No auth cookie")
		w.WriteHeader(http.StatusOK)
		return
	}

	token, err := jwt.ParseWithClaims(authCookie.Value, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))

	if err != nil {
		log.Println("Invalid JWT token", err)
		w.WriteHeader(http.StatusOK)
		return
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		resp, err := s.users.GetUserById(r.Context(), &pb.GetUserByIdRequest{Id: claims.UserID})
		if err != nil {
			log.Println("Failed to get user", err, "at path", r.Header.Get("X-Forwarded-Uri"))
			w.WriteHeader(http.StatusOK)
		} else {
			user := types.User{
				ID:           resp.User.Id,
				Name:         resp.User.Name,
				EmailAddress: resp.User.EmailAddress,
			}

			switch resp.User.Role {
			case pb.UserRole_USER_ROLE_ADMIN:
				user.Role = types.UserRoleAdmin
			default:
				user.Role = types.UserRoleUser
			}

			if resp.User.LastSeenAt != nil {
				t := resp.User.LastSeenAt.AsTime()
				user.LastSeenAt = &t
			}

			jsonUser, _ := json.Marshal(user)
			w.Header().Set("X-AppOS-User", string(jsonUser))
			w.WriteHeader(http.StatusOK)
		}
	} else {
		log.Println("Invalid JWT token")
		w.WriteHeader(http.StatusOK)
		return
	}
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	resp, err := s.users.GetUserByEmailAndPassword(
		r.Context(),
		&pb.GetUserByEmailAndPasswordRequest{
			EmailAddress: r.FormValue("email"),
			Password:     r.FormValue("password"),
		},
	)

	if err != nil {
		loginError(w, r, err)
		return
	}

	expiry := time.Now().Add(24 * time.Hour)
	claims := JWTClaims{
		UserID: resp.User.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "appos",
			ExpiresAt: jwt.NewNumericDate(expiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		loginError(w, r, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Path:     "/",
		Name:     "appos-auth",
		Value:    tokenString,
		Expires:  expiry,
		Secure:   r.URL.Scheme == "https",
		HttpOnly: true,
	})

	relativeReturnToUrl, err := url.Parse(r.FormValue("returnTo"))
	if err != nil || r.FormValue("returnTo") == "" {
		relativeReturnToUrl, _ = url.Parse("/")
	}
	returnToUrl := r.URL.ResolveReference(relativeReturnToUrl)
	log.Println("logged in", resp.User.Id)
	w.Header().Set("Location", returnToUrl.String())
	w.WriteHeader(http.StatusFound)
}

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Path:     "/",
		Name:     "appos-auth",
		Value:    "",
		Expires:  time.UnixMilli(0),
		Secure:   r.URL.Scheme == "https",
		HttpOnly: true,
	})

	relativeReturnToUrl, err := url.Parse(r.FormValue("returnTo"))
	if err != nil || r.FormValue("returnTo") == "" {
		relativeReturnToUrl, _ = url.Parse("/")
	}
	returnToUrl := r.URL.ResolveReference(relativeReturnToUrl)
	w.Header().Set("Location", returnToUrl.String())
	w.WriteHeader(http.StatusFound)
}

func RunHTTPServer(usersServer pb.UsersServiceServer) {
	s := NewServer(usersServer)

	mux := http.NewServeMux()

	mux.HandleFunc("/auth/check", s.check)

	log.Printf("Listening HTTP at %v", 3000)

	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}

func loginError(w http.ResponseWriter, r *http.Request, err error) {
	relativeLoginUrl, _ := url.Parse("/login")
	loginUrl := r.URL.ResolveReference(relativeLoginUrl)

	query := loginUrl.Query()
	query.Set("email", r.FormValue("email"))
	query.Set("error", err.Error())
	loginUrl.RawQuery = query.Encode()

	w.Header().Set("Location", loginUrl.String())
	w.WriteHeader(http.StatusFound)
}

type JWTClaims struct {
	UserID uint32 `json:"userId"`
	jwt.RegisteredClaims
}
