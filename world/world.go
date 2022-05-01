package world

type World interface {
	Resume()
	Pause()
	Init()
	GetName()
	RegisterService()
	RegisterStorage()
	Maintain()
}

type BaseWorld struct {
	dispatcher Dispatcher
}
