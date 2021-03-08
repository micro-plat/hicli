package tmpl

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/micro-plat/lib4go/types"
)

//Table 表名称
type Table struct {
	Name         string //表名
	Desc         string //表描述
	PKG          string //包名称
	Drop         bool   //创建表前是否先删除
	DBType       string //数据库类型
	DBLink       string //
	Rows         []*Row
	RawRows      []*Row
	Indexs       Indexs
	BasePath     string   //生成项目基本路径
	AllTables    []*Table //所有表
	Exclude      bool     //排除生成sql
	ELTableIndex int
}

//Row 行信息
type Row struct {
	Name       string //字段名
	Type       string //类型
	Def        string //默认值
	IsNull     string //为空
	Con        string //约束
	Desc       string //描述
	Len        int    //类型长度
	DecimalLen int    //小数长度
}

//Indexs 索引集
type Indexs map[string]*Index

//Index 索引
type Index struct {
	fields fields
	Name   string
	Type   string
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

//NewTable 创建表
func NewTable(name, desc string) *Table {
	return &Table{
		Name:    strings.TrimLeft(name, "^"),
		Desc:    desc,
		Rows:    make([]*Row, 0, 1),
		RawRows: make([]*Row, 0, 1),
		Exclude: strings.Contains(name, "^"),
	}
}

//AddRow 添加行信息
func (t *Table) AddRow(r *Row) error {
	t.Rows = append(t.Rows, r)
	t.RawRows = append(t.RawRows, r)
	return nil
}

//SetPkg 添加行信息
func (t *Table) SetPkg(path string) {
	t.PKG = getPKSName(path)
}

//SetBasePath 添加行信息
func (t *Table) SetBasePath(BasePath string) {
	t.BasePath = BasePath
}

//GetPKS 获取主键列表
func (t *Table) GetPKS() []string {
	indexs := t.GetIndexs()
	for _, index := range indexs {
		if index.Type == "pk" {
			return index.fields.List()
		}
	}
	return nil
}

//FilterRowByKW	过滤行信息
func (t *Table) FilterRowByKW(kwc string) {
	if kwc == "" {
		return
	}
	rows := make([]*Row, 0, 1)
	for _, row := range t.RawRows {
		if getKWCons(row.Con, kwc) {
			rows = append(rows, row)
		}
	}
	t.Rows = rows
}

//GetIndexs 获取所有索引信息
func (t *Table) GetIndexs() Indexs {
	if t.Indexs != nil {
		return t.Indexs
	}
	indexs := map[string]*Index{}
	for ri, r := range t.Rows {
		t.getIndex(indexs, r, ri, "idx")
		t.getIndex(indexs, r, ri, "unq")
		t.getIndex(indexs, r, ri, "pk")
	}
	for _, index := range indexs {
		sort.Sort(index.fields)
	}
	t.Indexs = indexs
	return t.Indexs
}
func (t *Table) getIndex(indexs map[string]*Index, row *Row, ri int, tp string) {
	ok, name, index := getIndex(row.Con, tp)
	if !ok {
		return
	}
	if name == "" {
		name = row.Name
	}
	index = types.DecodeInt(index, 0, ri)
	if v, ok := indexs[name]; ok {
		v.fields = append(v.fields, &Field{Name: row.Name, Index: index})
		return
	}
	indexs[name] = &Index{Name: name, Type: tp, fields: []*Field{{Name: row.Name, Index: index}}}
}

func (t *Table) String() string {
	buff := strings.Builder{}
	buff.WriteString(t.Name)
	buff.WriteString("(")
	buff.WriteString(t.Desc)
	buff.WriteString(")")
	buff.WriteString("\n")
	for _, c := range t.Rows {
		buff.WriteString(fmt.Sprintf("%s\t%s\t%s\t%s\t%v\t%s\n", c.Name, c.Type, c.Con, c.Def, c.IsNull, c.Desc))

	}
	return buff.String()
}

//Translate 翻译模板
func Translate(c string, tp string, input interface{}) (string, error) {
	var tmpl = template.New("table").Funcs(getfuncs(tp))
	np, err := tmpl.Parse(c)
	if err != nil {
		return "", err
	}
	buff := bytes.NewBufferString("")
	if err := np.Execute(buff, input); err != nil {
		return "", err
	}
	return strings.Replace(strings.Replace(buff.String(), "{###}", "`", -1), "&#39;", "'", -1), nil
}

//GetFilePath 获取文件路径
func GetFilePath(root string, name string, ext ...string) string {
	ex := types.GetStringByIndex(ext, 0, "vue")
	path, _ := Translate(fmt.Sprintf("{{.|rmhd|fpath}}.%s", ex), "", name)
	return filepath.Join(types.GetString(root, "."), path)
}

//GetFileName 获取文件名称
func GetFileName(root string, name string, prefix string, ext ...string) string {
	ex := types.GetStringByIndex(ext, 0, "go")
	path, _ := Translate(fmt.Sprintf("{{.|rmhd|l2d}}.%s", ex), "", name)
	return filepath.Join(types.GetString(root, "."), fmt.Sprintf("%s%s", prefix, path))
}
