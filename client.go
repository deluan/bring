package bring

import "image"

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

func (c *Client) LastUpdate() int64 {
	return c.display.lastUpdate
}

func (c *Client) Canvas() (image.Image, int64) {
	return c.display.Canvas()
}
