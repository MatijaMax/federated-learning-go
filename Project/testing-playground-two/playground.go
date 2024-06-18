package main

import (
	"project/actors"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
)

type Hello struct{ Who string }

func main() {
	system := actor.NewActorSystem()
	remoteConfig := remote.Configure("192.168.43.81", 8091)
	remoting := remote.NewRemote(system, remoteConfig)
	remoting.Start()
	remoting.Register("interfejs1", actor.PropsFromProducer(func() actor.Actor { return &actors.InterfaceActor{} }))
	remoting.Register("averager1", actor.PropsFromProducer(func() actor.Actor { return &actors.AveragerActor{} }))
	remoting.Register("trainer1", actor.PropsFromProducer(func() actor.Actor { return &actors.TrainerActor{} }))
	// context := system.Root

	// //var interfacePid *actor.PID = nil
	// var averagerPid *actor.PID = nil
	// var trainerPid *actor.PID = nil
	// var interfacePid *actor.PID = nil

	// // Spawn three local actors
	// for i := 0; i < 3; i++ {
	// 	if i == 0 {
	// 		props := actor.PropsFromProducer(func() actor.Actor { return &actors.InterfaceActor{} })
	// 		pid, err := context.SpawnNamed(props, "hello")
	// 		if err != nil {
	// 			panic(err)
	// 			return
	// 		}
	// 		interfacePid = pid
	// 		fmt.Print("INTERFEJS PID: ")
	// 		fmt.Println(interfacePid)

	// 	}
	// 	if i == 1 {
	// 		props := actor.PropsFromProducer(func() actor.Actor { return &actors.AveragerActor{} })
	// 		pid := context.Spawn(props)
	// 		averagerPid = pid
	// 		fmt.Print("AVERAGER PID: ")
	// 		fmt.Println(averagerPid)
	// 		go func() {
	// 			//message := &messages.Echo{Message: "Poruka init AVERAGER", Sender: pid}
	// 			//context.Send(pid, message)
	// 		}()
	// 	}
	// 	if i == 2 {
	// 		props := actor.PropsFromProducer(func() actor.Actor { return &actors.TrainerActor{} })
	// 		pid := context.Spawn(props)
	// 		trainerPid = pid
	// 		fmt.Print("TRENER PID: ")
	// 		fmt.Println(trainerPid)
	// 		go func() {
	// 			//message := &messages.Echo{Message: "Poruka init TRAINER", Sender: pid}
	// 			//context.Send(pid, message)
	// 		}()
	// 	}

	// }

	// time.Sleep(time.Second)
	// context.Send(interfacePid, actors.SpawnedAveragerPID{PID: averagerPid})
	// context.Send(averagerPid, actors.SpawnedTrainerPID{PID: trainerPid})
	// context.Send(trainerPid, actors.SpawnedAveragerPID{PID: averagerPid})
	// context.Send(trainerPid, actors.SpawnedInterfacePID{PID: interfacePid})

	// time.Sleep(time.Second * 10)
	// context.Send(interfacePid, &messages.InterInterfaceWeightsMessage{})

	time.Sleep(time.Hour)
}
