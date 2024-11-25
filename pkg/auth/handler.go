package auth

import (
	"fmt"
	"net/http"

	"github.com/eddycharly/generic-auth-server/pkg/policy"
)

func Handler(provider policy.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// fetch compiled policies
		policies, err := provider.CompiledPolicies(r.Context())
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
	}
}
