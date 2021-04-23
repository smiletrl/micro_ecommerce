package cart

import (
	"github.com/labstack/echo/v4"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/redis"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCartRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cart repository suite")
}

var _ = Describe("cart repository methods", func() {
	var (
		repo       Repository
		c          echo.Context
		customerID int64
		skuID      string
		quantity   int
		err        error
	)

	BeforeEach(func() {
		// initialize service
		stage := os.Getenv(constants.Stage)
		if stage == "" {
			stage = constants.StageLocal
		}
		config, err := config.Load(stage)
		if err != nil {
			panic(err)
		}
		repo = NewRepository(redis.Test(config))

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c = e.NewContext(req, rec)
	})

	Context("with Delete & Create & Get & Update cart item", func() {
		var items map[string]string

		BeforeEach(func() {
			// delete this cart item
			customerID = int64(12)
			skuID = "sku_abc"
			quantity = 8
			// clean this sku id firstly.
			err = repo.Delete(c, customerID, skuID)
			Expect(err).To(BeNil())

			err = repo.Create(c, customerID, skuID, quantity)
		})

		It("should create cart items successfully", func() {
			Expect(err).To(BeNil())
			items, err = repo.Get(c, customerID)
			Expect(err).To(BeNil())
			Expect(items).To(BeEquivalentTo(map[string]string{
				"sku_abc": "8",
			}))
		})

		Context("with more items to be increased with Create", func() {
			BeforeEach(func() {
				customerID = int64(12)
				skuID = "sku_abc"
				quantity = 10
				err = repo.Create(c, customerID, skuID, quantity)
			})

			It("should increase cart items successfully", func() {
				Expect(err).To(BeNil())
				items, err = repo.Get(c, customerID)
				Expect(err).To(BeNil())
				Expect(items).To(BeEquivalentTo(map[string]string{
					"sku_abc": "18",
				}))
			})

			// now update cart
			Context("with items to be updated with Create", func() {
				BeforeEach(func() {
					customerID = int64(12)
					skuID = "sku_abc"
					quantity = 29
					err = repo.Update(c, customerID, skuID, quantity)
				})

				It("should update cart items successfully", func() {
					Expect(err).To(BeNil())
					items, err = repo.Get(c, customerID)
					Expect(err).To(BeNil())
					Expect(items).To(BeEquivalentTo(map[string]string{
						"sku_abc": "29",
					}))
				})

				// now delete cart
				Context("with items to be updated with Create", func() {
					BeforeEach(func() {
						customerID = int64(12)
						skuID = "sku_abc"
						quantity = 29
						err = repo.Delete(c, customerID, skuID)
					})

					It("should delete cart items successfully", func() {
						Expect(err).To(BeNil())
						items, err = repo.Get(c, customerID)
						Expect(err).To(BeNil())
						Expect(items).To(BeEquivalentTo(map[string]string{}))
					})
				})
			})
		})
	})
})
