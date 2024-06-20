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
	remoteConfig := remote.Configure("127.0.0.1", 8090)
	remoting := remote.NewRemote(system, remoteConfig)
	remoting.Start()
	context := system.Root

	var interfacePid *actor.PID = nil
	var averagerPid *actor.PID = nil
	var trainerPid *actor.PID = nil

	var interfacePidRemote *actor.PID = nil
	var averagerPidRemote *actor.PID = nil
	var trainerPidRemote *actor.PID = nil

	var interfacePidRemote3 *actor.PID = nil
	var averagerPidRemote3 *actor.PID = nil
	var trainerPidRemote3 *actor.PID = nil

	var interfacePids []*actor.PID

	// Spawn actors
	for i := 0; i < 9; i++ {
		if i == 0 {
			props := actor.PropsFromProducer(func() actor.Actor { return &actors.InterfaceActor{} })
			pid := context.Spawn(props)
			interfacePid = pid
			interfacePids = append(interfacePids, interfacePid)
			fmt.Print("INTERFEJS PID: ")
			fmt.Println(interfacePid)

		}
		if i == 1 {
			spawnResponse, err := remoting.SpawnNamed("127.0.0.1:8092", "myactor1", "interfejs2", time.Second*12)
			if err != nil {
				panic(err)
				return
			}
			interfacePidRemote = spawnResponse.Pid
			interfacePids = append(interfacePids, interfacePidRemote)
		}
		if i == 2 {
			spawnResponse, err := remoting.SpawnNamed("127.0.0.1:8092", "myactor2", "trainer2", time.Second*12)
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
			spawnResponse, err := remoting.SpawnNamed("127.0.0.1:8092", "myactor3", "averager2", time.Second*12)
			if err != nil {
				panic(err)
				return
			}
			averagerPidRemote = spawnResponse.Pid
		}
		if i == 5 {
			props := actor.PropsFromProducer(func() actor.Actor { return &actors.TrainerActor{} })
			pid := context.Spawn(props)
			trainerPid = pid
			fmt.Print("TRENER PID: ")
			fmt.Println(trainerPid)

		}
		if i == 6 {
			spawnResponse, err := remoting.SpawnNamed("127.0.0.1:8093", "myactor13", "interfejs3", time.Second*12)
			if err != nil {
				panic(err)
				return
			}
			interfacePidRemote3 = spawnResponse.Pid
			interfacePids = append(interfacePids, interfacePidRemote3)
		}
		if i == 7 {
			spawnResponse, err := remoting.SpawnNamed("127.0.0.1:8093", "myactor23", "trainer3", time.Second*12)
			if err != nil {
				panic(err)
				return
			}
			trainerPidRemote3 = spawnResponse.Pid
		}
		if i == 8 {
			spawnResponse, err := remoting.SpawnNamed("127.0.0.1:8093", "myactor33", "averager3", time.Second*12)
			if err != nil {
				panic(err)
				return
			}
			averagerPidRemote3 = spawnResponse.Pid
		}
	}

	time.Sleep(time.Second * 1)
	context.Send(interfacePid, &messages.SpawnedAveragerPID{ThePid: averagerPid})
	context.Send(averagerPid, &messages.SpawnedTrainerPID{ThePid: trainerPid})
	context.Send(trainerPid, &messages.SpawnedAveragerPID{ThePid: averagerPid, DataPath: "../dataset/Diabetes31.csv"})
	context.Send(trainerPid, &messages.SpawnedInterfacePID{ThePid: interfacePid})

	context.Send(interfacePidRemote, &messages.RemoteIntegerPID{YourInterfacePid: interfacePidRemote, AllInterfacePids: interfacePids})
	context.Send(interfacePid, &messages.RemoteIntegerPID{YourInterfacePid: interfacePid, AllInterfacePids: interfacePids})
	context.Send(interfacePidRemote3, &messages.RemoteIntegerPID{YourInterfacePid: interfacePidRemote3, AllInterfacePids: interfacePids})

	context.Send(interfacePidRemote, &messages.SpawnedAveragerPID{ThePid: averagerPidRemote})
	context.Send(averagerPidRemote, &messages.SpawnedTrainerPID{ThePid: trainerPidRemote})
	context.Send(trainerPidRemote, &messages.SpawnedAveragerPID{ThePid: averagerPidRemote, DataPath: "../dataset/Diabetes32.csv"})
	context.Send(trainerPidRemote, &messages.SpawnedInterfacePID{ThePid: interfacePidRemote})

	context.Send(interfacePidRemote3, &messages.SpawnedAveragerPID{ThePid: averagerPidRemote3})
	context.Send(averagerPidRemote3, &messages.SpawnedTrainerPID{ThePid: trainerPidRemote3})
	context.Send(trainerPidRemote3, &messages.SpawnedAveragerPID{ThePid: averagerPidRemote3, DataPath: "../dataset/Diabetes33.csv"})
	context.Send(trainerPidRemote3, &messages.SpawnedInterfacePID{ThePid: interfacePidRemote3})

	// time.Sleep(time.Second * 10)
	// context.Send(interfacePid, &messages.InterInterfaceWeightsMessage{})

	time.Sleep(time.Hour)
}
