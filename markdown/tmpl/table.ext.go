package tmpl

import (
	"fmt"
	"reflect"
	"strings"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/lib4go/types"
)

type TabInfo struct {
	TabField         map[string]string //详情字段名---详情字段名
	TabListField     map[string]string //详情列表字段名---详情列表字段名
	TabTable         map[string]bool   //详情
	TabTableList     map[string]bool   //详情列表
	TabTablePreField map[string]string //前表关联字段 表名-字段名
	TabTableProField map[string]string //后表关联字段 表名-字段名
}

func newTableInfo() *TabInfo {
	return &TabInfo{
		TabField:         make(map[string]string),
		TabListField:     make(map[string]string),
		TabTable:         make(map[string]bool),
		TabTableList:     make(map[string]bool),
		TabTablePreField: make(map[string]string),
		TabTableProField: make(map[string]string),
	}
}

//DisposeELTab 处理前端详情页
func (t *Table) DisposeELTab() {
	if t.ExtInfo == "" {
		return
	}

	c := getBracketContent([]string{"el_tab"})(t.ExtInfo)
	tabs := strings.Split(c, "|")
	if len(tabs) == 0 || c == "" {
		return
	}
	for _, v := range tabs {
		tab := strings.Split(v, ",")
		if len(tab) < 1 || len(tab) > 3 {
			logs.Log.Warn("tab格式不正确：", v)
			continue
		}
		tabName := tab[0]
		tabList := ""
		tabField := make([]string, 2)

		if len(tab) > 1 {
			t := strings.Split(tab[1], "/")
			if len(t) == 1 {
				tabField = []string{t[0], t[0]}
			}
			if len(t) == 2 {
				tabField = []string{t[0], t[1]}
			}
		}
		if len(tab) > 2 {
			tabList = tab[2]
		}

		exist := false
		for _, tb := range t.AllTables {
			if tb.Name == tabName {
				if tabList == "list" { //列表
					tb.TabInfo.TabTableList[t.Name] = true
					tb.TabInfo.TabListField[tabField[1]] = tabField[1]
				} else {
					tb.TabInfo.TabTable[t.Name] = true
					tb.TabInfo.TabField[tabField[1]] = tabField[1]
				}
				tb.TabInfo.TabTablePreField[t.Name] = tabField[0]
				tb.TabInfo.TabTableProField[t.Name] = tabField[1] //后表关联字段
				t.TabTables = append(t.TabTables, tb)
				exist = true
				break
			}
		}
		if !exist {
			logs.Log.Warnf("[%s]详情tab表名不正确：[%s]", t.Name, tabName)
		}
	}
}

type BtnInfo struct {
	Name      string
	DESC      string
	KeyWord   string
	Confirm   string
	Condition string
	URL       string
	Rows      []Row
	FieldName string //详情页面按钮，对应的字段名
	IsQuery   bool
}

func newBtnInfo() *BtnInfo {
	return &BtnInfo{}
}

