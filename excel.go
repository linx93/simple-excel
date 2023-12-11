package simple_excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"reflect"
)

func center(f *excelize.File) int {
	style := new(excelize.Style)
	style.Alignment = &excelize.Alignment{Horizontal: "center"}
	styleId, _ := f.NewStyle(style)
	return styleId
}

type (
	//LeafNode 叶子节点，叶子节点就代表每个列
	LeafNode map[string]*Header //key是fieldName

	Cols struct {
		LeafNodes    LeafNode `json:"LeafNodes"`    //所有的叶子节点
		MaxTreeDepth int      `json:"maxTreeDepth"` //树的深度，决定数据是从第MaxTreeDepth+1行开始填充
	}

	// Sheet 一个工作区
	Sheet[T any] struct {
		data  []T   //数据
		Sheet sheet //metadata
	}

	sheet struct {
		name   string    //工作区名称
		header []*Header //表头
		cols   *Cols     //
	}
)

func (st sheet) getCell(fieldName string, dataRow int) (string, error) {
	h, ok := st.cols.LeafNodes[fieldName]
	if !ok {
		return "", nil
	}

	cell, err := excelize.CoordinatesToCellName(h.ColIndex+1, st.cols.MaxTreeDepth+dataRow+1)
	if err != nil {
		return "", err
	}

	return cell, nil
}

func (st sheet) setCell(f *excelize.File, fieldName string, dataRow int, val reflect.Value) error {
	cell, err := st.getCell(fieldName, dataRow)
	if err != nil {
		return err
	}

	if cell != "" {

		err = f.SetCellValue(st.name, cell, val.Interface())
		if err != nil {
			return err
		}

		//设置行居中样式
		f.SetCellStyle(st.name, cell, cell, center(f))
	}

	return nil
}

func (st *Sheet[T]) initHeads(file *excelize.File) error {
	//制造表头
	tabHead, cs, err := initHeads[T](file, st.Sheet)
	if err != nil {
		return err
	}

	st.Sheet.header = tabHead
	st.Sheet.cols = cs

	return nil
}

func (st *Sheet[T]) fullData(file *excelize.File) error {
	err := fullData[T](file, st)
	if err != nil {
		return err
	}
	return nil
}

func (st *Sheet[T]) createSheet(file *excelize.File) error {
	_, err2 := file.NewSheet(st.Sheet.name)
	if err2 != nil {
		return err2
	}

	err := st.initHeads(file)
	if err != nil {
		return err
	}

	err = st.fullData(file)
	if err != nil {
		return err
	}

	return nil
}

// 生成表头
func writeHeader(file *excelize.File, h Header, sheetName string) error {
	cell, err := excelize.CoordinatesToCellName(h.ColIndex+1, h.TreeLayer+1)
	if err != nil {
		log.Printf("坐标转换失败:row=%d,col=%d,err:%s\n", h.TreeLayer+1, h.ColIndex+1, err.Error())
		return fmt.Errorf("坐标转换失败:row=%d,col=%d,err:%s", h.TreeLayer+1, h.ColIndex+1, err.Error())
	}

	file.SetCellValue(sheetName, cell, h.Content)

	cell_ := ""
	//有子节点就需要横向合并单元格
	if h.HasChildren {
		cell_, err = excelize.CoordinatesToCellName(h.ColIndex+h.LeafNode, h.TreeLayer+1)
		if err != nil {
			log.Printf("坐标转换失败:row=%d,col=%d,err:%s\n", h.TreeLayer+1, h.ColIndex+h.LeafNode, err.Error())
			return fmt.Errorf("坐标转换失败:row=%d,col=%d,err:%s", h.TreeLayer+1, h.ColIndex+h.LeafNode, err.Error())
		}
		err = file.MergeCell(sheetName, cell, cell_)
		if err != nil {
			log.Printf("合并单元格失败:hCell=%s,vCell=%s,err:%s\n", cell, cell_, err.Error())
			return fmt.Errorf("合并单元格失败:hCell=%s,vCell=%s,err:%s", cell, cell_, err.Error())

		}
	} else {
		//没有子节点就需要纵向合并单元格
		cell_, err = excelize.CoordinatesToCellName(h.ColIndex+1, h.TreeLayer+h.TreeDepth)
		if err != nil {
			log.Printf("坐标转换失败:row=%d,col=%d,err:%s\n", h.TreeLayer+h.TreeDepth, h.ColIndex+1, err.Error())
			return fmt.Errorf("坐标转换失败:row=%d,col=%d,err:%s", h.TreeLayer+h.TreeDepth, h.ColIndex+1, err.Error())
		}
		err = file.MergeCell(sheetName, cell, cell_)
		if err != nil {
			log.Printf("合并单元格失败:hCell=%s,vCell=%s,err:%s\n", cell, cell_, err.Error())
			return fmt.Errorf("合并单元格失败:hCell=%s,vCell=%s,err:%s", cell, cell_, err.Error())
		}
	}

	//设置行居中样式
	file.SetCellStyle(sheetName, cell, cell_, center(file))

	//列宽
	if h.XlsxTag.Width > 0 {
		toName, _ := excelize.ColumnNumberToName(h.ColIndex + 1)
		file.SetColWidth(sheetName, toName, toName, float64(h.XlsxTag.Width))
	}

	if h.HasChildren {
		for _, item := range h.Children {
			err = writeHeader(file, item, sheetName)
			if err != nil {
				log.Printf("err:%s\n", err.Error())
				return fmt.Errorf("err:%s\n", err.Error())
			}
		}
	}

	return nil
}

func initHeads[T any](f *excelize.File, sht sheet) ([]*Header, *Cols, error) {
	tabHead, cs, err := buildTabHead[T]()
	if err != nil {
		log.Printf("%s\v", err.Error())
		return nil, nil, err
	}

	for _, h := range tabHead {
		err1 := writeHeader(f, *h, sht.name)
		if err1 != nil {
			log.Printf("%s\v", err1.Error())
			return nil, nil, err1
		}
	}

	return tabHead, cs, nil
}

// 填充数据
func fullData[T any](file *excelize.File, st *Sheet[T]) error {

	for i := 0; i < len(st.data); i++ {
		row := st.data[i]
		for j := 0; j < reflect.TypeOf(row).NumField(); j++ {
			err1 := full(file, &st.Sheet, i, reflect.TypeOf(row).Field(j), reflect.ValueOf(row).Field(j))
			if err1 != nil {
				return err1
			}
		}
	}

	return nil
}

// 递归写入
func full(file *excelize.File, sht *sheet, index int, rs reflect.StructField, rv reflect.Value) error {
	//xlsxTag, err := buildXlsxTagByStructField(rs)
	//if err != nil {
	//	return err
	//}
	//
	//if xlsxTag.Ignore {
	//	return nil
	//}

	if rs.Type.Kind() == reflect.Struct {
		for i := 0; i < rs.Type.NumField(); i++ {
			err := full(file, sht, index, rs.Type.Field(i), rv.Field(i))
			if err != nil {
				return err
			}
		}
	}

	//log.Println("name=", rs.Name, "   value=", rv, "    ignore=", xlsxTag.Ignore)

	err := sht.setCell(file, rs.Name, index, rv)
	if err != nil {
		return err
	}

	return nil
}
