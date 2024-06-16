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
	hasWeights bool
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

	case *messages.InterfaceToAveragerWeightsMessage:
		fmt.Println("WWWWWWWWWWWWWWWWWWWWWWWWWWW")
		if(state.hasWeights == false){
			fmt.Println("EEEEEEEEEEE")
			state.hasWeights = true
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
			
			
			fmt.Println("AVERAGER: dobio sam nove tezine od interfejsa")
		}else{
			//ovaj deo cu za sad da zakomentarisem gde prima tezine od interfejsa i usredni ih, testiram samo da prosledi opet do trenera i onda da ide sve u krug
			// fmt.Println("AAAAAAAAAA")
			// var weightsIH [][]float64
			// for _, floatArray := range msg.WeightsIH {
			// 	var row []float64
			// 	for _, value := range floatArray.Column {
			// 		row = append(row, value)
			// 	}
			// 	weightsIH = append(weightsIH, row)
			// }
			
			// var weightsHH [][]float64
			// for _, floatArray := range msg.WeightsHH {
			// 	var row []float64
			// 	for _, value := range floatArray.Column {
			// 		row = append(row, value)
			// 	}
			// 	weightsHH = append(weightsHH, row)
			// }
	
			// var weightsHO [][]float64
			// for _, floatArray := range msg.WeightsHO {
			// 	var row []float64
			// 	for _, value := range floatArray.Column {
			// 		row = append(row, value)
			// 	}
			// 	weightsHO = append(weightsHO, row)
			// }
	
			// biasH := msg.BiasH
			// biasH2 := msg.BiasH2
			// biasO  := msg.BiasO

			// // uprosecavanje jednostavno
			// for i := range state.weightsIH {
			// 	for j := range state.weightsIH[i] {
			// 		state.weightsIH[i][j] = (state.weightsIH[i][j] + weightsIH[i][j])/2
			// 	}
			// }
			// for i := range state.weightsHH {
			// 	for j := range state.weightsIH[i] {
			// 		state.weightsHH[i][j] = (state.weightsHH[i][j] + weightsHH[i][j])/2
			// 	}
			// }
			// for i := range state.weightsIH {
			// 	for j := range state.weightsHO[i] {
			// 		state.weightsHO[i][j] = (state.weightsHO[i][j] + weightsHO[i][j])/2
			// 	}
			// }
			// for j := range state.biasH {
			// 	state.biasH[j] = (state.biasH[j] + biasH[j])/2
			// }
			// for j := range state.biasH2 {
			// 	state.biasH2[j] = (state.biasH2[j] + biasH2[j])/2
			// }
			// for j := range state.biasO {
			// 	state.biasO[j] = (state.biasO[j] + biasO[j])/2
			// }
			// do ovde sam zakomentarisao sad, dacu mu one iste tezine, u sustini samo ce nastaviti trening

			var twoDArrayProtoIH []*messages.FloatArray;
			for _, row := range state.weightsIH {
				floatArray := &messages.FloatArray{}
				for _, value := range row {
					floatArray.Column = append(floatArray.Column, value)
				}
				twoDArrayProtoIH = append(twoDArrayProtoIH, floatArray)
			}

			var twoDArrayProtoHH []*messages.FloatArray;
			for _, row := range state.weightsHH {
				floatArray := &messages.FloatArray{}
				for _, value := range row {
					floatArray.Column = append(floatArray.Column, value)
				}
				twoDArrayProtoHH = append(twoDArrayProtoHH, floatArray)
			}

			var twoDArrayProtoHO []*messages.FloatArray;
			for _, row := range state.weightsHO {
				floatArray := &messages.FloatArray{}
				for _, value := range row {
					floatArray.Column = append(floatArray.Column, value)
				}
				twoDArrayProtoHO = append(twoDArrayProtoHO, floatArray)
			}


			myMessage := &messages.AveragerWeightsMessage{
				WeightsIH:   twoDArrayProtoIH,
				WeightsHH: twoDArrayProtoHH,
				WeightsHO:   twoDArrayProtoHO,
				BiasH:       state.biasH,
				BiasH2:  state.biasH2,
				BiasO:       state.biasO,
			}
			
			
			context.Send(state.spawnedTrainerPID, myMessage)
			fmt.Println("AVERAGER: dobio sam nove tezine od interfejsa")
		}

	case *messages.TrainerWeightsMessage:
		fmt.Println("BRALEE")
		time.Sleep(time.Second * 2)
		if(state.hasWeights == false){
			fmt.Println("BRALEE1")
			state.hasWeights = true
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
			// for i := range state.weightsIH {
			// 	for j := range state.weightsIH[i] {
			// 		fmt.Print(state.weightsIH[i][j])
			// 	}
			// }
		}else{
			fmt.Println("BRALEE2")
			var weightsIH [][]float64
			for _, floatArray := range msg.WeightsIH {
				var row []float64
				for _, value := range floatArray.Column {
					row = append(row, value)
				}
				weightsIH = append(weightsIH, row)
			}
			
			var weightsHH [][]float64
			for _, floatArray := range msg.WeightsHH {
				var row []float64
				for _, value := range floatArray.Column {
					row = append(row, value)
				}
				weightsHH = append(weightsHH, row)
			}
	
			var weightsHO [][]float64
			for _, floatArray := range msg.WeightsHO {
				var row []float64
				for _, value := range floatArray.Column {
					row = append(row, value)
				}
				weightsHO = append(weightsHO, row)
			}
	
			biasH := msg.BiasH
			biasH2 := msg.BiasH2
			biasO  := msg.BiasO

			fmt.Println("BRALrrrrEE2")

			// uprosecavanje jednostavno
			for i := range state.weightsIH {
				for j := range state.weightsIH[i] {
					state.weightsIH[i][j] = (state.weightsIH[i][j] + weightsIH[i][j])/2
				}
			}
			for i := range state.weightsHH {
				for j := range state.weightsHH[i] {
					state.weightsHH[i][j] = (state.weightsHH[i][j] + weightsHH[i][j])/2
				}
			}
			for i := range state.weightsHO {
				for j := range state.weightsHO[i] {
					state.weightsHO[i][j] = (state.weightsHO[i][j] + weightsHO[i][j])/2
				}
			}
			for j := range state.biasH {
				state.biasH[j] = (state.biasH[j] + biasH[j])/2
			}
			for j := range state.biasH2 {
				state.biasH2[j] = (state.biasH2[j] + biasH2[j])/2
			}
			for j := range state.biasO {
				state.biasO[j] = (state.biasO[j] + biasO[j])/2
			}

			fmt.Println("BRALEE2")

			var twoDArrayProtoIH []*messages.FloatArray;
			for _, row := range state.weightsIH {
				floatArray := &messages.FloatArray{}
				for _, value := range row {
					floatArray.Column = append(floatArray.Column, value)
				}
				twoDArrayProtoIH = append(twoDArrayProtoIH, floatArray)
			}

			var twoDArrayProtoHH []*messages.FloatArray;
			for _, row := range state.weightsHH {
				floatArray := &messages.FloatArray{}
				for _, value := range row {
					floatArray.Column = append(floatArray.Column, value)
				}
				twoDArrayProtoHH = append(twoDArrayProtoHH, floatArray)
			}

			var twoDArrayProtoHO []*messages.FloatArray;
			for _, row := range state.weightsHO {
				floatArray := &messages.FloatArray{}
				for _, value := range row {
					floatArray.Column = append(floatArray.Column, value)
				}
				twoDArrayProtoHO = append(twoDArrayProtoHO, floatArray)
			}


			myMessage := &messages.AveragerWeightsMessage{
				WeightsIH:   twoDArrayProtoIH,
				WeightsHH: twoDArrayProtoHH,
				WeightsHO:   twoDArrayProtoHO,
				BiasH:       state.biasH,
				BiasH2:  state.biasH2,
				BiasO:       state.biasO,
			}
			context.Send(state.spawnedTrainerPID, myMessage)
			fmt.Println("JA SAM AVERAGER: " + msg.NizFloatova)
		}
	}

}
