package transport

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wilianto/bababos-pricing-engine/price"
)

type HttpHandler struct {
	BasicPricing price.PricingStrategy
	SurgePricing price.PricingStrategy
}

func (h *HttpHandler) GetPrice(c echo.Context) error {
	var req price.PriceRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	// TODO: can change the strategy here, depend on city config for example:
	// Jakarta can use BasicPricing, but Surabaya can use SurgePricing
	price, err := h.BasicPricing.GetPrice(req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, price)
}
