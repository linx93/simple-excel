package simple_excel

import (
	"fmt"
	"github.com/linx93/simple-excel/sample"
	"reflect"
	"testing"
)

func Test_buildHeader(t *testing.T) {
	student := sample.Student{}

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
	headers, err := buildHeaders[sample.Student]()
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
	leafNodeMap := make(map[string]*Header, 0)
	for i := 0; i < len(headers); i++ {
		//设置叶子节点数
		setDepth(headers[i], maxDepthVal)
		maxLeafNode(headers[i])
		lastLeafNode, leafNodeMap = setColIndex(headers[i], lastLeafNode, leafNodeMap)
		//t.Logf("%#v\n", headers[i])
		t.Logf("lastLeafNode=%v\n", lastLeafNode)
		t.Logf("leafNodeMap=%v\n", leafNodeMap)
	}

}

func Test_maxDepth(t *testing.T) {
	headers, err := buildHeaders[sample.Student]()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	for _, h := range headers {
		depth := maxDepth(h)
		t.Logf("%s树的深度为%d\n", h.Content, depth)
	}
}

func Test_maxLeafNode(t *testing.T) {
	headers, err := buildHeaders[sample.Student]()
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
