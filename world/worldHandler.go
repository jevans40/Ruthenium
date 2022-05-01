package world

import (
	"github.com/jevans40/Ruthenium/ruthutil"
	log "github.com/sirupsen/logrus"
)

type TickChannelCommunication int

const (
	ResumeTick TickChannelCommunication = iota
	PauseTick
	MaintainTick
	InitTick
	KillTick
)

type WorldHandler struct {
	world  World
	paused bool
}

func (w *WorldHandler) RegisterWorld(world World) {
	w.world = world
}

func (w *WorldHandler) StartHandler(tickChannel chan TickChannelCommunication) {
	//Wait
	for {
		Tick, err := ruthutil.WaitChannel(tickChannel)
		if err != nil {
			return
		} else if Tick == ResumeTick {
			w.world.Resume()
		} else if Tick == PauseTick {
			w.world.Pause()
		} else if Tick == MaintainTick {
			w.world.Maintain()
		} else if Tick == InitTick {
			w.world.Init()
		} else if Tick == KillTick {
			return
		} else {
			log.WithFields(log.Fields{"Handler": w, "Tick": Tick}).Error("Invalid Tick Recieved in World Handler")
		}

	}

}
