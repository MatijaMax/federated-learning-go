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

	time.Sleep(time.Hour)
}
