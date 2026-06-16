package harness

import (
	"app/adapter/casbin"
	"app/adapter/casdoor"
	"app/adapter/featureflags"
	"app/adapter/health/health_module"
	apphttp "app/adapter/http"
	"app/adapter/logger"
	"app/adapter/metrics"
	"app/adapter/postgres"
	"app/adapter/redis"
	"app/adapter/time"
	"app/adapter/uuid"
	"app/adapter/validation"
	"app/adapter/watermill"
	"app/authn/authn_module"
	"app/authz/authz_http"
	"app/config"
	"app/internal/appcontext"
	"app/kv/kv_module"
	"app/test"
	"app/test/driver"
	"app/test/matchers"
	"app/test/req"
	"app/user/user_module"
	"fmt"

	gocasbin "github.com/casbin/casbin/v3"
	"github.com/onsi/ginkgo/v2"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

var defaultOptions = []fx.Option{
	logger.NopLoggerProvider,
	test.Queue,
	test.AvailablePortProvider,
	test.FakeAuthnProviders,
	test.TxIsolation,

	appcontext.Module,
	apphttp.Module,
	authz_http.Module,
	authn_module.HTTPModule,
	casbin.Module,
	casdoor.Module,
	config.Module,
	featureflags.Module,
	metrics.Module,
	postgres.Module,
	redis.Module,
	time.Module,
	uuid.Module,
	validation.Module,
	watermill.Module,

	user_module.AppModule,
	kv_module.HTTPModule,
	health_module.HTTPModule,
}

func Setup(options ...Option) *Harness {
	harness := &Harness{
		fxOptions:       defaultOptions,
		useTX:           true,
		deleteAuthUsers: true,
	}

	for _, option := range options {
		option.Apply(harness)
	}

	harness.Setup()

	return harness
}

type Harness struct {
	fxOptions []fx.Option
	app       *fxtest.App
	txPool    *test.TxPool
	authUsers *test.InMemoryUserProvider
	enforcer  *gocasbin.Enforcer
	port      int

	useTX           bool
	deleteAuthUsers bool
}

func (harness *Harness) Setup() {
	var (
		httpConfig apphttp.Config
		txPool     *test.TxPool
		authUsers  *test.InMemoryUserProvider
		enforcer   *gocasbin.Enforcer
	)

	harness.fxOptions = append(
		harness.fxOptions,
		fx.Populate(&httpConfig, &txPool, &authUsers, &enforcer),
	)

	harness.app = fxtest.New(ginkgo.GinkgoT(), harness.fxOptions...).RequireStart()
	harness.txPool = txPool
	harness.authUsers = authUsers
	harness.enforcer = enforcer
	harness.port = httpConfig.Port
}

func (harness *Harness) NewDriver() *driver.Driver {
	return driver.NewDriver(
		fmt.Sprintf("http://localhost:%d", harness.port),
		req.Headers{},
	)
}

func (harness *Harness) NewUser(email, password string) *driver.User {
	userDriver := harness.NewDriver()
	newUser := userDriver.Auth.MustRegister(email, password)
	userDriver.Login(newUser.Email, password)
	return &driver.User{
		App:    userDriver,
		Entity: newUser,
	}
}

func (harness *Harness) BeginTx() {
	matchers.Must(harness.txPool.Begin())
}

func (harness *Harness) RollbackTx() {
	matchers.Must(harness.txPool.Rollback())
}

func (harness *Harness) BeforeEach() {
	if harness.useTX {
		harness.BeginTx()
		// The casbin enforcer caches policy in memory; the rolled-back DB does
		// not clear it, so reload it to the clean committed state each test.
		matchers.Must(harness.enforcer.LoadPolicy())
	}

	if harness.deleteAuthUsers {
		harness.authUsers.Clear()
	}
}

func (harness *Harness) AfterEach() {
	if harness.useTX {
		harness.RollbackTx()
	}
}

func (harness *Harness) Teardown() {
	harness.app.RequireStop()
}
