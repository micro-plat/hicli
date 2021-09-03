package tmpl

import (
	"encoding/json"
	"path"
	"sort"
	"sync"

	"github.com/micro-plat/hicli/markdown/utils"
)

//SnippetConf 用于vue的路由,hydra服务的注册,import的路径等代码片段生成
type SnippetConf struct {
	Name      string `json:"name"`       //表名
	HasDetail bool   `json:"has_detail"` //是否有详情页
	HasList   bool   `json:"has_list"`   //是否有详情页
	BasePath  string `json:"base_path"`  //项目路径
	Desc      string `json:"desc"`       //描述
	PKG       string
	PkGAlias  string
	SortName  string `json:"sort_name"`
}

//NewSnippetConf .
func NewSnippetConf(t *Table) *SnippetConf {
	lrows := getRows("l")(t.Rows)
	drows := getRows("d")(t.Rows)
	return &SnippetConf{
		Name:      t.Name,
		SortName:  rmhd(t.Name),
		HasList:   len(lrows) > 0,
		HasDetail: len(drows) > 0,
		BasePath:  t.BasePath,
		Desc:      t.Desc,
	}
}

//SaveConf 保存配置
func (t *SnippetConf) SaveConf(confPath string) error {
	if confPath == "" {
		return nil
	}

	//读取配置
	conf := make(map[string]*SnippetConf)
	err := readConf(confPath, &conf)
	if err != nil {
		return err
	}

	//设置配置
	conf[t.Name] = t

	//写入配置
	return writeConf(confPath, conf)
}

//GetSnippetConf 获取配置
func GetSnippetConf(path string) (SnippetConfs, error) {

	conf := make(map[string]*SnippetConf)
	err := readConf(path, &conf)
	if err != nil {
		return nil, err
	}

	confs := make(SnippetConfs, 0)
	for _, v := range conf {
		confs = append(confs, v)
	}

	sort.Sort(confs)

	return confs, nil
}

//SnippetConfs 排序用
type SnippetConfs []*SnippetConf

func (t SnippetConfs) Len() int {
	return len(t)
}

//从低到高
func (t SnippetConfs) Less(i, j int) bool {
	return t[i].SortName < t[j].SortName
}

func (t SnippetConfs) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

//WebExtConf 自定义路由配置
type WebExtConf struct {
	Name        string `json:"name"`       //表名
	Path        string `json:"path"`       //路径
	Component   string `json:"component"`  //页面路径
	HasDetail   bool   `json:"has_detail"` //是否有详情页
	Desc        string `json:"desc"`
	Independent bool   `json:"independent"` //单独页面
}

//GetWebExtConf 获取配置
func GetWebExtConf(path string) ([]*WebExtConf, error) {

	confs := make([]*WebExtConf, 0)
	err := readConf(path, &confs)
	if err != nil {
		return nil, err
	}

	return confs, nil
}

//FieldConf 用于field文件生成
type FieldConf struct {
	Fields []*FieldItem `json:"fields"`
}

//FieldItem .
type FieldItem struct {
	Desc  string `json:"desc"`
	Name  string `json:"name"`
	Table string `json:"table"`
}

//NewFieldConf .
func NewFieldConf(t *Table) *FieldConf {
	fields := []*FieldItem{}
	for _, v := range t.Rows {
		item := &FieldItem{
			Desc:  v.Desc,
			Name:  v.Name,
			Table: t.Name,
		}
		fields = append(fields, item)
	}
	return &FieldConf{Fields: fields}
}

//GetFieldConf .
func GetFieldConf(path string) (map[string]*FieldItem, error) {
	conf := make(map[string]*FieldItem)
	err := readConf(path, &conf)
	return conf, err
}

func (t *FieldConf) SaveConf(confPath string) error {
	if confPath == "" {
		return nil
	}

	//读取配置
	conf := make(map[string]*FieldItem)
	err := readConf(confPath, &conf)
	if err != nil {
		return err
	}

	//设置配置
	for _, v := range t.Fields {
		if _, ok := conf[v.Name]; ok {
			continue
		}
		conf[v.Name] = v
	}

	//写入配置
	return writeConf(confPath, conf)
}

var mutex sync.Mutex

func writeConf(confPath string, conf interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()
	//创建文件
	fs, err := Create(confPath, true)
	if err != nil {
		return err
	}

	//写入
	r, err := json.Marshal(conf)
	if err != nil {
		return err
	}
	fs.WriteString(string(r))
	fs.Close()
	return nil
}

func readConf(path string, conf interface{}) error {
	//读取
	s, err := Read(path)
	if err != nil {
		return err
	}

	if len(s) > 0 {
		if err = json.Unmarshal(s, &conf); err != nil {
			return err
		}
	}

	return nil
}

func GetFieldConfPath(root string) string {
	projectPath := utils.GetProjectPath(root)
	if projectPath == "" {
		return ""
	}
	return path.Join(projectPath, ".hicli/server_filed.json")
}

func GetWebConfPath(root string) string {
	projectPath := utils.GetProjectPath(root)
	webPath, _ := utils.GetWebSrcPath(projectPath)
	if webPath == "" {
		return ""
	}
	return path.Join(webPath, ".hicli/web.json")
}

func GetWebExtConfPath(root string) string {
	projectPath := utils.GetProjectPath(root)
	webPath, _ := utils.GetWebSrcPath(projectPath)
	if webPath == "" {
		return ""
	}
	return path.Join(webPath, "public/router.ext.json")
}

func GetGoConfPath(root string) string {
	projectPath := utils.GetProjectPath(root)
	if projectPath == "" {
		return ""
	}
	return path.Join(projectPath, ".hicli/server.json")
}
