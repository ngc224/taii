package tail

import (
	"os"

	"github.com/fsnotify/fsnotify"
)

type Tail struct {
	filename        string
	file            *os.File
	fileInfoBeforef os.FileInfo
	Watcher         *fsnotify.Watcher
}

func NewTail(filename string) (*Tail, error) {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		return nil, err
	}

	if err := watcher.Add(filename); err != nil {
		return nil, err
	}

	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	t := &Tail{
		filename: filename,
		file:     file,
		Watcher:  watcher,
	}

	if err := t.reset(); err != nil {
		return nil, err
	}

	return t, nil
}

func (t *Tail) reset() (err error) {
	if _, err = t.file.Seek(0, os.SEEK_END); err != nil {
		return err
	}

	if t.fileInfoBeforef, err = t.file.Stat(); err != nil {
		return err
	}

	return nil
}

func (t *Tail) Close() {
	t.Watcher.Close()
	t.file.Close()
}

func (t *Tail) AddString() string {
	stat, err := t.file.Stat()

	if err != nil {
		return ""
	}

	readSize := stat.Size() - t.fileInfoBeforef.Size()
	t.fileInfoBeforef = stat

	if 0 > readSize {
		t.reset()
		return ""
	}

	buf := make([]byte, readSize)

	if _, err := t.file.Read(buf); err != nil {
		return ""
	}

	return string(buf)
}
