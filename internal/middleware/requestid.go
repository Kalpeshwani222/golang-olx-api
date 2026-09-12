package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ctxKey int
const (
	requestIDKey ctxKey = iota
)

const ( requestId = "X-Request-ID" )

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(requestId)
		
		if id == "" {
			id = uuid.NewString()
		}

		w.Header().Add(requestId,id)

		//requestIDKey:Create a custom key to avoid conflicts with context values from other packages
		ctx := context.WithValue(r.Context(),requestIDKey,id)

		//calling the next with the context and passing the id
		next.ServeHTTP(w,r.WithContext(ctx)) 
	}) 
}


func RequestIdFromContext (ctx context.Context) string{
	requestId := ctx.Value(requestIDKey).(string)
	return requestId
}