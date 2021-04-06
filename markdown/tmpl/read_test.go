package tmpl

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/micro-plat/lib4go/assert"
)

func TestRead(t *testing.T) {
	text := `## 一、商户信息

	###  1. 商户信息[ots_merchant_info]
	
	| 字段名      | 类型         | 默认值  | 为空  |     约束     | 描述     |
	| ----------- | ------------ | :-----: | :---: | :----------: | :------- |
	| mer_no      | varchar2(32) |         |  否   | PK,SEQ,RL,DI | 编号     |
	| mer_name    | varchar2(64) |         |  否   |   CRUQL,DN   | 名称     |
	| mer_crop    | varchar2(64) |         |  否   |   CRUQL,DN   | 公司     |
	| mer_type    | number(1)    |         |  否   |   CRUQL,DN   | 类型     |
	| bd_uid      | number(20)   |         |  否   |   CRUQL,DN   | 商务人员 |
	| status      | number(1)    |    0    |  否   |   RUQL,SL    | 状态     |
	| create_time | date         | sysdate |  否   |    RL,DTIME     | 创建时间 |`
	b := bytes.NewBuffer([]byte(text))
	lines, err := readMarkdownByReader(bufio.NewReader(b))
	assert.Equal(t, nil, err)
	assert.Equal(t, 11, len(lines))

	tl := line2TableLine(lines)
	assert.Equal(t, 1, len(tl.Lines))

	tb, err := tableLine2Table(tl)
	assert.Equal(t, nil, err)

	assert.Equal(t, 1, len(tb.Tbs))
	assert.Equal(t, 7, len(tb.Tbs[0].Rows))
	assert.Equal(t, "ots_merchant_info", tb.Tbs[0].Name)
	assert.Equal(t, "商户信息", tb.Tbs[0].Desc)
	assert.Equal(t, "create_time", tb.Tbs[0].Rows[6].Name)
	assert.Equal(t, "date", tb.Tbs[0].Rows[6].Type)
	assert.Equal(t, "sysdate", tb.Tbs[0].Rows[6].Def)
	assert.Equal(t, "否", tb.Tbs[0].Rows[6].IsNull)
	assert.Equal(t, "RL,DTIME", tb.Tbs[0].Rows[6].Con)
	assert.Equal(t, "创建时间", tb.Tbs[0].Rows[6].Desc)
}
func TestPKG(t *testing.T) {
	var cases = []struct {
		path string
		name string
	}{
		{path: "./modules/const/db/scheme", name: "scheme"},
		{path: "scheme", name: "scheme"},
		{path: "/", name: ""},
		{path: "/scheme", name: "scheme"},
		{path: "/modules/const/db", name: "db"},
	}
	for _, c := range cases {
		name := getPKSName(c.path)
		assert.Equal(t, c.name, name, c.path)
	}

}
