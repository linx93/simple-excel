package simple_excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"testing"
)

type Student struct {
	Age  int    `json:"age" xlsx:"ignore:true"`                                                        //年龄
	Name string `json:"name" xlsx:"head:名字;headCol:0;depth:3;width:20;merge:true;defValue:无名;index:0"` //名字
	Sub  Sub    `json:"sub" xlsx:"head:科目;headCol:2"`                                                  //科目
	//我的深度=我的兄弟的最大深度
}

type Sub struct {
	Chinese  string  `json:"chinese" xlsx:"head:语文;"`   //
	Math     string  `json:"math" xlsx:"head:数学;"`      //
	English  string  `json:"english" xlsx:"head:英语;"`   //
	CompSub  CompSub `json:"compSub" xlsx:"head:综合科目;"` //Comprehensive subjects 综合科目
	GymClass string  `json:"gymClass" xlsx:"head:体育课程"` //体育课
	//我的深度=我的兄弟的最大深度
}

type CompSub struct {
	Physics  string `json:"physics" xlsx:"head:物理;"`  //
	Chemical string `json:"chemical" xlsx:"head:化学;"` //
	Biology  string `json:"biology" xlsx:"head:生物;"`  //
}

func Test_write(t *testing.T) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	for _, h := range ss() {
		err := writeHeader(f, h)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	//stu := Student{}
	//tf := reflect.TypeOf(stu)
	//for i := 0; i < tf.NumField(); i++ {
	//	tf.f
	//}
	f.SaveAs("./linx17.xlsx")

	t.Log("success")
}
