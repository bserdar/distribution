package rhcert

import (
	"github.com/docker/distribution/context"
	"github.com/docker/distribution/registry/auth"
)

type accessController struct {
}

func (ac *accessController) Authorized(ctx context.Context, accessItems ...auth.Access) (context.Context, error) {
	req, err := context.GetRequest(ctx)
	if err != nil {
		return ctx, nil
	}
	context.GetLogger(ctx).Printf("Headers in auth: %v\n", req.Header)
	return ctx, nil
}

// newAccessController creates a new rhcert access controller with the given options
func newAccessController(options map[string]interface{}) (auth.AccessController, error) {
	return &accessController{}, nil
}

// Register the rhcert access controller
func init() {
	auth.Register("rhcert", auth.InitFunc(newAccessController))
}
