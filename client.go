package bring

import (
	"image"
	"strconv"
)

type Client struct {
	session *Session
	display *Display
	streams streams
	logger  Logger
}

func NewClient(session *Session, logger ...Logger) (*Client, error) {
	var log Logger
	if len(logger) > 0 {
		log = logger[0]
	} else {
		log = &DefaultLogger{}
	}

	c := &Client{
		session: session,
		display: newDisplay(log),
		streams: newStreams(),
		logger:  log,
	}
	return c, nil
}

func (c *Client) Start() {
	for {
		select {
		case ins := <-c.session.In:
			if h, ok := handlers[ins.opcode]; ok {
				err := h(c, ins.args)
				if err != nil {
					c.session.Terminate()
				}
				continue
			}
			c.logger.Errorf("Not implemented: %s", ins.opcode)
		}
	}
}

func (c *Client) Canvas() (image image.Image, lastUpdate int64) {
	return c.display.Canvas()
}

func (c *Client) MoveMouse(p image.Point, pressedButtons ...int) {
	if c.session.State != SessionActive {
		return
	}

	buttonMask := 0
	for _, b := range pressedButtons {
		buttonMask |= b
	}
	c.display.moveCursor(p.X, p.Y)
	err := c.session.Send(NewInstruction("mouse", strconv.Itoa(p.X), strconv.Itoa(p.Y), strconv.Itoa(buttonMask)))
	if err != nil {
		c.logger.Errorf("could not send mouse position: %s", err)
	}
}

func (c *Client) SendText(s string) {
	if c.session.State != SessionActive {
		return
	}

	for _, ch := range s {
		keycode := strconv.Itoa(int(ch))
		c.session.Send(NewInstruction("key", keycode, "1"))
		c.session.Send(NewInstruction("key", keycode, "0"))
	}
}

func (c *Client) SendKey(key KeyCode, pressed bool) {
	if c.session.State != SessionActive {
		return
	}

	var p string = "0"
	if pressed {
		p = "1"
	}
	for _, k := range key {
		keycode := strconv.Itoa(k)
		c.session.Send(NewInstruction("key", keycode, p))
	}
}
