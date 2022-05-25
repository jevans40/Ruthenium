package world

import (
	"reflect"

	"github.com/jevans40/Ruthenium/component"
	"github.com/jevans40/Ruthenium/render"
)

//TODO:: Documentation
//TODO:: Tests

type World interface {
	Resume()
	Pause()
	Init()
	GetName() string
	RegisterService(s Service)
	RegisterStorage(s component.ComponentStorage)
	Maintain() error
}

type BaseWorld struct {
	dispatcher Dispatcher
}

type WindowComponent struct {
	window *render.GoWindow
}

func (w WindowComponent) IsComponent()          {}
func (w WindowComponent) GetType() reflect.Type { return reflect.TypeOf(w) }
func (w WindowComponent) GetSize() (x, y int)   { return w.window.GetSize() }

func NewBaseWorld(renderChannel chan []float32, window *render.GoWindow) World {
	dispatcher := NewSimpleDispatcher()
	newWorld := BaseWorld{dispatcher: dispatcher}

	//Required Services
	renderService := NewRenderService(renderChannel)

	//Required Storages
	RenderableStorage := component.NewVectorStorage[Renderable]()

	//TODO:: Possibly Make a read only resource type for resources like this
	WindowResource := component.NewResourceStorage(WindowComponent{window: window})

	//Register Services and Storages to dispatcher
	newWorld.dispatcher.AddService(renderService)

	newWorld.dispatcher.AddStorage(RenderableStorage)
	newWorld.dispatcher.AddStorage(WindowResource)
	return &newWorld
}

func (b *BaseWorld) Resume() {
	//Other steps that need to be taken put here
	//TODO:: This
	//b.dispatcher.StartServices()
}

func (b *BaseWorld) Pause() {
	//Other steps that need to be taken put here
	//TODO:: This
	//b.dispatcher.StopServices()
}

func (b *BaseWorld) Init() {
	//Do any initialization stuff here that needs to be done once the game starts
}

func (b *BaseWorld) GetName() string {
	return "Base World"
}

func (b *BaseWorld) RegisterService(s Service) {
	b.dispatcher.AddService(s)
}

func (b *BaseWorld) RegisterStorage(s component.ComponentStorage) {
	b.dispatcher.AddStorage(s)
}

func (b *BaseWorld) Maintain() error {
	return b.dispatcher.Maintain()
}
