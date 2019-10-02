package bring

import (
	"image"
	"strconv"
)

// Guacamole protocol client. Given a Session, automatically handles incoming
// and outgoing Guacamole instructions via the provided session, updating its
// display using one or more graphic primitives.
type Client struct {
	session *Session
	display *Display
	streams streams
	logger  Logger
}

// Creates a new Client with the provided Session and Logger
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

// Starts the Client's main loop. It is a blocking call, so it
// should be called in its on go routine
func (c *Client) Start() {
	for {
		select {
		case ins := <-c.session.In:
			h, ok := handlers[ins.opcode]
			if !ok {
				c.logger.Errorf("Instruction not implemented: %s", ins.opcode)
				continue
			}
			err := h(c, ins.args)
			if err != nil {
				c.session.Terminate()
			}
		}
	}
}

// Returns a snapshot of the current screen, together with the last updated timestamp
func (c *Client) Screen() (image image.Image, lastUpdate int64) {
	return c.display.getCanvas()
}

// Returns the current session state
func (c *Client) State() SessionState {
	return c.session.State
}

// Send mouse events to the server. An event is composed by position of the
// cursor, and a list of any currently pressed MouseButtons
func (c *Client) SendMouse(p image.Point, pressedButtons ...MouseButton) error {
	if c.session.State != SessionActive {
		return ErrNotConnected
	}

	buttonMask := 0
	for _, b := range pressedButtons {
		buttonMask |= int(b)
	}
	c.display.moveCursor(p.X, p.Y)
	err := c.session.Send(NewInstruction("mouse", strconv.Itoa(p.X), strconv.Itoa(p.Y), strconv.Itoa(buttonMask)))
	if err != nil {
		return err
	}
	return nil
}

// Send the sequence of characters as they were typed. Only works with simple chars
// (no combination with control keys)
func (c *Client) SendText(sequence string) error {
	if c.session.State != SessionActive {
		return ErrNotConnected
	}

	for _, ch := range sequence {
		keycode := strconv.Itoa(int(ch))
		err := c.session.Send(NewInstruction("key", keycode, "1"))
		if err != nil {
			return nil
		}
		err = c.session.Send(NewInstruction("key", keycode, "0"))
		if err != nil {
			return nil
		}
	}
	return nil
}

// Send key presses and releases.
func (c *Client) SendKey(key KeyCode, pressed bool) error {
	if c.session.State != SessionActive {
		return ErrNotConnected
	}

	var p string = "0"
	if pressed {
		p = "1"
	}
	for _, k := range key {
		keycode := strconv.Itoa(k)
		err := c.session.Send(NewInstruction("key", keycode, p))
		if err != nil {
			return nil
		}
	}
	return nil
}
