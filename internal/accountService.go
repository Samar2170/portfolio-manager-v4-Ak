package internal

import "github.com/samar2170/portfolio-manager-v4-Ak/internal/models"

func RegisterBankAccount(ba BankAccountRequest, userCID string) error {
	_, err := models.GetUserByCID(userCID)
	if err != nil {
		return err
	}
	bankAccount := models.BankAccount{
		UserCID:   userCID,
		Bank:      ba.Bank,
		AccountNo: ba.AccountNo,
	}
	return models.CreateModelInstance(&bankAccount)

}

func RegisterDematAccount(da DematAccountRequest, userCID string) error {
	_, err := models.GetUserByCID(userCID)
	if err != nil {
		return err
	}
	dematAccount := models.DematAccount{
		UserCID:     userCID,
		AccountCode: da.AccountCode,
		Broker:      da.Broker,
	}
	return models.CreateModelInstance(&dematAccount)
}
