package tmpl

import (
	"fmt"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/micro-plat/hicli/markdown/const/enums"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/lib4go/types"
)

//MYSQL mysql数据库
const MYSQL = "mysql"

//ORACLE ORACLE数据库
const ORACLE = "oracle"

//webEnumComponents 前端涉及枚举的组件类型名称
var webEnumComponents = []string{"sl", "cb", "rd", "slm"}

type callHanlder func(string) string

func getfuncs(tp string) map[string]interface{} {
	return map[string]interface{}{
		//参数计算函数
		"add1":    add(1), //加1
		"mod":     getMod, //余数
		"mkSlice": mkSlice,
		"isTrue":  types.GetBool,

		//字符串处理函数
		"varName":   getVarName,            //获取pascal变量名称
		"snames":    getNames("/"),         //去掉首位斜线，并根据斜线分隔字符串
		"rmhd":      rmhd,                  //去除首段名称
		"rmhd2":     rmhd2,                 //去除首段名称
		"isNull":    isNull(),              //返回空语句
		"isMDNull":  isMDNull(),            //数据字典对应的为空判断
		"firstStr":  getStringByIndex(0),   //第一个字符
		"lastStr":   getLastStringByIndex,  //最后一个字符
		"l2d":       replaceUnderline("."), //下划线替换为.
		"hasPrefix": strings.HasPrefix,     //字符串前缀判定
		"lowerName": fGetLowerCase,         //小驼峰式命名
		"upperName": fGetUpperCase,         //大驼峰式命名

		//文件路径处理的函数
		"rpath":        getRouterPath,  //获取路由地址
		"fpath":        getFilePath,    //获取文件地址
		"webfpath":     getWebFilePath, //获取文件地址
		"parentPath":   getParentPath,  //获取文件夹地址
		"importPath":   getImportPath,  //go项目引用路径
		"pathPrefix":   pathPrefix,     //对项目前缀的预处理
		"fileBasePath": filepath.Base,  //文件基础路径

		//枚举处理函数
		"fIsEnumTB": hasKW("di", "dn"), //数据表的字段是否包含字典数据配置
		"fHasDT":    hasKW("dt"),       //数据表是否包含字典类型字段
		"fIsDI":     getKWS("di"),      //字段是否为字典ID
		"fIsDN":     getKWS("dn"),      //字段是否为字典Name
		"fIsDT":     getKWS("dt"),      //字段是否为字典Type
		"fIsDC":     getKWS("dc"),      //字段是否为字典column
		"fIsDPID":   getKWS("dpid"),    //字段是否为字典PID

		//数据库，sql，后端modules相关处理函数
		"shortName": shortName,                                          //获取特殊字段前的字符串
		"dbType":    dbType(tp),                                         //转换为SQL的数据类型
		"codeType":  codeType,                                           //转换为GO代码的数据类型
		"defValue":  defValue(tp),                                       //返回SQL中的默认值
		"seqTag":    getSEQTag(tp),                                      //获取SEQ的变量值
		"seqValue":  getSEQValue(tp),                                    //获取SEQ起始值
		"mysqlseq":  getSEQ(tp),                                         //获取SEQ
		"oracleseq": getSeqs,                                            //获取表的序列
		"pks":       getPKS,                                             //获取主键列表
		"indexs":    getDBIndex(tp),                                     //获取表的索引串
		"maxIndex":  getMaxIndex,                                        //最大索引值
		"lower":     strings.ToLower,                                    //获取变量的最小写字符
		"order":     getRows("order"),                                   //order排序
		"orderSort": sortByKw("order"),                                  //
		"orderCon":  getBracketContent([]string{"order"}, `(asc|desc)`), //
		"isTime":    isType("time.Time"),                                //是否是time
		"isDecimal": isType("types.Decimal"),                            //是否是decimal
		"isInt64":   isType("int64"),                                    //是否是int64
		"isInt":     isType("int"),                                      //是否是int
		"isString":  isType("string"),                                   //是否是string
		"replace":   replace(tp),
		"isPK":      judgeIndexType(enums.IndexPK),  //是否是主键
		"isUNQ":     judgeIndexType(enums.IndexUnq), //是否是唯一索引
		"isIndex":   judgeIndexType(enums.IndexNor), //是否是唯一索引
		"indexStr":  indexString(tp),

		//前后端约束处理函数
		"query":     getRows("q"),                                      //查询字段
		"list":      getRows("l"),                                      //列表展示字段
		"detail":    getRows("d"),                                      //详情展示字段
		"create":    getRows("c"),                                      //创建字段
		"delete":    getRows("del"),                                    //删除时判定字段
		"update":    getRows("u"),                                      //更新字段
		"export":    getRows("ept"),                                    //导出字段
		"tablist":   decodeRows("tabl", "l"),                           //详情tab页面列表字段
		"tabdetail": decodeRows("tabd", "d"),                           //详情tab页面字段
		"delCon":    getBracketContent([]string{"del"}),                //删除字段约束
		"sortCon":   getBracketContent([]string{"sort"}, `(asc|desc)`), //
		"sort":      getRows("sort"),                                   //查询字段
		"sortSort":  sortByKw("sort"),                                  //
		"qgroup":    getChildrenGroup("q"),                             //
		"cgroup":    getChildrenGroup("c"),                             //
		"ugroup":    getChildrenGroup("u"),                             //

		//前端约束处理函数
		"SL":            getKWS("sl"),                                                              //表单下拉框
		"SLM":           getKWS("slm"),                                                             //表单下拉框
		"CB":            getKWS("cb"),                                                              //表单复选框
		"RD":            getKWS("rd"),                                                              //表单单选框
		"TA":            getKWS("ta"),                                                              //表单文本域
		"DRANGE":        getKWS("drange"),                                                          //表单日期时间选择器
		"DTIME":         getKWS("dtime"),                                                           //表单日期时间选择器
		"DATE":          getKWS("date"),                                                            //表单日期选择器
		"UP":            getKWS("up"),                                                              //文件上传
		"CSCR":          getKWS("cscr"),                                                            //级联组件
		"dateType":      getDateType,                                                               //日期字段对应的组件的日期类型
		"dateFormat":    getDateFormat,                                                             //日期字段对应的组件的日期格式
		"dateFormatDef": getDateFormatDef,                                                          //日期字段对应的组件的日期默认值
		"CC":            getKWS("cc"),                                                              //表单颜色样式
		"FIXED":         getKWS("fixed"),                                                           //表单固定列
		"SORT":          getKWS("sort"),                                                            //表单固定列
		"fIsNofltr":     getKWS("nofltr"),                                                          //前端字段不格式化
		"LINK":          getKWS("link"),                                                            //表单点击跳转
		"linkCon":       getBracketContent([]string{"link"}),                                       //表单点击跳转约束
		"drangeCon":     getBracketContent([]string{"drange"}),                                     //表单点击跳转约束
		"cscrCon":       getBracketContent([]string{"cscr"}, "cscr\\(([\\w]+)"),                    //表单点击跳转约束
		"cscrDefault":   getBracketContent([]string{"cscr"}, "cscr\\([\\w]+\\,{0,1}([\\w,]+)?\\)"), //表单点击跳转约束
		"eptCon":        getBracketContent([]string{"ept"}),                                        //导出字段
		"lfCon":         getSubConContent("l", "f"),                                                //列表展示字段的过滤器子约束l(f:xx)
		"leCon":         getSubConContent("l", "e"),                                                //列表展示字段的枚举子约束l(e:xx)
		"qeCon":         getSubConContent("q", "e"),                                                //查询字段的枚举子约束q(e:xx)
		"qfCon":         getSubConContent("q", "f"),                                                //查询字段的枚举子约束q(f:xx)
		"rfCon":         getSubConContent("d", "f"),                                                //详情展示字段的过滤器子约束r(f:xx)
		"ufCon":         getSubConContent("u", "f"),                                                //详情展示字段的过滤器子约束r(f:xx)
		"cfCon":         getSubConContent("c", "f"),                                                //详情展示字段的过滤器子约束r(f:xx)
		"reCon":         getSubConContent("d", "e"),                                                //详情展示字段的枚举子约束r(e:xx)
		"ueCon":         getSubConContent("u", "e"),                                                //编辑字段的格式枚举子约束u(e:xx)
		"ceCon":         getSubConContent("c", "e"),                                                //添加字段的格式枚举子约束c(e:xx)
		"crCon":         getSubConContent("c", "r"),                                                //添加字段的格式枚举子约束c(r:xx)
		"udCon":         getSubConContent("u", "d"),                                                //添加字段的格式枚举子约束c(r:xx)
		"dicName":       getDicName(webEnumComponents...),                                          //编辑字段的格式枚举子约束u(d:true|false)
		"qDicCName":     getCascadeChildrenName("q", "e", webEnumComponents...),                    //查询下拉字段级联枚举对应的引用枚举名称
		"qDicPName":     getCascadeParentName("q", "e", webEnumComponents...),                      //查询下拉字段级联枚举对应的被引用枚举名称
		"cDicCName":     getCascadeChildrenName("c", "e", webEnumComponents...),                    //创建下拉字段级联枚举对应的引用枚举名称
		"cDicPName":     getCascadeParentName("c", "e", webEnumComponents...),                      //创建下拉字段级联枚举对应的被引用枚举名称
		"uDicCName":     getCascadeChildrenName("u", "e", webEnumComponents...),                    //更新下拉字段级联枚举对应的引用枚举名称
		"uDicPName":     getCascadeParentName("u", "e", webEnumComponents...),                      //更新下拉字段级联枚举对应的被引用枚举名称
		"qGroupCName":   getCascadeChildrenName("q", "g", webEnumComponents...),                    //创建下拉字段级联枚举对应的引用枚举名称
		"qGroupPName":   getCascadeParentName("q", "g", webEnumComponents...),                      //创建下拉字段级联枚举对应的被引用枚举名称
		"cGroupCName":   getCascadeChildrenName("c", "g", webEnumComponents...),                    //创建下拉字段级联枚举对应的引用枚举名称
		"cGroupPName":   getCascadeParentName("c", "g", webEnumComponents...),                      //创建下拉字段级联枚举对应的被引用枚举名称
		"uGroupCName":   getCascadeChildrenName("u", "g", webEnumComponents...),                    //更新下拉字段级联枚举对应的引用枚举名称
		"uGroupPName":   getCascadeParentName("u", "g", webEnumComponents...),                      //更新下拉字段级联枚举对应的被引用枚举名称

		"setIsInput":  setIsInput,
		"DMI":         getKWS("dmi"),  //dropdown menu+input 查询
		"dropmenurow": getRows("dmi"), //dropdown menu+input 查询
		"trimlist":    trimSuffix(".list"),

		"drangeValue": drangeValue, //表单日期时间选择器
		"ruleValue":   ruleValue,   //表单日期时间选择器
	}
}

