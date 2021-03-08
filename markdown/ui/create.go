package ui

import (
	"path/filepath"

	"github.com/codeskyblue/go-sh"
	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/hicli/markdown/tmpl"
)

var tmptls = map[string]string{
	"src/App.vue":                srcAppVUE,
	"src/main.js":                srcMainJS,
	"src/pages/system/menus.vue": srcPagesSystemMenus,
	"src/router/index.js":        srcRouterIndexJS,
	"src/utility/sys.js":         srcUtilitySysJS,
	"src/utility/http.js":        srcUtilityHTTPJS,
	"src/utility/enum.js":        srcUtilityEnumJS,
	"src/utility/filter.js":      srcUtilityFilter,
	"src/utility/utility.js":     srcUtilityUtilityJS,
	"src/utility/main.js":        srcUtilityMainJS,
	"src/utility/message.js":     srcUtilityMessageJS,
	"src/utility/env.js":         srcUtilityEnvJS,
	"src/utility/package.json":   srcUtilityPackageJson,
	"public/env.conf.json":       srcPublicEnvConfJson,
	"public/index.html":          publicIndexHTML,
	"package.json":               packageJSON,
	"babel.config.js":            babelConfigJS,
	"vue.config.js":              vueConfigJS,
	".gitignore":                 gitignore,
}

var ssoTmptls = map[string]string{
	"public/env.conf.json": srcSSOPublicEnvConfJson,
	"src/main.js":          srcSSOMainJS,
}

//CreateWeb 创建web项目
func CreateWeb(name string, sso bool) error {
	if sso {
		for k, v := range ssoTmptls {
			tmptls[k] = v
		}
	}

	return createFiles(name)
}

//Clear 清理缓存
func Clear(dir string) error {
	if err := run(dir, "npm", "install", "--no-optional", "--verbose"); err != nil {
		return err
	}
	return run(dir, "npm", "install")
}

func run(dir string, name string, args ...interface{}) error {
	session := sh.InteractiveSession()
	session.SetDir(filepath.Join("./", dir))
	logs.Log.Info(append([]interface{}{name}, args...)...)
	session.Command(name, args...)
	return session.Run()
}

//createFiles 创建文件
func createFiles(name string) error {
	for path, content := range tmptls {
		fs, err := tmpl.Create(filepath.Join(".", name, path), true)
		if err != nil {
			return err
		}
		fs.WriteString(content)
		fs.Close()
		logs.Log.Info("生成文件:", path)
	}
	return nil
}
