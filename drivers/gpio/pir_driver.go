package gpio

import (
	"time"

	"github.com/hybridgroup/gobot"
)

// PirDriver represents a PIR device
type PirDriver struct {
	pin        string
	name       string
	halt       chan bool
	connection DigitalReader
	Active     bool
	interval   time.Duration
	gobot.Eventer
}

const (
	// Error event
	PirTrigger = "trigger"
	// Release event
	PirKill = "kill"
)

// NewPirDriver returns a new PirDriver with polling interval of 10 Millisecond.
//
// pass in time Duration to change the polling interval of new informatoin
func NewPirDriver(a DigitalReader, pin string, t ...time.Duration) *PirDriver {
	p := &PirDriver{
		name:       "PIR",
		connection: a,
		pin:        pin,
		Active:     false,
		Eventer:    gobot.NewEventer(),
		interval:   10 * time.Millisecond,
		halt:       make(chan bool),
	}

	if len(t) > 0 {
		p.interval = t[0]
	}

	//p.AddEvent(Data)
	p.AddEvent(PirTrigger)
	p.AddEvent(PirKill)
	p.AddEvent(Error)

	return p
}

// Name return the PirDriver name
func (p *PirDriver) Name() string { return p.name }

// SetName sets the MakeyButtonDrivers name
func (p *PirDriver) SetName(n string) { p.name = n }

// Pin returns the PirDriver pin
func (p *PirDriver) Pin() string { return p.pin }

// Connection returns the PirDriver Connection
func (p *PirDriver) Connection() gobot.Connection { return p.connection.(gobot.Connection) }

// Start starts PirDriver and polls the state of the pin at the given interval
func (p *PirDriver) Start() (errs []error) {
	state := 0
	go func() {
		for {
			newValue, err := p.connection.DigitalRead(p.Pin())
			if err != nil {
				p.Publish(Error, err)
			} else if newValue != state && newValue != -1 {
				state = newValue
				if newValue == 0 {
					p.Active = false
					p.Publish(PirKill, newValue)
				} else {
					p.Active = true
					p.Publish(PirTrigger, newValue)
				}
			}
			select {
			case <-time.After(p.interval):
			case <-p.halt:
				return
			}
		}
	}()
	return
}

// Halt stops polling the PirDriver for new information
func (p *PirDriver) Halt() (err []error) {
	p.halt <- true
	return
}
