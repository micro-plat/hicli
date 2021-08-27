package tmpl

var tp2mysql = map[string]string{
	"^date$":                      "datetime",
	"^datetime$":                  "datetime",
	"^timestamp$":                 "datetime",
	"^decimal$":                   "decimal",
	"^float$":                     "float",
	"^int$":                       "int",
	"^number\\([1-2]{1}\\)$":      "tinyint",
	"^number\\([3-9]{1}\\)$":      "int",
	"^number\\(10\\)$":            "int",
	"^number\\(1[1-9]{1}\\)$":     "bigint",
	"^number\\(2[0-9]{1}\\)$":     "bigint",
	"^number\\((\\d+),(\\d+)\\)$": "decimal(*)",
	"^varchar\\((\\d+)\\)$":       "varchar(*)",
	"^varchar2\\((\\d+)\\)$":      "varchar(*)",
	"^string$":                    "tinytext",
	"^text$":                      "text",
	"^longtext$":                  "longtext",
	"^clob$":                      "longtext",
}

//参考资料 http://www.sqlines.com/mysql-to-oracle  https://docs.oracle.com/cd/E12151_01/doc.150/e12155/oracle_mysql_compared.htm#BABGACIF
var tp2oracle = map[string]string{
	"^bigint$":            "number(19)",
	"^binary$":            "raw",
	"^bit$":               "raw", //长度正确为（n/8）
	"^blob$":              "blob",
	"^boolean$":           "char(1)",
	"^bool$":              "char(1)",
	"^char$":              "char",
	"^character$":         "character",
	"^character varying$": "character varying",
	"^date$":              "date",
	"^datetime$":          "timestamp",
	"^decimal$":           "number",
	"^dec$":               "number",
	"^double$":            "binary_double",
	"^double precision$":  "binary_double",
	"^fixed$":             "number",
	"^float$":             "binary_double",
	"^float4$":            "binary_double",
	"^float8$":            "binary_double",
	"^int$":               "number",
	"^integer$":           "number",
	"^int1$":              "number(3)",
	"^int2$":              "number(5)",
	"^int3$":              "number(7)",
	"^int4$":              "number(10)",
	"^int8$":              "number(19)",
	"^longblob$":          "blob",
	"^longtext$":          "clob",
	"^long varbinary$":    "blob",
	"^long$":              "clob",
	"^long varchar$":      "clob",
	"^mediumblob$":        "blob",
	"^mediumint$":         "number(7)",
	"^mediumtext$":        "clob",
	"^middleint$":         "number(7)",
	"^nchar$":             "nchar",
	"^nvarchar$":          "nvarchar2",
	"^numeric$":           "number",
	"^real$":              "binary_double",
	"^smallint$":          "number(5)",
	"^text$":              "clob",
	"^time$":              "timestamp",
	"^timestamp$":         "timestamp",
	"^tinyblob$":          "raw(255)",
	"^tinyint$":           "number(3)",
	"^tinytext$":          "varchar2(255)",
	"^varbinary$":         "raw",
	"^varchar$":           "varchar2",
	"^year$":              "number(4)",
}

var def2mysql = []map[string]string{
	{
		"^$":                  "",
		"^-$":                 "default '-'",
		"^seq$":               "",
		"^current_timestamp$": "default current_timestamp",
		"^sysdate$":           "default current_timestamp",
		"^([0-9]+)$":          "default *",
	},
	{
		"^(.+)$": "default '*'",
	},
}

var def2oracle = []map[string]string{
	{
		"^$":                  "",
		"^-$":                 "default '-'",
		"^seq$":               "",
		"^current_timestamp$": "default sysdate",
		"^sysdate$":           "default sysdate",
		"^([0-9]+)$":          "default *",
	},
	{
		"^(.+)$": "default '*'",
	},
}

var any2code = map[string]string{
	"^number\\(([1-9]|10)\\)$":                      "int",
	"^number\\((1[1-9]|2[0-9])\\)$":                 "int64",
	"^number\\(\\d+,\\d+\\)$":                       "types.Decimal",
	"^(varchar2|varchar|nchar|nvarchar)\\(\\d+\\)$": "string",
	"^string$":                  "string",
	"^(bigint|int8)\\(\\d+\\)$": "int64",
	"^(binary|bit|blob|boolean|bool|char|character( varying)?)$":                          "string",
	"^(date|datetime|timestamp|time)$":                                                    "time.Time",
	"^(decimal|double( precision)?|float(4|8)?|real)$":                                    "types.Decimal",
	"^(dec|fixed|numeric|year|int(eger|[1-4])?|(medium|middle|small|tiny)int)\\(\\d+\\)$": "int",
	"^long(blob|text| varbinary| varchar)?$":                                              "string",
	"^(clob|text|medium(blob|text)|text|tinyblob|tinytext|varbinary)$":                    "string",
}

var keywordSubMatch = `\b%s\(([\w\s-:#,.|=/\p{Han}]+)\)`

var keywordMatch = []string{"^\\w*%s\\w*$", ",\\w*%s\\w*,", "^\\w*%s\\w*,", ",\\w*%s\\w*$"}

var cons = map[string][]string{
	"*":       {"\\b%s\\b"},
	"sl":      {"\\bsl(\\([\\w,]+\\)|\\b)"},
	"slm":     {"\\bslm(\\([\\w,]+\\)|\\b)"},
	"rd":      {"\\brd(\\([\\w,]+\\)|\\b)"},
	"cb":      {"\\bcb(\\([\\w,]+\\)|\\b)"},
	"ta":      {"\\bta(\\([\\w,]+\\)|\\b)"},
	"cc":      {"\\bcc(\\(\\w+\\)|\\b)"},
	"idx":     {"\\bidx(\\(([\\w]+)[,]?([\\d]?)\\)|\\b)"},
	"unq":     {"\\bunq(\\(([\\w]+)[,]?([\\d]?)\\)|\\b)"},
	"del":     {"\\bdel(\\([0-9]*\\)|\\b)"},
	"c":       {"\\bc(\\([\\w,:#]+\\)|\\b)"},
	"u":       {"\\bu(\\([\\w,:#]+\\)|\\b)"},
	"d":       {"\\bd(\\([\\w,:#]+\\)|\\b)"},
	"l":       {"\\bl(\\([\\w,:#]+\\)|\\b)"},
	"q":       {"\\bq(\\([\\w,:#]+\\)|\\b)"},
	"ept":     {"\\bept(\\([\\w,:#]+\\)|\\b)"},
	"sort":    {"\\bsort(\\((asc|desc)[,]?([\\d]?)\\)|\\b)"},
	"order":   {"\\border(\\((asc|desc)[,]?([\\d]?)\\)|\\b)"},
	"seq":     {"\\bseq(\\(([\\w]+)[,]?([\\d]?)[,]?([\\d]?)\\)|\\b)"},
	"replace": {"\\breplace(\\(([\\d]+),([\\d]+)[,]?([\\s\\S]+)\\))"},
}

var IsNull = map[string]string{
	"否":   "not null",
	"N":   "not null",
	"NO":  "not null",
	"是":   "",
	"":    "",
	"Y":   "",
	"YES": "",
}

var IsMDNull = map[string]string{
	"NO":  "否",
	"N":   "否",
	"n":   "否",
	"no":  "否",
	"YES": "是",
	"yes": "是",
	"Y":   "是",
	"y":   "是",
}
