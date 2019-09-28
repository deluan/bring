package bring

import (
	"fmt"
	"image"
	"image/draw"
	"sync"
	"time"

	"github.com/google/uuid"
)

var compositeOperations = map[byte]draw.Op{
	0xC: draw.Src,
	0xE: draw.Over,
}

type Display struct {
	logger         Logger
	cursor         *Layer
	cursorHotspotX int
	cursorHotspotY int
	cursorX        int
	cursorY        int
	tasks          chan task
	pendingTasks   sync.WaitGroup
	layers         layers
	canvas         *image.RGBA
	lastUpdate     int64
	canvasAccess   sync.RWMutex
}

func newDisplay(logger Logger) *Display {
	d := &Display{
		logger: logger,
		cursor: &Layer{image: image.NewRGBA(image.Rect(0, 0, 0, 0))},
		layers: newLayers(),
		canvas: image.NewRGBA(image.Rectangle{}),
		tasks:  make(chan task, 10),
	}
	go d.processTasks()
	return d
}

type taskFunc func() error

type task struct {
	taskFunc taskFunc
	name     string
	uuid     uuid.UUID
}

func (t *task) String() string {
	return fmt.Sprintf("%s [%s]", t.name, t.uuid)
}

func (d *Display) scheduleTask(name string, t taskFunc) {
	task := task{
		taskFunc: t,
		name:     name,
		uuid:     uuid.New(),
	}
	d.logger.Tracef("Adding new task: %s. Total: %d", task.String(), len(d.tasks)+1)
	d.pendingTasks.Add(1)
	d.tasks <- task
}

func (d *Display) processSingleTask(t task) {
	d.canvasAccess.Lock()
	d.logger.Tracef("Executing task %s", t.String())
	defer func() {
		d.pendingTasks.Done()
		d.canvasAccess.Unlock()
	}()
	err := t.taskFunc()
	if err != nil {
		d.logger.Errorf("Skipping task %s due to error. This can lead to invalid screen state! Error: %s", t.String(), err)
		return
	}
	defaultLayer := d.layers.getDefault()
	if !defaultLayer.modified {
		return
	}
	mr := defaultLayer.modifiedRect
	copyImage(d.canvas, mr.Min.X, mr.Min.Y, defaultLayer.image, mr, draw.Src)
	d.lastUpdate = time.Now().UnixNano()

	defaultLayer.resetModified()
}

func (d *Display) processTasks() {
	for {
		select {
		case t := <-d.tasks:
			d.processSingleTask(t)
		}
	}
}

func (d *Display) flush() error {
	id := uuid.New()
	d.logger.Tracef("Waiting for %d pending tasks [%s]", len(d.tasks), id.String())
	d.pendingTasks.Wait()
	d.logger.Tracef("All tasks completed [%s]", id.String())
	return nil
}

func (d *Display) Canvas() (image.Image, int64) {
	d.canvasAccess.RLock()
	defer func() {
		d.canvasAccess.RUnlock()
	}()
	return d.canvas, d.lastUpdate
}

func (d *Display) dispose(layerIdx int) {
	d.scheduleTask("dispose", func() error {
		d.layers.delete(layerIdx)
		return nil
	})
}

func (d *Display) copy(srcL, srcX, srcY, srcWidth, srcHeight,
	dstL, dstX, dstY int, compositeOperation byte) {
	srcLayer := d.layers.get(srcL)
	dstLayer := d.layers.get(dstL)
	op := compositeOperations[compositeOperation]
	d.scheduleTask("copy", func() error {
		dstLayer.Copy(srcLayer, srcX, srcY, srcWidth, srcHeight, dstX, dstY, op)
		return nil
	})
}

func (d *Display) draw(layerIdx, x, y int, compositeOperation byte, s *stream) {
	op := compositeOperations[compositeOperation]
	layer := d.layers.get(layerIdx)
	img, err := s.image()

	d.scheduleTask("draw", func() error {
		if err != nil {
			return err
		}
		layer.op = op
		layer.Draw(x, y, img)
		return nil
	})
}

func (d *Display) resize(layerIdx, w, h int) {
	layer := d.layers.get(layerIdx)
	d.scheduleTask("resize", func() error {
		layer.Resize(w, h)
		if layerIdx == 0 {
			d.canvas = image.NewRGBA(layer.image.Bounds())
			copyImage(d.canvas, 0, 0, layer.image, layer.image.Bounds(), draw.Src)
		}
		return nil
	})
}

func (d *Display) setCursor(cursorHotspotX, cursorHotspotY, srcL, srcX, srcY, srcWidth, srcHeight int) {
	layer := d.layers.get(srcL)
	d.scheduleTask("setCursor", func() error {
		d.cursorHotspotX = cursorHotspotX
		d.cursorHotspotY = cursorHotspotY
		d.cursor.Resize(srcWidth, srcHeight)
		d.cursor.Copy(layer, srcX, srcY, srcWidth, srcHeight, 0, 0, draw.Src)
		// TODO (?)
		//d.moveCursor(d.cursorX, d.cursorY)
		defaultLayer := d.layers.getDefault()
		defaultLayer.Copy(d.cursor, 0, 0, srcWidth, srcHeight, cursorHotspotX, cursorHotspotY, draw.Over)
		return nil
	})
}
