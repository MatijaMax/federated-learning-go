package main

import (
	"fmt"
	"project/actors"
	"project/messages"

	//"project/messages"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/cluster/clusterproviders/automanaged"
	"github.com/asynkron/protoactor-go/cluster/identitylookup/disthash"
	"github.com/asynkron/protoactor-go/remote"
)

func main() {

	system := actor.NewActorSystem()

	config := remote.Configure("192.168.43.151", 8090)

	// Configure a cluster on top of the above remote env
	provider := automanaged.NewWithConfig(1*time.Second, 6332, "192.168.43.81:6331")
	// provider, err := etcd.NewWithConfig("/protoactor", clientv3.Config{
	// 	Endpoints:   []string{"127.0.0.1:2379"},
	// 	DialTimeout: time.Second * 1,
	// })
	// if err != nil {
	// 	panic(err)
	// }
	lookup := disthash.New()
	clusterConfig := cluster.Configure("cluster-fed", provider, lookup, config)
	c := cluster.New(system, clusterConfig)

	c.StartMember()
	defer c.Shutdown(false)

	time.Sleep(1 * time.Second)

	//##########################

	context := system.Root

	var averagerPid *actor.PID = nil
	var trainerPid *actor.PID = nil

	var interfaceGrainPid *actor.PID = nil
	var interfaceGrainPidOther *actor.PID = nil
	var interfacePids []*actor.PID

	// Spawn three local actors
	for i := 0; i < 6; i++ {
		if i == 0 {
			// interfacePid = system.Root.Spawn(interfaceProps)
			// interfacePids = append(interfacePids, interfacePid)
		}
		if i == 1 {
			interfaceGrainPid = cluster.GetCluster(system).Get("remote-interface-1", "Interface")
			interfacePids = append(interfacePids, interfaceGrainPid)
			fmt.Println("XXXXXXXXXXXXX")
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
			interfaceGrainPidOther = cluster.GetCluster(system).Get("remote-interface-2", "Interface")
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

	//}

	time.Sleep(time.Second * 10)

	context.Send(interfaceGrainPid, &messages.SpawnedAveragerPID{ThePid: averagerPid})
	context.Send(averagerPid, &messages.SpawnedTrainerPID{ThePid: trainerPid})
	context.Send(trainerPid, &messages.SpawnedAveragerPID{ThePid: averagerPid, DataPath: "../dataset/Diabetes2.csv"})
	context.Send(trainerPid, &messages.SpawnedInterfacePID{ThePid: interfaceGrainPid})

	context.Send(interfaceGrainPid, &messages.RemoteIntegerPID{YourInterfacePid: interfaceGrainPid, AllInterfacePids: interfacePids})
	context.Send(interfaceGrainPidOther, &messages.RemoteIntegerPID{YourInterfacePid: interfaceGrainPidOther, AllInterfacePids: interfacePids})

	time.Sleep(time.Hour)
}
