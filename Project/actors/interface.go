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
		fmt.Println("INTERFEJS dobavio svoj PID:", msg.YourInterfacePid)
		fmt.Println("INTERFEJS REMOTE STARTERA:", msg.AllInterfacePids[0])
		// fmt.Println("INTERFEJS imaaa PID Averagera:", state.spawnedAveragerPID.Address)
		state.myPid = msg.YourInterfacePid
		state.queueInterfaces = msg.AllInterfacePids

	case *messages.Echo:
		var averagerPid *actor.PID = nil
		var trainerPid *actor.PID = nil
		for i := 0; i < 3; i++ {
		if i == 1 {
					props := actor.PropsFromProducer(func() actor.Actor { return &AveragerActor{} })
					pid := context.Spawn(props)
					averagerPid = pid
					fmt.Print("AVERAGER PID: ")
					fmt.Println(averagerPid)

				}
				if i == 2 {
					props := actor.PropsFromProducer(func() actor.Actor { return &TrainerActor{} })
					pid := context.Spawn(props)
					trainerPid = pid
					fmt.Print("TRENER PID: ")
					fmt.Println(trainerPid)
				}
		}
		context.Send(state.myPid, &messages.SpawnedAveragerPID{ThePid: averagerPid})
		context.Send(averagerPid, &messages.SpawnedTrainerPID{ThePid: trainerPid})
		context.Send(trainerPid, &messages.SpawnedAveragerPID{ThePid: averagerPid, DataPath: "../dataset/DiabetesNew1.csv"})
		context.Send(trainerPid, &messages.SpawnedInterfacePID{ThePid: state.myPid})
		fmt.Printf(msg.GetMessage() + "\n")
	case *messages.SpawnedAveragerPID:
		fmt.Println("INTERFEJS dobavio PID Averagera:", msg.ThePid)
		state.spawnedAveragerPID = msg.ThePid
	case *messages.TrainerWeightsMessage:
		time.Sleep(time.Second * 2)
		fmt.Println("JA SAM INTERFEJS: " + msg.NizFloatova)
		for _, pid := range state.queueInterfaces {
			//fmt.Println(state.queueInterfaces)
			// fmt.Printf("Index %d: %v\n", index, pid)
			if pid.Id != state.myPid.Id {
				fmt.Println(pid)
				context.Send(pid, &messages.InterInterfaceWeightsMessage{
					WeightsIH: msg.WeightsIH,
					WeightsHH: msg.WeightsHH,
					WeightsHO: msg.WeightsHO,
					BiasH:     msg.BiasH,
					BiasH2:    msg.BiasH2,
					BiasO:     msg.BiasO,
				})
			}
		}
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
