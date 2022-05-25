package render

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

//TODO:: Tests
//TODO:: Refactor

// GoWindow the struct that holds glfw window information
type GoWindow struct {
	lWindow *glfw.Window
	lSize   [2]int
}

// NewWindow Initalizes the information into a goWindow
// x int : Size in the x dimension
// y int : Size in the y dimension
func NewWindow(x int, y int) (*GoWindow, error) {
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	win, err := glfw.CreateWindow(x, y, "Testing", nil, nil)
	if err != nil {
		return nil, err
	}
	thisWindow := GoWindow{}
	thisWindow.lWindow = win
	thisWindow.lSize = [2]int{x, y}
	thisWindow.lWindow.SetSizeCallback(thisWindow.SizeCallback)
	return &thisWindow, nil
}

//GetWindow returns the lWindow feild associated with GoWindow
func (thisWindow *GoWindow) GetWindow() *glfw.Window {
	return thisWindow.lWindow
}

func (thisWindow *GoWindow) SizeCallback(w *glfw.Window, width int, height int) {
	thisWindow.lSize[0] = width
	thisWindow.lSize[1] = height
}

//GetSize returns the lSize feild associated with GoWindow
func (thisWindow *GoWindow) GetSize() (x int, y int) {
	x = thisWindow.lSize[0]
	y = thisWindow.lSize[1]
	//event.NotifyWindowResizeListeners(x, y)
	return
}
