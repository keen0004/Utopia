package tests

import (
	"os"
	"testing"
	"utopia/internal/excel"
)

var (
	excelfile = "./test.xlsx"
	sheetname = "test"
)

func TestWriteExcel(t *testing.T) {
	os.Remove(excelfile)

	excel, err := excel.NewExcel(excelfile)
	if err != nil {
		t.Errorf("New excel failed with error %v", err)
		return
	}

	err = excel.Open()
	if err != nil {
		t.Errorf("Open excel failed with error %v", err)
		return
	}
	defer excel.Close(true)

	err = excel.WriteCell(sheetname, "A1", "A1")
	if err != nil {
		t.Errorf("Write excel failed with error %v", err)
		return
	}

	err = excel.WriteCell(sheetname, "B1", "B1")
	if err != nil {
		t.Errorf("Write excel failed with error %v", err)
		return
	}

	err = excel.WriteCell(sheetname, "A2", "A2")
	if err != nil {
		t.Errorf("Write excel failed with error %v", err)
		return
	}

	err = excel.WriteCell(sheetname, "B2", "B2")
	if err != nil {
		t.Errorf("Write excel failed with error %v", err)
		return
	}
}

func TestReadCell(t *testing.T) {
	excel, err := excel.NewExcel(excelfile)
	if err != nil {
		t.Errorf("New excel failed with error %v", err)
		return
	}

	err = excel.Open()
	if err != nil {
		t.Errorf("Open excel failed with error %v", err)
		return
	}
	defer excel.Close(false)

	value, err := excel.ReadCell(sheetname, "A1")
	if err != nil {
		t.Errorf("Read excel failed with error %v", err)
		return
	} else if value != "A1" {
		t.Errorf("Expect A1 but %s", value)
		return
	}

	value, err = excel.ReadCell(sheetname, "B2")
	if err != nil {
		t.Errorf("Read excel failed with error %v", err)
		return
	} else if value != "B2" {
		t.Errorf("Expect B2 but %s", value)
		return
	}
}

func TestWriteAll(t *testing.T) {
	os.Remove(excelfile)

	excel, err := excel.NewExcel(excelfile)
	if err != nil {
		t.Errorf("New excel failed with error %v", err)
		return
	}

	err = excel.Open()
	if err != nil {
		t.Errorf("Open excel failed with error %v", err)
		return
	}
	defer excel.Close(true)

	data := [][]string{{"A1", "B1"}, {"A2", "B2"}}
	err = excel.WriteAll(sheetname, data)
	if err != nil {
		t.Errorf("Read excel failed with error %v", err)
		return
	}
}

func TestReadAll(t *testing.T) {
	excel, err := excel.NewExcel(excelfile)
	if err != nil {
		t.Errorf("New excel failed with error %v", err)
		return
	}
	defer os.Remove(excelfile)

	err = excel.Open()
	if err != nil {
		t.Errorf("Open excel failed with error %v", err)
		return
	}
	defer excel.Close(false)

	value, err := excel.ReadAll(sheetname)
	if err != nil {
		t.Errorf("Read excel failed with error %v", err)
		return
	} else if len(value) != 2 {
		t.Errorf("Expect 2 rows but %d", len(value))
		return
	}

	if value[0][0] != "A1" || value[0][1] != "B1" || value[1][0] != "A2" || value[1][1] != "B2" {
		t.Errorf("Not match values %v", value)
		return
	}
}
