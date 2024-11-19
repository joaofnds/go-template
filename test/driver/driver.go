package driver

import (
	"app/adapter/featureflags"
	"app/adapter/health"
	apphttp "app/adapter/http"
	"app/adapter/logger"
	"app/adapter/postgres"
	"app/adapter/redis"
	"app/adapter/time"
	"app/adapter/uuid"
	"app/adapter/validation"
	"app/config"
	"app/kv"
	"app/test"
	"app/test/matchers"
	usermodule "app/user/module"
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

		uuid.Module,
		time.Module,
		config.Module,
		featureflags.Module,
		apphttp.Module,
		validation.Module,
		postgres.Module,
		redis.Module,

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
		User: NewUserDriver(url),
		KV:   NewKVDriver(url),
	}
}

type Driver struct {
	app *fxtest.App
	db  *gorm.DB

	URL  string
	User *UserDriver
	KV   *KVDriver
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
