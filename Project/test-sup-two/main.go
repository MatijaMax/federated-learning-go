package main

import (
	"fmt"
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
	decider := func(reason interface{}) actor.Directive {
		fmt.Println("handling failure for child")
		return actor.StopDirective
	}
	supervisor := actor.NewOneForOneStrategy(10, 1000, decider)
	
	remoting.Register("interfejs1", actor.PropsFromProducer(func() actor.Actor { return &actors.InterfaceActor{} },actor.WithSupervisor(supervisor)))
	remoting.Register("averager1", actor.PropsFromProducer(func() actor.Actor { return &actors.AveragerActor{} }))
	remoting.Register("trainer1", actor.PropsFromProducer(func() actor.Actor { return &actors.TrainerActor{} }))
	
	

	time.Sleep(time.Hour)
}
