package actors

import (
	"fmt"
	"project/messages"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type InterfaceActor struct {
	count              int
	message            string
	spawnedAveragerPID *actor.PID
	myPid              *actor.PID
	queueInterfaces    []*actor.PID
}

func (state *InterfaceActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.Response:
		state.count++
		state.message = "Input" + string(state.count)
		fmt.Println(msg.GetSomeValue()+":", state.count)

	case *messages.RemoteIntegerPID:
		fmt.Println("INTERFEJS dobavio PID Remote Sistema:", msg.YourInterfacePid)
		fmt.Println("INTERFEJS REMOTE STARTERA:", msg.AllInterfacePids[0])
		// fmt.Println("INTERFEJS imaaa PID Averagera:", state.spawnedAveragerPID.Address)
		state.myPid = msg.YourInterfacePid
		state.queueInterfaces = msg.AllInterfacePids

	case *messages.Echo:
		fmt.Printf(msg.GetMessage() + "\n")
	case *messages.SpawnedAveragerPID:
		fmt.Println("INTERFEJS dobavio PID Averagera:", msg.ThePid)
		state.spawnedAveragerPID = msg.ThePid
	case *messages.TrainerWeightsMessage:
		time.Sleep(time.Second * 2)
		fmt.Println("JA SAM INTERFEJS: " + msg.NizFloatova)
	case *messages.InterInterfaceWeightsMessage:
		time.Sleep(time.Second * 2)
		//fmt.Println("Interfejs: Dobio sam tezine od brace iz klastera (mozda nekad proradi)")
		fmt.Println("Interfejs: Dobio sam tezine od brace iz klastera")
		context.Send(state.spawnedAveragerPID, &messages.InterfaceToAveragerWeightsMessage{
			WeightsIH: msg.WeightsIH,
			WeightsHH: msg.WeightsHH,
			WeightsHO: msg.WeightsHO,
			BiasH:     msg.BiasH,
			BiasH2:    msg.BiasH2,
			BiasO:     msg.BiasO,
		})

	}
}
