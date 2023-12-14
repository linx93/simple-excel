package simple_excel

import (
	"github.com/linx93/simple-excel/sample"
	"testing"
)

func TestExcelReader_Read(t *testing.T) {
	reader, err := NewExcelReader[sample.IntoRow]("./7-9事故与非事故已付明细 1.xlsx")
	if err != nil {
		t.Fatalf("err:%s", err.Error())
	}

	read, err := reader.Read()
	if err != nil {
		t.Fatalf("读err:%s", err.Error())
	}

	t.Log(read)
}
