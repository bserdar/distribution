package termsauth

import (
	"net/url"

	"github.com/docker/distribution/context"
	"github.com/docker/distribution/registry/auth"
)

type termsAuthData struct {
	// List of terms IDs required to grant access
	requiredTerms []int64
	// The name of the header field containing account id
	requestHeaderField string
	// Terms service URL
	termsSvc *url.URL
}

func newAccessController(options map[string]interface{}) (auth.AccessController, error) {
}

func init() {
	auth.Register("termsauth", auth.InitFunc(newAccessController))
}
