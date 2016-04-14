package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var errClosedWriter = errors.New("write to closed writer")

type periodicWriter struct {
	Prefix, Suffix string

	l sync.Mutex
	f *os.File
	y int
	m time.Month
	d int

	closed bool
}

func (w *periodicWriter) Write(p []byte) (int, error) {
	w.l.Lock()
	defer w.l.Unlock()
	if w.closed {
		return 0, errClosedWriter
	}
	if w.needNewFile() {
		if err := w.newFile(); err != nil {
			return 0, err
		}
	}
	return w.f.Write(p)
}

func (w *periodicWriter) Close() (err error) {
	w.l.Lock()
	defer w.l.Unlock()
	if w.f != nil {
		err = w.f.Close()
		w.f = nil
	}
	w.closed = true
	return err
}

func (w *periodicWriter) needNewFile() bool {
	if w.f == nil {
		return true
	}
	t := time.Now()
	y, m, d := t.Date()
	return d != w.d || m != w.m || y != w.y
}

func (w *periodicWriter) newFile() error {
	if w.f != nil {
		err := w.f.Close()
		if err != nil {
			return err
		}
		w.f = nil
		// rename a log file.
		newName := fmt.Sprintf("%s-%04d%02d%02d%s", w.Prefix, w.y, w.m, w.d,
			w.Suffix)
		err = os.Rename(w.Prefix+w.Suffix, newName)
		if err != nil {
			return err
		}
	}
	name := w.Prefix + w.Suffix
	dir := filepath.Dir(name)
	if err := os.MkdirAll(dir, 0777); err != nil {
		return err
	}
	f, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	w.f = f
	w.y, w.m, w.d = time.Now().Date()
	return nil
}
