package driver

import (
	"app/adapter/casdoor"
	"app/adapter/featureflags"
	"app/adapter/health"
	apphttp "app/adapter/http"
	"app/adapter/logger"
	"app/adapter/postgres"
	"app/adapter/redis"
	"app/adapter/time"
	"app/adapter/uuid"
	"app/adapter/validation"
	"app/authn/authn_http"
	"app/config"
	"app/kv"
	"app/test"
	"app/test/matchers"
	usermodule "app/user/user_module"
	"fmt"

	"github.com/onsi/ginkgo/v2"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"gorm.io/gorm"
)

func Setup() *Driver {
	var (
		httpConfig apphttp.Config
		db         *gorm.DB
	)
	app := fxtest.New(
		ginkgo.GinkgoT(),
		logger.NopLoggerProvider,
		test.Queue,
		test.AvailablePortProvider,

		casdoor.Module,
		uuid.Module,
		time.Module,
		validation.Module,
		config.Module,
		featureflags.Module,
		apphttp.Module,
		postgres.Module,
		redis.Module,
		authn_http.Module,

		usermodule.Module,
		kv.Module,
		health.Module,

		fx.Populate(&httpConfig, &db),
	).RequireStart()

	url := fmt.Sprintf("http://localhost:%d", httpConfig.Port)
	return &Driver{
		app: app,
		db:  db,

		URL:  url,
		Auth: NewAuthDriver(url),
		User: NewUserDriver(url),
		KV:   NewKVDriver(url),
	}
}

type Driver struct {
	app *fxtest.App
	db  *gorm.DB

	URL  string
	Auth *AuthDriver
	KV   *KVDriver
	User *UserDriver
}

func (driver *Driver) BeginTx() {
	matchers.Must(driver.db.Exec("BEGIN").Error)
}

func (driver *Driver) RollbackTx() {
	matchers.Must(driver.db.Exec("ROLLBACK").Error)
}

func (driver *Driver) Teardown() {
	driver.app.RequireStop()
}
