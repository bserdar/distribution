// Package chain implements an authentication scheme that uses
// multiple access controllers. The authorization request is passed to
// each controller one by one, until one of them accepts the request,
// or all fail.
//
// The configuration is as follows:
// auth:
//   chain:
//     htpasswd:
//       ... htpasswd configuratiun
//     token:
//       ... token configuration
//     etc.
package chain

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/docker/distribution/context"
	"github.com/docker/distribution/registry/auth"
)

// chainAccessController keeps the list of access controllers that
// will be called in order
type chainAccessController struct {
	controllers []auth.AccessController
}

// newAccessController parses the configuration options and creates a
// chain access controller. Each sub-controller is initialized using
// its section of the configuration
func newAccessController(options map[string]interface{}) (auth.AccessController, error) {
	var controller chainAccessController
	controller.controllers = make([]auth.AccessController, 0)
	for name, value := range options {
		subOptions, err := makeSubOptions(value)
		if err != nil {
			return nil, errors.New("Invalid options for " + name + ":" + err.Error())
		}
		acc, err := auth.GetAccessController(name, subOptions)
		if err != nil {
			return nil, err
		}
		controller.controllers = append(controller.controllers, acc)
	}
	return &controller, nil
}

// makeSubOptions copies the options for one of the sub-controllers
func makeSubOptions(value interface{}) (map[string]interface{}, error) {
	ret := make(map[string]interface{})
	if value != nil {
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Map {
			for _, key := range v.MapKeys() {
				ret[fmt.Sprint(key.Interface())] = v.MapIndex(key).Interface()
			}
		} else {
			return nil, errors.New(fmt.Sprintf("%v", value))
		}
	}
	return ret, nil
}

// Authorized calls each controller in order until one of them accepts, or all of them fail
func (ac *chainAccessController) Authorized(ctx context.Context, accessRecords ...auth.Access) (context.Context, error) {
	var composite ErrComposite

	for _, c := range ac.controllers {
		ctx, err := c.Authorized(ctx, accessRecords...)
		if err == nil {
			// Successful authorization, return
			return ctx, nil
		} else {
			// Add the error to composite error. This will be thrown
			// away if one of the subsequent controllers succeed
			composite.add(err)
		}
	}
	// All controllers failed, return error
	return ctx, composite
}

// init registers the chain auth backend.
func init() {
	auth.Register("chain", auth.InitFunc(newAccessController))
}

// ErrComposite stores the list of errors returned from controllers
type ErrComposite struct {
	components []error
}

func (e *ErrComposite) add(err error) {
	if e.components == nil {
		e.components = make([]error, 0)
	}
	e.components = append(e.components, err)
}

func (e ErrComposite) Error() string {
	return fmt.Sprintf("%v", e.components)
}
