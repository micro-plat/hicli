package tmpl

import (
	"fmt"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/lib4go/types"
)

//MYSQL mysql数据库
const MYSQL = "mysql"

//webEnumComponents 前端涉及枚举的组件类型名称
var webEnumComponents = []string{"sl", "cb", "rd", "slm"}

type callHanlder func(string) string

func getfuncs(tp string) map[string]interface{} {
	return map[string]interface{}{
		//参数计算函数
		"add1": add(1), //加1
		"mod":  getMod, //余数

		//字符串处理函数
		"varName":   getVarName,            //获取pascal变量名称
		"names":     getNames,              //去掉首位下划线，并根据下划线分隔字符串
		"rmhd":      rmhd,                  //去除首段名称
		"isNull":    isNull(tp),            //返回空语句
		"firstStr":  getStringByIndex(0),   //第一个字符
		"lastStr":   getLastStringByIndex,  //最后一个字符
		"l2d":       replaceUnderline("."), //下划线替换为.
		"hasPrefix": strings.HasPrefix,     //字符串前缀判定
		"lowerName": fGetLowerCase,         //小驼峰式命名
		"upperName": fGetUpperCase,         //大驼峰式命名

		//文件路径处理的函数
		"rpath":        getRouterPath, //获取路由地址
		"fpath":        getFilePath,   //获取文件地址
		"parentPath":   getParentPath, //获取文件夹地址
		"importPath":   getImportPath, //go项目引用路径
		"fileBasePath": filepath.Base, //文件基础路径

		//枚举处理函数
		"fIsEnumTB": hasKW("di", "dn"), //数据表的字段是否包含字典数据配置
		"fHasDT":    hasKW("dt"),       //数据表是否包含字典类型字段
		"fIsDI":     getKWS("di"),      //字段是否为字典ID
		"fIsDN":     getKWS("dn"),      //字段是否为字典Name
		"fIsDT":     getKWS("dt"),      //字段是否为字典Type

		//数据库，sql，后端modules相关处理函数
		"shortName": shortName,               //获取特殊字段前的字符串
		"dbType":    dbType(tp),              //转换为SQL的数据类型
		"codeType":  codeType,                //转换为GO代码的数据类型
		"defValue":  defValue(tp),            //返回SQL中的默认值
		"seqTag":    getSEQTag(tp),           //获取SEQ的变量值
		"seqValue":  getSEQValue(tp),         //获取SEQ起始值
		"pks":       getPKS,                  //获取主键列表
		"indexs":    getDBIndex(tp),          //获取表的索引串
		"maxIndex":  getMaxIndex,             //最大索引值
		"lower":     getLower,                //获取变量的最小写字符
		"order":     getOrderBy,              //order排序
		"ismysql":   stringsEqual("mysql"),   //是否是mysql
		"isoracle":  stringsEqual("oracle"),  //是否是oracle
		"isTime":    isType("time.Time"),     //是否是time
		"isDecimal": isType("types.Decimal"), //是否是decimal
		"isInt64":   isType("int64"),         //是否是int64
		"isInt":     isType("int"),           //是否是int
		"isString":  isType("string"),        //是否是string

		//前后端约束处理函数
		"query":  getRows("q"),           //查询字段
		"list":   getRows("l"),           //列表展示字段
		"detail": getRows("r"),           //详情展示字段
		"create": getRows("c"),           //创建字段
		"delete": getRows("d"),           //删除时判定字段
		"update": getRows("u"),           //更新字段
		"delCon": getBracketContent("d"), //删除字段约束

		//前端约束处理函数
		"SL":            getKWS("sl"),                                  //表单下拉框
		"SLM":           getKWS("slm"),                                 //表单下拉框
		"CB":            getKWS("cb"),                                  //表单复选框
		"RD":            getKWS("rd"),                                  //表单单选框
		"TA":            getKWS("ta"),                                  //表单文本域
		"DTIME":         getKWS("dtime"),                               //表单日期时间选择器
		"DATE":          getKWS("date"),                                //表单日期选择器
		"dateType":      getDateType,                                   //日期字段对应的组件的日期类型
		"dateFormat":    getDateFormat,                                 //日期字段对应的组件的日期格式
		"dateFormatDef": getDateFormatDef,                              //日期字段对应的组件的日期默认值
		"CC":            getKWS("cc"),                                  //表单颜色样式
		"FIXED":         getKWS("fixed"),                               //表单固定列
		"SORT":          getKWS("sort"),                                //表单固定列
		"lfCon":         getSubConContent("l", "f"),                    //列表展示字段的过滤器子约束l(f:xx)
		"leCon":         getSubConContent("l", "e"),                    //列表展示字段的枚举子约束l(e:xx)
		"qeCon":         getSubConContent("q", "e"),                    //查询字段的枚举子约束q(e:xx)
		"qfCon":         getSubConContent("q", "f"),                    //查询字段的枚举子约束q(f:xx)
		"rfCon":         getSubConContent("r", "f"),                    //详情展示字段的过滤器子约束r(f:xx)
		"reCon":         getSubConContent("r", "e"),                    //详情展示字段的枚举子约束r(e:xx)
		"ueCon":         getSubConContent("u", "e"),                    //编辑字段的格式枚举子约束u(e:xx)
		"ceCon":         getSubConContent("c", "e"),                    //添加字段的格式枚举子约束c(e:xx)
		"dicName":       getDicName(webEnumComponents...),              //字段的对应的枚举名称
		"qDicCName":     getDicChildrenName("q", webEnumComponents...), //查询下拉字段级联枚举对应的引用枚举名称
		"qDicPName":     getDicParentName("q", webEnumComponents...),   //查询下拉字段级联枚举对应的被引用枚举名称
		"cDicCName":     getDicChildrenName("c", webEnumComponents...), //创建下拉字段级联枚举对应的引用枚举名称
		"cDicPName":     getDicParentName("c", webEnumComponents...),   //创建下拉字段级联枚举对应的被引用枚举名称
		"uDicCName":     getDicChildrenName("u", webEnumComponents...), //更新下拉字段级联枚举对应的引用枚举名称
		"uDicPName":     getDicParentName("u", webEnumComponents...),   //更新下拉字段级联枚举对应的被引用枚举名称

	}
}

