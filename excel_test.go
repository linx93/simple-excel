package simple_excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"testing"
)

func Test_writeHeader(t *testing.T) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	headers, err := buildHeaders[Student]()
	if err != nil {
		t.Fatalf("buildHeaders:%s\n", err.Error())
	}
	for _, h := range headers {
		err := writeHeader(f, *h)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	f.SaveAs("./linx18.xlsx")

	t.Log("success")
}
