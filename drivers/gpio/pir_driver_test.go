package gpio

import (
	"errors"
	"testing"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/gobottest"
)

var _ gobot.Driver = (*PirDriver)(nil)

const PirDriver_TEST_DELAY = 30

func initTestPirDriver() *PirDriver {
	return NewPirDriver(newGpioTestAdaptor(), "1")
}

func TestPirDriverHalt(t *testing.T) {
	p := initTestPirDriver()
	go func() {
		<-p.halt
	}()
	gobottest.Assert(t, len(p.Halt()), 0)
}

func TestPirDriver(t *testing.T) {
	p := NewPirDriver(newGpioTestAdaptor(), "1")
	gobottest.Assert(t, p.Pin(), "1")
	gobottest.Refute(t, p.Connection(), nil)

	p = NewPirDriver(newGpioTestAdaptor(), "1", 30*time.Second)
	gobottest.Assert(t, p.interval, PirDriver_TEST_DELAY*time.Second)
}

func TestPirDriverStart(t *testing.T) {

	sem := make(chan bool, 0)
	p := initTestPirDriver()
	gobottest.Assert(t, len(p.Start()), 0)

	testAdaptorDigitalRead = func() (val int, err error) {
		val = 1
		return
	}

	p.Once(PirTrigger, func(data interface{}) {
		gobottest.Assert(t, p.Active, true)
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(PirDriver_TEST_DELAY * time.Millisecond):
		t.Errorf("PIR Event \"trigger\" was not published")
	}

	testAdaptorDigitalRead = func() (val int, err error) {
		val = 0
		return
	}

	p.Once(PirKill, func(data interface{}) {
		gobottest.Assert(t, p.Active, false)
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(PirDriver_TEST_DELAY * time.Millisecond):
		t.Errorf("PIR Event \"kill\" was not published")
	}

	testAdaptorDigitalRead = func() (val int, err error) {
		err = errors.New("digital read error")
		return
	}

	p.Once(Error, func(data interface{}) {
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(PirDriver_TEST_DELAY * time.Millisecond):
		t.Errorf("PIR Event \"error\" was not published")
	}

	// send a halt message
	p.Once(p.Event(Data), func(data interface{}) {
		sem <- true
	})

	testAdaptorDigitalRead = func() (val int, err error) {
		val = 0
		return
	}

	p.halt <- true

	select {
	case <-sem:
		t.Errorf("PirDriver Event should not published")
	case <-time.After(PirDriver_TEST_DELAY * time.Millisecond):
	}
}
