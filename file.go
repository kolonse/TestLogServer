package main

import (
	"log"
	"sync/atomic"
)

// IServer : server interface
type IServer interface {
	addPck(pck ICommand)
}

// IPck : package interface
type IPck interface {
	GetName() string
	GetData() []byte
}

// FileServer : manager of saving files
type FileServer struct {
	runningFlag     bool
	sessions        map[int32]IFileSession
	keyIndex        int32
	chanPck         chan ICommand
	chanRunning     chan bool
	chanWaitingExit chan bool
}

func (s *FileServer) run() {
	running := true
	for {
		if !running {
			break
		}

		select {
		case command := <-s.chanPck:
			// s.writer.Write(pck)
			switch command.GetType() {
			case commandWriteData:
				pck := command.(*writeCommand)
				if sess, ok := s.sessions[pck.GetID()]; ok {
					sess.write(pck)
				}
			case commandCreateSession:
				pck := command.(*createSessionCommand)
				s.sessions[pck.sess.GetID()] = pck.sess
			case commandDestorySession:
				pck := command.(*destorySessionCommand)
				if _, ok := s.sessions[pck.sess.GetID()]; ok {
					pck.sess.Close()
					delete(s.sessions, pck.sess.GetID())
				}
			}
		case running = <-s.chanRunning:
		}
	}

	s.chanWaitingExit <- true
}

// Start : start file server
func (s *FileServer) Start() {
	if s.runningFlag {
		return
	}

	go s.run()

	s.runningFlag = true
	log.Println("File Server Start")
}

// Stop : stop file server
func (s *FileServer) Stop() {
	s.chanRunning <- false
	log.Println("Log Server Waitting Exit")
	<-s.chanWaitingExit
	s.runningFlag = false
	log.Println("Log Server Exit Complete")
}

// AddPck : add a log to write
func (s *FileServer) addPck(pck ICommand) {
	s.chanPck <- pck
}

// CreateFileSession : create file session
func (s *FileServer) CreateFileSession() IFileSession {
	fileMgr := fileManager{
		files: make(map[string]IFileWriter),
	}

	sess := &FileSession{
		server:  s,
		id:      s.generateKey(),
		fileMgr: fileMgr,
	}

	pck := &createSessionCommand{
		sess: sess,
	}
	s.chanPck <- pck
	log.Printf("Create File Session : %v\n", sess.GetID())
	return sess
}

// DestoryFileSession : destory file session
func (s *FileServer) DestoryFileSession(sess IFileSession) {
	pck := &destorySessionCommand{
		sess: sess,
	}
	s.chanPck <- pck
	log.Printf("Destory File Session : %v\n", sess.GetID())
}

func (s *FileServer) generateKey() int32 {
	return atomic.AddInt32(&s.keyIndex, 1)
}

// NewFileServer : create file server
func NewFileServer() *FileServer {
	return &FileServer{
		runningFlag:     false,
		keyIndex:        0,
		sessions:        make(map[int32]IFileSession),
		chanPck:         make(chan ICommand, 1),
		chanRunning:     make(chan bool, 1),
		chanWaitingExit: make(chan bool, 1),
	}
}