func getLower(s string) string {
	return strings.ToLower(s)
}

func getMod(x int, y int) int {
	return x % y
}

//去掉首段名称
func rmhd(input string) string {
	index := strings.Index(input, "_")
	return input[index+1:]
}

//获取短文字
func shortName(input string) string {
	reg := regexp.MustCompile(`^[\p{Han}|\w]+`)
	return reg.FindString(input)
}

//获取短文字
func isNull(tp string) func(*Row) string {
	switch tp {
	case MYSQL:
		return func(row *Row) string {
			return mysqlIsNull[row.IsNull]
		}
	}
	return func(row *Row) string { return "" }
}

//首字母大写，并去掉下划线
func getVarName(input string) string {
	items := strings.Split(input, "_")
	nitems := make([]string, 0, len(items))
	for _, item := range items {
		if item == "" {
			continue
		}
		if len(item) == 1 {
			nitems = append(nitems, strings.ToUpper(item[0:1]))
			continue
		}
		if strings.EqualFold(item, "id") || strings.EqualFold(item, "url") {
			nitems = append(nitems, strings.ToUpper(item))
			continue
		}
		nitems = append(nitems, strings.ToUpper(item[0:1])+item[1:])
	}
	return strings.Join(nitems, "")
}

func getNames(input string) []string {
	items := strings.Split(strings.Trim(input, "_"), "_")
	return items
}

