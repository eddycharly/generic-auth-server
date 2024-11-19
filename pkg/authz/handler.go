package authz

import (
	"net/http"

	"github.com/eddycharly/generic-auth-server/pkg/authz/model"
)

type Handler = func(*http.Request) *model.Response
