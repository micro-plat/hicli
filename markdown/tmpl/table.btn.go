package tmpl

import (
	"fmt"
	"reflect"
	"strings"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/lib4go/types"
)

type BtnInfo struct {
	Method    string //方法名
	Name      string //按钮名称
	Alias     string //绑定字段别名
	Confirm   string
	Condition string
	Handler   string
	Rows      []Row
	FieldName string //详情页面按钮，对应的字段名
	IsQuery   bool
}

func newBtnInfo() *BtnInfo {
	return &BtnInfo{}
}

type Dialog struct {
	Name      string
	Method    string
	Path      string
	Condition string
}

type BatchInfo struct {
	Handler   string
	Name      string
	Method    string
	Condition string
	Confirm   string
}

func (a *BatchInfo) IsEmpty() bool {
	return a == nil || reflect.DeepEqual(a, &BatchInfo{})
}

var btnIndex = []string{"", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

//DispostELBtn
func (t *Table) DispostELBtnList() {
	if t.ExtInfo == "" {
		return
	}

	t.ListBtnInfo = make([]*BtnInfo, 0)
	t.ListDialogs = make([]*Dialog, 0)

	for _, v := range btnIndex {
		key := fmt.Sprintf("el_l_btn%s", v)
		if !strings.Contains(t.ExtInfo, key+"(") {
			break
		}
		//mode
		mode := types.GetString(getSubConContent(key, "mode")(t.ExtInfo), "btn")
		switch mode {
		case "btn":
			t.elListBtn(key)
		case "dialog":
			t.elListDialog(key)
		default:
			logs.Log.Warnf("%s的mode选项配置错误:", key, mode)
		}
	}
}

func (t *Table) elListBtn(key string) {
	info := newBtnInfo()

	//handler
	info.Handler = getSubConContent(key, "handler")(t.ExtInfo)
	if info.Handler == "" {
		logs.Log.Warn("列表页面btn的handler选项未配置:", key, t.ExtInfo)
		return
	}

	//name
	info.Name = getSubConContent(key, "name")(t.ExtInfo)
	if info.Name == "" {
		logs.Log.Warn("列表页面btn的name选项未配置:", t.ExtInfo)
		return
	}

	//method
	info.Method = getSubConContent(key, "method")(t.ExtInfo)

	//覆盖删除按钮
	cover := false
	if info.Method == "del" {
		t.BtnDel = true
		cover = true
	}

	//confirm
	info.Confirm = getSubConContent(key, "confirm")(t.ExtInfo)

	//condition
	info.Condition = translateCondition(getSubConContent(key, "condition")(t.ExtInfo))

	//cover
	if cover {
		t.ListBtnInfo = append(t.ListBtnInfo, info)
		return
	}

	//alias
	info.Alias = types.GetString(getSubConContent(key, "alias")(t.ExtInfo), key)

	//Rows
	for _, v := range getRows(info.Alias)(t.Rows) {
		info.Rows = append(info.Rows, *v)
	}

	if len(info.Rows) < 1 {
		logs.Log.Warn("列表页面btn的绑定的字段未配置")
	}

	t.ListBtnInfo = append(t.ListBtnInfo, info)
}

func (t *Table) elListDialog(key string) {

	info := &Dialog{}

	//Name
	info.Name = getSubConContent(key, "name")(t.ExtInfo)
	if info.Name == "" {
		logs.Log.Warn("列表按钮dialog的name选项未配置:", t.ExtInfo)
		return
	}

	//method
	info.Method = getSubConContent(key, "method")(t.ExtInfo)
	if info.Method == "Edit" {
		t.BtnShowEdit = true
	}
	if info.Method == "Detail" {
		t.BtnShowDetail = true
	}

	info.Path = getSubConContent(key, "path")(t.ExtInfo)
	if info.Path == "" {
		logs.Log.Warn("列表按钮dialog的path选项未配置:", key, t.ExtInfo)
		return
	}

	//confirm
	//info.Confirm = getSubConContent(key, "confirm")(t.ExtInfo)

	//condition
	info.Condition = translateCondition(getSubConContent(key, "condition")(t.ExtInfo))

	t.ListDialogs = append(t.ListDialogs, info)

}

//DispostELBtnDetail
func (t *Table) DispostELBtnDetail() {

	if t.ExtInfo == "" {
		return
	}

	for _, v := range btnIndex {
		key := fmt.Sprintf("el_d_btn%s", v)
		if !strings.Contains(t.ExtInfo, key+"(") {
			break
		}
		//mode
		mode := types.GetString(getSubConContent(key, "mode")(t.ExtInfo), "btn")
		switch mode {
		case "btn":
			t.elDetailBtn(key)
		case "dialog": //暂不处理
		//	t.elDetailDialog(key)
		default:
			logs.Log.Warnf("%s的mode选项配置错误:", key, mode)
		}
	}
}

func (t *Table) elDetailBtn(key string) {

	info := newBtnInfo()

	//handler
	info.Handler = getSubConContent(key, "handler")(t.ExtInfo)
	if info.Name == "" {
		logs.Log.Warn("详情按钮的handler选项未配置:", key, t.ExtInfo)
		return
	}

	//method
	info.Method = getSubConContent(key, "method")(t.ExtInfo)

	//desc
	info.Name = getSubConContent(key, "name")(t.ExtInfo)
	if info.Name == "" {
		logs.Log.Warn("详情按钮的name选项未配置:", t.ExtInfo)
		return
	}

	//confirm
	info.Confirm = getSubConContent(key, "confirm")(t.ExtInfo)

	//url
	info.Handler = getSubConContent(key, "handler")(t.ExtInfo)

	//condition
	info.Condition = translateCondition(getSubConContent(key, "condition")(t.ExtInfo))

	//key
	info.Alias = types.GetString(getSubConContent(key, "key")(t.ExtInfo), key)

	//Rows
	for _, v := range getRows(info.Alias)(t.Rows) {
		info.Rows = append(info.Rows, *v)
	}

	if info.FieldName == "" && len(info.Rows) == 1 {
		info.FieldName = info.Rows[0].Name
	}

	if info.FieldName == "" {
		logs.Log.Warn("详情页面btn未配置按钮绑定的字段")
		return
	}

	for _, v := range t.Rows {
		if v.Name == info.FieldName {
			if v.DetailBtnInfo == nil {
				v.DetailBtnInfo = make([]*BtnInfo, 0)
			}
			v.DetailBtnInfo = append(v.DetailBtnInfo, info)
		}
	}

}

func (t *Table) elDetailDialog(key string) {

	info := &Dialog{}

	//Name
	info.Name = getSubConContent(key, "name")(t.ExtInfo)
	if info.Name == "" {
		logs.Log.Warn("详情页面按钮dialog的name选项未配置:", t.ExtInfo)
		return
	}

	//method
	info.Method = getSubConContent(key, "method")(t.ExtInfo)
	// if info.Method == "Edit" {
	// 	t.BtnShowEdit = true
	// }
	// if info.Method == "Detail" {
	// 	t.BtnShowDetail = true
	// }

	info.Path = getSubConContent(key, "path")(t.ExtInfo)
	if info.Path == "" {
		logs.Log.Warn("详情页面dialog的path选项未配置:", key, t.ExtInfo)
		return
	}

	//confirm
	//info.Confirm = getSubConContent(key, "confirm")(t.ExtInfo)

	//condition
	info.Condition = translateCondition(getSubConContent(key, "condition")(t.ExtInfo))

	//@todo
	//t.ListDialogs = append(t.ListDialogs, info)

}

//DispostELBtnQuery
func (t *Table) DispostELBtnQuery() {

	if t.ExtInfo == "" {
		return
	}

	t.QueryBtnInfo = make([]*BtnInfo, 0)
	t.QueryDialogs = make([]*Dialog, 0)

	for _, v := range btnIndex {
		key := fmt.Sprintf("el_q_btn%s", v)
		if !strings.Contains(t.ExtInfo, key+"(") {
			break
		}
		//mode
		mode := types.GetString(getSubConContent(key, "mode")(t.ExtInfo), "btn")
		switch mode {
		case "btn", "query":
			t.elQueryBtn(key)
		case "batch":
			t.elQueryBatchBtn(key)
		case "dialog":
			t.elQueryDialog(key)
		default:
			logs.Log.Warnf("%s的mode选项配置错误:", key, mode)
		}
	}

}

func (t *Table) elQueryBtn(key string) {

	info := newBtnInfo()

	//method
	info.Method = getSubConContent(key, "method")(t.ExtInfo)

	//handler
	info.Handler = getSubConContent(key, "handler")(t.ExtInfo)
	if info.Name == "" {
		logs.Log.Warn("查询的handler选项未配置:", key, t.ExtInfo)
		return
	}

	//name
	info.Name = getSubConContent(key, "name")(t.ExtInfo)
	if info.Name == "" {
		logs.Log.Warn("查询按钮的name选项未配置:", t.ExtInfo)
		return
	}

	//confirm
	info.Confirm = getSubConContent(key, "confirm")(t.ExtInfo)

	//condition
	info.Condition = translateCondition(getSubConContent(key, "condition")(t.ExtInfo))

	if info.Method == "queryDatas" {
		t.BtnShowQuery = true
		t.QueryHandler = info.Handler
		info.IsQuery = true
	}

	t.QueryBtnInfo = append(t.QueryBtnInfo, info)
}

func (t *Table) elQueryBatchBtn(key string) {
	t.BatchInfo = &BatchInfo{}
	t.BatchInfo.Handler = getSubConContent(key, "handler")(t.ExtInfo)
	if t.BatchInfo.Handler == "" {
		logs.Log.Warn("批量操作的handler选项未配置:", key, t.ExtInfo)
		return
	}
	t.BatchInfo.Name = getSubConContent(key, "name")(t.ExtInfo)
	if t.BatchInfo.Name == "" {
		logs.Log.Warn("批量操作的name选项未配置:", t.ExtInfo)
		return
	}

	t.BatchInfo.Method = getSubConContent(key, "method")(t.ExtInfo)
	t.BatchInfo.Condition = translateCondition(getSubConContent(key, "condition")(t.ExtInfo))
	//confirm
	t.BatchInfo.Confirm = getSubConContent(key, "confirm")(t.ExtInfo)
}

func (t *Table) elQueryDialog(key string) {

	info := &Dialog{}

	//Name
	info.Name = getSubConContent(key, "name")(t.ExtInfo)
	if info.Name == "" {
		logs.Log.Warn("查询dialog的name选项未配置:", t.ExtInfo)
		return
	}

	//method
	info.Method = getSubConContent(key, "method")(t.ExtInfo)
	if info.Method == "Add" {
		t.BtnShowAdd = true
	}

	info.Path = getSubConContent(key, "path")(t.ExtInfo)
	if info.Path == "" {
		logs.Log.Warn("查询按钮dialog的path选项未配置:", key, t.ExtInfo)
		return
	}

	//confirm
	//info.Confirm = getSubConContent(key, "confirm")(t.ExtInfo)

	//condition
	info.Condition = translateCondition(getSubConContent(key, "condition")(t.ExtInfo))

	t.QueryDialogs = append(t.QueryDialogs, info)

}

type DownloadInfo struct {
	Title []string
}

func (t *Table) DispostELDownload() {
	t.DownloadInfo = &DownloadInfo{
		Title: make([]string, 0),
	}
	if t.ExtInfo == "" {
		return
	}

	c := getBracketContent([]string{"el_download"})(t.ExtInfo)
	if c == "" {
		return
	}
	t.DownloadInfo.Title = strings.Split(c, "|")
}

func translateCondition(c string) string {
	c = strings.Replace(c, " and ", " && ", -1)
	c = strings.Replace(c, " or ", " || ", -1)
	return c
}
