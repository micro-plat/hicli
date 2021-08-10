package tmpl

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/micro-plat/lib4go/types"
)

//Line 每一行信息
type Line struct {
	Text   string
	LineID int
}

//TableLine 表的每一行
type TableLine struct {
	Lines [][]*Line
}

//Markdown2DB 读取markdown文件并转换为MarkDownDB对象
func Markdown2DB(fn string) (*Tables, error) {
	lines, err := readMarkdown(fn)
	if err != nil {
		return nil, err
	}

	return tableLine2Table(line2TableLine(lines))
}

//Markdown2DB 读取markdown文件并转换为MarkDownDB对象
func Markdowns2DB(filePath string) (*Tables, error) {
	if !strings.Contains(filePath, "*") {
		return Markdown2DB(filePath)
	}

	fns := getAllMatchMD(filePath)
	//读取文件
	baseTable := &Tables{
		TableNames: make(map[string]bool),
	}
	for _, fn := range fns {
		newTable, err := Markdown2DB(fn)
		if err != nil {
			return nil, err
		}
		for key := range newTable.TableNames {
			if _, ok := baseTable.TableNames[key]; ok {
				return nil, fmt.Errorf("存在相同的表名：%s", key)
			}
			baseTable.TableNames[key] = true
		}
		baseTable.RawTables = append(baseTable.RawTables, newTable.Tbs...)
	}
	baseTable.Tbs = baseTable.RawTables
	return baseTable, nil
}

//readMarkdown 读取md文件
func readMarkdown(name string) ([]*Line, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return readMarkdownByReader(bufio.NewReader(f))
}

//readMarkdown 读取md文件
func readMarkdownByReader(rd *bufio.Reader) ([]*Line, error) {
	lines := make([]*Line, 0, 64)
	num := 0
	for {
		num++
		line, err := rd.ReadString('\n')
		if line == "" && (err != nil || io.EOF == err) {
			break
		}
		line = strings.Trim(strings.Trim(line, "\n"), "\t")
		if strings.TrimSpace(line) == "" {
			continue
		}
		lines = append(lines, &Line{Text: line, LineID: num})
	}
	return lines, nil
}

//lines2TableLine 数据行转变为以表为单个整体的数据行
func line2TableLine(lines []*Line) (tl TableLine) {
	dlines := []int{}
	for i, line := range lines {
		text := strings.TrimSpace(strings.Replace(line.Text, " ", "", -1))
		if text == "|字段名|类型|默认值|为空|约束|描述|" {
			dlines = append(dlines, i-1)
		}
		if len(dlines)%2 == 1 && strings.Count(text, "|") != 7 {
			dlines = append(dlines, i-1)
		}
	}
	if len(dlines)%2 == 1 {
		dlines = append(dlines, len(lines)-1)
	}
	//划分为以一张表为一个整体
	for i := 0; i < len(dlines); i = i + 2 {
		tl.Lines = append(tl.Lines, lines[dlines[i]:dlines[i+1]+1])
	}
	return tl
}

//tableLine2Table 表数据行变为表
func tableLine2Table(lines TableLine) (tables *Tables, err error) {
	tables = &Tables{Tbs: make([]*Table, 0, 1), TableNames: make(map[string]bool)}
	for _, tline := range lines.Lines {
		//markdown表格的表名，标题，标题数据区分行，共三行
		if len(tline) <= 3 {
			continue
		}
		var tb *Table
		for i, line := range tline {
			if i == 0 {
				//获取表名，描述名称
				name, err := getTableName(line)
				if err != nil {
					return nil, err
				}
				tb = NewTable(name, getTableDesc(line), getTableExtInfo(line))
				continue
			}
			if i < 3 {
				continue
			}
			c, err := line2TableRow(line)
			if err != nil {
				return nil, err
			}

			if err := tb.AddRow(c); err != nil {
				return nil, err
			}
		}
		if tb != nil {
			if _, ok := tables.TableNames[tb.Name]; ok {
				return nil, fmt.Errorf("存在相同的表名：%s", tb.Name)
			}
			tables.TableNames[tb.Name] = true
			tables.RawTables = append(tables.RawTables, tb)
		}
	}
	tables.Tbs = tables.RawTables
	return tables, nil
}

