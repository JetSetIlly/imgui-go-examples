//go:build glfw
// +build glfw

package platforms

import (
	"fmt"
	"math"
	"runtime"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/inkyblackness/imgui-go/v4"
)

// GLFWClientAPI identifies the render system that shall be initialized.
type GLFWClientAPI string

// This is a list of GLFWClientAPI constants.
const (
	GLFWClientAPIOpenGL2 GLFWClientAPI = "OpenGL2"
	GLFWClientAPIOpenGL3 GLFWClientAPI = "OpenGL3"
)

// GLFW implements a platform based on github.com/go-gl/glfw (v3.2).
type GLFW struct {
	imguiIO imgui.IO

	window *glfw.Window

	time             float64
	mouseJustPressed [3]bool
}

// NewGLFW attempts to initialize a GLFW context.
func NewGLFW(io imgui.IO, clientAPI GLFWClientAPI) (*GLFW, error) {
	runtime.LockOSThread()

	err := glfw.Init()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize glfw: %w", err)
	}

	switch clientAPI {
	case GLFWClientAPIOpenGL2:
		glfw.WindowHint(glfw.ContextVersionMajor, 2)
		glfw.WindowHint(glfw.ContextVersionMinor, 1)
	case GLFWClientAPIOpenGL3:
		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 2)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, 1)
	default:
		glfw.Terminate()
		return nil, ErrUnsupportedClientAPI
	}

	window, err := glfw.CreateWindow(windowWidth, windowHeight, "ImGui-Go GLFW+"+string(clientAPI)+" example", nil, nil)
	if err != nil {
		glfw.Terminate()
		return nil, fmt.Errorf("failed to create window: %w", err)
	}
	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	platform := &GLFW{
		imguiIO: io,
		window:  window,
	}
	platform.installCallbacks()

	return platform, nil
}

// Dispose cleans up the resources.
func (platform *GLFW) Dispose() {
	platform.window.Destroy()
	glfw.Terminate()
}

// ShouldStop returns true if the window is to be closed.
func (platform *GLFW) ShouldStop() bool {
	return platform.window.ShouldClose()
}

// ProcessEvents handles all pending window events.
func (platform *GLFW) ProcessEvents() {
	glfw.PollEvents()
}

// DisplaySize returns the dimension of the display.
func (platform *GLFW) DisplaySize() [2]float32 {
	w, h := platform.window.GetSize()
	return [2]float32{float32(w), float32(h)}
}

// FramebufferSize returns the dimension of the framebuffer.
func (platform *GLFW) FramebufferSize() [2]float32 {
	w, h := platform.window.GetFramebufferSize()
	return [2]float32{float32(w), float32(h)}
}

// NewFrame marks the begin of a render pass. It forwards all current state to imgui IO.
func (platform *GLFW) NewFrame() {
	// Setup display size (every frame to accommodate for window resizing)
	displaySize := platform.DisplaySize()
	platform.imguiIO.SetDisplaySize(imgui.Vec2{X: displaySize[0], Y: displaySize[1]})

	// Setup time step
	currentTime := glfw.GetTime()
	if platform.time > 0 {
		platform.imguiIO.SetDeltaTime(float32(currentTime - platform.time))
	}
	platform.time = currentTime

	// Setup inputs
	if platform.window.GetAttrib(glfw.Focused) != 0 {
		x, y := platform.window.GetCursorPos()
		platform.imguiIO.SetMousePosition(imgui.Vec2{X: float32(x), Y: float32(y)})
	} else {
		platform.imguiIO.SetMousePosition(imgui.Vec2{X: -math.MaxFloat32, Y: -math.MaxFloat32})
	}

	for i := 0; i < len(platform.mouseJustPressed); i++ {
		down := platform.mouseJustPressed[i] || (platform.window.GetMouseButton(glfwButtonIDByIndex[i]) == glfw.Press)
		platform.imguiIO.SetMouseButtonDown(i, down)
		platform.mouseJustPressed[i] = false
	}
}

// PostRender performs a buffer swap.
func (platform *GLFW) PostRender() {
	platform.window.SwapBuffers()
}

func (platform *GLFW) installCallbacks() {
	platform.window.SetMouseButtonCallback(platform.mouseButtonChange)
	platform.window.SetScrollCallback(platform.mouseScrollChange)
	platform.window.SetKeyCallback(platform.keyChange)
	platform.window.SetCharCallback(platform.charChange)
}

var glfwButtonIndexByID = map[glfw.MouseButton]int{
	glfw.MouseButton1: mouseButtonPrimary,
	glfw.MouseButton2: mouseButtonSecondary,
	glfw.MouseButton3: mouseButtonTertiary,
}

var glfwButtonIDByIndex = map[int]glfw.MouseButton{
	mouseButtonPrimary:   glfw.MouseButton1,
	mouseButtonSecondary: glfw.MouseButton2,
	mouseButtonTertiary:  glfw.MouseButton3,
}

