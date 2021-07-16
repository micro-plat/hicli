package tmpl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

//GetSchemePath 获取Scheme路径
func GetSchemePath(outpath string, name string, gofile bool) string {
	path := filepath.Join(outpath, fmt.Sprintf("%s.sql", name))
	if gofile {
		path = filepath.Join(outpath, fmt.Sprintf("%s.sql.go", name))

	}
	return path
}

//GetInstallPath 获取DB安装文件
func GetInstallPath(outpath string) string {
	return filepath.Join(outpath, "install.go")
}

//GetSEQFilePath 获取DB安装文件
func GetSEQFilePath(outpath string, gofile bool) string {
	if gofile {
		return filepath.Join(outpath, "sys.sequence.info.go")
	}
	return filepath.Join(outpath, "sys.sequence.info.sql")
}

//PathExists 文件是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

//Create 创建文件，文件夹 存在时写入则覆盖
func Create(path string, cover bool) (file *os.File, err error) {
	dir := filepath.Dir(path)
	if !PathExists(dir) {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("创建文件夹%s失败:%v", path, err)
		}
	}

	if !PathExists(path) {
		file, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("无法打开文件:%s(err:%v)", path, err)
		}
		return
	}

	if !cover {
		return nil, fmt.Errorf("文件:%s已经存在", path)
	}

	// 文件存在且文件头部设置不覆盖标识
	if checkCover(path) {
		return
	}

	file, err = os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件:%s(err:%v)", path, err)
	}
	return
}

//返回只可读的文件 外部调用注意Close()
func checkCover(path string) (cover bool) {
	file, _ := os.Open(path)
	//读取第一行的
	defer file.Close()
	br := bufio.NewReader(file)
	a, _, _ := br.ReadLine()
	return strings.HasPrefix(string(a), "//exclude") || strings.HasPrefix(string(a), "<!-- exclude -->")

}

func Read(path string) (s []byte, err error) {
	if !PathExists(path) {
		return
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}
