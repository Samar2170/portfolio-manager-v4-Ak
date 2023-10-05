package portfolio

import (
	"crypto/sha1"
	"fmt"

	portfoliobase "github.com/samar2170/portfolio-manager-v4-Ak/internal/portfolio/portfolio-base"
	"github.com/samar2170/portfolio-manager-v4-Ak/pkg/db"
	"github.com/samar2170/portfolio-manager-v4-Ak/pkg/utils/structs"
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
func createHashBlockForTrade(td *portfoliobase.TradeInterface) error {
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
