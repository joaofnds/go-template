package driver

import (
	"app/adapter/casdoor"
	"app/adapter/featureflags"
	"app/adapter/health"
	apphttp "app/adapter/http"
	"app/adapter/logger"
	"app/adapter/metrics"
	"app/adapter/postgres"
	"app/adapter/redis"
	"app/adapter/time"
	"app/adapter/uuid"
	"app/adapter/validation"
	"app/adapter/watermill"
	"app/authn/authn_http"
	"app/config"
	"app/internal/appcontext"
	"app/kv/kv_module"
	"app/test"
	"app/test/matchers"
	"app/test/req"
	"app/user/user_http"
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/onsi/ginkgo/v2"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"gorm.io/gorm"
)

func Setup(opts ...fx.Option) *Driver {
	var (
		httpConfig apphttp.Config
		db         *gorm.DB
	)
	allOpts := []fx.Option{
		logger.NopLoggerProvider,
		test.Queue,
		test.AvailablePortProvider,

		appcontext.Module,
		apphttp.Module,
		authn_http.Module,
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

		user_http.Module,
		kv_module.Module,
		health.Module,

		fx.Populate(&httpConfig, &db),
	}

	allOpts = append(allOpts, opts...)

	app := fxtest.New(ginkgo.GinkgoT(), allOpts...).RequireStart()

	url := fmt.Sprintf("http://localhost:%d", httpConfig.Port)
	headers := req.Headers{}
	return &Driver{
		app:     app,
		db:      db,
		url:     url,
		headers: headers,

		Auth:   NewAuthDriver(url, headers),
		Health: NewHealthDriver(url, headers),
		KV:     NewKVDriver(url, headers),
		Users:  NewUserDriver(url, headers),
	}
}

type Driver struct {
	app     *fxtest.App
	db      *gorm.DB
	url     string
	headers req.Headers

	Auth   *AuthDriver
	Health *HealthDriver
	KV     *KVDriver
	Users  *UserDriver
}

func (driver *Driver) SetHeader(key, value string) {
	driver.headers[key] = value
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

func marshal(v any) io.Reader {
	return bytes.NewReader(matchers.Must2(json.Marshal(v)))
}
