package game

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jevans40/Ruthenium/render"
	"github.com/jevans40/Ruthenium/world"
	log "github.com/sirupsen/logrus"
)

const Version = "0.1.0"

var _ Game = (*gameECS)(nil)

/*
Main game struct, manages all world layers.

*/
type Game interface {
	//Initialize the game and any rendering related functions.
	Init() error

	//Starts the game and render loop.
	//All world and logic layers should be added before this function is called.
	Start()

	//Add a world layer to the layerspace.
	AddWorld(*world.WorldHandler) (worldID int, err error)

	//Remove a layer from the layerspace.
	RemoveWorld(worldID int) error

	//Returns the game render channel
	GetRenderChannel() chan []float32
}

type gameECS struct {
	window     *render.GoWindow
	worlds     []*world.WorldHandler
	renderchan chan []float32
}

func NewGameECS() Game {
	return &gameECS{}
}

func (g *gameECS) Init() error {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)
	g.renderchan = make(chan []float32, runtime.NumCPU())
	return nil
}

func (g *gameECS) AddWorld(world *world.WorldHandler) (worldID int, err error) {
	g.worlds = append(g.worlds, world)
	return len(g.worlds), nil
}

//TODO: This is all sorts of wrong, it should pop worlds associated with a given ID,
//But here it pops the world at position ID.
//Fix later, not important until production but for now its kinda important.
func (g *gameECS) RemoveWorld(worldID int) error {
	g.worlds = append(g.worlds[:worldID], g.worlds[worldID+1:]...)
	return nil
}

func (g *gameECS) Start() {
	runtime.LockOSThread()

	//Initialize glfw
	err := glfw.Init()
	if err != nil {
		log.Panic(err)
	}
	defer glfw.Terminate()

	window, err := render.NewWindow(1920, 1080)
	if err != nil {
		log.Panic(err)
	}
	g.window = window
	go g.render()

	g.EventLoop()
}

//Returns the game render channel
func (g *gameECS) GetRenderChannel() chan []float32 {
	return g.renderchan
}

//NOTE: possibly move this to the systems category.

func (g *gameECS) render() {
	//Old vs new buffer

	//THIS HAS TO BE IN THE SAME THREAD AS THE OTHER RENDERING
	//Initialize open-gl
	runtime.LockOSThread()
	glfw.DetachCurrentContext()
	(g.window.GetWindow()).MakeContextCurrent()
	if err := gl.Init(); err != nil {
		log.Panic(err)
	}

	fmt.Print("Starting")

	//Generate the window for the game

	//Log Debug Info to the console
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Info("Starting the Rutherium game engine!")
	log.WithFields(log.Fields{"OpenGL version:": version}).Info()

	var nrAttributes int32
	gl.GetIntegerv(gl.MAX_VERTEX_ATTRIBS, &nrAttributes)
	log.WithFields(log.Fields{"Max Attributes Supported:": nrAttributes}).Info()

	//Get maximum texture size
	var MaximumTextureSize int32
	gl.GetIntegerv(gl.MAX_TEXTURE_SIZE, &MaximumTextureSize)
	log.WithFields(log.Fields{"Maximum Texture Size:": MaximumTextureSize}).Info()

	var MaximumTextureUnits int32
	gl.GetIntegerv(gl.MAX_TEXTURE_IMAGE_UNITS, &MaximumTextureUnits)
	log.WithFields(log.Fields{"Maximum Texture Units": MaximumTextureUnits}).Info()

	//Log game version
	log.WithFields(log.Fields{"Psychic Spork Version": Version}).Info()
	var Buffer []float32
	numObjects := 0
	var sprites []*render.VertexRenderable
	renderer := render.SpriteRendererFactory()
	for {
		//Create the renderer
		Buffer = <-g.renderchan

		numObjects = len(Buffer) / 28
		for len(sprites) < numObjects {
			sprites = append(sprites, render.VertexSpriteFactory(&renderer))
		}
		for len(sprites) > numObjects {
			sprites = sprites[0 : len(sprites)-2]
		}
		for i, v := range sprites {
			v.SetVerticies(Buffer[i*28 : (i+1)*28])
		}
		x, y := g.window.GetSize()
		renderer.Render(int32(x), int32(y))
		g.window.GetWindow().SwapBuffers()
	}
}

func (g *gameECS) EventLoop() {
	for {
		glfw.WaitEvents()
	}
}

func (g *gameECS) Update() {

}
