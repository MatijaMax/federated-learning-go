package main

import (
	"fmt"
	"os"
	"os/signal"
	"project/actors"
	"project/messages"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/cluster/clusterproviders/automanaged"
	"github.com/asynkron/protoactor-go/cluster/identitylookup/disthash"
	"github.com/asynkron/protoactor-go/remote"
)

func main() {
	// Set up the actor system
	system := actor.NewActorSystem()

	config := remote.Configure("192.168.43.81", 8081)

	// Configure a cluster on top of the above remote env
	provider := automanaged.NewWithConfig(1*time.Second, 6331, "localhost:6331")
	// provider, err := etcd.NewWithConfig("/protoactor", clientv3.Config{
	// 	Endpoints:   []string{"127.0.0.1:2379"},
	// 	DialTimeout: time.Second * 5,
	// })
	// if err != nil {
	// 	panic(err)
	// }
	lookup := disthash.New()
	clusterKind := cluster.NewKind(
		"Interface",
		actor.PropsFromProducer(func() actor.Actor {
			return &actors.InterfaceActor{}
		}))
	clusterConfig := cluster.Configure("cluster-fed", provider, lookup, config, cluster.WithKinds(clusterKind))
	c := cluster.New(system, clusterConfig)

	// Manage the cluster node's lifecycle
	c.StartMember()
	defer c.Shutdown(false)

	context := system.Root

	var interfaceGrainPid *actor.PID = nil
	var averagerPid *actor.PID = nil
	var trainerPid *actor.PID = nil

	var interfacePids []*actor.PID

	for i := 0; i < 6; i++ {
		if i == 1 {
			interfaceGrainPid = cluster.GetCluster(system).Get("remote-interface-2", "Interface")
			interfacePids = append(interfacePids, interfaceGrainPid)
			fmt.Println("SSSSSSSSSSSSSSSSSSSS")
			fmt.Println(interfaceGrainPid)
		}
		if i == 2 {

		}
		if i == 3 {
			props := actor.PropsFromProducer(func() actor.Actor { return &actors.AveragerActor{} })
			pid := context.Spawn(props)
			averagerPid = pid
		}
		if i == 4 {
			interfaceGrainPidOther := cluster.GetCluster(system).Get("remote-interface-1", "Interface")
			interfacePids = append(interfacePids, interfaceGrainPidOther)
			fmt.Println("EEEEEEEEEEEEEEEE")
			fmt.Println(interfaceGrainPidOther)
		}
		if i == 5 {
			props := actor.PropsFromProducer(func() actor.Actor { return &actors.TrainerActor{} })
			pid := context.Spawn(props)
			trainerPid = pid
		}
	}

	context.Send(interfaceGrainPid, &messages.SpawnedAveragerPID{ThePid: averagerPid})
	context.Send(averagerPid, &messages.SpawnedTrainerPID{ThePid: trainerPid})
	context.Send(trainerPid, &messages.SpawnedAveragerPID{ThePid: averagerPid, DataPath: "../dataset/Diabetes1.csv"})
	context.Send(trainerPid, &messages.SpawnedInterfacePID{ThePid: interfaceGrainPid})

	// context.Send(interfaceGrainPid, &messages.RemoteIntegerPID{YourInterfacePid: interfaceGrainPid, AllInterfacePids: interfacePids})

	// Run till a signal comes
	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt, os.Kill)
	<-finish
}