func fGetLowerCase(n string) string {
	items := strings.Split(n, "_")
	nitems := make([]string, 0, len(items))
	for k, i := range items {
		if strings.EqualFold(i, "id") || strings.EqualFold(i, "url") {
			nitems = append(nitems, strings.ToUpper(i))
			continue
		}
		if k == 0 {
			nitems = append(nitems, i)
		}
		if k > 0 {
			nitems = append(nitems, strings.ToUpper(i[0:1])+i[1:])
		}

	}
	return strings.Join(nitems, "")
}

func fGetUpperCase(n string) string {
	_, f := filepath.Split(n)
	f = strings.ReplaceAll(f, ".", "_")
	items := strings.Split(f, "_")
	nitems := make([]string, 0, len(items))
	for _, i := range items {
		if strings.EqualFold(i, "id") || strings.EqualFold(i, "url") {
			nitems = append(nitems, strings.ToUpper(i))
			continue
		}
		nitems = append(nitems, strings.ToUpper(i[0:1])+i[1:])
	}
	return strings.Join(nitems, "")
}

//通过正则表达式，转换正确的数据库类型
func dbType(tp string) callHanlder {
	switch tp {
	case MYSQL:
		return func(input string) string {
			buff := []byte(strings.ToLower(input))
			for k, v := range tp2mysql {
				reg := regexp.MustCompile(k)
				if reg.Match(buff) {
					if !strings.Contains(v, "*") {
						return v
					}
					value := reg.FindStringSubmatch(input)
					if len(value) > 1 {
						return strings.Replace(v, "*", strings.Join(value[1:], ","), -1)
					}
					return v
				}
			}
			return input
		}
	}
	return func(input string) string { return "" }
}

//通过正则表达式，转换正确的数据库类型
func codeType(input string) string {
	buff := []byte(strings.ToLower(input))
	for k, v := range any2code {
		reg := regexp.MustCompile(k)
		if reg.Match(buff) {
			return v
		}
	}
	return input
}

//通过正则表达式，转换正确的数据库类型
func defValue(tp string) func(*Row) string {
	switch tp {

	case MYSQL:
		return func(row *Row) string {
			if isCons(row.Con, "seq") {
				return ""
			}
			buff := []byte(strings.Trim(strings.ToLower(row.Def), "'"))
			for _, defs := range def2mysql {
				for k, v := range defs {
					reg := regexp.MustCompile(k)
					if reg.Match(buff) {
						if !strings.Contains(v, "*") {
							return v
						}
						value := reg.FindStringSubmatch(row.Def)
						if len(value) > 1 {
							return strings.Replace(v, "*", strings.Join(value[1:], ","), -1)
						}
						return row.Def
					}
				}
			}
			return row.Def
		}
	}
	return func(row *Row) string { return "" }
}

func getPKS(t *Table) []string {
	return t.GetPKS()
}

func getSEQTag(tp string) func(r *Row) string {
	switch tp {
	case MYSQL:
		return func(r *Row) string {
			if isCons(r.Con, "seq") {
				return "auto_increment"
			}
			return ""
		}
	}
	return func(r *Row) string { return "" }
}

func getSEQValue(tp string) func(r *Table) string {
	switch tp {
	case MYSQL:
		return func(r *Table) string {
			for _, r := range r.RawRows {
				if isCons(r.Con, "seq") {
					if v := types.GetInt(r.Def, 0); v != 0 {
						return fmt.Sprintf("auto_increment = %d", v)
					}

				}
			}
			return ""
		}
	}
	return func(r *Table) string { return "" }
}

func getMaxIndex(r interface{}) int {
	v := reflect.ValueOf(r)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array || v.Kind() == reflect.Map {
		return v.Len() - 1
	}
	return 0
}

