package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/eddycharly/generic-auth-server/pkg/policy"
	"github.com/eddycharly/generic-auth-server/pkg/server"
)

func NewHttpServer(addr, certFile, keyFile string, provider policy.Provider) server.ServerFunc {
	return func(ctx context.Context) error {
		// create mux
		mux := http.NewServeMux()
		// register auth
		mux.HandleFunc("GET /auth", func(w http.ResponseWriter, r *http.Request) {
			// fetch compiled policies
			policies, err := provider.CompiledPolicies(ctx)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println(err)
				return
			}
			// iterate over policies
			for _, policy := range policies {
				// execute policy
				response, err := policy(r)
				// return error if any
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Println(err)
					return
				}
				// if the reponse returned by the policy evaluation was not nil, return
				if response != nil {
					for k, v := range response.Header {
						w.Header()[k] = append(w.Header()[k], v...)
					}
					w.WriteHeader(response.StatusCode)
					if _, err := w.Write(response.Body); err != nil {
						fmt.Println(err)
					}
					return
				}
			}
			// we didn't have a response
			// TODO: default
			w.WriteHeader(http.StatusOK)
		})
		// create server
		s := &http.Server{
			Addr:    addr,
			Handler: mux,
		}
		// run server
		return server.RunHttp(ctx, s, certFile, keyFile)
	}
}
