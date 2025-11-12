package faker

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/sonni-a/wb-service/internal/models"
)

func ptr[T any](v T) *T { return &v }

func Init() {
	gofakeit.Seed(time.Now().UnixNano())
}

func GenerateFakeOrder() models.Order {
	orderUID := gofakeit.UUID()
	track := gofakeit.LetterN(12)

	return models.Order{
		OrderUID:          orderUID,
		TrackNumber:       track,
		Entry:             "WBIL",
		Locale:            "en",
		InternalSignature: ptr(""),
		CustomerID:        gofakeit.Username(),
		DeliveryService:   "meest",
		ShardKey:          gofakeit.Numerify("#"), // 1–9
		SmID:              gofakeit.Number(1, 100),
		DateCreated:       time.Now().UTC(),
		OofShard:          "1",

		Delivery: models.Delivery{
			OrderUID: orderUID,
			Name:     gofakeit.Name(),
			Phone:    "+1" + gofakeit.Numerify("##########"),
			Zip:      gofakeit.Numerify("#####"),
			City:     gofakeit.City(),
			Address:  gofakeit.Street(),
			Region:   gofakeit.State(),
			Email:    gofakeit.Email(),
		},

		Payment: models.Payment{
			OrderUID:     orderUID,
			Transaction:  orderUID,
			RequestID:    ptr(""),
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       gofakeit.Number(100, 3000),
			PaymentDt:    time.Now().Unix(),
			Bank:         "AlphaBank",
			DeliveryCost: gofakeit.Number(0, 500),
			GoodsTotal:   gofakeit.Number(100, 2000),
			CustomFee:    0,
		},

		Items: []models.Item{
			generateItem(orderUID, track),
			generateItem(orderUID, track),
		},
	}
}

func generateItem(orderUID, track string) models.Item {
	return models.Item{
		OrderUID:    orderUID,
		ChrtID:      int64(gofakeit.Number(100000, 999999)),
		TrackNumber: track,
		Price:       gofakeit.Number(10, 500),
		RID:         gofakeit.UUID(),
		Name:        gofakeit.ProductName(),
		Sale:        gofakeit.Number(0, 50),
		Size:        "M",
		TotalPrice:  gofakeit.Number(10, 1000),
		NmID:        int64(gofakeit.Number(1000000, 9999999)),
		Brand:       gofakeit.Company(),
		Status:      202,
	}
}
