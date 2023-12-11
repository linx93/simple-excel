package simple_excel

import (
	"github.com/xuri/excelize/v2"
	"log"
)

type ExcelWriter[T any] struct {
	file     *excelize.File
	savePath string
	Sheets   []*Sheet[T] //多个工作区
}

func NewSheet[T any](name string, data []T) *Sheet[T] {
	return &Sheet[T]{Sheet: sheet{name: name}, data: data}
}

func NewExcelWriter[T any](savePath string, file ...*excelize.File) *ExcelWriter[T] {
	var f *excelize.File
	if file == nil {
		f = excelize.NewFile()
		defer func() {
			if err := f.Close(); err != nil {
				log.Println(err)
			}
		}()
	} else {
		f = file[0]
	}

	return &ExcelWriter[T]{
		file:     f,
		savePath: savePath,
	}
}

func (eWriter *ExcelWriter[T]) AddSheet(sheet ...*Sheet[T]) {
	eWriter.Sheets = append(eWriter.Sheets, sheet...)
}

func (eWriter *ExcelWriter[T]) CreateTab() error {
	for _, sht := range eWriter.Sheets {
		err := sht.createSheet(eWriter.file)
		if err != nil {
			return err
		}
	}

	eWriter.file.DeleteSheet("Sheet1")

	err := eWriter.save()
	if err != nil {
		return err
	}

	return nil
}

func (eWriter *ExcelWriter[T]) save() error {
	err := eWriter.file.SaveAs(eWriter.savePath)
	if err != nil {
		log.Printf("%s\v", err.Error())
		return err
	}

	return nil
}
