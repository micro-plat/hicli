package tmpl

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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
		return filepath.Join(outpath, "seq_ids.sql.go")
	}
	return filepath.Join(outpath, "seq_ids.sql")
}

//PathExists 文件是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

//Create 创建文件，文件夹 存在时写入则覆盖
func Create(path string, append bool) (file *os.File, err error) {
	dir := filepath.Dir(path)
	if !PathExists(dir) {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("创建文件夹%s失败:%v", path, err)
		}
	}

	var srcf *os.File
	if !PathExists(path) {
		srcf, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("无法打开文件:%s(err:%v)", path, err)
		}
		return srcf, nil

	}
	if !append {
		return nil, fmt.Errorf("文件:%s已经存在", path)
	}
	srcf, err = os.OpenFile(path, os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件:%s(err:%v)", path, err)
	}
	return srcf, nil

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
