package main

import (
	"fmt"
	"project/actors"
	"project/messages"

	//"project/messages"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
)

type Hello struct{ Who string }

func main() {
	system := actor.NewActorSystem()
	remoteConfig := remote.Configure("192.168.43.151", 8090)
	remoting := remote.NewRemote(system, remoteConfig)
	remoting.Start()
	context := system.Root

	var interfacePid *actor.PID = nil
	var averagerPid *actor.PID = nil
	var trainerPid *actor.PID = nil

	var interfacePidRemote *actor.PID = nil
	var averagerPidRemote *actor.PID = nil
	var trainerPidRemote *actor.PID = nil
	var interfacePids []*actor.PID

	// Spawn three local actors
	for i := 0; i < 6; i++ {
		if i == 0 {
			props := actor.PropsFromProducer(func() actor.Actor { return &actors.InterfaceActor{} })
			pid := context.Spawn(props)
			interfacePid = pid
			interfacePids = append(interfacePids, interfacePid)
			fmt.Print("INTERFEJS PID: ")
			fmt.Println(interfacePid)

		}
		if i == 1 {
			spawnResponse, err := remoting.SpawnNamed("192.168.43.81:8091", "myactor1", "interfejs1", time.Second*12)
			if err != nil {
				panic(err)
				return
			}
			interfacePidRemote = spawnResponse.Pid
			interfacePids = append(interfacePids, interfacePidRemote)
		}
		if i == 2 {
			spawnResponse, err := remoting.SpawnNamed("192.168.43.81:8091", "myactor2", "trainer1", time.Second*12)
			if err != nil {
				panic(err)
				return
			}
			trainerPidRemote = spawnResponse.Pid
		}
		if i == 3 {
			props := actor.PropsFromProducer(func() actor.Actor { return &actors.AveragerActor{} })
			pid := context.Spawn(props)
			averagerPid = pid
			fmt.Print("AVERAGER PID: ")
			fmt.Println(averagerPid)

		}
		if i == 4 {
			spawnResponse, err := remoting.SpawnNamed("192.168.43.81:8091", "myactor3", "averager1", time.Second*12)
			if err != nil {
				panic(err)
				return
			}
			averagerPidRemote = spawnResponse.Pid
		}
		// "../dataset/Diabetes.csv"
		if i == 5 {
			props := actor.PropsFromProducer(func() actor.Actor { return &actors.TrainerActor{} })
			pid := context.Spawn(props)
			trainerPid = pid
			fmt.Print("TRENER PID: ")
			fmt.Println(trainerPid)

		}

	}

	time.Sleep(time.Second * 1)
	context.Send(interfacePid, &messages.SpawnedAveragerPID{ThePid: averagerPid})
	context.Send(averagerPid, &messages.SpawnedTrainerPID{ThePid: trainerPid})
	context.Send(trainerPid, &messages.SpawnedAveragerPID{ThePid: averagerPid, DataPath: "../dataset/Diabetes.csv1"})
	context.Send(trainerPid, &messages.SpawnedInterfacePID{ThePid: interfacePid})

	context.Send(interfacePidRemote, &messages.RemoteIntegerPID{YourInterfacePid: interfacePidRemote, AllInterfacePids: interfacePids})

	context.Send(interfacePidRemote, &messages.SpawnedAveragerPID{ThePid: averagerPidRemote})
	context.Send(averagerPidRemote, &messages.SpawnedTrainerPID{ThePid: trainerPidRemote})
	context.Send(trainerPidRemote, &messages.SpawnedAveragerPID{ThePid: averagerPidRemote, DataPath: "../dataset/Diabetes.csv1"})
	context.Send(trainerPidRemote, &messages.SpawnedInterfacePID{ThePid: interfacePidRemote})

	// time.Sleep(time.Second * 10)
	// context.Send(interfacePid, &messages.InterInterfaceWeightsMessage{})

	time.Sleep(time.Hour)
}
