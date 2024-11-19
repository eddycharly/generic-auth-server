package cel

import (
	"github.com/eddycharly/generic-auth-server/pkg/authz/cel/libs/http"
	"github.com/eddycharly/generic-auth-server/pkg/authz/cel/libs/jwt"
	"github.com/eddycharly/generic-auth-server/pkg/authz/cel/libs/model"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/ext"
	"k8s.io/apiserver/pkg/cel/library"
)

func NewEnv() (*cel.Env, error) {
	// create new cel env
	return cel.NewEnv(
		// configure env
		cel.HomogeneousAggregateLiterals(),
		cel.EagerlyValidateDeclarations(true),
		cel.DefaultUTCTimeZone(true),
		cel.CrossTypeNumericComparisons(true),
		// register common libs
		cel.OptionalTypes(),
		ext.Bindings(),
		ext.Encoders(),
		ext.Lists(),
		ext.Math(),
		ext.Protos(),
		ext.Sets(),
		ext.Strings(),
		// register kubernetes libs
		library.CIDR(),
		library.Format(),
		library.IP(),
		library.Lists(),
		library.Regex(),
		library.URLs(),
		// register our libs
		http.Lib(),
		model.Lib(),
		jwt.Lib(),
	)
}
