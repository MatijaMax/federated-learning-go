package main

import (
	"fmt"
	"time"

	//"project/actors"
	//"project/messages"
	console "github.com/asynkron/goconsole"
	"github.com/asynkron/protoactor-go/actor"
)

type Hello struct{ Who string }

func Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case Hello:
		context.Respond("Hello " + msg.Who)
	}
}

func main() {
	system := actor.NewActorSystem()
	rootContext := system.Root
	props := actor.PropsFromFunc(Receive)
	pid := rootContext.Spawn(props)
	result, _ := rootContext.RequestFuture(pid, Hello{Who: "World"}, 30*time.Second).Result()

	fmt.Println(result)
	_, _ = console.ReadLine()
}
