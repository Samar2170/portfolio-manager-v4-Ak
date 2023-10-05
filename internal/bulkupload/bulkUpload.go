package bulkupload

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strings"

	"github.com/samar2170/portfolio-manager-v4-Ak/internal"
	"github.com/samar2170/portfolio-manager-v4-Ak/internal/models"
	"github.com/samar2170/portfolio-manager-v4-Ak/pkg/utils/structs"
	"github.com/xuri/excelize/v2"
)

func CreateTradeTemplate() error {
	var err error
	f, err := excelize.OpenFile("assets/trade-template.xlsx", excelize.Options{})
	if err != nil {
		f = excelize.NewFile()
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()
	var str internal.StockTradeRequest
	var btr internal.BondTradeRequest
	var mtr internal.MutualFundTradeRequest
	var etr internal.ETSTradeRequest

	var names []string
	var types []string
	names, types = createRowFromApiRequest(str)
	fmt.Println(names, types)
	fmt.Println(len(names))
	for i, name := range names {
		fmt.Println(i, name)
	}

	fmt.Printf("%T\n", names)
	fmt.Printf("%T\n", types)
	f.NewSheet("Stock")
	err = f.SetSheetRow("Stock", "A1", &names)
	if err != nil {
		return err
	}
	err = f.SetSheetRow("Stock", "A2", &types)
	if err != nil {
		return err
	}
	names, types = createRowFromApiRequest(btr)
	f.NewSheet("Bond")
	f.SetSheetRow("Bond", "A1", &names)
	f.SetSheetRow("Bond", "A2", &types)

	names, types = createRowFromApiRequest(mtr)
	f.NewSheet("MutualFund")
	f.SetSheetRow("MutualFund", "A1", &names)
	f.SetSheetRow("MutualFund", "A2", &types)

	names, types = createRowFromApiRequest(etr)
	f.NewSheet("ETS")

	err = f.SetSheetRow("ETS", "A1", &names)
	if err != nil {
		return err
	}
	f.SetSheetRow("ETS", "A2", &types)

	if err := f.SaveAs("assets/trade-template.xlsx"); err != nil {
		return err
	}
	return nil
}

func createRowFromApiRequest(t interface{}) (names []string, types []string) {
	s := structs.New(t)
	m := s.MapWithType()
	for k, v := range m {
		names = append(names, k)
		types = append(types, v)
	}
	return
}

var UploadsDir = "uploads"

func SaveBulkUploadFile(file *multipart.FileHeader, userCID string) error {
	var err error
	user, err := models.GetUserByCID(userCID)
	if err != nil {
		return err
	}
	fileNameSplit := strings.Split(file.Filename, ".")
	fileExt := fileNameSplit[len(fileNameSplit)-1]
	newFileName := user.Username + "-" + "123456789" + "." + fileExt
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	newFileNameWithPath := fmt.Sprintf("%v/%v", UploadsDir, newFileName)
	dst, err := os.Create(newFileNameWithPath)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	bu := BulkUploadSheet{
		Name:    file.Filename,
		NewName: newFileName,
		UserCID: userCID,
		Path:    newFileNameWithPath,
	}
	err = bu.create()
	return err
}

func TestExcelize() {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	index := f.NewSheet("Sheet2")
	f.SetCellValue("Sheet2", "A2", "Hello world.")
	f.SetCellValue("Sheet1", "B2", 100)
	f.SetActiveSheet(index)
	if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}
