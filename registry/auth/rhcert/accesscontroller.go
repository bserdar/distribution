package rhcert

import (
	"github.com/docker/distribution/registry/auth"
)

// newAccessController creates a new rhcert access controller with the given options
func newAccessController(options map[string]interface{}) (auth.AccessController, error) {
}

// Register the rhcert access controller
func init() {
	auth.register("rhcert", auth.InitFunc(newAccessController))
}
