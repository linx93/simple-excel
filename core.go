package simple_excel

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

const (
	XLSX = "xlsx"

	index_    = "index:"
	name_     = "name:"
	defValue_ = "defValue:"
	head_     = "head:"
	width_    = "width:"
	replace_  = "replace:"
	ignore_   = "ignore:"
	cols_     = "cols:"

	index    = "index"
	name     = "name"
	defValue = "defValue"
	head     = "head"
	width    = "width"
	replace  = "replace"
	ignore   = "ignore"
	cols     = "cols"
)

var ModelErr = fmt.Errorf("模型定义异常")

type Header struct {
	Content     string   `json:"content"`     //表头的单元格内容
	FieldName   string   `json:"field"`       //字段名
	PFieldName  string   `json:"pFieldName"`  //父节点的字段名
	ColIndex    int      `json:"colIndex"`    //第几列,0开始
	Children    []Header `json:"children"`    //子节点集合
	HasChildren bool     `json:"hasChildren"` //存在子节点为true
	TreeLayer   int      `json:"treeLayer"`   //树的层,0开始
	TreeDepth   int      `json:"treeDepth"`   //树的深度
	LeafNode    int      `json:"leafNode"`    //叶子节点总数
}

type XlsxTag struct {
	DefValue string `json:"defValue"` //默认值
	Header   string `json:"header"`   //列的表头
	Index    int    `json:"index"`    //列的索引从0开始
	Width    int    `json:"width"`    //单元格的宽度
	Merge    bool   `json:"merge"`    //合并
	Replace  string `json:"replace"`  //替换
	Ignore   bool   `json:"ignore"`   //默认所有字段都会和excel去匹配，为true则会忽略该字段
	Cols     int    `json:"cols"`     //表头占用多少列，默认1
}

// 样例数据
func ss() []Header {
	hs := []Header{
		{
			Content:     "名字",
			FieldName:   "name",
			PFieldName:  "",
			ColIndex:    0,
			HasChildren: false,
			TreeLayer:   0,
			TreeDepth:   3,
			LeafNode:    1,
		},
		{
			Content:    "科目",
			FieldName:  "sub",
			PFieldName: "",
			ColIndex:   1,
			Children: []Header{
				{
					Content:     "语文",
					FieldName:   "chinese",
					PFieldName:  "sub",
					ColIndex:    1 + 0,
					HasChildren: false,
					TreeLayer:   0 + 1,
					TreeDepth:   3 - 1,
					LeafNode:    1,
				},
				{
					Content:     "数学",
					FieldName:   "math",
					PFieldName:  "sub",
					ColIndex:    1 + 1,
					HasChildren: false,
					TreeLayer:   0 + 1,
					TreeDepth:   3 - 1,
					LeafNode:    1,
				},
				{
					Content:     "英语",
					FieldName:   "english",
					PFieldName:  "sub",
					ColIndex:    1 + 2,
					HasChildren: false,
					TreeLayer:   0 + 1,
					TreeDepth:   3 - 1,
					LeafNode:    1,
				},
				{
					Content:    "理综",
					FieldName:  "compSub",
					PFieldName: "sub",
					ColIndex:   1 + 3,
					Children: []Header{
						{
							Content:     "物理",
							FieldName:   "physics",
							PFieldName:  "compSub",
							ColIndex:    1 + 3 + 0,
							HasChildren: false,
							TreeLayer:   0 + 1 + 1,
							TreeDepth:   3 - 2,
							LeafNode:    1,
						},
						{
							Content:     "化学",
							FieldName:   "chemical",
							PFieldName:  "compSub",
							ColIndex:    1 + 3 + 1,
							HasChildren: false,
							TreeLayer:   0 + 1 + 1,
							TreeDepth:   3 - 2,
							LeafNode:    1,
						},
						{
							Content:     "生物",
							FieldName:   "biology",
							PFieldName:  "compSub",
							ColIndex:    1 + 3 + 2,
							HasChildren: false,
							TreeLayer:   0 + 1 + 1,
							TreeDepth:   3 - 2,
							LeafNode:    1,
						},
					},
					HasChildren: true,
					TreeLayer:   0 + 1,
					TreeDepth:   3 - 1,
					LeafNode:    3,
				},
				{
					Content:     "体育",
					FieldName:   "sports",
					PFieldName:  "sub",
					ColIndex:    1 + 6,
					HasChildren: false,
					TreeLayer:   0 + 1,
					TreeDepth:   2,
					LeafNode:    1,
				},
			},
			HasChildren: true,
			TreeLayer:   0,
			TreeDepth:   3,
			LeafNode:    7,
		},
		{
			Content:     "身高",
			FieldName:   "height",
			PFieldName:  "",
			ColIndex:    8,
			HasChildren: false,
			TreeLayer:   0,
			TreeDepth:   3,
			LeafNode:    1,
		},
	}
	return hs
}

func BuildTabHead[T any]() ([]*Header, error) {
	//生成树
	headers, err := buildHeaders[T]()
	if err != nil {
		return nil, err
	}

	//设置深度
	maxDepthVal := 0
	for _, h := range headers {
		depth := maxDepth(h)
		if depth > maxDepthVal {
			maxDepthVal = depth
		}
	}

	for i := 0; i < len(headers); i++ {
		//设置树的深度
		setDepth(headers[i], maxDepthVal)
		//设置叶子节点数
		maxLeafNode(headers[i])
		//设置单元格的列坐标
		setColIndex(headers[i], i)
	}

	return headers, nil
}

