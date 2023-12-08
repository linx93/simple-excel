package simple_excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"reflect"
	"testing"
)

func Test_buildHeader(t *testing.T) {
	student := Student{}

	tf := reflect.TypeOf(student)

	headers := make([]*Header, 0)
	for i := 0; i < tf.NumField(); i++ {
		field := tf.Field(i)
		h, err := buildHeader(field, "", 0)
		if err != nil {
			panic(err)
		}
		if h == nil {
			continue
		}
		headers = append(headers, h)
	}

	fmt.Printf("headers长度=%d\n", len(headers))
	for _, h := range headers {
		fmt.Printf("%#v\n", *h)
	}
}

func Test_buildHeaders(t *testing.T) {

	//生成树
	headers, err := buildHeaders[Student]()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	//设置深度
	maxDepthVal := 0
	for _, h := range headers {
		depth := maxDepth(h)
		if depth > maxDepthVal {
			maxDepthVal = depth
		}
	}

	lastLeafNode := 0

	for i := 0; i < len(headers); i++ {
		//设置叶子节点数
		setDepth(headers[i], maxDepthVal)
		maxLeafNode(headers[i])
		lastLeafNode = setColIndex(headers[i], lastLeafNode)
		//t.Logf("%#v\n", headers[i])
		t.Logf("lastLeafNode=%v\n", lastLeafNode)
	}

}

func Test_maxDepth(t *testing.T) {
	headers, err := buildHeaders[Student]()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	for _, h := range headers {
		depth := maxDepth(h)
		t.Logf("%s树的深度为%d\n", h.Content, depth)
	}
}

func Test_maxLeafNode(t *testing.T) {
	headers, err := buildHeaders[Student]()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	for _, h := range headers {
		depth := maxLeafNode(h)
		t.Logf("%s树的叶子节点数为%d\n", h.Content, depth)
	}
	for _, h := range headers {
		t.Logf("%#v\n", *h)
	}
}

func TestBuildTabHead(t *testing.T) {
	tabHead, err := BuildTabHead[Student]()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	for _, h := range tabHead {
		t.Logf("%#v", *h)
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	for _, h := range tabHead {
		err1 := writeHeader(f, *h)
		if err1 != nil {
			log.Fatalf(err1.Error())
		}
	}
	f.SaveAs("./school.xlsx")

	t.Log("success")
}

func Test_calcColIndex(t *testing.T) {
	//生成树
	headers, err := buildHeaders[Student]()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	for i := 0; i < len(headers); i++ {
		colIndex := calcColIndex(headers[i], "名字")
		t.Logf("colIndex=%d", colIndex)
	}
}

func w(s []School) {

}