func getOrderBy(tb *Table) []map[string]interface{} {
	columns := make([]map[string]interface{}, 0, len(tb.Rows))
	fileds := []string{}
	orders := []string{}
	ob := map[string]string{}

	for _, v := range tb.Rows {
		if strings.Contains(v.Con, "OB") {
			if !strings.Contains(v.Con, "OB(") {
				fileds = append(fileds, v.Name)
				continue
			}
			for _, v1 := range strings.Split(v.Con, ",") {
				if !strings.Contains(v1, "OB(") {
					continue
				}
				s := strings.Index(v1, "OB(")
				e := strings.Index(v1, ")")
				orders = append(orders, v1[s+1:e])
				ob[v1[s+1:e]] = v.Name
			}
		}
	}

	if len(orders) > 0 {
		sort.Sort(sort.StringSlice(orders))
	}

	for _, v := range orders {
		fileds = append(fileds, ob[v])
	}

	for _, v := range fileds {
		row := map[string]interface{}{
			"name":  v,
			"comma": true,
		}
		columns = append(columns, row)
	}

	if len(columns) > 0 {
		columns[len(columns)-1]["comma"] = false
	}
	return columns
}

func getSeqs() func(tb *Table) []map[string]interface{} {
	return func(tb *Table) []map[string]interface{} {
		columns := make([]map[string]interface{}, 0, len(tb.Rows))

		for _, v := range tb.Rows {
			if strings.Contains(v.Con, "SEQ") {
				descsimple := getBracketContent("SEQ")(v.Desc)
				row := map[string]interface{}{
					"name":       v.Name,
					"descsimple": descsimple,
					"seqname":    "seqname", //  fmt.Sprintf("seq_%s_%s", fGetNName(tb.Name), getFilterName(tb.Name, v.Cname)),
					"desc":       v.Desc,
					"type":       v.Type,
					"len":        v.Len,
					"comma":      true,
				}
				columns = append(columns, row)
			}
		}
		return columns
	}
}

//去掉首段名称
func isCons(input string, tp string) bool {
	cks, ok := cons[strings.ToLower(tp)]
	if !ok {
		cks = cons["*"]
	}
	buff := []byte(strings.ToLower(input))
	for _, ck := range cks {
		nck := types.DecodeString(strings.Contains(ck, "%s"), true, fmt.Sprintf(ck, tp), ck)
		reg := regexp.MustCompile(nck)
		if reg.Match(buff) {
			return true
		}
	}
	return false
}

func getDBIndex(tp string) func(r *Table) string {
	switch tp {
	case MYSQL:
		return func(r *Table) string {
			indexs := r.GetIndexs()
			list := make([]string, 0, len(indexs))
			for _, index := range indexs {
				switch index.Type {
				case "idx":
					list = append(list, fmt.Sprintf("index %s(%s)", index.Name, index.fields.Join(",")))
				case "unq":
					list = append(list, fmt.Sprintf("unique index %s(%s)", index.Name, index.fields.Join(",")))
				case "pk":
					list = append(list, fmt.Sprintf("primary key (%s)", index.fields.Join(",")))
				}
			}
			if len(list) > 0 {
				return "," + strings.Join(list, "\n\t\t,")
			}
			return ""
		}
	}
	return func(r *Table) string { return "" }
}

func getIndex(input string, tp string) (bool, string, int) {
	buff := []byte(strings.Trim(strings.ToLower(input), "'"))
	for _, v := range cons[tp] {
		reg := regexp.MustCompile(v)
		if reg.Match(buff) {
			value := reg.FindStringSubmatch(strings.ToLower(input))
			if len(value) > 2 {
				return true, value[1], types.GetInt(value[2], 0)
			}
			if len(value) > 1 {
				return true, value[1], 0
			}
			return true, "", 0
		}
	}
	return false, "", 0
}

func getRows(tp ...string) func(row []*Row) []*Row {
	return func(row []*Row) []*Row {
		list := make([]*Row, 0, 1)
		for _, r := range row {
		NEXT:
			for _, t := range tp {
				if isCons(r.Con, t) {
					list = append(list, r)
					break NEXT
				}
			}
		}
		return list
	}
}

