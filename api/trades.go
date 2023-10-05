package api

import (
	"github.com/labstack/echo/v4"
	"github.com/samar2170/portfolio-manager-v4-Ak/internal/bulkupload"
	"github.com/samar2170/portfolio-manager-v4-Ak/pkg/response"
)

func downloadTradeTemplate(c echo.Context) error {
	return c.Attachment("assets/trade-template.xlsx", "Trade Template")
}

func uploadTradeSheet(c echo.Context) error {
	userCID := c.Get("user_cid").(string)
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(response.BadRequestResponseEcho(err.Error()))
	}
	err = bulkupload.SaveBulkUploadFile(file, userCID)
	if err != nil {
		return c.JSON(response.BadRequestResponseEcho(err.Error()))
	}
	return c.JSON(response.SuccessResponseEcho("Successfully Upload"))
}
