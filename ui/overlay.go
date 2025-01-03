package ui

import (
	"image"
	"math"
	"time"

	"image/draw"

	"golang.org/x/image/font/gofont/goregular"

	"github.com/BurntSushi/freetype-go/freetype/truetype"

	"github.com/jezek/xgbutil/ewmh"
	"github.com/jezek/xgbutil/icccm"
	"github.com/jezek/xgbutil/motif"
	"github.com/jezek/xgbutil/xevent"
	"github.com/jezek/xgbutil/xgraphics"
	"github.com/jezek/xgbutil/xwindow"

	"github.com/leukipp/cortile/v2/common"
	"github.com/leukipp/cortile/v2/desktop"
	"github.com/leukipp/cortile/v2/store"

	log "github.com/sirupsen/logrus"
)

var (
	fontSize   int = 16 // Size of text font
	fontMargin int = 4  // Margin of text font
	rectMargin int = 4  // Margin of layout rectangles
)

var (
	gui map[uint]*xwindow.Window = make(map[uint]*xwindow.Window) // Overlay window
)

func ShowLayout(ws *desktop.Workspace) {
	location := store.Location{Desktop: store.Workplace.CurrentDesktop}
	if ws == nil || ws.Location.Desktop != location.Desktop || common.Config.TilingGui <= 0 {
		return
	}

	// Wait for tiling events
	time.AfterFunc(150*time.Millisecond, func() {

		// Obtain layout name
		name := ws.ActiveLayout().GetName()
		if ws.TilingDisabled() {
			name = "disabled"
		}

		// Calculate scaled desktop dimensions
		dim := dimensions(ws)
		_, _, w, h := scale(dim.X, dim.Y, dim.Width, dim.Height)

		// Create an empty canvas image
		bg := bgra("gui_background")
		cv := xgraphics.New(store.X, image.Rect(0, 0, w+rectMargin, h+fontSize+2*fontMargin+2*rectMargin))
		cv.For(func(x int, y int) xgraphics.BGRA { return bg })

		// Draw client rectangles
		drawClients(cv, ws, name)

		// Draw layout name
		drawText(cv, name, bgra("gui_text"), cv.Rect.Dx()/2, cv.Rect.Dy()-2*fontMargin-rectMargin, fontSize)

		// Show the canvas graphics
		showGraphics(cv, ws, time.Duration(common.Config.TilingGui))
	})
}

func drawClients(cv *xgraphics.Image, ws *desktop.Workspace, layout string) {
	al := ws.ActiveLayout()
	mg := al.GetManager()
	clients := ws.VisibleClients()

	// Draw default rectangle
	dim := dimensions(ws)
	if len(clients) == 0 || layout == "disabled" {

		// Calculate scaled desktop dimensions
		x, y, w, h := scale(0, 0, dim.Width, dim.Height)

		// Draw client rectangle onto canvas
		color := bgra("gui_client_slave")
		drawImage(cv, &image.Uniform{color}, color, x+rectMargin, y+rectMargin, x+w, y+h)

		return
	}

	// Draw master and slave rectangle
	for _, c := range clients {
		if c == nil {
			continue
		}

		// Calculate scaled client dimensions
		cx, cy, cw, ch := c.OuterGeometry()
		x, y, w, h := scale(cx-dim.X, cy-dim.Y, cw, ch)

		// Calculate icon size
		iconSize := math.MaxInt
		if w < iconSize {
			iconSize = w
		}
		if h < iconSize {
			iconSize = h
		}
		iconSize /= 2

		// Obtain rectangle color
		color := bgra("gui_client_slave")
		if mg.IsMaster(c) || common.IsInList(layout, []string{"maximized", "fullscreen"}) {
			color = bgra("gui_client_master")
		}

		// Draw client rectangle onto canvas
		drawImage(cv, &image.Uniform{color}, color, x+rectMargin, y+rectMargin, x+w, y+h)

		// Draw client icon onto canvas
		ico, err := xgraphics.FindIcon(store.X, c.Window.Id, iconSize, iconSize)
		if err == nil {
			drawImage(cv, ico, color, x+rectMargin/2+w/2-iconSize/2, y+rectMargin/2+h/2-iconSize/2, x+w, y+h)
		}
	}
}

