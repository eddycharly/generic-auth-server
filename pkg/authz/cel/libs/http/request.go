package http

import (
	"net/http"

	"github.com/google/cel-go/common/types"
)

type Request = http.Request

var RequestType = types.NewObjectType("http.Request")
