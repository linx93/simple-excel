package sample

type School struct {
	SType string `xlsx:"head:学校类型"`
	SName string `xlsx:"head:学校名称"`
	SCode SCode  `xlsx:"head:财政教育经费投入(万元)"`
	Other Other  `xlsx:"head:其他投入(万元)"`
}

type SCode struct {
	Total   string `xlsx:"head:总计"`
	EduPay  EduPay `xlsx:"head:教育事业费"`
	BasePay string `xlsx:"head:基础拨款"`
}
type EduPay struct {
	Sum        Sum        `xlsx:"head:合计"`
	PeoplePay  string     `xlsx:"head:人员经费"`
	ComPay     string     `xlsx:"head:日常公共经费"`
	ProjectPay ProjectPay `xlsx:"head:项目经费"`
}
type ProjectPay struct {
	ProjectPayTotal string       `xlsx:"head:合计"`
	ProjectPayIn    ProjectPayIn `xlsx:"head:其中"`
}
type ProjectPayIn struct {
	ProjectPayIn1 string `xlsx:"head:标准化建设"`
	ProjectPayIn2 string `xlsx:"head:信息化建设"`
}

type Sum struct {
	Amount string `xlsx:"head:金额"`
	Incr   string `xlsx:"head:比上年增长（%）"`
}

type Other struct {
	Village Village `xlsx:"head:村投入"`
	Shehui  Shehui  `xlsx:"head:社会捐款"`
}
type Village struct {
	VillageToal string    `xlsx:"head:合计"`
	VillageIn   VillageIn `xlsx:"head:其中"`
}
type VillageIn struct {
	VillageInPeople  string `xlsx:"head:人员经费"`
	VillageInDaily   string `xlsx:"head:日常公用经费"`
	VillageInProject string `xlsx:"head:项目经费"`
	VillageInBase    string `xlsx:"head:基础投入"`
}

type Shehui struct {
	ShehuiToal string   `xlsx:"head:合计"`
	ShehuiIn   ShehuiIn `xlsx:"head:其中"`
}
type ShehuiIn struct {
	ShehuiInProject string `xlsx:"head:项目经费"`
	ShehuiInBase    string `xlsx:"head:基础投入"`
}