func getKWS(tp ...string) func(input string) bool {
	return func(input string) bool {
		for _, t := range tp {
			if isCons(input, t) {
				return true
			}
		}
		return false
	}
}

//getKWCons 获取关键字约束列
func getKWCons(input string, keyword string) bool {
	for _, kw := range keywordMatch {
		reg := regexp.MustCompile(fmt.Sprintf(kw, keyword))
		if reg.Match([]byte(input)) {
			return true
		}
	}
	return false
}

func isType(t string) func(input string) bool {
	return func(input string) bool {
		tp := codeType(input)
		return tp == t
	}
}

func stringsEqual(s string) func(s1 string) bool {
	return func(s1 string) bool {
		return strings.EqualFold(s, s1)
	}
}

func replaceUnderline(new string) func(s string) string {
	return func(s string) string {
		if s == "" {
			return ""
		}
		return strings.Replace(strings.ToLower(s), "_", new, -1)
	}
}

//getFilePath 获取文件地址
func getFilePath(tabName string) string {
	dir, _ := filepath.Split(replaceUnderline("/")(tabName))
	return path.Join(dir, replaceUnderline(".")(tabName))
}

//getParentPath 获取文件地址
func getParentPath(tabName string) string {
	dir, _ := filepath.Split(replaceUnderline("/")(tabName))
	return strings.TrimRight(dir, "/")
}

//getRouterPath 获取路由地址
func getRouterPath(tabName string) string {
	dir, f := filepath.Split(replaceUnderline("/")(tabName))
	return dir + f
}

func getStringByIndex(index int) func(s []string) string {
	return func(s []string) string {
		return types.GetStringByIndex(s, index)
	}
}

func getLastStringByIndex(s []string) string {
	if len(s) == 0 {
		return ""
	}
	return types.GetStringByIndex(s, len(s)-1)
}

func getDicName(keys ...string) func(con string, subcon string, tb *Table) string {
	return func(con string, subcon string, tb *Table) string {
		tp := subcon
		if tp == "" || strings.HasPrefix(subcon, "#") { //子约束为空或指定级联
			tp = getBracketContent(keys...)(con)        //获取组件的约束
			if tp == "" || strings.HasPrefix(tp, "#") { //约束不是表名，不是指定枚举名称
				return ""
			}
		}
		for _, tb := range tb.AllTables { //查看是否匹配表名
			if tb.Name == tp {
				if hasKW("di", "dn")(tb) && hasKW("dt")(tb) {
					for _, v := range tb.Rows {
						if getKWS("dt")(v.Con) {
							return v.Name
						}
					}
				}
				return rmhd(tb.Name)
			}
		}

		return tp
	}
}

func getImportPath(s []*SnippetConf) map[string]*SnippetConf {
	r := make(map[string]*SnippetConf, 0)
	for _, v := range s {
		path, _ := Translate("{{.Name|rmhd|parentPath}}", "", v)
		tpath := fmt.Sprintf("%s/services/%s", v.BasePath, path)
		if _, ok := r[tpath]; !ok {
			r[tpath] = v
		}
	}
	return r
}

func getDateFormat(con, subCon string) string {
	if subCon == "" {
		if getKWS("dtime")(con) {
			return "yyyy-MM-dd HH:mm:ss"
		}
		if getKWS("date")(con) {
			return "yyyy-MM-dd"
		}
		return "yyyy-MM-dd"
	}

	return subCon
}

func getDateFormatDef(con, subCon string) string {
	f := getDateFormat(con, subCon)
	f = strings.ReplaceAll(f, "H", "0")
	f = strings.ReplaceAll(f, "h", "0")
	f = strings.ReplaceAll(f, "m", "0")
	f = strings.ReplaceAll(f, "s", "0")
	return f
}

