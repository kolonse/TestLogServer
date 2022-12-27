package main

import "path"

// fileManager : file manager
type fileManager struct {
	files map[string]IFileWriter
}

func (f *fileManager) write(pck *writeCommand) {
	wrapper, ok := f.files[pck.GetName()]
	if !ok {
		wrapper = newFileWrapper()
		f.files[pck.GetName()] = wrapper
	}

	if len(pck.GetData()) == 0 {
		wrapper.close()
		delete(f.files, pck.GetName())
	} else {
		wrapper.write(pck)
	}
}

// Close :
func (f *fileManager) close() {
	for k, wrapper := range f.files {
		wrapper.close()
		delete(f.files, k)
	}
}

// IFileSession : file session
type IFileSession interface {
	SetDir(...string)
	Write(IPck)
	Close()
	GetID() int32

	write(*writeCommand)
}

// FileSession : file session
type FileSession struct {
	server  IServer
	baseDir string
	id      int32

	fileMgr fileManager
}

// GetID :
func (s *FileSession) GetID() int32 {
	return s.id
}

// SetDir ： set dir
func (s *FileSession) SetDir(dirs ...string) {
	s.baseDir = path.Join(dirs...)
}

// Write : write log
func (s *FileSession) Write(pck IPck) {
	filePck := &writeCommand{}
	filePck.id = s.id
	filePck.name = path.Join(s.baseDir, pck.GetName())
	filePck.data = pck.GetData()

	s.server.addPck(filePck)
}

// Write : write log
func (s *FileSession) write(pck *writeCommand) {
	s.fileMgr.write(pck)
}

// Close : close file session
func (s *FileSession) Close() {
	s.server = nil
	s.fileMgr.close()
}
