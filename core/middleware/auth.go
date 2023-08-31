package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Season-04/app-os/core/types"
)

const UserHeaderName = "X-AppOS-User"

func ToContext(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userStr := r.Header.Get(UserHeaderName)
		if userStr != "" {
			user := &types.User{}
			err := json.Unmarshal([]byte(userStr), user)
			if err == nil {
				r = r.WithContext(context.WithValue(r.Context(), UserHeaderName, user))
			}
		}

		next(w, r)
	}
}

func UserFromContext(ctx context.Context) *types.User {
	v := ctx.Value(UserHeaderName)
	if user, ok := v.(*types.User); ok {
		return user
	}

	return nil
}
