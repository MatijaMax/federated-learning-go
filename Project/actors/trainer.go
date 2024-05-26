package actors

import (
	"fmt"
	"project/messages"

	"github.com/asynkron/protoactor-go/actor"
)

type SpawnedInterfacePID struct{ PID *actor.PID }

type TrainerActor struct {
	count               int
	message             string
	spawnedInterfacePID *actor.PID
	spawnedAveragerPID  *actor.PID
}

func (state *TrainerActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.Response:
		state.count++
		state.message = "Input" + string(state.count)
		fmt.Println(msg.GetSomeValue()+":", state.count)
	case *messages.Echo:
		fmt.Printf(msg.GetMessage() + "\n")

	case SpawnedAveragerPID:
		fmt.Print("TRENER dobavio PID Averagera \n")
		state.spawnedAveragerPID = msg.PID
	case SpawnedInterfacePID:
		fmt.Print("TRENER dobavio PID Interfejsa \n")
		state.spawnedInterfacePID = msg.PID

	}
}
