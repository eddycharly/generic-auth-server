package auth

import (
	"net/http"

	"github.com/eddycharly/generic-auth-server/pkg/auth/model"
)

type Handler = func(*http.Request) *model.Response
