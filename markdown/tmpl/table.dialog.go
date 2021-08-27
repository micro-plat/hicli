package tmpl

import (
	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/lib4go/types"
)

type Dialog struct {
	Name      string
	Method    string
	Path      string
	Condition string
}

func (t *Table) elListDialog(key string) {

	info := &Dialog{}

	//Name
	info.Name = getSubConContent(key, "name")(t.ExtInfo)
	if info.Name == "" {
		logs.Log.Warn("列表按钮dialog的name选项未配置:", t.ExtInfo)
		return
	}

	//method
	info.Method = types.GetString(getSubConContent(key, "method")(t.ExtInfo), getClickFunc(info.Name, key))
	if info.Method == "Edit" {
		t.BtnShowEdit = true
	}
	if info.Method == "Detail" {
		t.BtnShowDetail = true
	}

	info.Path = getSubConContent(key, "path")(t.ExtInfo)
	if info.Path == "" {
		logs.Log.Warn("列表按钮dialog的path选项未配置:", key, t.ExtInfo)
		return
	}

	//confirm
	//info.Confirm = getSubConContent(key, "confirm")(t.ExtInfo)

	//condition
	info.Condition = translateCondition(getSubConContent(key, "condition")(t.ExtInfo))

	t.ListDialogs = append(t.ListDialogs, info)

}

func (t *Table) elDetailDialog(key string) {

	info := &Dialog{}

	//Name
	info.Name = getSubConContent(key, "name")(t.ExtInfo)
	if info.Name == "" {
		logs.Log.Warn("详情页面按钮dialog的name选项未配置:", t.ExtInfo)
		return
	}

	//method
	info.Method = types.GetString(getSubConContent(key, "method")(t.ExtInfo), getClickFunc(info.Name, key))
	// if info.Method == "Edit" {
	// 	t.BtnShowEdit = true
	// }
	// if info.Method == "Detail" {
	// 	t.BtnShowDetail = true
	// }

	info.Path = getSubConContent(key, "path")(t.ExtInfo)
	if info.Path == "" {
		logs.Log.Warn("详情页面dialog的path选项未配置:", key, t.ExtInfo)
		return
	}

	//confirm
	//info.Confirm = getSubConContent(key, "confirm")(t.ExtInfo)

	//condition
	info.Condition = translateCondition(getSubConContent(key, "condition")(t.ExtInfo))

	//@todo
	//t.ListDialogs = append(t.ListDialogs, info)

}

func (t *Table) elQueryDialog(key string) {

	info := &Dialog{}

	//Name
	info.Name = getSubConContent(key, "name")(t.ExtInfo)
	if info.Name == "" {
		logs.Log.Warn("查询dialog的name选项未配置:", t.ExtInfo)
		return
	}

	//method
	info.Method = types.GetString(getSubConContent(key, "method")(t.ExtInfo), getClickFunc(info.Name, key))
	if info.Method == "Add" {
		t.BtnShowAdd = true
	}

	info.Path = getSubConContent(key, "path")(t.ExtInfo)
	if info.Path == "" {
		logs.Log.Warn("查询按钮dialog的path选项未配置:", key, t.ExtInfo)
		return
	}

	//confirm
	//info.Confirm = getSubConContent(key, "confirm")(t.ExtInfo)

	//condition
	info.Condition = translateCondition(getSubConContent(key, "condition")(t.ExtInfo))

	t.QueryDialogs = append(t.QueryDialogs, info)

}
