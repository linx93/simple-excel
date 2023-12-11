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

//type LeafNode struct {
//	ColIndex int    `json:"column"` //列索引，从0开始
//	Name     string `json:"name"`   //列名，也是字段名
//}

// LeafNode 叶子节点，叶子节点就代表每个列
type LeafNode map[string]int //key是fieldName => value是colIndex，列索引，从0开始

type Cols struct {
	LeafNodes    LeafNode `json:"LeafNodes"`    //所有的叶子节点
	MaxTreeDepth int      `json:"maxTreeDepth"` //树的深度，决定数据是从第MaxTreeDepth+1行开始填充
}

func (cs Cols) getCell(fieldName string, dataRow int) (string, error) {
	colIndex, ok := cs.LeafNodes[fieldName]
	if !ok {
		return "", nil
	}

	cell, err := excelize.CoordinatesToCellName(colIndex+1, cs.MaxTreeDepth+dataRow+1)
	if err != nil {
		return "", err
	}

	return cell, nil
}

func (cs Cols) setCell(f *excelize.File, fieldName string, dataRow int, val any) error {
	cell, err := cs.getCell(fieldName, dataRow)
	if err != nil {
		return err
	}

	if cell != "" {
		err = f.SetCellValue("Sheet1", cell, val)
		if err != nil {
			return err
		}
	}

	return nil
}

// 生成表头
func writeHeader(file *excelize.File, h Header) error {
	cell, err := excelize.CoordinatesToCellName(h.ColIndex+1, h.TreeLayer+1)
	if err != nil {
		log.Printf("坐标转换失败:row=%d,col=%d,err:%s\n", h.TreeLayer+1, h.ColIndex+1, err.Error())
		return fmt.Errorf("坐标转换失败:row=%d,col=%d,err:%s", h.TreeLayer+1, h.ColIndex+1, err.Error())
	}

	file.SetCellValue("Sheet1", cell, h.Content)

	cell_ := ""
	//有子节点就需要横向合并单元格
	if h.HasChildren {
		cell_, err = excelize.CoordinatesToCellName(h.ColIndex+h.LeafNode, h.TreeLayer+1)
		if err != nil {
			log.Printf("坐标转换失败:row=%d,col=%d,err:%s\n", h.TreeLayer+1, h.ColIndex+h.LeafNode, err.Error())
			return fmt.Errorf("坐标转换失败:row=%d,col=%d,err:%s", h.TreeLayer+1, h.ColIndex+h.LeafNode, err.Error())
		}
		err = file.MergeCell("Sheet1", cell, cell_)
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
		err = file.MergeCell("Sheet1", cell, cell_)
		if err != nil {
			log.Printf("合并单元格失败:hCell=%s,vCell=%s,err:%s\n", cell, cell_, err.Error())
			return fmt.Errorf("合并单元格失败:hCell=%s,vCell=%s,err:%s", cell, cell_, err.Error())
		}
	}

	//设置行居中样式
	file.SetCellStyle("Sheet1", cell, cell_, center(file))

	if h.HasChildren {
		for _, item := range h.Children {
			err = writeHeader(file, item)
			if err != nil {
				log.Printf("err:%s\n", err.Error())
				return fmt.Errorf("err:%s\n", err.Error())
			}
		}
	}

	return nil
}

func CreateTab[T any](saveAs string, rows []T) error {
	file, cs, err := CreateTabHeadFile[T]()
	if err != nil {
		return err
	}

	//开始填充数据
	err = fullData[T](file, cs, rows)
	if err != nil {
		return err
	}

	err = file.SaveAs(saveAs)
	if err != nil {
		log.Printf("%s\v", err.Error())
		return err
	}

	return nil
}

func CreateTabHeadFile[T any]() (*excelize.File, *Cols, error) {
	tabHead, cs, err := BuildTabHead[T]()
	if err != nil {
		log.Printf("%s\v", err.Error())
		return nil, nil, err
	}

	f := excelize.NewFile()
	defer func() {
		if err1 := f.Close(); err1 != nil {
			log.Printf("%s\v", err.Error())
		}
	}()

	for _, h := range tabHead {
		err1 := writeHeader(f, *h)
		if err1 != nil {
			log.Printf("%s\v", err1.Error())
			return nil, nil, err1
		}
	}

	return f, cs, nil
}

// CreateTabHead 生成一个复杂表头
func CreateTabHead[T any](saveAs string) error {
	f, _, err := CreateTabHeadFile[T]()
	if err != nil {
		return err
	}

	err = f.SaveAs(saveAs)
	if err != nil {
		log.Printf("%s\v", err.Error())
		return err
	}
	return nil
}

// 填充数据
func fullData[T any](file *excelize.File, cols *Cols, rows []T) error {

	for i := 0; i < len(rows); i++ {
		row := rows[i]
		for j := 0; j < reflect.TypeOf(row).NumField(); j++ {
			err1 := full(file, cols, i, reflect.TypeOf(row).Field(j), reflect.ValueOf(row).Field(j))
			if err1 != nil {
				return err1
			}
		}
	}

	return nil
}

// 递归写入
func full(file *excelize.File, cols *Cols, index int, rs reflect.StructField, rv reflect.Value) error {
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
			err := full(file, cols, index, rs.Type.Field(i), rv.Field(i))
			if err != nil {
				return err
			}
		}
	}

	//log.Println("name=", rs.Name, "   value=", rv, "    ignore=", xlsxTag.Ignore)

	err := cols.setCell(file, rs.Name, index, rv)
	if err != nil {
		return err
	}

	return nil
}
