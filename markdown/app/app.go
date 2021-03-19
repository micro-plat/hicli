package app

import (
	"fmt"
	"path/filepath"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/hicli/markdown/tmpl"
	"github.com/micro-plat/hicli/markdown/utils"
)

var tmptls = map[string]string{
	"main.go": tmplMainGo,
	"conf.go": tmplConfGo,
	"go.mod":  tmplGoMod,
}

var tmptlsSSO = map[string]string{
	"conf.go": tmplConfSSoGo,
}

//CreateApp 创建web项目
func CreateApp(name string, sso bool) error {
	projectPath := utils.GetProjectPath(name)
	tmptls := tmptls
	if sso {
		for k, v := range tmptlsSSO {
			tmptls[k] = v
		}
	}
	for file, template := range tmptls {
		//翻译文件
		param := map[string]interface{}{
			"projectPath": projectPath,
			"router":      true,
			"basePath":    utils.GetProjectBasePath(projectPath),
		}
		content, err := tmpl.Translate(template, "", param)
		if err != nil {
			return fmt.Errorf("翻译%s模板出错:%+v", file, err)
		}
		fs, err := tmpl.Create(filepath.Join(projectPath, file), true)
		if err != nil {
			return err
		}
		fs.WriteString(content)
		fs.Close()
		logs.Log.Info("生成文件:", file)
	}
	return nil
}
