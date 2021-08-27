package tmpl

import (
	"strings"

	"github.com/micro-plat/hicli/markdown/const/enums"
)

//Tables markdown中的所有表信息
type Tables struct {
	PKG        string
	RawTables  []*Table
	Tbs        []*Table
	TableNames map[string]bool
	Drop       bool
	SEQFile    bool
}

//FilterByKW 过滤行信息
func (t *Tables) FilterByKW(kwc string) {
	if kwc == "" {
		return
	}
	kws := strings.Split(kwc, ",")
	tbs := make([]*Table, 0, len(kws))
	for _, v := range kws {
		for _, tb := range t.RawTables {
			if strings.Contains(tb.Name, v) {
				tbs = append(tbs, tb)
			}
		}
	}
	t.Tbs = tbs
}

func (t *Tables) Exclude() {
	tbs := make([]*Table, 0, 1)
	for _, tb := range t.Tbs {
		if !tb.Exclude {
			tbs = append(tbs, tb)
		}
	}
	t.Tbs = tbs
}

//BuildSEQFile 生成seq_ids脚本
func (t *Tables) BuildSEQFile(d bool) {
	t.SEQFile = d
}

//DropTable 如果表存在是否删除
func (t *Tables) DropTable(d bool) {
	t.Drop = d
	for _, tb := range t.RawTables {
		tb.Drop = d
	}
	for _, tb := range t.Tbs {
		tb.Drop = d
	}
}

//SetPkg 添加行信息
func (t *Tables) SetPkg(path string) {

	t.PKG = getPKSName(path)

	for _, tb := range t.RawTables {
		tb.SetPkg(t.PKG)
	}
	for _, tb := range t.Tbs {
		tb.SetPkg(t.PKG)
	}
}

func (source *Tables) getTableMap() map[string]*Table {
	m := make(map[string]*Table, len(source.RawTables))
	for _, t := range source.RawTables {
		m[t.Name] = t
	}
	return m
}

//Diff 比较两个不同
func (source *Tables) Diff(target *Tables) *Tables {
	sourceM := source.getTableMap()
	targetM := target.getTableMap()

	diff := &Tables{
		Tbs: make([]*Table, 0),
	}

	//新增
	for tname, table := range sourceM {
		if _, ok := targetM[tname]; !ok {
			table.Operation = enums.DiffInsert
			diff.Tbs = append(diff.Tbs, table)
			delete(sourceM, tname)
		}
	}

	//减少
	for tname, table := range targetM {
		if _, ok := sourceM[tname]; !ok {
			table.Operation = enums.DiffDelete
			diff.Tbs = append(diff.Tbs, table)
			delete(targetM, tname)
		}
	}

	//变动
	for name, sourceTable := range sourceM {
		targetTable := targetM[name]
		//字段变化
		diffRows := sourceTable.Rows.Diff(targetTable.Rows)
		if len(diffRows) > 0 {
			sourceTable.DiffRows = diffRows
		}
		//索引变化
		sourceTable.GetIndexs()
		targetTable.GetIndexs()
		diffIndex := sourceTable.Indexs.Diff(targetTable.Indexs)
		if len(diffIndex) > 0 {
			sourceTable.DiffIndexs = diffIndex
		}
		if len(diffRows) > 0 || len(diffIndex) > 0 {
			sourceTable.Operation = enums.DiffModify
			diff.Tbs = append(diff.Tbs, sourceTable)
		}
	}

	return diff
}
