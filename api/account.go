package api

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/samar2170/portfolio-manager-v4/internal"
	"github.com/samar2170/portfolio-manager-v4/internal/models"
	"github.com/samar2170/portfolio-manager-v4/pkg/response"
)

func registerAccount(c echo.Context) error {
	accountType := c.Param("accountType")
	user := c.Get("user").(models.User)
	log.Println("user: ", user.Username, "user_cid: ", user.UserCID)
	switch accountType {
	case "bank":
		bankAccountRequest := new(internal.BankAccountRequest)
		if err := c.Bind(&bankAccountRequest); err != nil {
			return c.JSON(response.BadRequestResponseEcho(err.Error()))
		}
		err := internal.RegisterBankAccount(*bankAccountRequest, user.UserCID)
		if err != nil {
			return c.JSON(response.BadRequestResponseEcho(err.Error()))
		}
		return c.JSON(response.SuccessResponseEcho("Bank Account Registered Successfully"))
	case "demat":
		dematAccountRequest := new(internal.DematAccountRequest)
		if err := c.Bind(&dematAccountRequest); err != nil {
			return c.JSON(response.BadRequestResponseEcho(err.Error()))
		}
		err := internal.RegisterDematAccount(*dematAccountRequest, user.UserCID)
		if err != nil {
			return c.JSON(response.BadRequestResponseEcho(err.Error()))
		}
		return c.JSON(response.SuccessResponseEcho("Demat Account Registered Successfully"))
	default:
		return c.JSON(response.BadRequestResponseEcho("Invalid Account Type"))
	}
}

func listAccounts(c echo.Context) error {
	accountType := c.Param("accountType")
	user := c.Get("user").(models.User)
	switch accountType {
	case "bank":
		bankAccounts, err := models.GetBankAccountsByUserCID(user.UserCID)
		if err != nil {
			return c.JSON(response.BadRequestResponseEcho(err.Error()))
		}
		return c.JSON(response.JSONResponseEcho(bankAccounts))
	case "demat":
		dematAccounts, err := models.GetDematAccountsByUserCID(user.UserCID)
		if err != nil {
			return c.JSON(response.BadRequestResponseEcho(err.Error()))
		}
		return c.JSON(response.JSONResponseEcho(dematAccounts))
	case "general":
		generalAccounts, err := models.GetGeneralAccountsByUserCID(user.UserCID)
		if err != nil {
			return c.JSON(response.BadRequestResponseEcho(err.Error()))
		}
		return c.JSON(response.JSONResponseEcho(generalAccounts))
	default:
		bankAccounts, dematAccounts, generalAccounts, err := models.GetAccountsByUserCID(user.UserCID)
		if err != nil {
			return c.JSON(response.BadRequestResponseEcho(err.Error()))
		}
		return c.JSON(response.JSONResponseEcho(map[string]interface{}{
			"bank":    bankAccounts,
			"demat":   dematAccounts,
			"general": generalAccounts,
		}))
	}
}
