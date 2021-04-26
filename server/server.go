package server

import (
	"fmt"
	ifs "io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/codeskyblue/go-sh"
	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/hicli/markdown/utils"
	"github.com/urfave/cli"

	"github.com/fsnotify/fsnotify"
)

type server struct {
	session    *sh.Session
	serverName string
	path       string
	fs         *fs
	running    bool
	hasNotify  bool
	notifyChan chan int
	closeChan  chan int
	watchers   map[string]string
	ticker     *time.Ticker
	errChan    chan error
	startFlag  map[string][]interface{}
}

func newServer(c *cli.Context) (*server, error) {
	//path
	path, watchpath, err := checkPath(c)
	if err != nil {
		return nil, err
	}

	session := sh.InteractiveSession()
	session.SetDir(path)

	r, err := NewFileSystem(watchpath)
	if err != nil {
		return nil, err
	}
	return &server{
		serverName: filepath.Base(path),
		path:       watchpath,
		fs:         r,
		session:    session,
		notifyChan: make(chan int, 1),
		closeChan:  make(chan int, 1),
		errChan:    make(chan error, 1),
		watchers:   make(map[string]string),
		ticker:     time.NewTicker(time.Millisecond * 500),
		startFlag:  getStartFlag(c),
	}, nil
}

var tags = map[string][]string{
	"install": {"tags", "mod"},
	"run":     {"fixed", "registry", "plat", "system", "server_types", "cluster", "trace", "tport", "mask", "debug"},
}

func checkPath(c *cli.Context) (projectPath, watchPath string, err error) {

	//判断项目是否存在
	projectPath = utils.GetProjectPath(c.Args().Get(0))
	if !utils.PathExists(filepath.Join(projectPath, "main.go")) {
		return "", "", fmt.Errorf("未指定的运行应用程序的项目路径:%s", projectPath)
	}

	//指定监控当前目录
	watchPath = projectPath
	if c.Bool("work") {
		watchPath = utils.GetProjectPath("./")
	}

	if c.String("exclude") != "" {
		defExcludePath = append(defExcludePath, strings.Split(c.String("exclude"), ",")...)
	}

	return
}

func getStartFlag(c *cli.Context) map[string][]interface{} {
	startFlag := make(map[string][]interface{})
	for k, v := range tags {
		startFlag[k] = []interface{}{k}

		for _, v1 := range strings.Split(c.String(k), " ") {
			if v1 != "" {
				startFlag[k] = append(startFlag[k], v1)
			}
		}

		for _, v1 := range v {
			if c.Bool(v1) {
				startFlag[k] = append(startFlag[k], fmt.Sprintf("-%s", v1))
				continue
			}
			if c.String(v1) != "" {
				startFlag[k] = append(startFlag[k], fmt.Sprintf("-%s", v1))
				startFlag[k] = append(startFlag[k], c.String(v1))
			}
		}

	}

	return startFlag
}

//Reset 拉取项目
func (s *server) resume() {

	s.fs.Start()
	go s.start()

	for {
		select {
		case <-s.notifyChan:
			s.pause()
			go s.start()
		case <-s.closeChan:
			s.pause()
			return
		}
	}
}

func (s *server) start() {
	if s.running {
		return
	}
	s.running = true

	//开启文件监控
	go s.watch()

	//文件打包
	err := s.session.Command("go", s.startFlag["install"]...).Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	//程序启动
	s.session.Command(s.serverName, s.startFlag["run"]...).Run()
}

func (s *server) pause() {
	if s.running {
		s.running = false
		s.session.Kill(os.Interrupt)
	}
}

func (s *server) close() (err error) {

	var sigChan = make(chan os.Signal, 3)
	signal.Notify(sigChan, syscall.SIGTERM, os.Interrupt)

	select {
	case <-sigChan:
	case err = <-s.errChan:
		s.running = false
		close(s.closeChan)
		close(s.notifyChan)
		close(s.errChan)
		time.Sleep(time.Second)
	}

	return err
}

var defExcludePath = []string{"vendor", "node_modules", ".gitignore", ".hicli", ".git"}

func (s *server) isExclude(path string) bool {
	if path == s.path {
		return false
	}

	for _, v := range defExcludePath {
		if strings.Contains(path, v) {
			return true
		}
	}

	return false
}

func (s *server) watch() {
	filepath.WalkDir(s.path, func(path string, d ifs.DirEntry, err error) error {
		if !s.isExclude(path) && d.IsDir() {
			if _, ok := s.watchers[path]; !ok {
				s.watchers[path] = path
				go s.watchChildren(path)
			}
		}
		return nil
	})
}

func (s *server) watchChildren(path string) {
	//监控子节点变化
	ch, err := s.fs.WatchChildren(path)
	if err != nil && !s.fs.done {
		s.fs.Close()
		s.errChan <- err
		return
	}

	for {
		select {
		case <-s.ticker.C:
			if s.hasNotify {
				s.notifyChan <- 1
				s.hasNotify = false
				delete(s.watchers, path)
				return
			}
		case <-s.closeChan:
			s.fs.Close()
			return
		case cldWatcher := <-ch:
			if cldWatcher.GetOp() == fsnotify.Remove {
				logs.Log.Info("----------------------项目文件被删除，应用程序重启----------------------")
				s.notifyChan <- 1
				delete(s.watchers, path)
				return
			}
			if cldWatcher.GetError() != nil && s.running {
				s.errChan <- fmt.Errorf("监控项目文件发生错误：%+v", cldWatcher.GetError())
				return
			}
			if !s.isExclude(cldWatcher.GetPath()) {
				logs.Log.Info("----------------------项目发生变化，应用程序重启----------------------")
				s.hasNotify = true
			}
		LOOP:
			ch, err = s.fs.WatchChildren(path)
			if err != nil {
				if s.running {
					s.errChan <- fmt.Errorf("应用程序运行中，未获取到文件监控")
					return
				}
				logs.Log.Errorf("文件监控错误%+v", err)
				time.Sleep(time.Second * 3)
				goto LOOP
			}
		}
	}
}
