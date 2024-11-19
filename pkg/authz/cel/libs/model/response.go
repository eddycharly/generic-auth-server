package model

import (
	"github.com/eddycharly/generic-auth-server/pkg/authz/model"
	"github.com/google/cel-go/common/types"
)

type Response = model.Response

var ResponseType = types.NewObjectType("model.Response")
