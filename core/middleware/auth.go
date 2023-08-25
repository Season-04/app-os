package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const UserIDHeaderName = "X-AppOS-User-ID"

func ToContext(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.Header.Get(UserIDHeaderName)

		if userIDStr != "" {
			userID, err := strconv.ParseUint(userIDStr, 10, 32)
			if err != nil {
				log.Println(fmt.Errorf("invalid user ID, %w", err))
			} else {
				r = r.WithContext(context.WithValue(r.Context(), UserIDHeaderName, uint32(userID)))
			}
		}

		next(w, r)
	}
}

func UserIDFromContext(ctx context.Context) uint32 {
	v := ctx.Value(UserIDHeaderName)
	if userID, ok := v.(uint32); ok {
		return userID
	}

	return 0
}
