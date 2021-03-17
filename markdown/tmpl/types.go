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
var def2mysql = []map[string]string{
	{
		"^$":         "",
		"^-$":        "default '-'",
		"^seq$":      "",
		"^sysdate$":  "default current_timestamp",
		"^([0-9]+)$": "default *",
	},
	{
		"^(.+)$": "default '*'",
	},
}

var any2code = map[string]string{
	"^date$":                      "time.Time",
	"^datetime$":                  "time.Time",
	"^timestamp$":                 "time.Time",
	"^decimal$":                   "types.Decimal",
	"^float$":                     "types.Decimal",
	"^int$":                       "int",
	"^number\\([1-2]{1}\\)$":      "int",
	"^number\\([3-9]{1}\\)$":      "int",
	"^number\\(10\\)$":            "int",
	"^number\\(1[1-9]{1}\\)$":     "int64",
	"^number\\(2[0-9]{1}\\)$":     "int64",
	"^number\\((\\d+),(\\d+)\\)$": "types.Decimal",
	"^varchar\\(\\d+\\)$":         "string",
	"^varchar2\\(\\d+\\)$":        "string",
	"^string$":                    "string",
	"^text$":                      "string",
	"^longtext$":                  "string",
	"^clob$":                      "string",
}

var keywordSubMatch = `\b%s\(([\w\s-:#,|/]+)\)`

var keywordMatch = []string{"^\\w*%s\\w*$", ",\\w*%s\\w*,", "^\\w*%s\\w*,", ",\\w*%s\\w*$"}

var cons = map[string][]string{
	"pk":    {"\\bpk\\b"},
	"seq":   {"\\bseq\\b"},
	"di":    {"\\bdi\\b"},
	"dn":    {"\\bdn\\b"},
	"sl":    {"\\bsl(\\([\\w,]+\\)|\\b)"},
	"slm":   {"\\bslm(\\([\\w,]+\\)|\\b)"},
	"rd":    {"\\brd(\\([\\w,]+\\)|\\b)"},
	"cb":    {"\\bcb(\\([\\w,]+\\)|\\b)"},
	"ta":    {"\\bta(\\([\\w,]+\\)|\\b)"},
	"cc":    {"\\bcc(\\(\\w+\\)|\\b)"},
	"idx":   {"\\bidx(\\(([\\w]+)[,]?([\\d]?)\\)|\\b)"},
	"unq":   {"\\bunq(\\(([\\w]+)[,]?([\\d]?)\\)|\\b)"},
	"d":     {"\\bd(\\([0-9]*\\)|\\b)"},
	"c":     {"\\bc(\\([\\w,:#]+\\)|\\b)"},
	"u":     {"\\bu(\\([\\w,:#]+\\)|\\b)"},
	"r":     {"\\br(\\([\\w,:#]+\\)|\\b)"},
	"l":     {"\\bl(\\([\\w,:#]+\\)|\\b)"},
	"q":     {"\\bq(\\([\\w,:#]+\\)|\\b)"},
	"sort":  {"\\bsort(\\((asc|desc)[,]?([\\d]?)\\)|\\b)"},
	"order": {"\\border(\\((asc|desc)[,]?([\\d]?)\\)|\\b)"},
	"*":     {"\\b%s\\b"},
}
var mysqlIsNull = map[string]string{
	"否":   "not null",
	"N":   "not null",
	"NO":  "not null",
	"是":   "",
	"":    "",
	"Y":   "",
	"YES": "",
}
