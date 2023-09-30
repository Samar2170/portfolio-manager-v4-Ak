package portfolio

import (
	"crypto/sha1"
	"fmt"

	"github.com/samar2170/portfolio-manager-v4/internal/models"
	"github.com/samar2170/portfolio-manager-v4/internal/portfolio/pbond"
	"github.com/samar2170/portfolio-manager-v4/internal/portfolio/pets"
	"github.com/samar2170/portfolio-manager-v4/internal/portfolio/pmutualfund"
	"github.com/samar2170/portfolio-manager-v4/internal/portfolio/pstock"
	"github.com/samar2170/portfolio-manager-v4/pkg/db"
	"github.com/samar2170/portfolio-manager-v4/pkg/utils/structs"
	"gorm.io/gorm"
)

type BlockHash struct {
	*gorm.Model
	UserCID      string
	Hash         string
	PreviousHash string
}

func (b *BlockHash) create() error {
	return db.DB.Create(b).Error
}
func GetLatestBlockHash(userCID string) (BlockHash, error) {
	var b BlockHash
	err := db.DB.Where("user_cid = ?", userCID).Last(&b).Error
	return b, err
}

type TradeInterface interface {
	GetAccount() models.DematAccount
	GetInvestedValue() float64
}

// lets do it blockchain style
func RegisterTrade(td TradeInterface) error {
	var err error
	switch td := td.(type) {
	case *pstock.StockTrade:
		err = pstock.RegisterStockTrade(td)
	case *pbond.BondTrade:
		err = pbond.RegisterBondTrade(td)
	case *pmutualfund.MutualFundTrade:
		err = pmutualfund.RegisterMutualFundTrade(td)
	case *pets.ETSTrade:
		err = pets.RegisterETSTrade(td)
	}
	// createHashBlockForTrade(&td)
	return err
}

func createHashBlockForTrade(td *TradeInterface) error {
	s := structs.New(td)
	m := fmt.Sprint(s.Map())
	accountCID := (*td).GetAccount().UserCID
	latestBlock, err := GetLatestBlockHash(accountCID)
	if err != nil {
		return err
	}
	hasher := sha1.New()
	hasher.Write([]byte(m))
	hash := hasher.Sum(nil)
	hashString := fmt.Sprintf("%x", hash)
	blockHash := BlockHash{
		UserCID:      accountCID,
		Hash:         hashString,
		PreviousHash: latestBlock.Hash,
	}
	err = blockHash.create()
	if err != nil {
		return err
	}
	return nil
}
