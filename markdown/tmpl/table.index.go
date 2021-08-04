package tmpl

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/micro-plat/hicli/markdown/const/enums"
)

//Indexs 索引集
type Indexs map[string]*Index

//Index 索引
type Index struct {
	fields    fields
	Name      string
	TableName string
	Type      enums.IndexType
	Operation enums.Operation
}
type fields []*Field

//Field 字段信息
type Field struct {
	Name  string
	Index int
}

//List 获取所有字段名的列表
func (t fields) List() []string {
	list := make([]string, 0, len(t))
	for _, fi := range t {
		list = append(list, fi.Name)
	}
	return list
}

//Len 字段个数
func (t fields) Len() int {
	return len(t)
}

//Join 指定连接符，将字段名连接为一个长字符串
func (t fields) Join(s string) string {
	list := t.List()
	return strings.Join(list, s)
}

//从低到高
func (t fields) Less(i, j int) bool {
	if t[i].Index < t[j].Index {
		return true
	}
	if t[i].Index == t[j].Index {
		return true
	}
	return false
}

func (t fields) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (s *Index) JudgeType(t enums.IndexType) bool {
	return s.Type == t
}

func (s *Index) String(tp string) string {
	switch tp {
	case MYSQL:
		switch s.Type {
		case enums.IndexNor:
			return fmt.Sprintf("index %s(%s)", s.Name, s.fields.Join(","))
		case enums.IndexUnq:
			return fmt.Sprintf("unique index %s(%s)", s.Name, s.fields.Join(","))
		case enums.IndexPK:
			return fmt.Sprintf("primary key (%s)", s.fields.Join(","))
		}
	case ORACLE:
		switch s.Type {
		case enums.IndexNor:
			return fmt.Sprintf("index %s on %s(%s)", s.Name, s.TableName, s.fields.Join(","))
		case enums.IndexUnq:
			return fmt.Sprintf("constraint %s unique(%s)", s.Name, s.fields.Join(","))
		case enums.IndexPK:
			return fmt.Sprintf("constraint pk_%s primary key (%s)", s.Name, s.fields.Join(","))
		}
	}
	return ""
}

func (s *Index) Equal(t *Index) bool {
	if s.Name != t.Name {
		return false
	}
	if s.Type != t.Type {
		return false
	}
	return reflect.DeepEqual(s.fields, t.fields)
}

func (source Indexs) Diff(target Indexs) []*Index {
	tempSource := source
	diff := make([]*Index, 0)

	//减少,索引要处理删除
	for name, index := range target {
		if _, ok := tempSource[name]; !ok {
			index.Operation = enums.DiffDelete
			diff = append(diff, index)
		}
	}

	//新增
	for name, index := range tempSource {
		if _, ok := target[name]; !ok {
			index.Operation = enums.DiffInsert
			diff = append(diff, index)
			delete(tempSource, name)
		}
	}

	//变动
	for name, sourceIndex := range tempSource {
		if !sourceIndex.Equal(target[name]) {
			sourceIndex.Operation = enums.DiffModify
			diff = append(diff, sourceIndex)
		}
	}

	return diff
}
