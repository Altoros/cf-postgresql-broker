package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	"code.cloudfoundry.org/lager"
	"github.com/Altoros/pg-puppeteer-go"
	"github.com/pivotal-cf/brokerapi"
)

var (
	ErrAsyncNotSupported    = errors.New("async operaions are not supported")
	ErrUpdatingNotSupported = errors.New("updating is not supported")
)

// Represents requests Handler
type Handler struct {
	*pgp.PGPuppeteer

	services []brokerapi.Service
}

// Returns the list of provided services
func (h Handler) Services(ctx context.Context) []brokerapi.Service {
	return h.services
}

// Creates a DB and return its name as DashboardURL
func (h Handler) Provision(ctx context.Context, instanceID string, _ brokerapi.ProvisionDetails, _ bool) (brokerapi.ProvisionedServiceSpec, error) {
	dbname, err := h.CreateDB(instanceID)

	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, err
	}

	return brokerapi.ProvisionedServiceSpec{
		IsAsync:      false,
		DashboardURL: dbname,
	}, nil
}

// Drops a DB
func (h Handler) Deprovision(ctx context.Context, instanceID string, _ brokerapi.DeprovisionDetails, _ bool) (brokerapi.DeprovisionServiceSpec, error) {
	return brokerapi.DeprovisionServiceSpec{}, h.DropDB(instanceID)
}

// Creates a DB user for specified DB
func (h Handler) Bind(ctx context.Context, instanceID, bindingID string, _ brokerapi.BindDetails) (brokerapi.Binding, error) {
	creds, err := h.CreateUser(instanceID, bindingID)
	if err != nil {
		return brokerapi.Binding{}, err
	}

	return brokerapi.Binding{
		Credentials: creds,
	}, nil
}

// Drops a DB user
func (h Handler) Unbind(ctx context.Context, instanceID, bindingID string, _ brokerapi.UnbindDetails) error {
	return h.DropUser(instanceID, bindingID)
}

// Not supported
func (Handler) LastOperation(ctx context.Context, instanceID, operationData string) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{}, ErrAsyncNotSupported
}

// Not supported
func (Handler) Update(ctx context.Context, instanceID string, _ brokerapi.UpdateDetails, _ bool) (brokerapi.UpdateServiceSpec, error) {
	return brokerapi.UpdateServiceSpec{}, ErrUpdatingNotSupported
}

// Creates new requests handler
// Connects it to the database and parses services JSON string
func newHandler(source string, servicesJSON string, GUID string) (*Handler, error) {
	conn, err := pgp.New(source)

	if err != nil {
		return nil, err
	}

	services := make([]brokerapi.Service, 0)

	// Parse services list
	if err := json.Unmarshal([]byte(servicesJSON), &services); err != nil {
		return nil, err
	}

	replace := func(str string) string {
		return strings.Replace(str, "{GUID}", GUID, 1)
	}

	// Replace GUID with runtime value
	for i := 0; i < len(services); i++ {
		services[i].ID = replace(services[i].ID)

		for j := 0; j < len(services[i].Plans); j++ {
			services[i].Plans[j].ID = replace(services[i].Plans[j].ID)
		}
	}

	return &Handler{conn, services}, nil
}

// Determines the application name or returns a default value
func appName(envJSON string, defaultName string) string {
	env := struct {
		ApplicationName string `json:"application_name"`
	}{}

	if envJSON == "" {
		goto DEFAULT
	}

	if err := json.Unmarshal([]byte(envJSON), &env); err != nil {
		panic(err)
	}

	if env.ApplicationName == "" {
		goto DEFAULT
	}

	return env.ApplicationName

DEFAULT:
	return defaultName
}

func main() {
	// Set up logger
	logger := lager.NewLogger(appName(os.Getenv("VCAP_APPLICATION"), "cf-postgresql-broker"))
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	// Set up authentication
	credentials := brokerapi.BrokerCredentials{
		Username: os.Getenv("AUTH_USER"),
		Password: os.Getenv("AUTH_PASSWORD"),
	}

	// Create requests handler
	handler, err := newHandler(
		os.Getenv("PG_SOURCE"),
		os.Getenv("PG_SERVICES"),
		os.Getenv("CF_INSTANCE_GUID"))

	if err != nil {
		logger.Fatal("handler", err)
	}

	// Register requests handler
	http.Handle("/", brokerapi.New(handler, logger, credentials))

	// Boot up
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	logger.Info("boot-up", lager.Data{"port": port})

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Fatal("listen-and-serve", err)
	}
}
