package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/cucumber/messages-go/v16"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
)

type apiFeature struct {
	resp     *http.Response
	migrator *migrator
}

type DbConfig struct {
	dsn  string
	host string
	port string
}

type migrator struct {
	config *DbConfig
}

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "calendar_config.toml", "Path to configuration file")
}

func (a *apiFeature) resetResponse(*godog.Scenario) {
	a.resp = nil
}

func (d *migrator) loadDump(ctx context.Context) {
	dumpFile := filepath.Join("dump", "dump.sql")

	if _, err := os.Stat(dumpFile); err != nil {
		fmt.Println("Dump file does not exist")
		os.Exit(1)
	}

	conn, err := pgx.Connect(ctx, d.config.dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	command, err := os.ReadFile(dumpFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read file: %v\n", err)
		os.Exit(1)
	}

	_, err = conn.Exec(ctx, string(command))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to persist to DB: %v\n", err)
		os.Exit(1)
	}
}

func (d *migrator) truncate(ctx context.Context) {
	conn, err := pgx.Connect(ctx, d.config.dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	if _, err := conn.Exec(ctx, `TRUNCATE TABLE events`); err != nil {
		fmt.Fprintf(os.Stderr, "can't to truncate table: %v\n", err)
		os.Exit(1)
	}

}

func (a *apiFeature) iSendRequestTo(method, endpoint string) error {
	return a.iSendRequestToWithJsonBody(method, endpoint, nil)
}

func (a *apiFeature) iWaitS(delay int) {
	time.Sleep(time.Duration(delay) * time.Second)
}

func (a *apiFeature) theResponseCodeShouldBe(code int) error {
	if code != a.resp.StatusCode {
		var actual interface{}
		err := json.NewDecoder(a.resp.Body).Decode(&actual)
		if err != nil {
			return err
		}
		fmt.Println(actual)

		return fmt.Errorf("expected response code to be: %d, but actual is: %d", code, a.resp.StatusCode)
	}

	return nil
}

func (a *apiFeature) iSendRequestToWithJsonBody(method, endpoint string, message *messages.PickleDocString) error {
	var body io.Reader
	if message != nil {
		body = strings.NewReader(message.Content)
	}

	endpointWithPrefix := "http://" + net.JoinHostPort(a.migrator.config.host, a.migrator.config.port) + endpoint
	req, err := http.NewRequest(method, endpointWithPrefix, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	a.resp, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (a *apiFeature) theResponseBodyHasItems(count int) (err error) {
	var actual map[string]interface{}
	err = json.NewDecoder(a.resp.Body).Decode(&actual)

	var i int
	for _, t := range actual {
		i += len(t.([]interface{}))
	}

	if i != count {
		return fmt.Errorf("expected response body count to be: %d, but actual is: %d", count, i)
	}

	return nil
}

func (a *apiFeature) theResponseShouldBeContainsKey(key string) (err error) {
	var actual map[string]interface{}
	err = json.NewDecoder(a.resp.Body).Decode(&actual)

	if _, ok := actual[key]; !ok {
		return fmt.Errorf("expected response body should be contain: %s", key)
	}

	return nil
}

func (a *apiFeature) theResponseShouldMatchJSON(resp *godog.DocString) error {
	var expected, actual interface{}

	err := json.NewDecoder(strings.NewReader(resp.Content)).Decode(&expected)
	if err != nil {
		return err
	}

	err = json.NewDecoder(a.resp.Body).Decode(&actual)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("expected JSON does not match actual, %v vs. %v", expected, actual)
	}

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	cfg, err := newConfig(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	migrator := &migrator{
		cfg,
	}

	api := &apiFeature{
		migrator: migrator,
	}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		api.migrator.truncate(ctx)
		api.migrator.loadDump(ctx)
		api.resetResponse(sc)
		return ctx, nil
	})

	ctx.Step(`^I send "(GET|POST|PUT|DELETE)" request to "([^"]*)"$`, api.iSendRequestTo)
	ctx.Step(`^the response code should be (\d+)$`, api.theResponseCodeShouldBe)
	ctx.Step(`^the response body should be contains key: "(\w+)"$`, api.theResponseShouldBeContainsKey)
	ctx.Step(`^the response body has (\d+) items$`, api.theResponseBodyHasItems)
	ctx.Step(`^the response should match json:$`, api.theResponseShouldMatchJSON)
	ctx.Step(`^I wait (\d+)s$`, api.iWaitS)
	ctx.Step(`^I send "(GET|POST|PUT|DELETE)" request to "([^"]*)" with json body:$`, api.iSendRequestToWithJsonBody)
}

func TestMain(m *testing.M) {
	flag.Parse()
	format := "progress"
	for _, arg := range os.Args[1:] {
		if arg == "-test.v=true" { // go test transforms -v option
			format = "pretty"
			break
		}
	}

	opts := godog.Options{
		Format:    format,
		Output:    colors.Colored(os.Stdout),
		Paths:     []string{"features"},
		Randomize: 0,
	}

	status := godog.TestSuite{
		Name:                "integration-tests",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()

	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}

func newConfig(configFile string) (*DbConfig, error) {
	viper.SetConfigFile(configFile)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("can't read config file: %w", err)
	}

	return &DbConfig{
		dsn:  viper.GetString("storage.DSN"),
		host: viper.GetString("http.host"),
		port: viper.GetString("http.port"),
	}, nil
}
