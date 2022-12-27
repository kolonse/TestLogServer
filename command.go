package main

const (
	commandWriteData      int = 0
	commandCreateSession  int = 1
	commandDestorySession int = 2
)

// ICommand :
type ICommand interface {
	GetType() int
}

type createSessionCommand struct {
	sess IFileSession
}

func (p *createSessionCommand) GetType() int {
	return commandCreateSession
}

type destorySessionCommand struct {
	sess IFileSession
}

func (p *destorySessionCommand) GetType() int {
	return commandDestorySession
}

type writeCommand struct {
	id   int32
	name string
	data []byte
}

func (p *writeCommand) GetType() int {
	return commandWriteData
}

func (p *writeCommand) GetID() int32 {
	return p.id
}

func (p *writeCommand) GetName() string {
	return p.name
}
func (p *writeCommand) GetData() []byte {
	return p.data
}