func drawImage(cv *xgraphics.Image, img image.Image, color xgraphics.BGRA, x0 int, y0 int, x1 int, y1 int) {

	// Draw rectangle
	draw.Draw(cv, image.Rect(x0, y0, x1, y1), img, image.Point{}, draw.Src)

	// Blend background
	xgraphics.BlendBgColor(cv, color)
}

func drawText(cv *xgraphics.Image, txt string, color xgraphics.BGRA, x int, y int, size int) {
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Error("Parsing font failed: ", err)
		return
	}

	// Obtain maximum font size
	w, _ := xgraphics.Extents(font, float64(size), txt)
	if w > 2*(x-fontMargin-rectMargin) {
		drawText(cv, txt, color, x, y, size-1)
		return
	}

	// Draw text onto canvas
	cv.Text(x-w/2, y-size, color, float64(size), font, txt)
}

func showGraphics(img *xgraphics.Image, ws *desktop.Workspace, duration time.Duration) *xwindow.Window {
	win, err := xwindow.Generate(img.X)
	if err != nil {
		log.Error("Graphics generation failed: ", err)
		return nil
	}

	// Calculate window dimensions
	dim := dimensions(ws)
	w, h := img.Rect.Dx(), img.Rect.Dy()
	x, y := dim.X+dim.Width/2-w/2, dim.Y+dim.Height/2-h/2

	// Create the graphics window
	win.Create(img.X.RootWin(), x, y, w, h, 0)

	// Set class and name
	icccm.WmClassSet(win.X, win.Id, &icccm.WmClass{
		Instance: common.Build.Name,
		Class:    common.Build.Name,
	})
	icccm.WmNameSet(win.X, win.Id, common.Build.Name)

	// Set states for modal like behavior
	icccm.WmStateSet(win.X, win.Id, &icccm.WmState{
		State: icccm.StateNormal,
	})
	ewmh.WmStateSet(win.X, win.Id, []string{
		"_NET_WM_STATE_SKIP_TASKBAR",
		"_NET_WM_STATE_SKIP_PAGER",
		"_NET_WM_STATE_ABOVE",
		"_NET_WM_STATE_MODAL",
	})

	// Set hints for size and decorations
	icccm.WmNormalHintsSet(img.X, win.Id, &icccm.NormalHints{
		Flags:     icccm.SizeHintPPosition | icccm.SizeHintPMinSize | icccm.SizeHintPMaxSize,
		X:         x,
		Y:         y,
		MinWidth:  uint(w),
		MinHeight: uint(h),
		MaxWidth:  uint(w),
		MaxHeight: uint(h),
	})
	motif.WmHintsSet(img.X, win.Id, &motif.Hints{
		Flags:      motif.HintFunctions | motif.HintDecorations,
		Function:   motif.FunctionNone,
		Decoration: motif.DecorationNone,
	})

	// Ensure the window closes gracefully
	win.WMGracefulClose(func(w *xwindow.Window) {
		xevent.Detach(w.X, w.Id)
		xevent.Quit(w.X)
		w.Destroy()
	})

	// Paint the image and map the window
	img.XSurfaceSet(win.Id)
	img.XDraw()
	img.XPaint(win.Id)
	win.Map()

	// Move focus to active window
	store.ActiveWindowSet(store.X, &store.Windows.Active)

	// Close previous opened window
	if v, ok := gui[ws.Location.Screen]; ok {
		v.Destroy()
	}
	gui[ws.Location.Screen] = win

	// Close window after given duration
	if duration > 0 {
		time.AfterFunc(duration*time.Millisecond, win.Destroy)
	}

	return win
}

func dimensions(ws *desktop.Workspace) *common.Geometry {
	dim := store.DesktopGeometry(ws.Location.Screen)

	// Ignore desktop margins on fullscreen mode
	if ws.ActiveLayout().GetName() == "fullscreen" {
		dim = store.ScreenGeometry(ws.Location.Screen)
	}

	return dim
}

func scale(x, y, w, h int) (sx, sy, sw, sh int) {
	s := 10

	// Rescale dimensions by factor s
	sx, sy, sw, sh = x/s, y/s, w/s, h/s

	return
}
