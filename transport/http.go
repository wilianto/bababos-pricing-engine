package transport

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/wilianto/bababos-pricing-engine/price"
)

type HttpHandler struct {
	BasicPricing price.PricingStrategy
	SurgePricing price.PricingStrategy
}

func (h *HttpHandler) GetPrice(c echo.Context) error {
	customerIDStr := c.QueryParam("customer_id")
	skuID := c.QueryParam("sku_id")
	qtyStr := c.QueryParam("qty")

	if customerIDStr == "" || skuID == "" || qtyStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing required params"})
	}

	customerID, err := strconv.ParseInt(customerIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid customer_id"})
	}

	qty, err := strconv.ParseInt(qtyStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid qty"})
	}

	req := price.PriceRequest{
		CustomerID: customerID,
		SkuID:      skuID,
		Qty:        qty,
	}

	// TODO: can change the strategy here, depend on city config for example:
	// Jakarta can use BasicPricing, but Surabaya can use SurgePricing
	price, err := h.BasicPricing.GetPrice(req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, price)
}
