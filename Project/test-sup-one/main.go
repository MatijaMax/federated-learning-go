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
	

	var interfacePid *actor.PID = nil
	

	var interfacePidRemote *actor.PID = nil
	
	var interfacePids []*actor.PID

	
	decider := func(reason interface{}) actor.Directive {
		fmt.Println("handling failure for child")
		return actor.StopDirective
	}
	supervisor := actor.NewOneForOneStrategy(10, 1000, decider)
	context := system.Root

	for i := 0; i < 2; i++ {
		if i == 0 {
			
			props := actor.PropsFromProducer(func() actor.Actor { return &actors.InterfaceActor{} },actor.WithSupervisor(supervisor))
			pid := context.Spawn(props)
			interfacePid = pid
			interfacePids = append(interfacePids, interfacePid)
			
		}
		if i == 1 {
			spawnResponse, err := remoting.SpawnNamed("127.0.0.1:8091", "myactor1", "interfejs1", time.Second*12)
			if err != nil {
				panic(err)
				return
			}
			
			interfacePidRemote = spawnResponse.Pid
			interfacePids = append(interfacePids, interfacePidRemote)
		}
		
	}

	time.Sleep(time.Second * 1)
	

	context.Send(interfacePidRemote, &messages.RemoteIntegerPID{YourInterfacePid: interfacePidRemote, AllInterfacePids: interfacePids})
	context.Send(interfacePid, &messages.RemoteIntegerPID{YourInterfacePid: interfacePid, AllInterfacePids: interfacePids})

	time.Sleep(time.Second * 1)

	context.Send(interfacePidRemote, &messages.Echo{Message: "Kreni 1"})
	context.Send(interfacePid, &messages.Echo{Message: "Kreni 2"})

	

	time.Sleep(time.Hour)
}
