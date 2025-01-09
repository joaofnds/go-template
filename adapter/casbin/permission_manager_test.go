package casbin_test

import (
	"app/adapter/casbin"
	"app/adapter/logger"
	"app/adapter/postgres"
	"app/adapter/validation"
	"app/authz"
	"app/config"
	"app/internal/ref"
	. "app/test/matchers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"gorm.io/gorm"
)

var _ = Describe("casbin enforcer", func() {
	user := ref.New("user", "111")
	post := ref.New("post", "222")
	anyPost := ref.New("post", "*")
	admin := ref.New("role", "admin")

	adminAnyPostDelete := authz.NewRequest(admin, anyPost, "delete")
	adminAnyPostDeletePolicy := authz.NewAllowPolicy(admin, anyPost, "delete")

	userPostDelete := authz.NewRequest(user, post, "delete")
	userPostDeletePolicy := authz.NewAllowPolicy(user, post, "delete")
	userAnyPostDeletePolicy := authz.NewAllowPolicy(user, anyPost, "delete")

	denyUserPostDeletePolicy := authz.NewDenyPolicy(user, post, "delete")

	var (
		app   *fxtest.App
		db    *gorm.DB
		sut   *casbin.PermissionManager
		roles *casbin.RoleManager
	)

	BeforeEach(func() {
		app = fxtest.New(
			GinkgoT(),
			logger.NopLoggerProvider,
			validation.Module,
			config.Module,
			postgres.Module,
			casbin.Module,
			fx.Populate(&db, &roles, &sut),
		)
		app.RequireStart()
		Must(db.Exec("BEGIN").Error)
	})

	AfterEach(func() {
		Must(db.Exec("ROLLBACK").Error)
		app.RequireStop()
	})

	It("grants permission", func(ctx SpecContext) {
		Expect(sut.Check(userPostDelete)).To(BeFalse())

		Must(sut.Add(userPostDeletePolicy))

		Expect(sut.Check(userPostDelete)).To(BeTrue())
	})

	It("grants multiple permissions", func(ctx SpecContext) {
		Expect(sut.Check(userPostDelete)).To(BeFalse())
		Expect(sut.Check(adminAnyPostDelete)).To(BeFalse())

		Must(sut.Add(userPostDeletePolicy, adminAnyPostDeletePolicy))

		Expect(sut.Check(userPostDelete)).To(BeTrue())
		Expect(sut.Check(adminAnyPostDelete)).To(BeTrue())
	})

	Describe("after direct grant", func() {
		It("has permission", func(ctx SpecContext) {
			Expect(sut.Check(userPostDelete)).To(BeFalse())

			Must(sut.Add(userPostDeletePolicy))

			Expect(sut.Check(userPostDelete)).To(BeTrue())
		})
	})

	Describe("after role grant", func() {
		It("has permission", func(ctx SpecContext) {
			Must(sut.Add(adminAnyPostDeletePolicy))
			Expect(sut.Check(userPostDelete)).To(BeFalse())

			Must(roles.Assign(user, admin))

			Expect(sut.Check(userPostDelete)).To(BeTrue())
		})

		When("deny override", func() {
			It("does not have permission", func(ctx SpecContext) {
				Must(sut.Add(adminAnyPostDeletePolicy))
				Must(roles.Assign(user, admin))
				Expect(sut.Check(userPostDelete)).To(BeTrue())

				Must(sut.Add(denyUserPostDeletePolicy))
				Expect(sut.Check(userPostDelete)).To(BeFalse())
			})
		})
	})

	Describe("after role revoke", func() {
		It("looses access", func(ctx SpecContext) {
			Must(sut.Add(userPostDeletePolicy))
			Expect(sut.Check(userPostDelete)).To(BeTrue())

			Must(sut.Remove(userPostDeletePolicy))
			Expect(sut.Check(userPostDelete)).To(BeFalse())
		})
	})

	Describe("after removing permissions for a subject", func() {
		It("looses access", func(ctx SpecContext) {
			Must(sut.Add(userPostDeletePolicy, userAnyPostDeletePolicy))
			Expect(sut.Check(userPostDelete)).To(BeTrue())

			Must(sut.RemoveBySubject(user))
			Expect(sut.Check(userPostDelete)).To(BeFalse())
		})
	})
})
