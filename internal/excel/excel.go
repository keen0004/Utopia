package excel

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

type Excel struct {
	path string         // The abs path of excel file
	fd   *excelize.File // The fd of open file
}

func NewExcel(path string) (*Excel, error) {
	if path == "" {
		return nil, errors.New("Invalid excel file path")
	}

	// change to abs path
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, errors.New("Cant not get abs of path")
	}

	return &Excel{path: path, fd: nil}, nil
}

func (excel *Excel) Open() error {
	// open excl file, create it if not exist
	var err error
	_, err = os.Stat(excel.path)
	if os.IsNotExist(err) {
		excel.fd = excelize.NewFile()
		excel.fd.Path = excel.path
	} else if err != nil {
		return err
	} else {
		excel.fd, err = excelize.OpenFile(excel.path)
		if err != nil {
			return err
		}
	}

	if excel.fd == nil {
		return errors.New("Open excel file failed")
	}

	return nil
}

func (excel *Excel) Close(save bool) {
	// auto save on close excel
	if excel.fd != nil {
		if save {
			excel.fd.SaveAs(excel.path)
		}

		excel.fd.Close()
	}

	excel.fd = nil
}

// read one cell value, index like "A2", "B10"
func (excel *Excel) ReadCell(sheet string, index string) (string, error) {
	if excel.fd == nil {
		return "", errors.New("Excel file not opened")
	}

	return excel.fd.GetCellValue(sheet, index)
}

// write one cell value, index like "A2", "B10"
func (excel *Excel) WriteCell(sheet string, index string, value string) error {
	if excel.fd == nil {
		return errors.New("Excel file not opened")
	}

	// create new sheet if not exist
	if excel.fd.GetSheetIndex(sheet) == -1 {
		excel.fd.NewSheet(sheet)
	}

	return excel.fd.SetCellValue(sheet, index, value)
}

// read all cell value in sheet
func (excel *Excel) ReadAll(sheet string) ([][]string, error) {
	if excel.fd == nil {
		return nil, errors.New("Excel file not opened")
	}

	if excel.fd.GetSheetIndex(sheet) == -1 {
		return nil, errors.New("Invalid sheet")
	}

	return excel.fd.GetRows(sheet)
}

// write all cell value in sheet
func (excel *Excel) WriteAll(sheet string, values [][]string) error {
	if excel.fd == nil {
		return errors.New("Excel file not opened")
	}

	// create a tmp sheet for write values
	tmp := uuid.NewString()
	excel.fd.NewSheet(tmp)

	rowsize := len(values)
	colsize := len(values[0])

	for i := 0; i < rowsize; i++ {
		for j := 0; j < colsize; j++ {
			err := excel.fd.SetCellValue(tmp, fmt.Sprintf("%c%d", 65+j, i+1), values[i][j])
			if err != nil {
				excel.fd.DeleteSheet(tmp)
				return err
			}
		}
	}

	// delete old sheet
	if excel.fd.GetSheetIndex(sheet) == -1 {
		excel.fd.DeleteSheet(sheet)
	}

	// remove to dst sheet
	excel.fd.SetSheetName(tmp, sheet)
	excel.fd.SetActiveSheet(excel.fd.GetSheetIndex(sheet))
	excel.fd.SaveAs(excel.path)

	return nil
}
