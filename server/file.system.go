package server

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

var fileMode = os.FileMode(0664)
var dirMode = os.FileMode(0755)

type fsChildrenWatcher struct {
	watcher  chan ChildrenWatcher
	event    chan fsnotify.Event
	syncChan chan fsnotify.Event
}

type fs struct {
	watcher             *fsnotify.Watcher
	childrenWatcherMaps map[string]*fsChildrenWatcher
	watchLock           sync.Mutex
	tempNodes           map[string]bool
	tempNodeLock        sync.Mutex
	closeCh             chan struct{}
	rootDir             string
	done                bool
}

//NewFileSystem 文件系统的注册中心
func NewFileSystem(rootDir string) (*fs, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	rootDir = strings.TrimRight(rootDir, "/")
	if strings.HasPrefix(rootDir, "./") {
		rootDir = rootDir[2:]
	}
	registryfs := &fs{
		rootDir:             rootDir,
		watcher:             w,
		childrenWatcherMaps: make(map[string]*fsChildrenWatcher),
		tempNodes:           make(map[string]bool),
		closeCh:             make(chan struct{}),
	}
	return registryfs, nil
}

//Start 启动文件监控
func (l *fs) Start() {
	go func() {
		for {
			select {
			case <-l.closeCh:
				l.watcher.Close()
				return
			case event := <-l.watcher.Events:
				if l.done {
					return
				}
				func(event fsnotify.Event) {
					l.watchLock.Lock()
					defer l.watchLock.Unlock()
					dataPath := l.formatPath(event.Name)
					path := filepath.Dir(dataPath)
					l.bubblingChildrenEvent(path, event)
				}(event)

			}
		}
	}()
}

//bubblingChildrenEvent 冒泡父节点的监控事件
func (l *fs) bubblingChildrenEvent(path string, event fsnotify.Event) {
	for len(path) > 1 {
		childrenWatcher, ok := l.childrenWatcherMaps[path]
		if ok {
			childrenWatcher.event <- event
			return
		}
		path = filepath.Dir(path)
	}
}

func (l *fs) replaceColon(path string) string {
	return strings.ReplaceAll(path, ":", "@@@")
}

func (l *fs) restoreColon(path string) string {
	return strings.ReplaceAll(path, "@@@", ":")
}

//formatPath 将rootDir 构建到路径中去
func (l *fs) formatPath(path string) string {
	if !strings.HasPrefix(path, l.rootDir) {
		return l.rootDir + filepath.Join("/", path)
	}
	return path
}

func (l *fs) GetChildren(path string) (paths []string, version int32, err error) {
	rpath := l.replaceColon(l.formatPath(path))
	fs, err := os.Stat(rpath)
	if os.IsNotExist(err) {
		return nil, 0, errors.New(path + "不存在")
	}
	version = int32(fs.ModTime().Unix())
	children, err := os.ReadDir(rpath)
	if err != nil {
		return nil, 0, err
	}
	paths = make([]string, 0, len(children))
	for _, f := range children {
		if strings.HasSuffix(f.Name(), ".swp") || strings.HasPrefix(f.Name(), "~") || strings.HasPrefix(f.Name(), ".init") {
			continue
		}
		paths = append(paths, l.restoreColon(f.Name()))
	}
	return paths, version, nil
}

func (l *fs) WatchChildren(path string) (data chan ChildrenWatcher, err error) {

	realPath := l.replaceColon(l.formatPath(path))
	_, err = os.Stat(realPath)
	if os.IsNotExist(err) {
		err = fmt.Errorf("Watch path:%s 不存在", path)
		return
	}

	l.watchLock.Lock()
	defer l.watchLock.Unlock()
	v, ok := l.childrenWatcherMaps[realPath]
	if ok {
		return v.watcher, nil
	}
	l.childrenWatcherMaps[realPath] = &fsChildrenWatcher{
		event:    make(chan fsnotify.Event),
		watcher:  make(chan ChildrenWatcher),
		syncChan: make(chan fsnotify.Event, 100),
	}

	go func(rpath string, v *fsChildrenWatcher) {
		rpath = l.formatPath(rpath)
		if err := l.watcher.Add(rpath); err != nil {
			v.watcher <- &valuesEntity{path: rpath, Err: err}
		}
		go func(evtw *fsChildrenWatcher) {
			ticker := time.NewTicker(time.Second * 2)
			for {
				select {
				case <-ticker.C:
					path := ""
					var op fsnotify.Op
				INFOR:
					for {
						select {
						case p := <-evtw.syncChan:
							path = p.Name
							op = p.Op
						default:
							break INFOR
						}
					}
					if len(path) > 0 {
						vals, version, err := l.GetChildren(rpath)
						ett := &valuesEntity{
							path:    path,
							OP:      op,
							values:  vals,
							version: version,
							Err:     err,
						}
						evtw.watcher <- ett
					}
				}
			}
		}(v)

		for {
			select {
			case <-l.closeCh:
				return
			case event := <-v.event:
				if event.Op == fsnotify.Chmod {
					break
				}
				v.syncChan <- event
			}
		}
	}(realPath, l.childrenWatcherMaps[realPath])

	return l.childrenWatcherMaps[realPath].watcher, nil
}

func (l *fs) Close() error {
	l.tempNodeLock.Lock()
	defer l.tempNodeLock.Unlock()
	if l.done {
		return nil
	}
	l.done = true
	close(l.closeCh)
	for path := range l.tempNodes {
		os.RemoveAll(path)
	}
	return nil
}

type valueEntity struct {
	Value   []byte
	version int32
	path    string
	Err     error
}
type valuesEntity struct {
	values  []string
	version int32
	path    string
	Err     error
	OP      fsnotify.Op
}

func (v *valueEntity) GetPath() string {
	return v.path
}
func (v *valueEntity) GetValue() ([]byte, int32) {
	return v.Value, v.version
}
func (v *valueEntity) GetError() error {
	return v.Err
}

func (v *valuesEntity) GetValue() ([]string, int32) {
	return v.values, v.version
}
func (v *valuesEntity) GetError() error {
	return v.Err
}
func (v *valuesEntity) GetPath() string {
	return v.path
}

func (v *valuesEntity) GetOp() fsnotify.Op {
	return v.OP
}

type ChildrenWatcher interface {
	GetValue() ([]string, int32)
	GetPath() string
	GetError() error
	GetOp() fsnotify.Op
}
