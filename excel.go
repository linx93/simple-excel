package simple_excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
)

func center(f *excelize.File) int {
	style := new(excelize.Style)
	style.Alignment = &excelize.Alignment{Horizontal: "center"}
	styleId, _ := f.NewStyle(style)
	return styleId
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
