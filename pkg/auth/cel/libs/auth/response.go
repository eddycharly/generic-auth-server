package auth

import (
	"net/http"

	"github.com/eddycharly/generic-auth-server/pkg/auth/model"
	"github.com/google/cel-go/common/types"
)

type Response = model.Response

var ResponseType = types.NewObjectType("model.Response")

func makeResponse(code int, body []byte) Response {
	return Response{
		StatusCode: code,
		Body:       body,
		Header:     http.Header{},
	}
}
