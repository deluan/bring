package bring

import (
	"errors"
	"image"
	"strconv"

	"github.com/deluan/bring/protocol"
)

var ErrInvalidKeyCode = errors.New("invalid key code")

// OnSyncFunc is the signature for OnSync event handlers. It will receive the current screen image and the
// timestamp of the last update.
type OnSyncFunc = func(image image.Image, lastUpdate int64)

// Guacamole protocol client. Automatically handles incoming and outgoing Guacamole instructions,
// updating its display using one or more graphic primitives.
type Client struct {
	session *session
	display *display
	streams streams
	logger  Logger
	onSync  OnSyncFunc
}

// NewClient creates a Client and connects it to the guacd server with the provided configuration. Logger is optional
func NewClient(addr string, remoteProtocol string, config map[string]string, logger ...Logger) (*Client, error) {
	var log Logger
	if len(logger) > 0 {
		log = logger[0]
	} else {
		log = &DefaultLogger{}
	}

	s, err := newSession(addr, remoteProtocol, config, log)
	if err != nil {
		return nil, err
	}

	c := &Client{
		session: s,
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
			h, ok := handlers[ins.Opcode]
			if !ok {
				c.logger.Errorf("Instruction not implemented: %s", ins.Opcode)
				continue
			}
			err := h(c, ins.Args)
			if err != nil {
				c.session.Terminate()
			}
		}
	}
}

// OnSync sets a function that will be called on every sync instruction received. This event
// usually happens after a batch of updates are received from the guacd server, making it a
// perfect way to get the current screenshot without having to poll with Screen().
// The handler is expected to be called frequently, so avoid adding any blocking behaviour.
// If your handler is slow, consider using a concurrent pattern (using goroutines)
func (c *Client) OnSync(f OnSyncFunc) {
	c.onSync = f
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
	err := c.session.Send(protocol.NewInstruction("mouse", strconv.Itoa(p.X), strconv.Itoa(p.Y), strconv.Itoa(buttonMask)))
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
		err := c.session.Send(protocol.NewInstruction("key", keycode, "1"))
		if err != nil {
			return nil
		}
		err = c.session.Send(protocol.NewInstruction("key", keycode, "0"))
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

	p := "0"
	if pressed {
		p = "1"
	}
	keySym, ok := keySyms[key]
	if !ok {
		return ErrInvalidKeyCode
	}
	for _, k := range keySym {
		keycode := strconv.Itoa(k)
		err := c.session.Send(protocol.NewInstruction("key", keycode, p))
		if err != nil {
			return nil
		}
	}
	return nil
}
