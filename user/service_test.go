package user_test

import (
	"testing"

	"app/adapter/logger"
	"app/adapter/postgres"
	"app/adapter/time"
	"app/adapter/uuid"
	"app/config"
	"app/test"
	. "app/test/matchers"
	"app/user"
	usermodule "app/user/user_module"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"gorm.io/gorm"
)

func TestUserService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "user service suite")
}

var _ = Describe("user service", func() {
	var (
		app         *fxtest.App
		userService *user.Service
		db          *gorm.DB
	)

	BeforeEach(func() {
		app = fxtest.New(
			GinkgoT(),
			logger.NopLoggerProvider,
			test.Queue,
			config.Module,
			postgres.Module,
			usermodule.Module,
			uuid.Module,
			time.Module,
			fx.Populate(&userService, &db),
		)
		app.RequireStart()

		Must(db.Exec("BEGIN").Error)
	})

	AfterEach(func() {
		Must(db.Exec("ROLLBACK").Error)
		app.RequireStop()
	})

	Describe("DeleteAll", func() {
		It("removes all users", func(ctx SpecContext) {
			Must2(userService.CreateUser(ctx, "joao"))
			Expect(userService.List(ctx)).NotTo(BeEmpty())

			Must(userService.DeleteAll(ctx))
			Expect(userService.List(ctx)).To(BeEmpty())
		})
	})

	It("created users can be found by name", func(ctx SpecContext) {
		user := Must2(userService.CreateUser(ctx, "joao"))

		found := Must2(userService.FindByName(ctx, user.Name))
		Expect(found).To(Equal(user))
	})

	It("lists created users", func(ctx SpecContext) {
		user := Must2(userService.CreateUser(ctx, "joao"))
		Expect(userService.List(ctx)).To(ContainElement(user))
	})

	It("removed users are not listed", func(ctx SpecContext) {
		user := Must2(userService.CreateUser(ctx, "joao"))
		Must(userService.Remove(ctx, user))

		users := Must2(userService.List(ctx))
		Expect(users).NotTo(ContainElement(user))
	})
})
