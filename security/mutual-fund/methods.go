package mutualfund

import "github.com/samar2170/portfolio-manager-v4-Ak/pkg/db"

func SearchMutualFund(query string) ([]MutualFund, error) {
	var mfs []MutualFund
	mf, err := GetMutualFundBySchemeNavName(query)
	if err == nil {
		mfs = append(mfs, mf)
	}
	err = db.DB.Where("scheme_nav_name LIKE ?", "%"+query+"%").Find(&mfs).Error
	if err != nil {
		return nil, err
	}
	return mfs, nil
}

func GetMutualFundBySchemeNavName(schemeNavName string) (MutualFund, error) {
	var mf MutualFund
	err := db.DB.First(&mf, "scheme_nav_name = ?", schemeNavName).Error
	return mf, err
}

func GetMutualFundByID(id int) (MutualFund, error) {
	var mf MutualFund
	err := db.DB.First(&mf, id).Error
	return mf, err
}
