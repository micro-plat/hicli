package tmpl

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/hicli/markdown/const/enums"
	"github.com/micro-plat/hicli/markdown/utils"
	"github.com/micro-plat/lib4go/types"
)

//Table 表名称
type Table struct {
	Name          string          //表名
	Desc          string          //表描述
	ExtInfo       string          //扩展信息
	PKG           string          //包名称
	Drop          bool            //创建表前是否先删除
	DBType        string          //数据库类型
	DBLink        string          //dblink
	Rows          TableColumn     //行
	RawRows       []*Row          //原始行
	Indexs        Indexs          //序列
	BasePath      string          //生成项目基本路径
	AllTables     []*Table        //所有表
	DiffRows      TableColumn     //差异对比行
	DiffIndexs    []*Index        //差异对比序列
	Operation     enums.Operation //差异对比
	Exclude       bool            //排除生成sql
	DBObjectName  string          //数据库对象名
	ELTableIndex  int             //前端数据序列
	TabTables     []*Table        //前端详情切换的tab页对应表
	TabInfo       *TabInfo        //前端
	ListBtnInfo   []*BtnInfo      //前端
	QueryBtnInfo  []*BtnInfo      //前端
	DownloadInfo  *DownloadInfo   //前端
	BatchInfo     BatchInfos      //前端
	ListDialogs   []*Dialog       //前端
	QueryDialogs  []*Dialog       //前端
	BtnShowEdit   bool            //前端
	BtnShowQuery  bool            //前端
	QueryHandler  string          //前端
	BtnShowAdd    bool            //前端
	BtnShowDetail bool            //前端
	BtnDel        bool            //前端\
}

//NewTable 创建表
func NewTable(name, desc, extinfo string) *Table {
	return &Table{
		Name:    strings.TrimLeft(name, "^"),
		Desc:    desc,
		Rows:    make([]*Row, 0, 1),
		RawRows: make([]*Row, 0, 1),
		Exclude: strings.Contains(name, "^"),
		ExtInfo: extinfo,
		TabInfo: newTableInfo(),
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
	if path == "" {
		path = utils.GetProjectPath(path)
	}

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
	logs.Log.Errorf("[%s]主键未配置", t.Name)
	return nil
}

//SetELTableIndex 设置前端列表索引
func (t *Table) SetELTableIndex() {
	if t.ExtInfo == "" {
		return
	}
	c := getBracketContent([]string{"el_index"})(t.ExtInfo)
	t.ELTableIndex = types.GetInt(c)

}

//SetAllTables 关联所有表
func (t *Table) SetAllTables(tbs []*Table) {
	t.AllTables = tbs
}

//SortRows 行排序
func (t *Table) SortRows() {
	sorts := make(map[string]int, len(t.Rows))
	cons := make(map[string]string, len(t.Rows))
	for _, v := range t.Rows {
		v.Sort = v.LineID * 1000 //取1000倍
		sorts[v.Name] = v.Sort
		cons[v.Name] = v.Con
	}
	for k, v := range t.Rows {
		after := getBracketContent([]string{"after"})(v.Con)
		if after == "" {
			continue
		}
		if after == "0" {
			t.Rows[k].Sort = 0
			continue
		}
		if s, ok := sorts[after]; ok {
			t.Rows[k].Sort = t.getSortByAfter(cons, sorts, after, s)
		}
	}
	sort.Sort(t.Rows)
}

func (t *Table) getSortByAfter(cons map[string]string, sorts map[string]int, rowName string, sort int) int {
	con, ok := cons[rowName]
	if !ok {
		return sort + 1
	}

	rowName = getBracketContent([]string{"after"})(con)
	if rowName == "" {
		return sort + 1
	}
	if rowName == "0" {
		return 0 + 1
	}

	if s, ok := sorts[rowName]; ok {
		ts := t.getSortByAfter(cons, sorts, rowName, s)
		return ts + 1
	}

	return sort + 1
}

//FilterRowByKW 过滤行信息
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
		t.getIndex(indexs, r, ri, enums.IndexNor)
		t.getIndex(indexs, r, ri, enums.IndexUnq)
		t.getIndex(indexs, r, ri, enums.IndexPK)
	}
	for _, index := range indexs {
		sort.Sort(index.fields)
	}
	t.Indexs = indexs
	return t.Indexs
}
func (t *Table) getIndex(indexs map[string]*Index, row *Row, ri int, tp enums.IndexType) {
	ok, name, i, _ := getCapturingGroup(row.Con, string(tp))
	if !ok {
		return
	}
	name = types.DecodeString(name, "", row.Name)
	index := types.DecodeInt(types.GetInt(i), 0, ri)
	if v, ok := indexs[name]; ok {
		v.fields = append(v.fields, &Field{Name: row.Name, Index: index})
		return
	}
	indexs[name] = &Index{
		Name:      name,
		TableName: t.Name,
		Type:      tp,
		fields:    []*Field{{Name: row.Name, Index: index}},
	}
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

//GetFilePath 获取后端文件路径
func GetFilePath(root string, name string, ext ...string) string {
	ex := types.GetStringByIndex(ext, 0, "go")
	path, _ := Translate(fmt.Sprintf("{{.|rmhd|fpath}}.%s", ex), "", name)
	return filepath.Join(types.GetString(root, "."), path)
}

//GetWebFilePath 获取前端文件路径
func GetWebFilePath(root string, name string, ext ...string) string {
	ex := types.GetStringByIndex(ext, 0, "vue")
	path, _ := Translate(fmt.Sprintf("{{.|rmhd|webfpath}}.%s", ex), "", name)
	return filepath.Join(types.GetString(root, "."), path)
}

//GetFileName 获取文件名称
func GetFileName(root string, name string, prefix string, ext ...string) string {
	ex := types.GetStringByIndex(ext, 0, "go")
	path, _ := Translate(fmt.Sprintf("{{.|rmhd|l2d}}.%s", ex), "", name)
	return filepath.Join(types.GetString(root, "."), fmt.Sprintf("%s%s", prefix, path))
}
