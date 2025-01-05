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

	"github.com/onsi/ginkgo/v2"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"gorm.io/gorm"
)

var defaultOptions = []fx.Option{
	logger.NopLoggerProvider,
	test.Queue,
	test.AvailablePortProvider,
	test.FakeAuthnProviders,

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
	db        *gorm.DB
	authUsers *test.InMemoryUserProvider
	port      int

	useTX           bool
	deleteAuthUsers bool
}

func (harness *Harness) Setup() {
	var (
		httpConfig apphttp.Config
		db         *gorm.DB
		authUsers  *test.InMemoryUserProvider
	)

	harness.fxOptions = append(
		harness.fxOptions,
		fx.Populate(&httpConfig, &db, &authUsers),
	)

	harness.app = fxtest.New(ginkgo.GinkgoT(), harness.fxOptions...).RequireStart()
	harness.db = db
	harness.authUsers = authUsers
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
	matchers.Must(harness.db.Exec("BEGIN").Error)
}

func (harness *Harness) RollbackTx() {
	matchers.Must(harness.db.Exec("ROLLBACK").Error)
}

func (harness *Harness) BeforeEach() {
	if harness.useTX {
		harness.BeginTx()
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