func mkSlice(args ...interface{}) []interface{} {
	return args
}

func trimSuffix(suffix string) func(s string) string {
	return func(s string) string {
		return strings.TrimSuffix(s, suffix)
	}
}

func setIsInput(r *Row) string {
	r.IsInput = true
	return ""
}

func getMod(x int, y int) int {
	return x % y
}

//去掉首段名称
func rmhd(input string) string {
	if !trimPrefix {
		return strings.TrimPrefix(input, "_")
	}

	index := strings.Index(input, "_")
	return input[index+1:]
}

//去掉首段名称2
func rmhd2(input string) string {
	index := strings.Index(input, "_")
	return input[index+1:]
}

//获取短文字
func shortName(input string) string {
	reg := regexp.MustCompile(`^[\p{Han}|\w]+`)
	return reg.FindString(input)
}

//获取短文字
func isNull() func(*Row) string {
	return func(row *Row) string {
		return IsNull[row.IsNull]
	}
}

func isMDNull() func(*Row) string {
	return func(row *Row) string {
		return IsMDNull[row.IsNull]
	}
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

func getNames(kw string) func(input string) []string {
	return func(input string) []string {
		items := strings.Split(strings.Trim(input, kw), kw)
		return items
	}
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

func getSeqs(tb *Table) []map[string]interface{} {
	columns := make([]map[string]interface{}, 0, len(tb.Rows))

	for _, v := range tb.Rows {
		ok, seqName, start, increament := getCapturingGroup(v.Con, "seq")
		if !ok {
			continue
		}
		seqName = types.DecodeString(seqName, "", fmt.Sprintf("seq_%s_id", rmhd(tb.Name)))
		if len(seqName) > 64 {
			logs.Log.Errorf("自动生成或配置%s的序列名长度不正确(%s),请重新配置", v.Name, seqName)
			return nil
		}
		row := map[string]interface{}{
			"name":      v.Name,
			"seqname":   seqName,
			"desc":      v.Desc,
			"type":      v.Type,
			"len":       v.Len,
			"increment": types.DecodeInt(types.GetInt(increament, 0), 0, 1),
			"min":       types.DecodeInt(types.GetInt(start, 0), 0, 1),
			"max":       99999999999,
		}
		columns = append(columns, row)
	}

	return columns
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
	case ORACLE:
		return func(input string) string {
			reg := regexp.MustCompile(`[\w]+`)
			tps := reg.FindAllString(strings.ToLower(input), -1)
			if len(tps) == 0 {
				return input
			}
			t := tps[0]
			for k, v := range tp2oracle {
				reg2 := regexp.MustCompile(k)
				if !reg2.Match([]byte(strings.ToLower(t))) {
					continue
				}
				oracleType := v
				index := strings.Index(oracleType, "(")
				if len(tps) == 2 {
					if index >= 0 {
						return fmt.Sprintf("%s(%s)", oracleType[:index], tps[1])
					}
					return fmt.Sprintf("%s(%s)", oracleType, tps[1])
				}
				if len(tps) == 3 {
					if index >= 0 {
						return fmt.Sprintf("%s(%s,%s)", oracleType[:index], tps[1], tps[2])
					}
					return fmt.Sprintf("%s(%s,%s)", oracleType, tps[1], tps[2])
				}
				return t
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
	defMaps := []map[string]string{}
	switch tp {
	case MYSQL:
		defMaps = def2mysql
	case ORACLE:
		defMaps = def2oracle
	}
	return func(row *Row) string {
		if isCons(row.Con, "seq") {
			return ""
		}
		buff := []byte(strings.Trim(strings.ToLower(row.Def), "'"))
		for _, defs := range defMaps {
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
					return ""
				}
			}
			return ""
		}
	}
	return func(r *Table) string { return "" }
}

func getSEQ(tp string) func(r *Table) bool {
	switch tp {
	case MYSQL:
		return func(r *Table) bool {
			for _, r := range r.RawRows {
				if isCons(r.Con, "seq") {
					return true
				}
			}
			return false
		}
	}
	return func(r *Table) bool { return false }
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

func sortByKw(kw string) func(rows TableColumn) []*Row {
	return func(rows TableColumn) []*Row {
		result := make(TableColumn, 0, len(rows))
		for _, v := range rows {
			ok, _, sort, _ := getCapturingGroup(v.Con, kw)
			if !ok {
				continue
			}
			v.Sort = types.GetInt(sort, 0)
			result = append(result, v)
		}
		sort.Sort(result)
		return result
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
				case enums.IndexNor:
					list = append(list, index.String(tp))
				case enums.IndexUnq:
					list = append(list, index.String(tp))
				case enums.IndexPK:
					list = append(list, index.String(tp))
				}
			}
			if len(list) > 0 {
				return "," + strings.Join(list, "\n\t\t,")
			}
			return ""
		}
	case ORACLE:
		return func(r *Table) string {
			indexs := r.GetIndexs()
			list := make([]string, 0, len(indexs))
			for _, index := range indexs {
				switch index.Type {
				case enums.IndexNor:
					list = append(list, fmt.Sprintf("create %s;\n\t", index.String(tp)))
				case enums.IndexUnq:
					list = append(list, fmt.Sprintf("alter table %s add %s;\n\t", r.Name, index.String(tp)))
				case enums.IndexPK:
					list = append(list, fmt.Sprintf("alter table %s add %s;\n\t", r.Name, index.String(tp)))
				}
			}
			if len(list) > 0 {
				return strings.Join(list, "")
			}
			return ""
		}
	}
	return func(r *Table) string { return "" }
}

