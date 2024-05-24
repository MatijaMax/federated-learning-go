package actors

import (
	"fmt"
	"project/messages"

	"github.com/asynkron/protoactor-go/actor"
)

type AveragerActor struct {
	count   int
	message string
	//spawnedPID *actor.PID
}

func (state *AveragerActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.Response:
		state.count++
		state.message = "Input" + string(state.count)
		fmt.Println(msg.GetSomeValue()+":", state.count)
		/*
				go func() {
					fmt.Println("Please enter something:")
					var input string
					fmt.Scanln(&input)
					message := &messages.Echo{Message: input, Sender: context.Self()}
					context.Send(state.spawnedPID, message)

				}()
			/*
			case *messages.Echo:
				state.spawnedPID = msg.GetSender()
		*/
	}
}
