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
	TabTable         map[string]bool   //主页
	TabTableList     map[string]bool   //详情tab关联字段
	TabTablePreField map[string]string //详情tab关联字段
	TabTableProField map[string]string //详情tab关联字段
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

//DisposeTabTables 处理前端详情页
func (t *Table) DisposeTabTables() {
	if t.ExtInfo == "" {
		return
	}

	c := getBracketContent([]string{"el_tab"})(t.ExtInfo)
	tabs := strings.Split(c, "|")
	if len(tabs) == 0 {
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
				if tabList == "list" {
					tb.TabInfo.TabTableList[t.Name] = true
					tb.TabInfo.TabListField[tabField[1]] = tabField[1]
				} else {
					tb.TabInfo.TabTable[t.Name] = true
					tb.TabInfo.TabField[tabField[1]] = tabField[1]
				}
				tb.TabInfo.TabTablePreField[t.Name] = tabField[0]
				tb.TabInfo.TabTableProField[t.Name] = tabField[1]
				t.TabTables = append(t.TabTables, tb)
				exist = true
				break
			}
		}
		if !exist {
			logs.Log.Warn("tab表名不正确：", tabName)
		}
	}
}

type BtnInfo struct {
	Name               string
	DESC               string
	VIF                []*VIF
	KeyWord            string
	Confirm            string
	URL                string
	Table              []*Table
	Rows               []Row
	RelativeShelfFiled map[string]string
	RelativeFiled      map[string]string
	LastRowIndex       string
}
type VIF struct {
	IfName string
	IfDESC string
}

func newBtnInfo() *BtnInfo {
	return &BtnInfo{
		RelativeShelfFiled: make(map[string]string),
		RelativeFiled:      make(map[string]string),
	}
}

//DispostBtnTables {el_btn(name:funcName,desc:1-启用|2-禁用,confirm:你确定进行修改吗,table:adas/iqe,key:sa)}
func (t *Table) DispostBtnTables() {
	if t.ExtInfo == "" {
		return
	}
	a := []string{"", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

	for _, v := range a {
		key := fmt.Sprintf("el_btn%s", v)
		if !strings.Contains(t.ExtInfo, key) {
			break
		}

		info := newBtnInfo()

		//name
		info.Name = getSubConContent(key, "name")(t.ExtInfo)
		if info.Name == "" {
			logs.Log.Warn("列表页面btn的name选项未配置:", key, t.ExtInfo)
			continue
		}

		//desc and if
		desc := getSubConContent(key, "desc")(t.ExtInfo)
		if desc == "" {
			logs.Log.Warn("列表页面btn的desc选项未配置:", t.ExtInfo)
			continue
		}

		if strings.Contains(desc, "|") {
			for _, v := range strings.Split(desc, "|") {
				pos := strings.Index(v, "-")
				if pos < 0 {
					logs.Log.Warn("列表页面btn的if选项不正确:", desc)
					continue
				}
				info.VIF = append(info.VIF, &VIF{
					IfName: v[:pos],
					IfDESC: v[pos+1:],
				})
			}

			if len(info.VIF) < 2 {
				logs.Log.Warn("列表页面btn的if选项最少为2个：", desc)
			}
		} else {
			info.DESC = desc
		}

		//confirm
		info.Confirm = getSubConContent(key, "confirm")(t.ExtInfo)

		info.URL = getSubConContent(key, "url")(t.ExtInfo)

		//key
		info.KeyWord = types.GetString(getSubConContent(key, "key")(t.ExtInfo), key)

		//table
		tabs := getSubConContent(key, "table")(t.ExtInfo)

		for _, v := range strings.Split(tabs, "|") {

			tabName := v
			tabField := make([]string, 2)
			if pos := strings.Index(v, ":"); pos > 0 {
				tabName = v[0:pos]
				t := strings.Split(v[pos+1:], "/")
				if len(t) == 1 {
					tabField = []string{t[0], t[0]}
				}
				if len(t) == 2 {
					tabField = []string{t[0], t[1]}
				}
			}

			for _, tb := range t.AllTables {
				if tb.Name == tabName {
					info.RelativeShelfFiled[tb.Name] = tabField[0]
					info.RelativeFiled[tb.Name] = tabField[1]
					info.Table = append(info.Table, tb)
				}
			}
		}

		//Rows
		for _, v := range getRows(info.KeyWord)(t.Rows) {
			v.BelongTable = t
			info.Rows = append(info.Rows, *v)
			info.LastRowIndex = v.Name
		}

		for k, v := range info.Table {
			for _, v1 := range getRows(info.KeyWord)(v.Rows) {
				v1.BelongTable = v
				v1.Disable = true
				v1.SQLAliasName = fmt.Sprintf("t%d", k)
				info.Rows = append(info.Rows, *v1)
			}
		}
		if len(info.Rows) < 1 {
			logs.Log.Warn("列表页面btn的更新的字段未配置")
			continue
		}

		t.BtnInfo = append(t.BtnInfo, info)
	}

}

type DownloadInfo struct {
	Title []string
}

func (t *Table) DispostDownloadTables() {
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
	t.DownloadInfo.Title = strings.Split(c, ",")
}

type SelectInfo struct {
	URL string
}

func (t *Table) DispostSelectTables() {
	if t.ExtInfo == "" {
		return
	}
	t.SelectInfo = &SelectInfo{}
	t.SelectInfo.URL = getSubConContent("el_select", "url")(t.ExtInfo)
}

func (a *SelectInfo) IsEmpty() bool {

	return reflect.DeepEqual(a, &SelectInfo{})

}
