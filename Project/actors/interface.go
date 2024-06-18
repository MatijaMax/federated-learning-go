package actors

import (
	"fmt"
	"project/messages"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type SpawnedAveragerPID struct{ PID *actor.PID }

type RemoteIntegerPID struct{ PID *actor.PID }

type InterfaceActor struct {
	count              int
	message            string
	spawnedAveragerPID *actor.PID
	queueInterfaces    []*actor.PID
}

func (state *InterfaceActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.Response:
		state.count++
		state.message = "Input" + string(state.count)
		fmt.Println(msg.GetSomeValue()+":", state.count)

	case RemoteIntegerPID:
		fmt.Println("INTERFEJS dobavio PID Remote Sistema:", msg.PID)
		state.queueInterfaces = append(state.queueInterfaces, msg.PID)

	case *messages.Echo:
		fmt.Printf(msg.GetMessage() + "\n")
	case SpawnedAveragerPID:
		fmt.Println("INTERFEJS dobavio PID Averagera:", msg.PID)
		state.spawnedAveragerPID = msg.PID
	case *messages.TrainerWeightsMessage:
		time.Sleep(time.Second * 2)
		fmt.Println("JA SAM INTERFEJS: " + msg.NizFloatova)
	case *messages.InterInterfaceWeightsMessage:
		time.Sleep(time.Second * 2)
		fmt.Println("Interfejs: Dobio sam tezine od brace iz klastera (mozda nekad proradi)")
		context.Send(state.spawnedAveragerPID, &messages.InterfaceToAveragerWeightsMessage{})

	}
}