func getDateType(con, subCon string) string {
	if subCon == "" {
		if getKWS("dtime")(con) {
			return "datetime"
		}
		if getKWS("date")(con) {
			return "date"
		}
	}

	if strings.Contains(subCon, "h") || strings.Contains(subCon, "H") || strings.Contains(subCon, "m") || strings.Contains(subCon, "s") {
		return "datetime"
	}
	return "date"
}

func getDicChildrenName(tp string, keys ...string) func(name string, t *Table) string {
	return func(name string, t *Table) string {
		kw := fmt.Sprintf("#%s", name)
		for _, v := range t.Rows {
			subCon := getSubConContent(tp, "e")(v.Con) //该字段枚举子约束
			if kw == subCon {
				return v.Name
			}
			if subCon != "" { //字段标识配置配置了对应枚举,不再处理组件标识的级联枚举
				return ""
			}
			con := getBracketContent(keys...)(v.Con)
			if strings.Contains(con, kw) {
				return v.Name
			}
		}
		return ""
	}
}

func getDicParentName(tp string, keys ...string) func(con string, t *Table) string {
	return func(con string, t *Table) string {
		subCon := getSubConContent(tp, "e")(con)             //该字段枚举子约束
		if subCon != "" && !strings.HasPrefix(subCon, "#") { //字段标识配置配置了对应枚举,不再处理组件标识的级联枚举
			return ""
		}

		parentName := ""
		if strings.HasPrefix(subCon, "#") { ///该字段设置有级联枚举子约束
			parentName = strings.TrimPrefix(subCon, "#")
		}

		if parentName == "" {
			//查找组件约束的级联
			c := getBracketContent(keys...)(con)
			if strings.Index(c, "#") < 0 { //该字段组件约束没有级联
				return ""
			}
			for _, v := range strings.Split(c, ",") {
				if strings.HasPrefix(v, "#") {
					parentName = strings.TrimPrefix(v, "#")
					break
				}
			}
		}

		if parentName == "" {
			return ""
		}

		for _, v := range t.Rows {
			if v.Name == parentName {
				return parentName
			}
		}

		return ""
	}
}

func getSubConContent(tp, kw string) func(con string) string {
	return func(con string) string {
		c := getBracketContent(tp)(con)
		if c == "" {
			return ""
		}
		subConMap := map[string]string{}
		for _, v := range strings.Split(c, ",") {
			sub := strings.Index(v, ":")
			if sub < 1 {
				logs.Log.Warn("约束格式不正确：", con)
				continue
			}
			subConMap[v[0:sub]] = v[sub+1 : len(v)]
		}
		//		fmt.Println("con:", con, "map:", subConMap)
		if v, ok := subConMap[kw]; ok {
			return v
		}
		return ""
	}
}

func getBracketContent(keys ...string) func(con string) string {
	return func(con string) string {
		s := ""
		for _, key := range keys {
			kw := ""
			for k := range key {
				kw += fmt.Sprintf("[%s%s]", strings.ToLower(key[k:k+1]), strings.ToUpper(key[k:k+1]))
			}
			rex := regexp.MustCompile(fmt.Sprintf(`%s\((.+?)\)`, kw))
			strs := rex.FindAllString(con, -1)
			if len(strs) < 1 {
				continue
			}
			str := strs[0]
			str = str[strings.Index(str, "(")+1 : len(str)]
			str = strings.TrimRight(str, ")")
			s = fmt.Sprintf("%s,%s", s, str)
		}
		if s == "" {
			return ""
		}
		s = strings.TrimLeft(s, ",")
		return s
	}
}

func hasKW(tp ...string) func(t *Table) bool {
	return func(t *Table) bool {
		ext := map[string]bool{}
		for _, r := range t.Rows {
			for _, t := range tp {
				if isCons(r.Con, t) {
					ext[t] = true
				}
			}
		}
		for _, t := range tp {
			if _, ok := ext[t]; !ok {
				return false
			}
		}
		return true
	}
}

func add(a int) func(b int) int {
	return func(b int) int { return a + b }
}
