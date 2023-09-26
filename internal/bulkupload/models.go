package bulkupload

import (
	"github.com/samar2170/portfolio-manager-v4/pkg/db"
	"gorm.io/gorm"
)

func init() {
	db.DB.AutoMigrate(BulkUploadSheet{})
}

type BulkUploadSheet struct {
	*gorm.Model
	ID          uint
	Name        string
	NewName     string
	Path        string
	Parsed      bool
	Error       bool
	ErrorString string
	UserCID     string
}

func (s *BulkUploadSheet) create() error {
	return db.DB.Create(s).Error
}

func GetBulkUploadSheetByID(id uint) (BulkUploadSheet, error) {
	var s BulkUploadSheet
	err := db.DB.Where("id =?", id).First(&s).Error
	return s, err
}

func getUnparsedBulkUploadSheets(limit uint) ([]uint, error) {
	var ids []uint
	err := db.DB.Model(&BulkUploadSheet{}).Where("parsed =?", false).Limit(int(limit)).Select("id").Find(&ids).Error
	return ids, err
}
