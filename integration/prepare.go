package integration

import (
	"donation-mgmt/src/config"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/net/context"
)

const (
	EnvironmentKey string = "APP_ENVIRONMENT"
)

var smokeTestRegex = regexp.MustCompile(`^Test_Smoke_.+$`)
var activeTests atomic.Int32
var initSync sync.Once
var app *IntegrationApp

func Prepare(t *testing.T) (*IntegrationApp, bool) {
	t.Helper()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Integration test preparation panicked: %v", r)
			t.FailNow()
		}
	}()

	if !isRunningIntTests() && !isRunningDirectTest() {
		t.Skip("Skipping integration tests")
		return nil, false
	}

	isSmokeTest := smokeTestRegex.MatchString(t.Name())
	if testing.Short() && !isSmokeTest {
		t.Skip("Test is not a smoke test. Smoke tests begin with 'Test_Smoke_'")
		return nil, false
	}

	activeTests.Add(1)

	t.Cleanup(func() {
		next := activeTests.Add(-1)

		if next == 0 {
			fmt.Println("Last test in package finished. Cleaning up")

			err := app.Stop()
			if err != nil {
				t.Errorf("Failed to stop integration app: %v", err)
				return
			}
		}
	})

	t.Parallel()

	initSync.Do(func() {
		var err error
		app = NewIntegrationApp()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = app.Start(ctx)
		if err != nil {
			fmt.Printf("Could not start integration app: %v\n", err)
			app = nil
		}
	})

	if app == nil {
		t.Errorf("Integration app was not initialized. Cannot run test")
		t.FailNow()
		return nil, false
	}

	return app, true
}

func isRunningIntTests() bool {
	if env, ok := os.LookupEnv(EnvironmentKey); ok {
		return strings.EqualFold(env, string(config.IntegrationTest))
	}

	return false
}

func isRunningDirectTest() bool {
	runFlag := flag.Lookup("test.run")
	return runFlag != nil && runFlag.Value.String() != ""
}
