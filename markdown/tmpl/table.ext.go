package tmpl

import (
	"strings"

	logs "github.com/lib4dev/cli/logger"
)

type TabInfo struct {
	TabField         map[string]string
	TabListField     map[string]string //详情是否生成list
	TabTable         map[string]bool   //详情tab关联字段
	TabTableList     map[string]bool   //详情tab关联字段
	TabTablePreField map[string]string //详情tab关联字段
	TabTableProField map[string]string //详情tab关联字段
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
				tb.TabInfo = &TabInfo{
					TabField:         make(map[string]string),
					TabListField:     make(map[string]string),
					TabTable:         make(map[string]bool),
					TabTableList:     make(map[string]bool),
					TabTablePreField: make(map[string]string),
					TabTableProField: make(map[string]string),
				}
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
