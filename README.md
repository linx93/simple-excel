# excel的简单读写
## 读
```func TestExcelReader_Read(t *testing.T) {
    reader, err := NewExcelReader[sample.Student]("./test5.xlsx")
    if err != nil {
        t.Fatalf("err:%s", err.Error())
    }
    read, err := reader.Read()
    if err != nil {
        t.Fatalf("读err:%s", err.Error())
        }
        t.Log(read)
    }
```

## 写
```
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
```


## 标签的使用规则,只是针对写excel
|    标签    | 是否必须 |  默认值  |  类型  |     作用      | 进度  |
|:--------:|:----:|:-----:|:----:|:-----------:|:---:|
|   head   |  是   | 无默认值  | bool |    设置表头列名     | 已实现 |
|  ignore  |  否   | false | bool |    字段忽略     | 已实现 |
|  width   |  否   |  10   | uint |    设置列宽度    | 已实现 |
| defValue |  否   |       |      |    设置默认值    | 待实现 |
| replace  |  否   |       |      |    设置替换     | 待实现 |
|  merge   |  否   | false | bool | 纵向合并值相同的单元格 | 待实现 |


##  注意事项
1. **定义结构体时，都是用类型而不是指针，可以参考sample目录下的两个结构体**
2. **标签如果配置了ignore:true时，读写都会忽略此字段，如果此字段是结构体类型那它的所有的子字段也都将被忽略**
3. **定义结构体时，字段的先后顺序就是写入excel列的先后顺序，如果想要调整excel中列的先后顺序，更改对应字段的顺序即可**
4. **~~读excel根据excel的表头创建结构体model时，表头列的顺序必须和model的属性字段顺序保持一直，否则读取出来的数据是不对的~~(已优化，不需要保持一致)**
5. **读取excel时，excel中的列可以多于model中的属性（未忽略的属性）但不能少。excel中列的顺序方面无要求**