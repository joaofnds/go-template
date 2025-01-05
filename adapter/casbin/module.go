package casbin

import (
	"app/authz"
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
		fx.Provide(newCasbinModel),
		fx.Provide(newCasbinPersistAdapter),
		fx.Provide(newCasbinEnforcer),

		fx.Provide(NewEnforcer),
		fx.Provide(func(e *Enforcer) authz.Enforcer { return e }),
		fx.Provide(NewRoleManager),
		fx.Provide(func(r *RoleManager) authz.RoleManager { return r }),
	)
	Invokes = fx.Options(
		fx.Invoke(loadPolicy),
	)

	//go:embed model.conf
	modelStr string
)

func newCasbinModel() (model.Model, error) {
	return model.NewModelFromString(modelStr)
}

func newCasbinPersistAdapter(db *gorm.DB) (persist.Adapter, error) {
	return gormadapter.NewAdapterByDBUseTableName(db, "casbin", "policies")
}

func newCasbinEnforcer(model model.Model, adapter persist.Adapter) (*casbin.Enforcer, error) {
	return casbin.NewEnforcer(model, adapter)
}

func loadPolicy(lifecycle fx.Lifecycle, enforcer *casbin.Enforcer) {
	enforcer.EnableAutoSave(true)
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error { return enforcer.LoadPolicy() },
	})
}
