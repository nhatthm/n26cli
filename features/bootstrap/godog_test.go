package bootstrap

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/godogx/aferosteps"
	"github.com/nhatthm/clockdog"
	"github.com/nhatthm/consoledog"
	"github.com/nhatthm/n26godog"
	"github.com/nhatthm/surveydog"
	"github.com/nhatthm/surveyexpect"
	"github.com/stretchr/testify/assert"
)

// Used by init().
//
//nolint:gochecknoglobals
var (
	runGoDogTests bool

	out = new(surveyexpect.Buffer)
	opt = godog.Options{
		Strict: true,
		Output: out,
	}
)

// This has to run on init to define -godog flag, otherwise "undefined flag" error happens.
//
//nolint:gochecknoinits
func init() {
	flag.BoolVar(&runGoDogTests, "godog", false, "Set this flag is you want to run godog BDD tests")
	godog.BindFlags("godog.", flag.CommandLine, &opt) // nolint: staticcheck
}

func TestIntegration(t *testing.T) {
	if !runGoDogTests {
		t.Skip(`Missing "-godog" flag, skipping integration test.`)
	}

	server := n26godog.New(t)
	clock := clockdog.New()
	am := newAppManager(t, server.URL(), clock)
	fsManager := aferosteps.NewManager()
	console := consoledog.New(t)
	survey := surveydog.New(t).
		WithConsole(console).
		WithStarter(am.WithStdio)

	RunSuite(t, "..", func(_ *testing.T, ctx *godog.ScenarioContext) {
		am.registerContext(ctx)
		clock.RegisterContext(ctx)
		fsManager.RegisterContext(t, ctx)
		console.RegisterContext(ctx)
		survey.RegisterContext(ctx)
		server.RegisterContext(ctx)
	})
}

func RunSuite(t *testing.T, path string, featureContext func(t *testing.T, ctx *godog.ScenarioContext)) {
	t.Helper()

	flag.Parse()

	if opt.Randomize == 0 {
		opt.Randomize = rand.Int63()
	}

	var paths []string

	files, err := ioutil.ReadDir(filepath.Clean(path))
	assert.NoError(t, err)

	paths = make([]string, 0, len(files))

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".feature") {
			paths = append(paths, filepath.Join(path, f.Name()))
		}
	}

	for _, path := range paths {
		path := path

		t.Run(path, func(t *testing.T) {
			opt.Paths = []string{path}
			suite := godog.TestSuite{
				Name:                 "Integration",
				TestSuiteInitializer: nil,
				ScenarioInitializer: func(s *godog.ScenarioContext) {
					featureContext(t, s)
				},
				Options: &opt,
			}
			status := suite.Run()

			if status != 0 {
				fmt.Println(out.String())
				assert.Fail(t, "one or more scenarios failed in feature: "+path)
			}
		})
	}
}
