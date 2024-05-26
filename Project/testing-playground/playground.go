package main

import (
	"project/actors"
	"project/messages"
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

	// Spawn three local actors
	for i := 0; i < 7; i++ {
		if i == 0 {
			props := actor.PropsFromProducer(func() actor.Actor { return &actors.InterfaceActor{} })
			pid := context.Spawn(props)
			interfacePid = pid
			go func() {
				message := &messages.Echo{Message: "Poruka init INTERFACE", Sender: pid}
				context.Send(pid, message)
			}()

		}
		if i == 1 {
			props := actor.PropsFromProducer(func() actor.Actor { return &actors.AveragerActor{} })
			pid := context.Spawn(props)
			averagerPid = pid
			go func() {
				message := &messages.Echo{Message: "Poruka init AVERAGER", Sender: pid}
				context.Send(pid, message)
			}()
		}
		if i == 2 {
			props := actor.PropsFromProducer(func() actor.Actor { return &actors.TrainerActor{} })
			pid := context.Spawn(props)
			trainerPid = pid
			go func() {
				message := &messages.Echo{Message: "Poruka init TRAINER", Sender: pid}
				context.Send(pid, message)
			}()
		}

	}

	context.Send(interfacePid, actors.SpawnedAveragerPID{PID: averagerPid})
	context.Send(averagerPid, actors.SpawnedTrainerPID{PID: trainerPid})
	context.Send(trainerPid, actors.SpawnedAveragerPID{PID: averagerPid})
	context.Send(trainerPid, actors.SpawnedInterfacePID{PID: interfacePid})

	time.Sleep(time.Hour)

}
