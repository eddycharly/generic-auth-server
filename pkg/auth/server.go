package auth

import (
	"context"
	"net/http"

	"github.com/eddycharly/generic-auth-server/pkg/policy"
	"github.com/eddycharly/generic-auth-server/pkg/server"
)

func NewHttpServer(addr, certFile, keyFile string, provider policy.Provider) server.ServerFunc {
	return func(ctx context.Context) error {
		// create mux
		mux := http.NewServeMux()
		// register auth
		mux.HandleFunc("GET /auth", Handler(provider))
		// create server
		s := &http.Server{
			Addr:    addr,
			Handler: mux,
		}
		// run server
		return server.RunHttp(ctx, s, certFile, keyFile)
	}
}