// 类似 kw(内容,内容,内容)
func getCapturingGroup(input string, kw string) (bool, string, string, string) {
	buff := []byte(strings.Trim(strings.ToLower(input), "'"))
	cks, ok := cons[strings.ToLower(kw)]
	if !ok {
		cks = cons["*"]
	}
	for _, v := range cks {
		nck := types.DecodeString(strings.Contains(v, "%s"), true, fmt.Sprintf(v, kw), v)
		reg := regexp.MustCompile(nck)
		if reg.Match(buff) {
			value := reg.FindStringSubmatch(strings.ToLower(input))
			if len(value) == 5 {
				return true, value[2], value[3], value[4]
			}
			if len(value) == 4 {
				return true, value[2], value[3], ""
			}
			if len(value) == 3 {
				return true, value[2], "", ""
			}
			return true, "", "", ""
		}
	}
	return false, "", "", ""
}

func decodeRows(tp, deftp string) func(row []*Row) []*Row {
	return func(row []*Row) []*Row {
		list := make([]*Row, 0, 1)
		for _, r := range row {
			if isCons(r.Con, tp) {
				list = append(list, r)
			}
		}
		if len(list) > 0 {
			return list
		}

		return getRows(deftp)(row)
	}
}

func HasRow(row []*Row, tp ...string) bool {
	return len(getRows(tp...)(row)) > 0
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

func replaceUnderline(new string) func(s string) string {
	return func(s string) string {
		if s == "" {
			return ""
		}
		return strings.Replace(strings.ToLower(s), "_", new, -1)
	}
}

func replace(tp string) func(row *Row) string {
	return func(row *Row) string {
		ok, start, end, s := getCapturingGroup(row.Con, "replace")
		if !ok {
			return ""
		}
		str := types.DecodeString(s, "", "****") //默认替换成’*‘
		switch tp {
		case MYSQL:
			format := "concat(left(t.%s,%d),'%s',right(t.%s,%d))"
			return fmt.Sprintf(format, row.Name, types.GetInt(start, 0), str, row.Name, types.GetInt(end, 0))
		case ORACLE:
			format := "substr(t.%s,0,%d)|| '%s' || substr(t.%s,-%d,%d)"
			return fmt.Sprintf(format, row.Name, types.GetInt(start, 0), str, row.Name, types.GetInt(end, 0), types.GetInt(end, 0))
		}
		return ""
	}
}

//getFilePath 获取文件地址
func getFilePath(tabName string) string {
	dir, _ := filepath.Split(replaceUnderline("/")(tabName))
	if strings.HasSuffix(dir, "/main/") {
		dir = strings.TrimSuffix(dir, "/main/") + "/mainx/"
	}

	return path.Join(dir, replaceUnderline(".")(tabName))
}

//getFilePath 获取文件地址
func getWebFilePath(tabName string) string {
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

func getLastbutoneStringByIndex(s []string) string {
	if len(s) < 2 {
		return "x"
	}
	return types.GetStringByIndex(s, len(s)-2)
}

func getDicName(keys ...string) func(con string, subcon string, tb *Table) string {
	return func(con string, subcon string, tb *Table) string {
		tp := subcon
		if tp == "" || strings.HasPrefix(subcon, "#") { //子约束为空或指定级联
			tp = getBracketContent(keys)(con)           //获取组件的约束
			if tp == "" || strings.HasPrefix(tp, "#") { //约束不是表名，不是指定枚举名称
				return ""
			}
		}
		for _, tb := range tb.AllTables { //查看是否匹配表名
			if tb.Name == tp {
				if hasKW("di", "dn")(tb) && hasKW("dt")(tb) {
					logs.Log.Warn("约束%s指定表无法判断具体类型，需指定具体的枚举类型", con)
					return ""
					// for _, v := range tb.Rows {
					// 	if getKWS("dt")(v.Con) {
					// 		return v.Name
					// 	}
					// }
				}
				return strings.ToLower(rmhd(tb.Name))
			}
		}

		return tp
	}
}

func pathPrefix(v *SnippetConf) string {
	if trimPrefix {
		return v.Name
	}

	if !v.TrimPrefix {
		return v.Name
	}
	index := strings.Index(v.Name, "_")
	return v.Name[index+1:]
}

func getImportPath(s []*SnippetConf) map[string]*SnippetConf {
	r := make(map[string]*SnippetConf)
	t := make(map[string]string)

	for _, v := range s {
		path, _ := Translate("{{.|pathPrefix|rmhd|parentPath}}", "", v)
		if strings.HasSuffix(path, "/main") {
			path = strings.TrimSuffix(path, "/main") + "/mainx"
		}
		tpath := filepath.Join(fmt.Sprintf("%s/services", v.BasePath), path)
		alias := getLastStringByIndex(getNames("/")(path))
		if path == "" {
			alias = "services"
			v.PKG = alias
			if _, ok := r[tpath]; !ok {
				r[tpath] = v
			}
			continue
		}
		if p, ok := t[alias]; ok && p != tpath {
			pre := getLastbutoneStringByIndex(getNames("/")(path))
			alias = fmt.Sprintf("%s%s", pre, alias)
			v.PkGAlias = fmt.Sprintf("%s ", alias)
		}
		v.PKG = alias
		t[alias] = tpath
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
		if getKWS("drange")(con) {
			return "yyyy-MM-dd"
		}
		if getKWS("date")(con) {
			return "yyyy-MM-dd"
		}
		return "yyyy-MM-dd HH:mm:ss"
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

func drangeValue(con string) []string {
	if con == "" {
		return []string{"30"}
	}
	r := strings.Split(con, ",")
	if len(r) == 1 {
		return []string{r[0], r[0]}
	}
	if len(r) == 2 {
		return r
	}
	return []string{"30"}
}

func ruleValue(con string) []string {
	if con == "" {
		return []string{}
	}
	r := strings.Split(con, ",")
	if len(r) == 1 {
		return []string{"0", r[0]}
	}
	if len(r) == 2 {
		return r
	}
	return []string{}
}

func getCascadeChildrenName(tp, tkey string, keys ...string) func(name string, t *Table) string {
	return func(name string, t *Table) string {
		kw := fmt.Sprintf("#%s", name)
		for _, v := range t.Rows {
			subCon := getSubConContent(tp, tkey)(v.Con) //该字段枚举子约束
			if kw == subCon {
				return v.Name
			}
			if subCon != "" { //字段标识配置配置了对应枚举,不再处理组件标识的级联枚举
				continue
			}
			con := getBracketContent(keys)(v.Con)
			if strings.Contains(con, kw) {
				return v.Name
			}
		}
		return ""
	}
}

func getCascadeParentName(tp, tkey string, keys ...string) func(con string, t *Table) string {
	return func(con string, t *Table) string {
		subCon := getSubConContent(tp, tkey)(con)            //该字段枚举子约束
		if subCon != "" && !strings.HasPrefix(subCon, "#") { //字段标识配置配置了对应枚举,不再处理组件标识的级联枚举
			return ""
		}

		parentName := ""
		if strings.HasPrefix(subCon, "#") { ///该字段设置有级联枚举子约束
			parentName = strings.TrimPrefix(subCon, "#")
		}

		if parentName == "" {
			//查找组件约束的级联
			c := getBracketContent(keys)(con)
			if !strings.Contains(c, "#") { //该字段组件约束没有级联
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

func getChildrenGroup(tp string) func(name string, t *Table) []*Row {
	return func(name string, t *Table) []*Row {
		r := []*Row{}
		var wg sync.WaitGroup
		var lock sync.Mutex
		for _, v := range t.Rows {
			wg.Add(1)
			go func(tp string, v *Row, t *Table) {
				defer wg.Done()
				pname := getCascadeParentName(tp, "g", webEnumComponents...)(v.Con, t)
				if pname == name {
					lock.Lock()
					defer lock.Unlock()
					r = append(r, v)
				}
			}(tp, v, t)
		}
		wg.Wait()

		return r
	}
}

func getSubConContent(tp, kw string) func(con string) string {
	return func(con string) string {
		c := getBracketContent([]string{tp})(con)
		if c == "" {
			return ""
		}
		subConMap := map[string]string{}
		for _, v := range strings.Split(c, ",") {
			sub := strings.Index(v, ":")
			if sub < 1 {
				logs.Log.Warn("约束格式不正确：", con, tp, kw)
				continue
			}
			subConMap[v[0:sub]] = v[sub+1:]
		}
		if v, ok := subConMap[kw]; ok {
			return v
		}
		return ""
	}
}

// func getTest(con string) string {
// 	return getBracketContent([]string{"cscr"}, "cscr\\(([\\w]+)")(con)
// }

//pattern 取第一个分组
func getBracketContent(keys []string, pattern ...string) func(con string) string {
	return func(con string) string {
		s := make([]string, 0)
		for _, key := range keys {
			kw := ""
			for k := range key {
				kw += fmt.Sprintf("[%s%s]", strings.ToLower(key[k:k+1]), strings.ToUpper(key[k:k+1]))
			}
			rex := regexp.MustCompile(fmt.Sprintf(keywordSubMatch, kw))
			if len(pattern) > 0 {
				rex = regexp.MustCompile(pattern[0])
			}
			value := rex.FindStringSubmatch(con)
			if len(value) == 2 {
				s = append(s, value[1])
			}
		}
		return strings.Join(s, ",")
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

func indexString(tp string) func(index *Index) string {
	return func(index *Index) string {
		return index.String(tp)
	}
}

func judgeIndexType(t enums.IndexType) func(index *Index) bool {
	return func(index *Index) bool {
		return index.JudgeType(t)
	}
}
