package healthz

import (
	"context"
	"net/http"

	"github.com/eddycharly/generic-auth-server/pkg/server"
	"github.com/eddycharly/generic-auth-server/pkg/server/handlers"
)

func NewServer(addr, certFile, keyFile string) server.ServerFunc {
	return func(ctx context.Context) error {
		// create mux
		mux := http.NewServeMux()
		// register health check
		mux.Handle("GET /livez", handlers.Healthy(handlers.True))
		// register ready check
		mux.Handle("GET /readyz", handlers.Ready(handlers.True))
		// create server
		s := &http.Server{
			Addr:    addr,
			Handler: mux,
		}
		// run server
		return server.RunHttp(ctx, s, certFile, keyFile)
	}
}
