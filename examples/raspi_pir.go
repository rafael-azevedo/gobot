package main

import (
	"fmt"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/drivers/gpio"
	"github.com/hybridgroup/gobot/platforms/raspi"
)

func main() {
	r := raspi.NewAdapter()
	led := gpio.NewLedDriver(r, "12")
	pir := gpio.NewPirDriver(r, "11")

	work := func() {

		pir.On(gpio.PirTrigger, func(data interface{}) {
			fmt.Println("Turning LED on")
			led.On()
			fmt.Println("Should be on")
			time.Sleep(1 * time.Second)
		})

		pir.On(gpio.PirKill, func(data interface{}) {
			fmt.Println("Turning LED off")
			led.Off()
			fmt.Println("Should be off")
			time.Sleep(1 * time.Second)

		})

	}

	robot := gobot.NewRobot("pirbot",
		[]gobot.Connection{r},
		[]gobot.Device{led, pir},
		work,
	)

	robot.Start()

}
