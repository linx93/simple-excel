package sample

type Student struct {
	Hobby  Hobby  `json:"hobby" xlsx:"head:爱好;headCol:2"` //爱好
	Weight int    `json:"weight" xlsx:"head:体重"`
	Age    int    `json:"age" xlsx:"ignore:true"`                                                        //年龄
	Name   string `json:"name" xlsx:"head:名字;headCol:0;depth:3;width:20;merge:true;defValue:无名;index:0"` //名字
	Sub    Sub    `json:"sub" xlsx:"head:科目;headCol:2"`                                                  //科目
	Height int    `json:"height" xlsx:"head:身高"`
	//我的深度=我的兄弟的最大深度
}
type Hobby struct {
	Sing       string `json:"sing" xlsx:"head:唱歌;"`       //
	Basketball string `json:"basketball" xlsx:"head:篮球;"` //
}

type Sub struct {
	CompSub  CompSub  `json:"compSub" xlsx:"head:理科;"`   //Comprehensive subjects 综合科目 理科
	Chinese  string   `json:"chinese" xlsx:"head:语文;"`   //
	Math     string   `json:"math" xlsx:"head:数学;"`      //
	English  string   `json:"english" xlsx:"head:英语;"`   //
	CompSub1 CompSub1 `json:"compSub1" xlsx:"head:文科;"`  //Comprehensive subjects 综合科目 文科
	GymClass string   `json:"gymClass" xlsx:"head:体育课程"` //体育课

	//我的深度=我的兄弟的最大深度
}

type CompSub struct {
	Physics  string `json:"physics" xlsx:"head:物理;"`  //
	Chemical string `json:"chemical" xlsx:"head:化学;"` //
	Biology  string `json:"biology" xlsx:"head:生物;"`  //
}

type CompSub1 struct {
	Politics  string `json:"politics" xlsx:"head:政治;"` //政治
	Choose    Choose `json:"choose" xlsx:"head:文科选修"`
	History   string `json:"history" xlsx:"head:历史;"`   //历史
	Geography string `json:"geography" xlsx:"head:地理;"` //地理

}
type Choose struct {
	ChooseA       string              `json:"chooseA" xlsx:"head:选修A"`
	ChooseB       string              `json:"chooseB" xlsx:"head:选修B"`
	ModernHistory ChooseModernHistory `json:"modernHistory" xlsx:"head:选修近代史"`
	ChooseC       string              `json:"chooseC" xlsx:"head:选修C"`
}

type ChooseModernHistory struct {
	ModernChineseHistory string `json:"modernChineseHistory" xlsx:"head:中国近代史"`
	ModernWorldHistory   string `json:"modernWorldHistory" xlsx:"head:世界近代史"`
}
