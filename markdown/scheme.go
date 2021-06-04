package markdown

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/hicli/markdown/tmpl"
	"github.com/urfave/cli"
)

const gofile = "gofile"

//createScheme 生成数据库结构
func createScheme(c *cli.Context) (err error) {
	if len(c.Args()) == 0 {
		return fmt.Errorf("未指定markdown文件")
	}
	if len(c.Args()) < 2 {
		return fmt.Errorf("未指定输出路径")
	}

	createfile := func(filepath, appendDir string) (err error) {
		//读取文件
		tbs, err := tmpl.Markdown2DB(filepath)
		if err != nil {
			return err
		}
		//设置包名称
		if c.Bool(gofile) {
			tbs.SetPkg(c.Args().Get(1))
		}
		//是否删除表
		tbs.DropTable(c.Bool("drop"))
		//过滤数据表
		tbs.FilterByKW(c.String("table"))
		tbs.Exclude()
		tbs.BuildSEQFile(c.Bool("seqfile"))

		//循环创建表
		outpath := c.Args().Get(1)
		if appendDir != "" {
			outpath += "/" + appendDir + ""
		}
		for _, tb := range tbs.Tbs {
			//创建文件
			path := tmpl.GetSchemePath(outpath, tb.Name, c.Bool(gofile))

			//翻译文件
			content, err := tmpl.Translate(tmpl.SQLTmpl, dbtp, tb)
			if err != nil {
				return err
			}
			fs, err := tmpl.Create(path, c.Bool("cover"))
			if err != nil {
				return err
			}
			logs.Log.Info("生成文件:", path)
			if _, err := fs.Write([]byte(content)); err != nil {
				return err
			}
		}
		if tbs.SEQFile {
			content, err := tmpl.Translate(tmpl.CreateSEQTable, dbtp, tbs)
			if err != nil {
				return err
			}
			path := tmpl.GetSEQFilePath(outpath, c.Bool(gofile))
			fs, err := tmpl.Create(path, c.Bool("cover"))
			if err != nil {
				return err
			}
			logs.Log.Info("生成文件:", path)
			fs.WriteString(content)
			fs.Close()
		}

		//生成安装文件
		if c.Bool(gofile) {
			content, err := tmpl.Translate(tmpl.InstallTmpl, dbtp, tbs)
			if err != nil {
				return err
			}
			path := tmpl.GetInstallPath(outpath)
			fs, err := tmpl.Create(path, c.Bool("cover"))
			if err != nil {
				return err
			}
			logs.Log.Info("生成文件:", path)
			fs.WriteString(content)
			fs.Close()
		}
		return nil
	}

	filePath := c.Args().First()
	if !strings.Contains(filePath, "*") {
		return createfile(filePath, "")
	}

	//找到符合的md文件

	var wg sync.WaitGroup
	files := getAllMatchMD(filePath)
	for _, v := range files {
		_, f := filepath.Split(v)
		wg.Add(1)
		go func(v, f string) {
			defer wg.Done()
			err := createfile(v, strings.TrimRight(f, ".md"))
			if err != nil {
				logs.Log.Errorf("[%s]发生错误：%+v", v, err)
			}
		}(v, f)
	}

	wg.Wait()

	return nil
}

func getAllMatchMD(path string) (paths []string) {

	//路径是的具体文件
	_, err := os.Stat(path)
	if err == nil {
		return []string{path}
	}
	//查找匹配的文件
	dir, f := filepath.Split(path)
	fmt.Println("dir", dir, f)

	regexName := strings.Replace(strings.Replace(f, ".md", "\\.md", -1), "*", ".+", -1)
	reg := regexp.MustCompile(regexName)

	fmt.Println("regexName：", regexName)
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		fname := f.Name()
		if strings.HasPrefix(fname, ".") || f.IsDir() {
			continue
		}
		if reg.Match([]byte(fname)) {
			paths = append(paths, filepath.Join(dir, fname))
		}
	}
	return paths
}
