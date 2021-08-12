package markdown

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/micro-plat/hicli/markdown/const/enums"
	"github.com/micro-plat/hicli/markdown/tmpl"
	"github.com/micro-plat/lib4go/types"
	"github.com/urfave/cli"
)

//createDataDiff 比较差异文件
func createDataDiff(c *cli.Context) (err error) {

	if len(c.Args()) < 2 {
		return fmt.Errorf("未指定需要对比的源xlsx文件和目标xlsx文件")
	}

	//读取excel
	sheet := types.GetString(c.String("sheet"), "Sheet1")
	sourcePath := c.Args().First()
	sourceFile, err := tmpl.GetExcel(sourcePath, sheet)
	if err != nil {
		return fmt.Errorf("读取%s出错：%+v", sourcePath, err)
	}
	sourceRows := sourceFile.GetRows(sheet)
	if len(sourceRows) < 1 {
		return fmt.Errorf("%s为空", sourcePath)
	}
	targetPath := c.Args().Get(1)
	targetFile, err := tmpl.GetExcel(targetPath, sheet)
	if err != nil {
		return fmt.Errorf("读取%s出错：%+v", targetPath, err)
	}
	targetRows := targetFile.GetRows(sheet)
	if len(targetRows) < 1 {
		return fmt.Errorf("%s为空", targetPath)
	}

	//处理比较列
	titles := make([]string, 0)
	if c.String("titles") != "" {
		titles = strings.Split(c.String("titles"), ",")
	}
	sourceTitle := sourceRows[0]
	targetTitle := targetRows[0]
	if len(titles) == 0 { //全部列比较
		if !reflect.DeepEqual(sourceTitle, targetTitle) {
			return fmt.Errorf("%s和%s标题或者标题顺序不一致", sourcePath, targetPath)
		}
	}
	sourceDataMap, err := getDataMap(sourceRows, titles)
	if err != nil {
		return fmt.Errorf("%s：%+v", sourcePath, err)
	}
	targetDataMap, err := getDataMap(targetRows, titles)
	if err != nil {
		return fmt.Errorf("%s：%+v", targetPath, err)
	}

	//比较并将不一样的数据样式填充为红色
	diff := dataDiff(sourceDataMap, targetDataMap)
	style, err := sourceFile.NewStyle(`{"fill":{"type":"gradient","color":["#FF0000","#FF0000"],"shading":1}}`)
	if err != nil {
		return fmt.Errorf("设置样式错误：%+v", err)
	}
	for _, v := range diff { //修改excel
		fmt.Println("diff", v.index, v.key)
		sourceFile.SetCellStyle(sheet, fmt.Sprintf("A%d", v.index), fmt.Sprintf("%s%d", getEnd(len(sourceTitle)), v.index), style)
	}

	//保存样式
	err = sourceFile.Save()
	if err != nil {
		return fmt.Errorf("保存样式错误：%+v", err)
	}
	return nil
}

//getDataMap 获取excel对应map，以_拼接对比列的数据
func getDataMap(rows [][]string, titles []string) (dataMap map[string]int, err error) {
	dataMap = make(map[string]int, len(rows))
	if len(titles) == 0 { //全部列
		for index, row := range rows[1:] {
			key := strings.Join(row, "_")
			dataMap[key] = index
		}
		return
	}

	compareTitlesMap := make(map[string]int, len(titles))
	for k, title := range titles {
		compareTitlesMap[title] = k
	}

	dataTitlesMap := make(map[string]int, len(rows[0]))
	for k, title := range rows[0] {
		dataTitlesMap[title] = k
	}

	indexs := make([]int, 0, len(titles)) //对比例所在的顺序位置
	for k := range compareTitlesMap {
		v, ok := dataTitlesMap[k]
		if !ok {
			err = fmt.Errorf("不存在比较的列%s", k)
			return
		}
		indexs = append(indexs, v)
	}
	sort.Ints(indexs)
	for index, row := range rows[1:] {
		tempRow := make([]string, 0, len(indexs))
		for _, v := range indexs { //对比列的数据
			tempRow = append(tempRow, row[v])
		}
		key := strings.Join(tempRow, "_")
		dataMap[key] = index
	}
	return
}

type rowDiff struct {
	key       string
	index     int
	operation enums.Operation
}

func dataDiff(sourceMap, targetMap map[string]int) []*rowDiff {
	r := make([]*rowDiff, 0)
	//新增
	for name, index := range sourceMap {
		if _, ok := targetMap[name]; !ok {
			diff := &rowDiff{
				key:       name,
				index:     index + 2, //加上标题
				operation: enums.DiffInsert,
			}
			r = append(r, diff)
		}
	}
	return r
}

func getEnd(count int) string {
	count = count - 1
	if count < 26 {
		return string(65 + count)
	}
	return fmt.Sprintf("%s%s", string(64+count/26), string(65+count%26))
}
