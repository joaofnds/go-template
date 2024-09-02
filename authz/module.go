package authz

import (
	"context"
	_ "embed"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var (
	Module = fx.Module("authz", Providers, Invokes)

	Providers = fx.Options(
		fx.Provide(newModel),
		fx.Provide(newAdapter),
		fx.Provide(newEnforcer),
		fx.Provide(NewService),
	)
	Invokes = fx.Options(
		fx.Invoke(loadPolicy),
	)

	//go:embed model.conf
	modelStr string
)

func newModel() (model.Model, error) {
	return model.NewModelFromString(modelStr)
}

func newAdapter(db *gorm.DB) (persist.Adapter, error) {
	return gormadapter.NewAdapterByDBUseTableName(db, "casbin", "policies")
}

func newEnforcer(model model.Model, adapter persist.Adapter) (*casbin.Enforcer, error) {
	return casbin.NewEnforcer(model, adapter)
}

func loadPolicy(lifecycle fx.Lifecycle, enforcer *casbin.Enforcer) {
	enforcer.EnableAutoSave(true)
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error { return enforcer.LoadPolicy() },
	})
}
