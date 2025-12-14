package watcher

import (
	"log"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	watcher  *fsnotify.Watcher
	onChange func()
}

func New(dirs ...interface{}) (*Watcher, error) {
	// Extract callback (last arg) and dirs
	onChange := dirs[len(dirs)-1].(func())
	watchDirs := make([]string, len(dirs)-1)
	for i := 0; i < len(dirs)-1; i++ {
		watchDirs[i] = dirs[i].(string)
	}

	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	w := &Watcher{
		watcher:  fsw,
		onChange: onChange,
	}

	for _, dir := range watchDirs {
		if _, err := os.Stat(dir); err == nil {
			if err := fsw.Add(dir); err != nil {
				log.Printf("Warning: failed to watch %s: %v\n", dir, err)
			}
		}
	}

	go w.watch()
	return w, nil
}

func (w *Watcher) watch() {
	var debounce <-chan time.Time

	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove) != 0 {
				debounce = time.After(100 * time.Millisecond)
			}
		case <-debounce:
			w.onChange()
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v\n", err)
		}
	}
}

func (w *Watcher) Close() error {
	return w.watcher.Close()
}
