package driver

import (
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
	"app/authn"
	"app/authn/authn_module"
	"app/config"
	"app/internal/appcontext"
	"app/kv/kv_module"
	"app/test"
	"app/test/matchers"
	"app/test/req"
	"app/user/user_module"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/user"

	"github.com/onsi/ginkgo/v2"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"gorm.io/gorm"
)

var defaultOptions = []fx.Option{
	logger.NopLoggerProvider,
	test.Queue,
	test.AvailablePortProvider,

	appcontext.Module,
	apphttp.Module,
	authn_module.HTTPModule,
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

func Setup(options ...DriverOption) *Driver {
	driver := &Driver{
		fxOptions:       defaultOptions,
		useTX:           true,
		deleteAuthUsers: true,
	}

	for _, option := range options {
		option.Apply(driver)
	}

	driver.Setup()

	return driver
}

type Driver struct {
	fxOptions []fx.Option
	app       *fxtest.App
	db        *gorm.DB
	authUsers authn.UserProvider

	url     string
	headers req.Headers

	useTX           bool
	deleteAuthUsers bool

	Auth   *AuthDriver
	Health *HealthDriver
	KV     *KVDriver
	Users  *UserDriver
}

func (driver *Driver) Setup() {
	var (
		httpConfig apphttp.Config
		db         *gorm.DB
		authUsers  authn.UserProvider
	)

	driver.fxOptions = append(
		driver.fxOptions,
		fx.Populate(&httpConfig, &db, &authUsers),
	)

	driver.app = fxtest.New(ginkgo.GinkgoT(), driver.fxOptions...).RequireStart()
	driver.db = db
	driver.authUsers = authUsers
	driver.url = fmt.Sprintf("http://localhost:%d", httpConfig.Port)
	driver.headers = req.Headers{}

	driver.Auth = NewAuthDriver(driver.url, driver.headers)
	driver.Health = NewHealthDriver(driver.url, driver.headers)
	driver.KV = NewKVDriver(driver.url, driver.headers)
	driver.Users = NewUserDriver(driver.url, driver.headers)
}

func (driver *Driver) SetHeader(key, value string) {
	driver.headers[key] = value
}

func (driver *Driver) DeleteAuthUsers() {
	var emails = []string{}
	matchers.Must(driver.db.Model(&user.User{}).Pluck("email", &emails).Error)

	for _, email := range emails {
		matchers.Must(driver.authUsers.Delete(context.Background(), email))
	}
}

func (driver *Driver) BeginTx() {
	matchers.Must(driver.db.Exec("BEGIN").Error)
}

func (driver *Driver) RollbackTx() {
	matchers.Must(driver.db.Exec("ROLLBACK").Error)
}

func (driver *Driver) BeforeEach() {
	if driver.useTX {
		driver.BeginTx()
	}

	if driver.deleteAuthUsers {
		driver.DeleteAuthUsers()
	}
}

func (driver *Driver) AfterEach() {
	if driver.useTX {
		driver.RollbackTx()
	}
}

func (driver *Driver) Teardown() {
	driver.app.RequireStop()
}

func marshal(v any) io.Reader {
	return bytes.NewReader(matchers.Must2(json.Marshal(v)))
}
