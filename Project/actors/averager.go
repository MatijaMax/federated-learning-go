package actors

import (
	"fmt"
	"project/messages"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type SpawnedTrainerPID struct{ PID *actor.PID }

type AveragerActor struct {
	count             int
	message           string
	spawnedTrainerPID *actor.PID
	weightsIH   [][]float64 // Weights between input and hidden layer
    weightsHH   [][]float64 // Weights between input and hidden layer
    weightsHO   [][]float64 // Weights between hidden and output layer
    biasH       []float64   // Bias for the hidden layer
    biasH2       []float64   // Bias for the hidden layer
    biasO       []float64   // Bias for the output layer
	
}

func (state *AveragerActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.Response:
		state.count++
		state.message = "Input" + string(state.count)
		fmt.Println(msg.GetSomeValue()+":", state.count)
		/*
				go func() {
					fmt.Println("Please enter something:")
					var input string
					fmt.Scanln(&input)
					message := &messages.Echo{Message: input, Sender: context.Self()}
					context.Send(state.spawnedPID, message)

				}()
			/*
			case *messages.Echo:
				state.spawnedPID = msg.GetSender()
		*/
	case *messages.Echo:
		fmt.Printf(msg.GetMessage() + "\n")

	case SpawnedTrainerPID:
		fmt.Println("AVERAGER dobavio PID Trenera: ", msg.PID)
		state.spawnedTrainerPID = msg.PID

	case *messages.TrainerWeightsMessage:
		time.Sleep(time.Second * 2)
		for _, floatArray := range msg.WeightsIH {
			var row []float64
			for _, value := range floatArray.Column {
				row = append(row, value)
			}
			state.weightsIH = append(state.weightsIH, row)
		}
		
		for _, floatArray := range msg.WeightsHH {
			var row []float64
			for _, value := range floatArray.Column {
				row = append(row, value)
			}
			state.weightsHH = append(state.weightsHH, row)
		}

		for _, floatArray := range msg.WeightsHO {
			var row []float64
			for _, value := range floatArray.Column {
				row = append(row, value)
			}
			state.weightsHO = append(state.weightsHO, row)
		}

		state.biasH = msg.BiasH
		state.biasH2 = msg.BiasH2
		state.biasO = msg.BiasO
		fmt.Println("JA SAM AVERAGER: " + msg.NizFloatova)
		for i := range state.weightsIH {
			for j := range state.weightsIH[i] {
				fmt.Print(state.weightsIH[i][j])
			}
		}

	}

}