func buildXlsxTag(xlsxTag *XlsxTag, tf reflect.Type) error {
	for i := 0; i < tf.NumField(); i++ {
		tag, ok := tf.Field(i).Tag.Lookup(XLSX)

		if !ok {
			continue
		}

		if tag == "-" || tag == "" {
			continue
		}

		err := matchTag(xlsxTag, tag)
		if err != nil {
			return err
		}

	}

	return nil
}

func matchTag(xlsxTag *XlsxTag, tag string) error {
	split := strings.Split(tag, ";")

	for _, st := range split {
		if strings.HasPrefix(st, defValue_) {
			xlsxTag.DefValue = strings.ReplaceAll(st, defValue_, "")
		} else if strings.HasPrefix(st, head_) {
			xlsxTag.Header = strings.ReplaceAll(st, head_, "")
		} else if strings.HasPrefix(st, index_) {

			idx, err := strconv.ParseInt(strings.ReplaceAll(st, index_, ""), 10, 64)
			if err != nil {
				return fmt.Errorf("%w:%s的值必须是整数类型,err:%s", ModelErr, index, err.Error())
			}
			xlsxTag.Index = int(idx)
		} else if strings.HasPrefix(st, width_) {
			wh, err := strconv.ParseUint(strings.ReplaceAll(st, width_, ""), 10, 64)
			if err != nil {
				return fmt.Errorf("%w:%s的值必须是正整数类型,err:%s", ModelErr, width, err.Error())
			}
			xlsxTag.Width = int(wh)
		} else if strings.HasPrefix(st, replace_) {
			xlsxTag.Replace = strings.ReplaceAll(st, replace_, "")
		} else if strings.HasPrefix(st, ignore_) {
			ig, err := strconv.ParseBool(strings.ReplaceAll(st, ignore_, ""))
			if err != nil {
				return fmt.Errorf("%w:%s的值必须是bool类型,err:%s", ModelErr, ignore, err.Error())
			}
			xlsxTag.Ignore = ig
		} else if strings.HasPrefix(st, cols_) {
			c, err := strconv.ParseUint(strings.ReplaceAll(st, cols_, ""), 10, 64)
			if err != nil {
				return fmt.Errorf("%w:%s的值必须是正整数类型,err:%s", ModelErr, cols, err.Error())
			}
			xlsxTag.Cols = int(c)
		}
	}

	return nil
}

// 根据结构体类型获取表头集合
func buildHeaders[T any]() ([]*Header, error) {
	t := new(T)

	tf := reflect.TypeOf(*t)

	headers := make([]*Header, 0)
	for i := 0; i < tf.NumField(); i++ {
		field := tf.Field(i)
		h, err := buildHeader(field, "", 0)
		if err != nil {
			return nil, err
		}
		if h == nil {
			continue
		}
		headers = append(headers, h)
	}
	return headers, nil
}

// 生成树
func buildHeader(field reflect.StructField, pFieldName string, pTreeLayer int) (*Header, error) {
	tag, ok := field.Tag.Lookup(XLSX)
	if !ok || tag == "-" || tag == "" {
		//跳过
		return nil, nil
	}

	xlsxTag := &XlsxTag{}
	err := matchTag(xlsxTag, tag)
	if err != nil {
		return nil, err
	}

	if xlsxTag.Ignore == true {
		//跳过
		return nil, nil
	}

	h := Header{
		Content:    xlsxTag.Header,
		FieldName:  field.Name,
		PFieldName: pFieldName,
		ColIndex:   xlsxTag.Index,
		Children:   make([]Header, 0),
		TreeLayer:  pTreeLayer,
	}

	if field.Type.Kind() == reflect.Struct {
		h.HasChildren = true
		for i := 0; i < field.Type.NumField(); i++ {
			children, err1 := buildHeader(field.Type.Field(i), field.Name, pTreeLayer+1)
			if err1 != nil {
				log.Printf("构建头失败:%s\n", err1.Error())
				return nil, err1
			}
			h.Children = append(h.Children, *children)
		}

	}

	return &h, nil
}

// 计算树的最大深度
func maxDepth(tree *Header) int {
	//空树深度为0
	if tree == nil {
		return 0
	}

	maxChildDepth := 0

	for _, child := range tree.Children {
		childDepth := maxDepth(&child)

		//取最大的子树深度
		if childDepth > maxChildDepth {
			maxChildDepth = childDepth
		}

	}

	//最大的字数深度+1就是本树的最大深度
	return maxChildDepth + 1
}

// 计算树的叶子节点数,找到树的节点中没有子节点的总数
func maxLeafNode(tree *Header) int {
	if tree == nil {
		return 0
	}

	leafNode := 1

	if tree.HasChildren {
		leafNode = leafNode - 1
		for i := 0; i < len(tree.Children); i++ {
			node := maxLeafNode(&tree.Children[i])
			leafNode = leafNode + node
		}
	}

	//计算做设置
	tree.LeafNode = leafNode

	return leafNode
}

// 设置树的深度
func setDepth(tree *Header, maxDepth int) {
	tree.TreeDepth = maxDepth
	for i := 0; i < len(tree.Children); i++ {
		setDepth(&tree.Children[i], maxDepth-1)
	}
}

// 设置colIndex
func setColIndex(tree *Header, startColIndex int) {

	tree.ColIndex = startColIndex
	log.Println(tree.Content, " ", tree.LeafNode)

	if tree.HasChildren {
		for i := 0; i < len(tree.Children); i++ {
			if tree.Children[i].HasChildren {
				setColIndex(&tree.Children[i], startColIndex+i)
			} else {
				setColIndex(&tree.Children[i], startColIndex+i)
			}
		}
	} else {

	}
}
