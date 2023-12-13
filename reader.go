package simple_excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"reflect"
	"strconv"
)

type ExcelReader[T any] struct {
	readPath string
	file     *excelize.File
	sht      []Sheet[T]
}

func NewExcelReader[T any](readPath string, sheetNames ...string) (*ExcelReader[T], error) {

	f, err := excelize.OpenFile(readPath)
	if err != nil {
		return nil, err
	}

	tabHead, cs, err := buildTabHead[T]()
	if err != nil {
		return nil, err
	}

	sheets := make([]Sheet[T], 0)

	if sheetNames == nil || len(sheetNames) == 0 {
		sheetNames = f.GetSheetList()
	}

	for _, sheetName := range sheetNames {
		st := Sheet[T]{
			Sheet: sheet{
				name:   sheetName,
				cols:   cs,
				header: tabHead,
			},
		}
		sheets = append(sheets, st)
	}

	return &ExcelReader[T]{
		readPath: readPath,
		file:     f,
		sht:      sheets,
	}, nil
}

func (er *ExcelReader[T]) Read() (map[string][]T, error) {
	m := make(map[string][]T, 0)
	for _, st := range er.sht {
		sheetName := st.Sheet.name
		leafNodes := st.Sheet.cols.LeafNodes

		rows, err := er.file.Rows(sheetName)
		if err != nil {
			log.Printf("读取工作区%s出错\v", sheetName)
			continue
		}

		rowIndex := 0
		skipHead := er.sht[0].Sheet.cols.MaxTreeDepth

		ts := make([]T, 0)

		for rows.Next() {
			rowIndex++
			if rowIndex <= skipHead {
				//跳过表头
				continue
			}

			colVals, err1 := rows.Columns()
			if err1 != nil {
				log.Printf("读取工作区%s的行出错: %s", sheetName, err1)
				continue
			}

			//用空串补充长度
			lLen := len(leafNodes)
			cLen := len(colVals)
			if cLen < lLen {
				for i := 0; i < lLen-cLen; i++ {
					colVals = append(colVals, "")
				}
			}

			//log.Println("   colVals=", colVals)
			t := new(T)
			tv := reflect.ValueOf(t)
			for _, leafNode := range leafNodes {
				//log.Println("   Field=", leafNode.Content, "   setColVal=", colVals[leafNode.ColIndex], "   ColIndex=", leafNode.ColIndex)
				setStruct(tv.Elem(), leafNode.FieldName, colVals[leafNode.ColIndex])
			}
			ts = append(ts, *t)
		}

		st.data = ts
		m[sheetName] = ts
	}

	return m, nil
}

func setVal(field reflect.Value, xlsxValue string) {
	switch field.Type().Kind() {
	case reflect.String:
		field.SetString(xlsxValue)
	case
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		intValue, err := strconv.ParseInt(xlsxValue, 10, 64)
		if err != nil {
			return
		}
		field.SetInt(intValue)
	case
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		intValue, err := strconv.ParseUint(xlsxValue, 10, 64)
		if err != nil {
			return
		}
		field.SetUint(intValue)
	case reflect.Float64:
		floatValue, err := strconv.ParseFloat(xlsxValue, 10)
		if err != nil {
			return
		}
		field.SetFloat(floatValue)
	default:
		field.SetString(fmt.Sprintf("%#v", xlsxValue))
	}
}

func setStruct(val reflect.Value, fieldName string, setV string) {
	switch val.Kind() {
	case reflect.Ptr:
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		setStruct(val.Elem(), fieldName, setV)
	case reflect.Struct:
		fieldVal := val.FieldByName(fieldName)
		if fieldVal.IsValid() {
			//fieldVal.Set(reflect.ValueOf(setV).Convert(fieldVal.Type()))
			setVal(fieldVal, setV)
			return
		} else {
			for i := 0; i < val.NumField(); i++ {
				field := val.Field(i)
				setStruct(field, fieldName, setV)
			}
		}

	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			elem := val.Index(i)
			setStruct(elem, fieldName, setV)
		}
	case reflect.String:
		//log.Println(val.Type(), "  fieldName=", fieldName, " val=", setV)

	}
}
