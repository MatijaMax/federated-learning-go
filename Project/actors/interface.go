package actors

import (
	"fmt"
	"project/messages"

	"github.com/asynkron/protoactor-go/actor"
)

type SpawnedAveragerPID struct{ PID *actor.PID }

type InterfaceActor struct {
	count              int
	message            string
	spawnedAveragerPID *actor.PID
}

func (state *InterfaceActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.Response:
		state.count++
		state.message = "Input" + string(state.count)
		fmt.Println(msg.GetSomeValue()+":", state.count)
	case *messages.Echo:
		fmt.Printf(msg.GetMessage() + "\n")
	case SpawnedAveragerPID:
		fmt.Println("INTERFEJS dobavio PID Averagera:", msg.PID)
		state.spawnedAveragerPID = msg.PID
	case *messages.TrainerWeightsMessage:
		fmt.Println("JA SAM INTERFEJS: " + msg.NizFloatova)

	}
}
