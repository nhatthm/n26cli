package bootstrap

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/cucumber/godog"
	"github.com/nhatthm/clockdog"
	"github.com/nhatthm/n26cli/internal/app"
	"github.com/nhatthm/n26cli/internal/cli"
	"github.com/nhatthm/n26cli/internal/io"
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
)

var (
	uuidPattern        = `\b[0-9a-f]{8}\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\b[0-9a-f]{12}\b`
	credentialsService = "n26api.credentials.test"
	tokenService       = "n26api.token.test" // nolint: gosec
)

type appManager struct {
	fs    afero.Fs
	clock *clockdog.Clock
	stdio terminal.Stdio
	test  *testing.T

	homeDir string
	keys    map[string][]string
	baseURL string

	mu sync.Mutex
}

// WithStdio configures stdio for a given scenario.
func (m *appManager) WithStdio(_ *godog.Scenario, stdio terminal.Stdio) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.stdio = stdio
}

func (m *appManager) registerContext(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(m.init)

	ctx.AfterScenario(func(sc *godog.Scenario, _ error) {
		m.cleanup()
	})

	ctx.Step(`run command "([^"]*)"`, m.runCommandSimple)
	ctx.Step(`run command (\[[^\]]*\])`, m.runCommandArgs)
	ctx.Step(`create a file "([^"]+)" with content:`, m.createFileContent)
	ctx.Step(`create a credentials "([^"]+)" in keychain with content:`, m.createKeychainKey)
	ctx.Step(`delete token "([^"]+)" in keychain`, m.deleteKeychainToken)
	ctx.Step(`there is a file "([^"]+)" with content:`, m.hasFileContent)
	ctx.Step(`configured device is not "([^"]*)"`, m.isNotDevice)
	ctx.Step(`keychain has no credentials "([^"]*)"`, m.hasNoCredentialsInKeychain)
	ctx.Step(`keychain has username "([^"]*)" and password "([^"]*)"`, m.hasCredentialsInKeychain)
}

func (m *appManager) init(sc *godog.Scenario) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.homeDir = filepath.Join(m.test.TempDir(), "n26", sc.Id)
}

func (m *appManager) cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()

	err := m.fs.RemoveAll(m.homeDir)
	require.NoError(m.test, err)

	for svc, usernames := range m.keys {
		for _, username := range usernames {
			err := keyring.Delete(svc, username)
			if err != nil && !errors.Is(err, keyring.ErrNotFound) {
				require.NoError(m.test, err)
			}
		}
	}
}

func (m *appManager) config() (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigFile(filepath.Join(m.homeDir, ".n26", "config.toml"))

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return v, nil
}

func (m *appManager) device() (string, error) {
	cfg, err := m.config()
	if err != nil {
		return "", err
	}

	device := cfg.GetString("n26.device")

	if device == "" {
		return "", errors.New("device id is empty") // nolint: goerr113
	}

	return device, nil
}

func (m *appManager) registerKeyring(service, username string) {
	if _, ok := m.keys[service]; !ok {
		m.keys[service] = make([]string, 0)
	}

	m.keys[service] = append(m.keys[service], username)
}

func (m *appManager) keyring(service, username string) (string, error) {
	m.registerKeyring(service, username)

	return keyring.Get(service, username)
}

func (m *appManager) runCommandSimple(params string) error {
	return m.runCommand(strings.Split(params, " "))
}

func (m *appManager) runCommandArgs(params string) error {
	var args []string

	if err := json.Unmarshal([]byte(params), &args); err != nil {
		return err
	}

	return m.runCommand(args)
}

func (m *appManager) runCommand(args []string) (err error) {
	l := app.NewServiceLocator()

	l.ClockProvider = m.clock
	l.StdioProvider = io.Stdio(m.stdio.In, m.stdio.Out, m.stdio.Err)
	l.N26.BaseURL = m.baseURL
	l.N26.MFAWait = 5 * time.Millisecond
	l.N26.MFATimeout = time.Second

	cmd := cli.NewApp(l, m.homeDir)

	cmd.SetIn(m.stdio.In)
	cmd.SetOut(m.stdio.Out)
	cmd.SetArgs(args)

	pflag.CommandLine = nil

	doneCh := make(chan struct{})

	go func() {
		defer func() {
			if r := recover(); r != nil {
				_, _ = fmt.Fprintf(m.stdio.Out, "panic: %s\n", r)
			}

			close(doneCh)
		}()

		err = cmd.Execute()
	}()

	select {
	case <-time.After(time.Second):
		return errors.New("command timed out") // nolint: goerr113

	case <-doneCh:
		return
	}
}

func (m *appManager) createFileContent(filePath string, expectedBody *godog.DocString) error {
	filePath = filepath.Join(m.homeDir, filePath)

	if err := m.fs.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	return afero.WriteFile(m.fs, filePath, []byte(expectedBody.Content), os.ModePerm)
}

func (m *appManager) createKeychainKey(key string, content *godog.DocString) error {
	m.registerKeyring(credentialsService, key)

	return keyring.Set(credentialsService, key, content.Content)
}

func (m *appManager) deleteKeychainToken(key string) error {
	err := keyring.Delete(tokenService, key)
	if err != nil && !errors.Is(keyring.ErrNotFound, err) {
		return err
	}

	return nil
}

func (m *appManager) hasFileContent(filePath string, expectedBody *godog.DocString) error {
	filePath = filepath.Join(m.homeDir, filePath)

	_, err := m.fs.Stat(filePath)
	if err != nil {
		return err
	}

	actual, err := afero.ReadFile(m.fs, filePath)
	if err != nil {
		return err
	}

	expected := regexp.QuoteMeta(expectedBody.Content)
	expected = strings.ReplaceAll(expected, "<uuid>", uuidPattern)
	expected = fmt.Sprintf("^%s$", expected)

	t := t()
	assert.Regexp(t, expected, string(actual))

	return t.LastError()
}

func (m *appManager) isNotDevice(expected string) error {
	actual, err := m.device()
	if err != nil {
		return err
	}

	t := t()
	assert.NotEqual(t, expected, actual)

	return t.LastError()
}

func (m *appManager) hasNoCredentialsInKeychain(device string) error {
	_, err := m.keyring(credentialsService, device)
	if err == nil {
		return fmt.Errorf("(service=%q, key=%q) exists in keychain", credentialsService, device) // nolint: goerr113
	}

	t := t()
	assert.Equal(t, keyring.ErrNotFound, err)

	return t.LastError()
}

func (m *appManager) hasCredentialsInKeychain(username, password string) error {
	device, err := m.device()
	if err != nil {
		return err
	}

	data, err := m.keyring(credentialsService, device)
	if err != nil {
		return err
	}

	var actual map[string]string

	if err := json.Unmarshal([]byte(data), &actual); err != nil {
		return err
	}

	expected := map[string]string{
		"username": username,
		"password": password,
	}

	t := t()
	assert.Equal(t, expected, actual)

	return t.LastError()
}

func newAppManager(t *testing.T, baseURL string, clock *clockdog.Clock) *appManager { // nolint: thelper
	return &appManager{
		fs:      afero.NewOsFs(),
		clock:   clock,
		test:    t,
		keys:    make(map[string][]string),
		baseURL: baseURL,
	}
}
