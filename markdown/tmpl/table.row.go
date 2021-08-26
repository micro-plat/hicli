package tmpl

import (
	"github.com/micro-plat/hicli/markdown/const/enums"
)

//Row 行信息
type Row struct {
	Name          string //字段名
	Type          string //类型
	Def           string //默认值
	IsNull        string //为空
	Con           string //约束
	Desc          string //描述
	Len           int    //类型长度
	LenStr        string
	DecimalLen    int //小数长度
	LineID        int
	Sort          int //字段在列表中排序位置
	Disable       bool
	SQLAliasName  string //SQL别名
	IsInput       bool
	Operation     enums.Operation
	DetailBtnInfo []*BtnInfo //详情页按钮
}

func (s *Row) Equal(t *Row) bool {
	return s.Name == t.Name && s.Type == t.Type && s.Def == t.Def && s.IsNull == t.IsNull && s.Desc == t.Desc && s.Len == t.Len
}

//TableColumn 表的列排序用
type TableColumn []*Row

func (t TableColumn) Len() int {
	return len(t)
}

//从低到高
func (t TableColumn) Less(i, j int) bool {
	return t[i].Sort < t[j].Sort
}

func (t TableColumn) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (source TableColumn) getColumnMap() map[string]*Row {
	m := make(map[string]*Row, len(source))
	for _, c := range source {
		m[c.Name] = c
	}
	return m
}

func (source TableColumn) Diff(target TableColumn) []*Row {
	sourceM := source.getColumnMap()
	targetM := target.getColumnMap()

	diff := make([]*Row, 0)

	//新增
	for name, column := range sourceM {
		if _, ok := targetM[name]; !ok {
			column.Operation = enums.DiffInsert
			diff = append(diff, column)
			delete(sourceM, name)
		}
	}

	//减少
	for name, column := range targetM {
		if _, ok := sourceM[name]; !ok {
			column.Operation = enums.DiffDelete
			diff = append(diff, column)
			delete(targetM, name)
		}
	}

	//变动
	for name, sourceColumn := range sourceM {
		if !sourceColumn.Equal(targetM[name]) {
			sourceColumn.Operation = enums.DiffModify
			diff = append(diff, sourceColumn)
		}
	}

	return diff
}