func line2TableRow(line *Line) (*Row, error) {
	if strings.Count(line.Text, "|") != 7 {
		return nil, fmt.Errorf("表结构有误(行:%d)", line.LineID)
	}
	colums := strings.Split(strings.Trim(line.Text, "|"), "|")
	if colums[0] == "" {
		return nil, fmt.Errorf("字段名称不能为空 %s(行:%d)", line.Text, line.LineID)
	}

	tp, len, decimalLen, err := getType(line)
	if err != nil {
		return nil, err
	}
	c := &Row{
		LineID:     line.LineID,
		Name:       strings.TrimSpace(strings.Replace(colums[0], "&#124;", "|", -1)),
		Type:       tp,
		Len:        len,
		DecimalLen: decimalLen,
		Def:        strings.TrimSpace(strings.Replace(colums[2], "&#124;", "|", -1)),
		IsNull:     strings.TrimSpace(colums[3]),
		Con:        strings.TrimSpace(colums[4]), // strings.Replace(strings.TrimSpace(colums[4]), " ", "", -1),
		Desc:       strings.TrimSpace(strings.Replace(colums[5], "&#124;", "|", -1)),
	}
	return c, nil
}

func getTableDesc(line *Line) string {
	reg := regexp.MustCompile(`[^\d\.|\s]?[^\x00-\xff]+[^\[]+`)
	names := reg.FindAllString(line.Text, -1)
	if len(names) == 0 {
		return ""
	}
	return strings.TrimSpace(names[0])
}

func getTableExtInfo(line *Line) string {
	reg := regexp.MustCompile(`{(.*?)}$`)
	names := reg.FindStringSubmatch(line.Text)
	if len(names) == 0 {
		return ""
	}
	return strings.TrimSpace(names[1])
}

func getTableName(line *Line) (string, error) {
	if !strings.HasPrefix(line.Text, "###") {
		return "", fmt.Errorf("%d行表名称标注不正确，请以###开头:(%s)", line.LineID, line.Text)
	}

	reg := regexp.MustCompile(`\[[\^]?[\w]+[,]?[\p{Han}A-Za-z0-9_]+\]`)
	names := reg.FindAllString(line.Text, -1)
	if len(names) == 0 {
		return "", fmt.Errorf("未设置表名称或者格式不正确:%s(行:%d)，格式：### 描述[表名,菜单名]，菜单名可选", line.Text, line.LineID)
	}
	s := strings.Split(strings.TrimRight(strings.TrimLeft(names[0], "["), "]"), ",")
	return s[0], nil
}

func getPKSName(path string) string {
	ext := filepath.Ext(path)
	dir := path
	if ext != "" {
		dir = filepath.Dir(path)
	}
	_, name := filepath.Split(dir)
	return name
}

//类型，长度，小数长度，错误
func getType(line *Line) (string, int, int, error) {
	colums := strings.Split(strings.Trim(line.Text, "|"), "|")
	if colums[0] == "" {
		return "", 0, 0, fmt.Errorf("字段名称不能为空 %s(行:%d)", line.Text, line.LineID)
	}

	t := strings.TrimSpace(colums[1])
	reg := regexp.MustCompile(`[\w]+`)
	names := reg.FindAllString(t, -1)
	if len(names) == 0 || len(names) > 4 {
		return "", 0, 0, fmt.Errorf("未设置字段类型:%v(行:%d)", names, line.LineID)
	}
	if len(names) == 1 {
		return t, 0, 0, nil
	}
	if len(names) == 2 {
		return t, types.GetInt(names[1]), 0, nil
	}
	return t, types.GetInt(names[1]), types.GetInt(names[2]), nil
}

func getAllMatchMD(path string) (paths []string) {

	//路径是的具体文件
	_, err := os.Stat(path)
	if err == nil {
		return []string{path}
	}
	//查找匹配的文件
	dir, f := filepath.Split(path)

	regexName := fmt.Sprintf("^%s$", strings.Replace(strings.Replace(f, ".md", "\\.md", -1), "*", "(.+)?", -1))
	reg := regexp.MustCompile(regexName)

	if dir == "" {
		dir = "./"
	}
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		fname := f.Name()
		if strings.HasPrefix(fname, ".") || f.IsDir() {
			continue
		}
		if reg.Match([]byte(fname)) {
			paths = append(paths, filepath.Join(dir, fname))
		}
	}
	return paths
}
