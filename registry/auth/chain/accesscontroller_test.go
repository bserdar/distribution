package chain

import (
	"errors"
	"strings"
	"testing"

	"github.com/docker/distribution/configuration"
	"github.com/docker/distribution/context"
	"github.com/docker/distribution/registry/auth"
	_ "github.com/docker/distribution/registry/auth/htpasswd"
	_ "github.com/docker/distribution/registry/auth/silly"
)

const (
	testCfg1 = "version: 0.1\n" +
		"auth:\n" +
		"  chain:\n" +
		"    htpasswd:\n" +
		"       realm: htest-realm\n" +
		"       path: htest-path\n" +
		"    silly:\n" +
		"       realm: test-realm\n" +
		"       service: test-service\n" +
		"storage:\n" +
		"  cache:\n" +
		"    blobdescriptor: redis\n" +
		"  filesystem:\n" +
		"    rootdirectory: /var/lib/registry-cache\n"

	testCfg2 = "version: 0.1\n" +
		"auth:\n" +
		"  chain:\n" +
		"    test1:\n" +
		"    test2:\n" +
		"storage:\n" +
		"  cache:\n" +
		"    blobdescriptor: redis\n" +
		"  filesystem:\n" +
		"    rootdirectory: /var/lib/registry-cache\n"
)

type ChainController interface {
	GetControllers() []auth.AccessController
}

func (c chainAccessController) GetControllers() []auth.AccessController { return c.controllers }

func TestCfg(t *testing.T) {
	cfg, err := configuration.Parse(strings.NewReader(testCfg1))
	if err != nil {
		t.Errorf("Cannot parse configuration " + err.Error())
	}
	ctr, err := auth.GetAccessController("chain", cfg.Auth.Parameters())
	if err != nil {
		t.Errorf("Cannot get chain controller:" + err.Error())
	}
	ch, ok := ctr.(ChainController)
	if !ok {
		t.Errorf("Wrong kind of controller")
	}
	controllers := ch.GetControllers()
	if len(controllers) != 2 {
		t.Errorf("Expecting 2 controllers")
	}
}

type testController struct {
	returnError bool
	called      bool
}

func (t *testController) Authorized(ctx context.Context, accessRecords ...auth.Access) (context.Context, error) {
	t.called = true
	if t.returnError {
		return ctx, errors.New("Error")
	} else {
		return context.WithValue(ctx, "auth", true), nil
	}
}

func TestAuth(t *testing.T) {

	var t1, t2 testController

	auth.Register("test1", func(options map[string]interface{}) (auth.AccessController, error) {
		return &t1, nil
	})
	auth.Register("test2", func(options map[string]interface{}) (auth.AccessController, error) {
		return &t2, nil
	})

	cfg, err := configuration.Parse(strings.NewReader(testCfg2))
	if err != nil {
		t.Errorf("Cannot parse configuration " + err.Error())
	}
	ctr, err := auth.GetAccessController("chain", cfg.Auth.Parameters())
	if err != nil {
		t.Errorf("Cannot get chain controller:" + err.Error())
	}

	// Both fail
	t1.returnError, t2.returnError = true, true
	t1.called, t2.called = false, false
	ctx, err := ctr.Authorized(context.Background())
	if err == nil || ctx.Value("auth") != nil || !t1.called || !t2.called {
		t.Errorf("Expecting fail err:%v auth:%v t1: %v t2:%v", err, ctx.Value("auth"), t1.called, t2.called)
	}

	// t1 fails, t2 authorizes
	t1.returnError, t2.returnError = true, false
	t1.called, t2.called = false, false
	ctx, err = ctr.Authorized(context.Background())
	if err != nil || ctx.Value("auth").(bool) != true || !t1.called || !t2.called {
		t.Errorf("Expecting success err:%v auth:%v t1: %v t2:%v", err, ctx.Value("auth"), t1.called, t2.called)
	}

	// t1 authorizes, t2 isn't called
	t1.returnError, t2.returnError = false, true
	t1.called, t2.called = false, false
	ctx, err = ctr.Authorized(context.Background())
	if err != nil || ctx.Value("auth").(bool) != true || !t1.called || t2.called {
		t.Errorf("Expecting success err:%v auth:%v t1: %v t2:%v", err, ctx.Value("auth"), t1.called, t2.called)
	}
}
