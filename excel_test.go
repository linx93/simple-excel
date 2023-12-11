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
		err := writeHeader(f, *h)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	f.SaveAs("./linx18.xlsx")

	t.Log("success")
}

func TestCreateTab(t *testing.T) {
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
			Chinese:  "60",
			Math:     "100",
			English:  "0",
			CompSub1: sample.CompSub1{}, /*sample.CompSub1{
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
			},*/
			//GymClass: "80",
		},
		Height: 170,
	}

	stus := []sample.Student{stu, stu, stu, stu, stu, stu, stu, stu, stu, stu, stu, stu}

	err := CreateTab[sample.Student]("./linx-student9.xlsx", stus)
	if err != nil {
		t.Fatalf("CreateTab失败:%s", err.Error())
	}
	t.Log("success")
}
