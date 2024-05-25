package actors

import (
	"fmt"
	"project/messages"

	"github.com/asynkron/protoactor-go/actor"
)

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

	}
}
