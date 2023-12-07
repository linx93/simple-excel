package simple_excel

/*
import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"reflect"
	"strings"
)

func write[T any](data []*T, path string) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	t := new(T)
	tf := reflect.TypeOf(t)
	if tf.Kind() != reflect.Struct {
		return fmt.Errorf("%w:%s", ModelErr, "模型不是一个结构体")
	}

	if tf.NumField() == 0 {
		return fmt.Errorf("%w:%s", ModelErr, "模型结构体中没有属性")
	}

	setHeader(f, tf)
	setDate(f, tf, data)

	// Save spreadsheet by the given path.
	if err1 := f.SaveAs("test.xlsx"); err1 != nil {
		fmt.Println(err1)
	}
	return nil
}

func setDate[T any](f *excelize.File, tf reflect.Type, data []*T) {
	for rowIndex, row := range data {
		row
		tf := reflect.TypeOf(t)
	}
}

func setHeader(sheetName string, rowIndex int, f *excelize.File, tf reflect.Type, isHeader bool) error {
	for i := 0; i < tf.NumField(); i++ {
		tag, ok := tf.Field(i).Tag.Lookup(XLSX)

		if !ok {
			continue
		}

		if tag == "-" || tag == "" {
			continue
		}

		columnIndex := ""
		tagName := ""
		split := strings.Split(tag, ";")

		for _, st := range split {
			if strings.HasPrefix(st, colIndex_) {
				columnIndex = strings.ReplaceAll(st, colIndex_, "")
			}

			if strings.HasPrefix(st, name_) {
				tagName = strings.ReplaceAll(st, name_, "")
			}
		}

		f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnIndex, rowIndex), func() any {
			if isHeader {
				return tagName
			} else {
				return nil
			}
		})
	}
	return nil
}
*/
