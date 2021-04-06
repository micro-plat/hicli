package gitlabs

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/codeskyblue/go-sh"
	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/lib4go/envs"
)

//Repository 仓库信息
type Repository struct {
	Name     string `json:"name"`
	Desc     string `json:"description"`
	Path     string `json:"relative_path"`
	Type     string `json:"type"`
	FullPath string `json:"-"`
}

//NewRepository 创建分组/仓库信息
func NewRepository(fullPath string) *Repository {
	u, _ := url.Parse(fullPath)
	return &Repository{FullPath: fullPath, Path: u.Path, Type: "project"}
}

//String 输出内容
func (r *Repository) String() string {
	if buff, err := json.Marshal(&r); err == nil {
		return string(buff)
	}
	return ""
}

//GetLocalPath 获取本地路径
func (r *Repository) GetLocalPath() string {
	u, _ := url.Parse(r.FullPath)
	gopath := envs.GetString("GOPATH")
	return filepath.Join(gopath, "src", u.Host, r.Path)
}

//Exists 本地仓库是否是存在
func (r *Repository) Exists() bool {
	rpath := filepath.Join(r.GetLocalPath(), ".git")
	if _, err := os.Stat(rpath); err != nil {
		return os.IsExist(err)
	}
	return true
}

//Reset 拉取项目
func (r *Repository) Reset(branch ...string) error {
	session := sh.InteractiveSession()
	session.SetDir(r.GetLocalPath())
	for _, b := range branch {
		session.Command("git", "branch")
		buff, err := session.Output()
		if err != nil {
			return err
		}
		if hasBranch(string(buff), b) {
			logs.Log.Info(r.GetLocalPath(), ">", "git", "reset", "--hard")
			session.Command("git", "reset", "--hard")
			if err := session.Run(); err != nil {
				return err
			}
			logs.Log.Info(r.GetLocalPath(), ">", "git", "checkout", b)
			session.Command("git", "checkout", b)
			if err := session.Run(); err != nil {
				return err
			}
			logs.Log.Info(r.GetLocalPath(), ">", "git", "reset", "--hard")
			session.Command("git", "reset", "--hard")
			if err := session.Run(); err != nil {
				return err
			}
		}
	}
	return nil
}

//Update 更新项目
func (r *Repository) Update() error {
	session := sh.InteractiveSession()
	session.SetDir(r.GetLocalPath())
	logs.Log.Info("hicli", "pull", r.FullPath)
	session.Command("hicli", "pull", r.FullPath)
	if err := session.Run(); err != nil {
		return err
	}
	logs.Log.Info(r.GetLocalPath(), ">", "go", "install")
	session.Command("go", "install")
	if err := session.Run(); err != nil {
		return err
	}

	logs.Log.Info(r.GetLocalPath(), ">", "hicli", "--version")
	session.Command("hicli", "--version")
	if err := session.Run(); err != nil {
		return err
	}
	return nil
}

//Pull 拉取项目
func (r *Repository) Pull(branch ...string) error {
	session := sh.InteractiveSession()
	session.SetDir(r.GetLocalPath())
	session.Command("git", "branch")
	buff, err := session.Output()
	if err != nil {
		return err
	}
	for _, b := range branch {
		if hasBranch(string(buff), b) {
			logs.Log.Info(r.GetLocalPath(), ">", "git", "pull", "origin", b)
			session.Command("git", "pull", "origin", b)
		} else {
			logs.Log.Info(r.GetLocalPath(), ">", "git", "fetch", "origin", b+":"+b)
			session.Command("git", "fetch", "origin", b+":"+b)
		}
		if err := session.Run(); err != nil {
			return err
		}
	}
	return nil
}

//Clone 克隆项目
func (r *Repository) Clone() error {
	session := sh.InteractiveSession()
	session.Command("git", "clone", r.FullPath, r.GetLocalPath())
	logs.Log.Info("git", "clone", r.FullPath, r.GetLocalPath())
	err := session.Run()
	if err != nil && strings.Contains(err.Error(), "exit status 128") {
		return fmt.Errorf("fatal: 目标路径 '%s' 已经存在，并且不是一个空目录。", r.GetLocalPath())
	}
	return err
}

//hasBranch 本地是否包含指定分支
func hasBranch(s string, b string) bool {
	items := strings.Split(s, "\n")
	for _, i := range items {
		if strings.Contains(i, b) {
			return true
		}
	}
	return false
}
