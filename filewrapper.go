package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

const (
	fileMaxSize int = 10 * 1024 * 1024
)

// IFileWriter : writer interface
type IFileWriter interface {
	write(*writeCommand)
	close()
}

type fileWrapper struct {
	key       string
	path      string
	handle    *os.File
	index     int
	writeSize int
}

func (w *fileWrapper) isExistPath(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func (w *fileWrapper) isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func (w *fileWrapper) formatPath() {
	w.index++
	w.path = fmt.Sprintf("%v.%v", w.path, w.index)
}

func (w *fileWrapper) createFile(id int32) bool {
	// check file exist
	var err error
	for {
		if w.isExistPath(w.path) {
			w.formatPath()
		} else {
			break
		}
	}

	dir := path.Dir(w.path)
	if w.isDir(dir) == false {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			log.Println("Create Dir Fail : ", err.Error())
			return false
		}
	}
	// create file handle
	w.handle, err = os.Create(w.path)
	if err != nil {
		log.Println("Create File Fail : ", err.Error())
		return false
	}
	log.Printf("Create File Success : Session Id [ %v ] Path [ %v ]\n", id, w.path)
	return true
}

func (w *fileWrapper) getNowStr() string {
	return "[ go write time: " + time.Now().Format("2006-01-02 15:04:05.000") + " ]"

}

func (w *fileWrapper) writeText(pck *writeCommand) int {
	buf := pck.GetData()
	if len(buf) == 0 {
		return 0
	}
	if buf[len(buf)-1] == 0 {
		buf = buf[0 : len(buf)-1]
	}
	str := string(buf)
	if len(str) == 0 {
		return 0
	}

	if str[len(str)-1] != '\n' {
		str += "\n"
	}

	w.handle.WriteString(w.getNowStr() + str)
	return len(str)
}

func (w *fileWrapper) write(pck *writeCommand) {
	w.key = pck.GetName()
	w.path = pck.GetName()
	var err error
	var writeSize int
	if w.handle == nil && !w.createFile(pck.GetID()) {
		return
	}
	if strings.HasSuffix(w.key, ".log") {
		writeSize = w.writeText(pck)
	} else {
		writeSize, err = w.handle.Write(pck.GetData())
		if err != nil {
			log.Println("Write File Fail : ", w.path, err.Error())
			return
		}
	}

	w.writeSize += writeSize
	if w.writeSize > fileMaxSize {
		w.close()
		w.formatPath()
		w.createFile(pck.GetID())
		w.writeSize = 0
	}
}

func (w *fileWrapper) close() {
	if w.handle != nil {
		w.handle.Close()
		w.handle = nil
	}

	log.Println("Close File : ", w.path)
}

func newFileWrapper() *fileWrapper {
	return &fileWrapper{
		index:     0,
		writeSize: 0,
	}
}
