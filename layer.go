package gocamole

import (
	"image"
	"image/draw"

	"github.com/llgcode/draw2d/draw2dimg"
)

type Layer struct {
	x, y         int
	width        int
	height       int
	op           draw.Op
	image        *image.RGBA
	gc           draw2dimg.GraphicContext
	visible      bool
	modified     bool
	modifiedRect image.Rectangle
}

func (l *Layer) updateModifiedRect(modArea image.Rectangle) {
	before := l.modifiedRect
	l.modifiedRect = l.modifiedRect.Union(modArea)
	l.modified = l.modified || !before.Eq(l.modifiedRect)
}

func (l *Layer) resetModified() {
	l.modifiedRect = image.Rectangle{}
	l.modified = false
}

func copyImage(dest draw.Image, x, y int, src image.Image, sr image.Rectangle, op draw.Op) {
	dp := image.Pt(x, y)
	r := image.Rectangle{Min: dp, Max: dp.Add(sr.Size())}
	draw.Draw(dest, r, src, sr.Min, op)
}

func (l *Layer) Draw(x, y int, src image.Image) {
	l.updateModifiedRect(image.Rect(x, y, x+src.Bounds().Max.X, y+src.Bounds().Max.Y))
	copyImage(l.image, x, y, src, src.Bounds(), l.op)
}

func (l *Layer) Resize(w int, h int) {
	original := l.image.Bounds()
	if w == l.width && h == l.height {
		return
	}
	newImage := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(newImage, l.image.Bounds(), l.image, image.Pt(0, 0), l.op)
	l.image = newImage
	l.width = w
	l.height = h
	l.updateModifiedRect(original.Union(l.image.Bounds()))
}

type layers map[int]*Layer

func newLayers() layers {
	ls := make(layers)
	ls[0] = &Layer{image: image.NewRGBA(image.Rect(0, 0, 0, 0))}
	return ls
}

func newLayer(l0 *Layer, visible bool) *Layer {
	var l *Layer
	// If it is an invisible layer (aka buffer)
	if visible {
		l = &Layer{
			width:  l0.width,
			height: l0.height,
			image:  image.NewRGBA(image.Rect(0, 0, l0.width, l0.height)),
		}
	} else {
		l = &Layer{image: image.NewRGBA(image.Rect(0, 0, 0, 0))}
	}
	l.visible = visible
	return l
}

func (ls layers) getDefault() *Layer {
	return ls[0]
}

func (ls layers) get(id int) *Layer {
	if l, ok := ls[id]; ok {
		return l
	}
	l := newLayer(ls.getDefault(), id > 0)
	ls[id] = l
	return l
}

func (ls layers) delete(id int) {
	if id == 0 {
		return
	}
	ls[0].updateModifiedRect(ls[id].image.Bounds())
	ls[id].image = nil
	ls[id] = nil
	delete(ls, id)
}