func (platform *GLFW) mouseButtonChange(window *glfw.Window, rawButton glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	buttonIndex, known := glfwButtonIndexByID[rawButton]

	if known && (action == glfw.Press) {
		platform.mouseJustPressed[buttonIndex] = true
	}
}

func (platform *GLFW) mouseScrollChange(window *glfw.Window, x, y float64) {
	platform.imguiIO.AddMouseWheelDelta(float32(x), float32(y))
}

func (platform *GLFW) keyChange(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	k := glfwKeyEventToImguiKey(key, scancode)
	if action == glfw.Press {
		platform.imguiIO.AddKeyEvent(k, true)
	}
	if action == glfw.Release {
		platform.imguiIO.AddKeyEvent(k, false)
	}
	glfwSetImguiModKey(platform.imguiIO, mods)
}

func glfwSetImguiModKey(io imgui.IO, mod glfw.ModifierKey) {
	io.AddKeyEvent(imgui.KeyModCtrl, (mod&glfw.ModControl) != 0)
	io.AddKeyEvent(imgui.KeyModShift, (mod&glfw.ModShift) != 0)
	io.AddKeyEvent(imgui.KeyModAlt, (mod&glfw.ModAlt) != 0)
	io.AddKeyEvent(imgui.KeyModSuper, (mod&glfw.ModSuper) != 0)
}

