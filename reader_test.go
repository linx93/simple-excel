package simple_excel

import (
	"github.com/linx93/simple-excel/sample"
	"testing"
)

func TestExcelReader_Read(t *testing.T) {
	reader, err := NewExcelReader[sample.Student]("D:\\code\\go\\linx\\simple-excel\\test5.xlsx")
	if err != nil {
		t.Fatalf("err:%s", err.Error())
	}
	read, err := reader.Read()
	if err != nil {
		t.Fatalf("读err:%s", err.Error())
	}
	t.Log(read)
}