//DispostBtnTables {el_btn(name:funcName,desc:1-启用|2-禁用,confirm:你确定进行修改吗,condition:condition,key:sa)}
func (t *Table) DispostELBtn() {
	if t.ExtInfo == "" {
		return
	}
	a := []string{"", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

	for _, v := range a {
		key := fmt.Sprintf("el_btn%s", v)
		if !strings.Contains(t.ExtInfo, key+"(") {
			break
		}

		info := newBtnInfo()

		//name
		info.Name = getSubConContent(key, "name")(t.ExtInfo)
		if info.Name == "" {
			logs.Log.Warn("列表页面btn的name选项未配置:", key, t.ExtInfo)
			continue
		}

		//覆盖删除按钮
		cover := false
		if info.Name == "del" {
			t.BtnDel = true
			cover = true
		}

		//desc and if
		info.DESC = getSubConContent(key, "desc")(t.ExtInfo)
		if info.DESC == "" {
			logs.Log.Warn("列表页面btn的desc选项未配置:", t.ExtInfo)
			continue
		}

		//confirm
		info.Confirm = getSubConContent(key, "confirm")(t.ExtInfo)

		//url
		info.URL = getSubConContent(key, "url")(t.ExtInfo)

		//condition
		info.Condition = translateCondition(getSubConContent(key, "condition")(t.ExtInfo))

		//cover
		if cover {
			t.BtnInfo = append(t.BtnInfo, info)
			continue
		}

		//key
		info.KeyWord = types.GetString(getSubConContent(key, "key")(t.ExtInfo), key)

		//Rows
		for _, v := range getRows(info.KeyWord)(t.Rows) {
			info.Rows = append(info.Rows, *v)
		}

		if len(info.Rows) < 1 {
			logs.Log.Warn("列表页面btn的绑定的字段未配置")
		}

		t.BtnInfo = append(t.BtnInfo, info)
	}

}

//DispostELBtnDetail {el_btn_detail(name:funcName,desc:1-启用|2-禁用,confirm:你确定进行修改吗,field_name:name,condition:condition,key:sa)}
func (t *Table) DispostELBtnDetail() {
	if t.ExtInfo == "" {
		return
	}
	a := []string{"", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
Next:
	for _, v := range a {
		key := fmt.Sprintf("el_btn_detail%s", v)
		if !strings.Contains(t.ExtInfo, key+"(") {
			break
		}

		info := newBtnInfo()

		//name
		info.Name = getSubConContent(key, "name")(t.ExtInfo)
		if info.Name == "" {
			logs.Log.Warn("详情页面btn的name选项未配置:", key, t.ExtInfo)
			continue
		}

		//desc
		info.DESC = getSubConContent(key, "desc")(t.ExtInfo)
		if info.DESC == "" {
			logs.Log.Warn("详情页面btn的desc选项未配置:", t.ExtInfo)
			continue
		}

		//confirm
		info.Confirm = getSubConContent(key, "confirm")(t.ExtInfo)

		//url
		info.URL = getSubConContent(key, "url")(t.ExtInfo)

		//condition
		info.Condition = translateCondition(getSubConContent(key, "condition")(t.ExtInfo))

		//key
		info.KeyWord = types.GetString(getSubConContent(key, "key")(t.ExtInfo), key)

		//field_name
		info.FieldName = getSubConContent(key, "field_name")(t.ExtInfo)

		//Rows
		for _, v := range getRows(info.KeyWord)(t.Rows) {
			info.Rows = append(info.Rows, *v)
		}

		if info.FieldName == "" && len(info.Rows) == 1 {
			info.FieldName = info.Rows[0].Name
		}

		for _, v := range t.Rows {
			if v.Name == info.FieldName {
				v.DetailBtnInfo = append(v.DetailBtnInfo, info)
				break Next
			}
		}

		logs.Log.Warn("详情页面btn未配置按钮字段")
	}

}

//DispostELBtnQuery {el_btn_query(name:funcName,desc:desc,confirm:你确定进行修改吗,url:xxxx,condition:condition)}
func (t *Table) DispostELBtnQuery() {
	if t.ExtInfo == "" {
		return
	}
	a := []string{"", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

	for _, v := range a {
		key := fmt.Sprintf("el_btn_query%s", v)
		if !strings.Contains(t.ExtInfo, key+"(") {
			break
		}

		info := newBtnInfo()

		//name
		info.Name = getSubConContent(key, "name")(t.ExtInfo)
		if info.Name == "" {
			logs.Log.Warn("查询的btn的name选项未配置:", key, t.ExtInfo)
			continue
		}

		//desc
		info.DESC = getSubConContent(key, "desc")(t.ExtInfo)
		if info.DESC == "" {
			logs.Log.Warn("查询的btn的desc选项未配置:", t.ExtInfo)
			continue
		}

		//confirm
		info.Confirm = getSubConContent(key, "confirm")(t.ExtInfo)

		//url
		info.URL = getSubConContent(key, "url")(t.ExtInfo)

		//condition
		info.Condition = translateCondition(getSubConContent(key, "condition")(t.ExtInfo))

		if info.Name == "queryDatas" {
			t.BtnShowQuery = true
			t.QueryURL = info.URL
			info.IsQuery = true
		}

		t.QueryBtnInfo = append(t.QueryBtnInfo, info)
	}

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

type SelectInfo struct {
	URL       string
	Name      string
	Desc      string
	Condition string
	Confirm   string
}

func (t *Table) DispostELSelect() {
	key := "el_select"
	if t.ExtInfo == "" || !strings.Contains(t.ExtInfo, key) {
		return
	}
	t.SelectInfo = &SelectInfo{}
	t.SelectInfo.URL = getSubConContent(key, "url")(t.ExtInfo)
	t.SelectInfo.Name = getSubConContent(key, "name")(t.ExtInfo)
	t.SelectInfo.Desc = getSubConContent(key, "desc")(t.ExtInfo)
	t.SelectInfo.Condition = translateCondition(getSubConContent(key, "condition")(t.ExtInfo))
	//confirm
	t.SelectInfo.Confirm = getSubConContent(key, "confirm")(t.ExtInfo)
}

func (a *SelectInfo) IsEmpty() bool {
	return a == nil || reflect.DeepEqual(a, &SelectInfo{})
}

type ListComponents struct {
	Name      string
	Path      string
	BtnName   string
	Condition string
}

func (t *Table) DispostELListComponents() {
	key := "el_components"
	if t.ExtInfo == "" || !strings.Contains(t.ExtInfo, key) {
		return
	}
	c := getBracketContent([]string{key})(t.ExtInfo)
	tabs := strings.Split(c, "|")
	if len(tabs) == 0 {
		return
	}
	t.ListComponents = make([]*ListComponents, 0)
	for _, v := range tabs {
		info := &ListComponents{}
		tab := strings.Split(v, ",")
		if len(tab) < 2 {
			logs.Log.Warn("列表页面components的选项配置不正确(name,path,btn_name,condition|name,path...)", key, t.ExtInfo)
			continue
		}

		//name
		info.Name = tab[0]
		if info.Name == "" {
			logs.Log.Warn("列表页面components的name选项未配置:", key, t.ExtInfo)
			continue
		}
		if info.Name == "Edit" {
			t.BtnShowEdit = true
		}
		if info.Name == "Detail" {
			t.BtnShowDetail = true
		}

		info.Path = tab[1]
		if info.Path == "" {
			logs.Log.Warn("列表页面components的path选项未配置:", key, t.ExtInfo)
			continue
		}
		if len(tab) > 2 {
			info.BtnName = tab[2]
		}
		if len(tab) > 3 {
			info.Condition = translateCondition(tab[3])
		}
		t.ListComponents = append(t.ListComponents, info)
	}
}

type QueryComponents struct {
	Name    string
	Path    string
	BtnName string
}

func (t *Table) DispostELQueryComponents() {
	key := "el_query_components"
	if t.ExtInfo == "" || !strings.Contains(t.ExtInfo, key) {
		return
	}
	c := getBracketContent([]string{key})(t.ExtInfo)
	tabs := strings.Split(c, "|")
	if len(tabs) == 0 {
		return
	}
	t.QueryComponents = make([]*QueryComponents, 0)
	for _, v := range tabs {
		info := &QueryComponents{}
		tab := strings.Split(v, ",")
		if len(tab) < 2 {
			logs.Log.Warn("页面查询components的选项配置不正确(name,path,btn_name|name,path...)", key, t.ExtInfo)
			continue
		}

		//name
		info.Name = tab[0]
		if info.Name == "" {
			logs.Log.Warn("页面查询components的name选项未配置:", key, t.ExtInfo)
			continue
		}
		info.Path = tab[1]
		if info.Path == "" {
			logs.Log.Warn("页面查询components的path选项未配置:", key, t.ExtInfo)
			continue
		}
		if len(tab) > 2 {
			info.BtnName = tab[2]
		}
		t.QueryComponents = append(t.QueryComponents, info)
	}
}

func translateCondition(c string) string {
	c = strings.Replace(c, " and ", " && ", -1)
	c = strings.Replace(c, " or ", " || ", -1)
	return c
}
