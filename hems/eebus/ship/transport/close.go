package transport

import (
	"errors"
	"time"

	"github.com/andig/evcc/hems/eebus/ship/message"
	"github.com/andig/evcc/hems/eebus/ship/ship"
)

// AcceptClose accepts connection close
func (c *Transport) AcceptClose() error {
	err := c.WriteJSON(message.CmiTypeEnd, ship.CmiConnectionClose{
		ConnectionClose: ship.ConnectionClose{
			Phase: ship.ConnectionClosePhaseTypeConfirm,
		},
	})

	// stop read/write pump
	close(c.closeC)

	return err
}

// Close closes the connection
func (c *Transport) Close() error {
	err := c.WriteJSON(message.CmiTypeEnd, ship.CmiConnectionClose{
		ConnectionClose: ship.ConnectionClose{
			Phase: ship.ConnectionClosePhaseTypeAnnounce,
			// MaxTime: int(ship.CmiCloseTimeout / time.Millisecond),
		},
	})

	timer := time.NewTimer(message.CmiCloseTimeout)
	for err == nil {
		msg, err := c.ReadMessage(timer.C)
		if err != nil {
			break
		}

		if typed, ok := msg.(ship.ConnectionClose); ok && typed.Phase == ship.ConnectionClosePhaseTypeConfirm {
			break
		}

		err = errors.New("close: invalid response")
	}

	// stop read/write pump
	close(c.closeC)

	return err
}
