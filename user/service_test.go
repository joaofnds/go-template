package user_test

import (
	"testing"

	"app/adapter/casbin"
	"app/adapter/logger"
	"app/adapter/postgres"
	"app/adapter/time"
	"app/adapter/uuid"
	"app/adapter/validation"
	"app/adapter/watermill"
	"app/config"
	"app/internal/appcontext"
	"app/test"
	. "app/test/matchers"
	"app/user"
	"app/user/user_module"

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

var _ = Describe("user service", Ordered, func() {
	var (
		app         *fxtest.App
		userService *user.Service
		db          *gorm.DB
	)

	BeforeAll(func() {
		app = fxtest.New(
			GinkgoT(),
			logger.NopLoggerProvider,
			test.Queue,
			casbin.Module,
			appcontext.Module,
			watermill.Module,
			validation.Module,
			config.Module,
			postgres.Module,
			user_module.ServiceModule,
			uuid.Module,
			time.Module,
			fx.Populate(&userService, &db),
		)
		app.RequireStart()
	})

	BeforeEach(func() {
		Must(db.Exec("BEGIN").Error)
	})

	AfterEach(func() {
		Must(db.Exec("ROLLBACK").Error)
	})

	AfterAll(func() {
		app.RequireStop()
	})

	Describe("DeleteAll", func() {
		It("removes all users", func(ctx SpecContext) {
			Must2(userService.CreateUser(ctx, "joao@template.com"))
			Expect(userService.List(ctx)).NotTo(BeEmpty())

			Must(userService.DeleteAll(ctx))
			Expect(userService.List(ctx)).To(BeEmpty())
		})
	})

	It("created users can be found by id", func(ctx SpecContext) {
		user := Must2(userService.CreateUser(ctx, "joao@template.com"))

		found := Must2(userService.FindByID(ctx, user.ID))
		Expect(found).To(Equal(user))
	})

	It("created users can be found by email", func(ctx SpecContext) {
		user := Must2(userService.CreateUser(ctx, "joao@template.com"))

		found := Must2(userService.FindByEmail(ctx, user.Email))
		Expect(found).To(Equal(user))
	})

	It("lists created users", func(ctx SpecContext) {
		user := Must2(userService.CreateUser(ctx, "joao@template.com"))
		Expect(userService.List(ctx)).To(ContainElement(user))
	})

	It("removed users are not listed", func(ctx SpecContext) {
		user := Must2(userService.CreateUser(ctx, "joao@template.com"))
		Must(userService.Remove(ctx, user))

		users := Must2(userService.List(ctx))
		Expect(users).NotTo(ContainElement(user))
	})
})
