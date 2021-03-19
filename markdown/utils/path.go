package utils

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

//GetHicliHomePath 获取用户home目录
func GetHicliHomePath() string {
	user, err := user.Current()
	if nil == err {
		return filepath.Join(user.HomeDir, ".hicli")
	}

	// cross compile support
	var home string
	if "windows" == runtime.GOOS {
		home, err = homeWindows()
	} else {
		home, err = homeUnix()
	}

	if err != nil {
		panic(err)
	}

	return filepath.Join(home, ".hicli")
}

func homeUnix() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}

//GetProjectPath 获取项目路径
func GetProjectPath(root string) string {
	if !strings.HasPrefix(root, "./") && !strings.HasPrefix(root, "/") && !strings.HasPrefix(root, "../") {
		srcPath, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		root = filepath.Join(srcPath, root)
	}

	aPath, err := filepath.Abs(root)
	if err != nil {
		panic(fmt.Errorf("不是有效的项目路径:%s,%+v", root, err))
	}
	return aPath
}

//GetWebSrcPath 获取web项目src目录
//判断路径下是否有src目录且src下有App.vue,有则返回src目录和项目目录
//默认返回空
func GetWebSrcPath(projectPath string) (string, string) {
	n := strings.LastIndex(projectPath, "src")
	if n < 0 {
		return "", ""
	}
	parentDir := projectPath[0:n]
	srcPath := path.Join(parentDir, "src")
	appVuePath := path.Join(srcPath, "App.vue")
	if PathExists(appVuePath) { //存在返回
		return parentDir, srcPath
	}
	return GetWebSrcPath(parentDir)
}

func getGOENV(key string) string {
	envs, err := exec.Command("go", "env").Output()
	if err != nil {
		panic(fmt.Errorf("执行go env出错，%+v", err))
	}
	rex := regexp.MustCompile(fmt.Sprintf(`%s=(.*)`, key))
	strs := rex.FindAllString(string(envs), -1)
	if len(strs) < 1 {
		return ""
	}
	env := strs[0]
	env = strings.TrimPrefix(env, fmt.Sprintf(`%s=`, key))
	env = strings.TrimPrefix(env, `"`)
	env = strings.TrimRight(env, `"`)
	return env
}

//GetGOMOD .
func GetGOMOD() string {
	return getGOENV("GOMOD")
}

//GetProjectBasePath 如果开启了gomod 则返回module名
//未使用gomod则判断path中是否存在$GOPATH，存在则返回$GOPATH下面的名字
//默认返回空
func GetProjectBasePath(projectPath string) string {
	gomod := getGOENV("GOMOD")
	basePath := ""
	if gomod != "" && strings.Contains(gomod, projectPath) {
		f, err := os.Open(gomod)
		if err != nil {
			panic(fmt.Errorf("打开%s文件出错，%+v", gomod, err))
		}
		defer f.Close()

		br := bufio.NewReader(f)
		for {
			a, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			line := string(a)
			if strings.HasPrefix(line, "module ") {
				basePath = strings.TrimPrefix(line, "module ")
				break
			}
		}
		return basePath
	}

	gopath := getGOENV("GOPATH")
	if gopath != "" {
		root := filepath.Join(gopath, "src")
		if strings.HasPrefix(strings.ToLower(projectPath), strings.ToLower(root)) {
			basePath = projectPath[len(root)+1:]
		}
		return basePath
	}
	return ""
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
