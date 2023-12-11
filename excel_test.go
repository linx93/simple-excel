package simple_excel

import (
	"fmt"
	"github.com/linx93/simple-excel/sample"
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
	headers, err := buildHeaders[sample.Student]()
	if err != nil {
		t.Fatalf("buildHeaders:%s\n", err.Error())
	}
	for _, h := range headers {
		err1 := writeHeader(f, *h, "Sheet1")
		if err1 != nil {
			log.Fatalf(err1.Error())
		}
	}

	f.SaveAs("./linx18.xlsx")

	t.Log("success")
}

func TestCreateTab(t *testing.T) {
	savePath := "./test5.xlsx"

	stu := sample.Student{
		Hobby: sample.Hobby{
			Sing:       "0分",
			Basketball: "80分",
		},
		Weight: 120,
		Age:    18,
		Name:   "熊林",
		Sub: sample.Sub{
			CompSub: sample.CompSub{
				Physics:  "90",
				Chemical: "80",
				Biology:  "80",
			},
			Chinese: "60",
			Math:    "100",
			English: "0",
			CompSub1: sample.CompSub1{
				Politics: "40",
				Choose: sample.Choose{
					ChooseA: "40",
					ChooseB: "40",
					ModernHistory: sample.ChooseModernHistory{
						ModernChineseHistory: "50",
						ModernWorldHistory:   "50",
					},
					ChooseC: "40",
				},
				History:   "40",
				Geography: "80",
			},
			//GymClass: "80",
		},
		Height: 170,
	}

	data := []sample.Student{stu, stu, stu, stu, stu, stu, stu, stu, stu, stu, stu, stu}
	sheetStu1 := NewSheet("学生1", data)
	sheetStu2 := NewSheet("学生2", data)
	writerStu := NewExcelWriter[sample.Student](savePath)
	writerStu.AddSheet(sheetStu1, sheetStu2)
	writerStu.CreateTab()

	//继续追加写入school数据

	file, _ := excelize.OpenFile(savePath)
	school := sample.School{SType: "好学校", SName: "某某中学"}
	schData := []sample.School{school}
	sheetSch := NewSheet("学校", schData)
	writerSch := NewExcelWriter[sample.School](savePath, file)
	writerSch.AddSheet(sheetSch)
	writerSch.CreateTab()

	t.Log("success")
}
