package transport

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/wilianto/bababos-pricing-engine/price"
)

type HttpHandlerTestSuite struct {
	suite.Suite
}

func TestHttpHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HttpHandlerTestSuite))
}

func (s *HttpHandlerTestSuite) TestGetPrice() {
	req := httptest.NewRequest(http.MethodGet, "/price?customer_id=1&sku_id=2&qty=120", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(req, rec)
	handler := &HttpHandler{
		BasicPricing: &mockPricingStrategy{},
	}

	err := handler.GetPrice(ctx)
	require.NoError(s.T(), err)
	require.Equal(s.T(), http.StatusOK, rec.Code)
	require.Equal(s.T(), "{\"suggested_price\":1000,\"recommended_supplier_id\":1234}\n", rec.Body.String())
}

type mockPricingStrategy struct {
}

func (m *mockPricingStrategy) GetPrice(req price.PriceRequest) (price.PriceResponse, error) {
	return price.PriceResponse{
		SuggestedPrice:        1000,
		RecommendedSupplierID: 1234,
	}, nil
}