func glfwKeyEventToImguiKey(key glfw.Key, scancode int) imgui.ImguiKey {
	switch key {
	case glfw.KeyTab:
		return imgui.KeyTab
	case glfw.KeyLeft:
		return imgui.KeyLeftArrow
	case glfw.KeyRight:
		return imgui.KeyRightArrow
	case glfw.KeyUp:
		return imgui.KeyUpArrow
	case glfw.KeyDown:
		return imgui.KeyDownArrow
	case glfw.KeyPageUp:
		return imgui.KeyPageUp
	case glfw.KeyPageDown:
		return imgui.KeyPageDown
	case glfw.KeyHome:
		return imgui.KeyHome
	case glfw.KeyEnd:
		return imgui.KeyEnd
	case glfw.KeyInsert:
		return imgui.KeyInsert
	case glfw.KeyDelete:
		return imgui.KeyDelete
	case glfw.KeyBackspace:
		return imgui.KeyBackspace
	case glfw.KeySpace:
		return imgui.KeySpace
	case glfw.KeyEnter:
		return imgui.KeyEnter
	case glfw.KeyEscape:
		return imgui.KeyEscape
	case glfw.KeyApostrophe:
		return imgui.KeyApostrophe
	case glfw.KeyComma:
		return imgui.KeyComma
	case glfw.KeyMinus:
		return imgui.KeyMinus
	case glfw.KeyPeriod:
		return imgui.KeyPeriod
	case glfw.KeySlash:
		return imgui.KeySlash
	case glfw.KeySemicolon:
		return imgui.KeySemicolon
	case glfw.KeyEqual:
		return imgui.KeyEqual
	case glfw.KeyLeftBracket:
		return imgui.KeyLeftBracket
	case glfw.KeyBackslash:
		return imgui.KeyBackslash
	case glfw.KeyWorld1:
		return imgui.KeyOem102
	case glfw.KeyWorld2:
		return imgui.KeyOem102
	case glfw.KeyRightBracket:
		return imgui.KeyRightBracket
	case glfw.KeyGraveAccent:
		return imgui.KeyGraveAccent
	case glfw.KeyCapsLock:
		return imgui.KeyCapsLock
	case glfw.KeyScrollLock:
		return imgui.KeyScrollLock
	case glfw.KeyNumLock:
		return imgui.KeyNumLock
	case glfw.KeyPrintScreen:
		return imgui.KeyPrintScreen
	case glfw.KeyPause:
		return imgui.KeyPause
	case glfw.KeyKP0:
		return imgui.KeyKeypad0
	case glfw.KeyKP1:
		return imgui.KeyKeypad1
	case glfw.KeyKP2:
		return imgui.KeyKeypad2
	case glfw.KeyKP3:
		return imgui.KeyKeypad3
	case glfw.KeyKP4:
		return imgui.KeyKeypad4
	case glfw.KeyKP5:
		return imgui.KeyKeypad5
	case glfw.KeyKP6:
		return imgui.KeyKeypad6
	case glfw.KeyKP7:
		return imgui.KeyKeypad7
	case glfw.KeyKP8:
		return imgui.KeyKeypad8
	case glfw.KeyKP9:
		return imgui.KeyKeypad9
	case glfw.KeyKPDecimal:
		return imgui.KeyKeypadDecimal
	case glfw.KeyKPDivide:
		return imgui.KeyKeypadDivide
	case glfw.KeyKPMultiply:
		return imgui.KeyKeypadMultiply
	case glfw.KeyKPSubtract:
		return imgui.KeyKeypadSubtract
	case glfw.KeyKPAdd:
		return imgui.KeyKeypadAdd
	case glfw.KeyKPEnter:
		return imgui.KeyKeypadEnter
	case glfw.KeyKPEqual:
		return imgui.KeyKeypadEqual
	case glfw.KeyLeftShift:
		return imgui.KeyLeftShift
	case glfw.KeyLeftControl:
		return imgui.KeyLeftCtrl
	case glfw.KeyLeftAlt:
		return imgui.KeyLeftAlt
	case glfw.KeyLeftSuper:
		return imgui.KeyLeftSuper
	case glfw.KeyRightShift:
		return imgui.KeyRightShift
	case glfw.KeyRightControl:
		return imgui.KeyRightCtrl
	case glfw.KeyRightAlt:
		return imgui.KeyRightAlt
	case glfw.KeyRightSuper:
		return imgui.KeyRightSuper
	case glfw.KeyMenu:
		return imgui.KeyMenu
	case glfw.Key0:
		return imgui.Key0
	case glfw.Key1:
		return imgui.Key1
	case glfw.Key2:
		return imgui.Key2
	case glfw.Key3:
		return imgui.Key3
	case glfw.Key4:
		return imgui.Key4
	case glfw.Key5:
		return imgui.Key5
	case glfw.Key6:
		return imgui.Key6
	case glfw.Key7:
		return imgui.Key7
	case glfw.Key8:
		return imgui.Key8
	case glfw.Key9:
		return imgui.Key9
	case glfw.KeyA:
		return imgui.KeyA
	case glfw.KeyB:
		return imgui.KeyB
	case glfw.KeyC:
		return imgui.KeyC
	case glfw.KeyD:
		return imgui.KeyD
	case glfw.KeyE:
		return imgui.KeyE
	case glfw.KeyF:
		return imgui.KeyF
	case glfw.KeyG:
		return imgui.KeyG
	case glfw.KeyH:
		return imgui.KeyH
	case glfw.KeyI:
		return imgui.KeyI
	case glfw.KeyJ:
		return imgui.KeyJ
	case glfw.KeyK:
		return imgui.KeyK
	case glfw.KeyL:
		return imgui.KeyL
	case glfw.KeyM:
		return imgui.KeyM
	case glfw.KeyN:
		return imgui.KeyN
	case glfw.KeyO:
		return imgui.KeyO
	case glfw.KeyP:
		return imgui.KeyP
	case glfw.KeyQ:
		return imgui.KeyQ
	case glfw.KeyR:
		return imgui.KeyR
	case glfw.KeyS:
		return imgui.KeyS
	case glfw.KeyT:
		return imgui.KeyT
	case glfw.KeyU:
		return imgui.KeyU
	case glfw.KeyV:
		return imgui.KeyV
	case glfw.KeyW:
		return imgui.KeyW
	case glfw.KeyX:
		return imgui.KeyX
	case glfw.KeyY:
		return imgui.KeyY
	case glfw.KeyZ:
		return imgui.KeyZ
	case glfw.KeyF1:
		return imgui.KeyF1
	case glfw.KeyF2:
		return imgui.KeyF2
	case glfw.KeyF3:
		return imgui.KeyF3
	case glfw.KeyF4:
		return imgui.KeyF4
	case glfw.KeyF5:
		return imgui.KeyF5
	case glfw.KeyF6:
		return imgui.KeyF6
	case glfw.KeyF7:
		return imgui.KeyF7
	case glfw.KeyF8:
		return imgui.KeyF8
	case glfw.KeyF9:
		return imgui.KeyF9
	case glfw.KeyF10:
		return imgui.KeyF10
	case glfw.KeyF11:
		return imgui.KeyF11
	case glfw.KeyF12:
		return imgui.KeyF12
	case glfw.KeyF13:
		return imgui.KeyF13
	case glfw.KeyF14:
		return imgui.KeyF14
	case glfw.KeyF15:
		return imgui.KeyF15
	case glfw.KeyF16:
		return imgui.KeyF16
	case glfw.KeyF17:
		return imgui.KeyF17
	case glfw.KeyF18:
		return imgui.KeyF18
	case glfw.KeyF19:
		return imgui.KeyF19
	case glfw.KeyF20:
		return imgui.KeyF20
	case glfw.KeyF21:
		return imgui.KeyF21
	case glfw.KeyF22:
		return imgui.KeyF22
	case glfw.KeyF23:
		return imgui.KeyF23
	case glfw.KeyF24:
		return imgui.KeyF24
	default:
		return imgui.KeyNone
	}
}

func (platform *GLFW) charChange(window *glfw.Window, char rune) {
	platform.imguiIO.AddInputCharacters(string(char))
}

// ClipboardText returns the current clipboard text, if available.
func (platform *GLFW) ClipboardText() (string, error) {
	return platform.window.GetClipboardString()
}

// SetClipboardText sets the text as the current clipboard text.
func (platform *GLFW) SetClipboardText(text string) {
	platform.window.SetClipboardString(text)
}
