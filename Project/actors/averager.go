package actors

import (
	"fmt"
	"project/messages"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type WeightsBiases struct {
	weightsIH [][]float64 // Weights between input and hidden layer
	weightsHH [][]float64 // Weights between input and hidden layer
	weightsHO [][]float64 // Weights between hidden and output layer
	biasH     []float64   // Bias for the hidden layer
	biasH2    []float64   // Bias for the hidden layer
	biasO     []float64   // Bias for the output layer
}

type SpawnedTrainerPID struct{ PID *actor.PID }

type AveragerActor struct {
	count             int
	message           string
	spawnedTrainerPID *actor.PID
	queueTrainersWB            []WeightsBiases // First dynamic queue
	queueInterfacesWB            []WeightsBiases // Second dynamic queue
	hasWeights bool
}

func (state *AveragerActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.Response:
		state.count++
		state.message = "Input" + string(state.count)
		fmt.Println(msg.GetSomeValue()+":", state.count)
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
			var weightsBiases WeightsBiases;
			for _, floatArray := range msg.WeightsIH {
				var row []float64
				for _, value := range floatArray.Column {
					row = append(row, value)
				}
				weightsBiases.weightsIH = append(weightsBiases.weightsIH, row)
			}
			
			for _, floatArray := range msg.WeightsHH {
				var row []float64
				for _, value := range floatArray.Column {
					row = append(row, value)
				}
				weightsBiases.weightsHH = append(weightsBiases.weightsHH, row)
			}
	
			for _, floatArray := range msg.WeightsHO {
				var row []float64
				for _, value := range floatArray.Column {
					row = append(row, value)
				}
				weightsBiases.weightsHO = append(weightsBiases.weightsHO, row)
			}
	
			weightsBiases.biasH = msg.BiasH
			weightsBiases.biasH2 = msg.BiasH2
			weightsBiases.biasO = msg.BiasO
			state.queueInterfacesWB = append(state.queueInterfacesWB, weightsBiases)
			
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

			weightsBiasesRES, hasRes := state.AverageFirstN();
			if(hasRes == true){
				var twoDArrayProtoIH []*messages.FloatArray;
				for _, row := range weightsBiasesRES.weightsIH {
					floatArray := &messages.FloatArray{}
					for _, value := range row {
						floatArray.Column = append(floatArray.Column, value)
					}
					twoDArrayProtoIH = append(twoDArrayProtoIH, floatArray)
				}

				var twoDArrayProtoHH []*messages.FloatArray;
				for _, row := range weightsBiasesRES.weightsHH {
					floatArray := &messages.FloatArray{}
					for _, value := range row {
						floatArray.Column = append(floatArray.Column, value)
					}
					twoDArrayProtoHH = append(twoDArrayProtoHH, floatArray)
				}

				var twoDArrayProtoHO []*messages.FloatArray;
				for _, row := range weightsBiasesRES.weightsHO {
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
					BiasH:       weightsBiasesRES.biasH,
					BiasH2:  weightsBiasesRES.biasH2,
					BiasO:       weightsBiasesRES.biasO,
				}
				context.Send(state.spawnedTrainerPID, myMessage)
				fmt.Println("AVERAGER: dobio sam nove tezine od interfejsa")
			}else{
				fmt.Println("Nema sta da prosecim")
				//sad samo ovako al menjacu
				weightsBiasesRES := state.queueTrainersWB[0]
				var twoDArrayProtoIH []*messages.FloatArray;
				for _, row := range weightsBiasesRES.weightsIH {
					floatArray := &messages.FloatArray{}
					for _, value := range row {
						floatArray.Column = append(floatArray.Column, value)
					}
					twoDArrayProtoIH = append(twoDArrayProtoIH, floatArray)
				}

				var twoDArrayProtoHH []*messages.FloatArray;
				for _, row := range weightsBiasesRES.weightsHH {
					floatArray := &messages.FloatArray{}
					for _, value := range row {
						floatArray.Column = append(floatArray.Column, value)
					}
					twoDArrayProtoHH = append(twoDArrayProtoHH, floatArray)
				}

				var twoDArrayProtoHO []*messages.FloatArray;
				for _, row := range weightsBiasesRES.weightsHO {
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
					BiasH:       weightsBiasesRES.biasH,
					BiasH2:  weightsBiasesRES.biasH2,
					BiasO:       weightsBiasesRES.biasO,
				}
				context.Send(state.spawnedTrainerPID, myMessage)
			}
		}

	case *messages.TrainerWeightsMessage:
		fmt.Println("BRALEE")
		time.Sleep(time.Second * 2)
		if(state.hasWeights == false){
			fmt.Println("BRALEE1")
			state.hasWeights = true
			var weightsBiases WeightsBiases;
			for _, floatArray := range msg.WeightsIH {
				var row []float64
				for _, value := range floatArray.Column {
					row = append(row, value)
				}
				weightsBiases.weightsIH = append(weightsBiases.weightsIH, row)
			}
			
			for _, floatArray := range msg.WeightsHH {
				var row []float64
				for _, value := range floatArray.Column {
					row = append(row, value)
				}
				weightsBiases.weightsHH = append(weightsBiases.weightsHH, row)
			}
	
			for _, floatArray := range msg.WeightsHO {
				var row []float64
				for _, value := range floatArray.Column {
					row = append(row, value)
				}
				weightsBiases.weightsHO = append(weightsBiases.weightsHO, row)
			}
	
			weightsBiases.biasH = msg.BiasH
			weightsBiases.biasH2 = msg.BiasH2
			weightsBiases.biasO = msg.BiasO
			state.queueTrainersWB = append(state.queueTrainersWB, weightsBiases)
			fmt.Println("JA SAM AVERAGER: " + msg.NizFloatova)
			// for i := range state.weightsIH {
			// 	for j := range state.weightsIH[i] {
			// 		fmt.Print(state.weightsIH[i][j])
			// 	}
			// }
		}else{
			fmt.Println("BRALEE2")
			var weightsBiases WeightsBiases;
			for _, floatArray := range msg.WeightsIH {
				var row []float64
				for _, value := range floatArray.Column {
					row = append(row, value)
				}
				weightsBiases.weightsIH = append(weightsBiases.weightsIH, row)
			}
			
			for _, floatArray := range msg.WeightsHH {
				var row []float64
				for _, value := range floatArray.Column {
					row = append(row, value)
				}
				weightsBiases.weightsHH = append(weightsBiases.weightsHH, row)
			}
	
			for _, floatArray := range msg.WeightsHO {
				var row []float64
				for _, value := range floatArray.Column {
					row = append(row, value)
				}
				weightsBiases.weightsHO = append(weightsBiases.weightsHO, row)
			}
	
			weightsBiases.biasH = msg.BiasH
			weightsBiases.biasH2 = msg.BiasH2
			weightsBiases.biasO = msg.BiasO

			state.queueTrainersWB = append(state.queueTrainersWB, weightsBiases)
			fmt.Println("BRALrrrrEE2")

			

			fmt.Println(len(state.queueTrainersWB))
			weightsBiasesRES, hasRes := state.AverageFirstN();
			if(hasRes == true){
				var twoDArrayProtoIH []*messages.FloatArray;
				for _, row := range weightsBiasesRES.weightsIH {
					floatArray := &messages.FloatArray{}
					for _, value := range row {
						floatArray.Column = append(floatArray.Column, value)
					}
					twoDArrayProtoIH = append(twoDArrayProtoIH, floatArray)
				}

				var twoDArrayProtoHH []*messages.FloatArray;
				for _, row := range weightsBiasesRES.weightsHH {
					floatArray := &messages.FloatArray{}
					for _, value := range row {
						floatArray.Column = append(floatArray.Column, value)
					}
					twoDArrayProtoHH = append(twoDArrayProtoHH, floatArray)
				}

				var twoDArrayProtoHO []*messages.FloatArray;
				for _, row := range weightsBiasesRES.weightsHO {
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
					BiasH:       weightsBiasesRES.biasH,
					BiasH2:  weightsBiasesRES.biasH2,
					BiasO:       weightsBiasesRES.biasO,
				}
				context.Send(state.spawnedTrainerPID, myMessage)
				fmt.Println("JA SAM AVERAGER: " + msg.NizFloatova)
			}else{
				fmt.Println("Nema sta da prosecim")
			}

			
		}
	}

}

func addWeightsBiases(wb1, wb2 WeightsBiases) WeightsBiases {
	result := WeightsBiases{
		weightsIH: make([][]float64, len(wb1.weightsIH)),
		weightsHH: make([][]float64, len(wb1.weightsHH)),
		weightsHO: make([][]float64, len(wb1.weightsHO)),
		biasH:     make([]float64, len(wb1.biasH)),
		biasH2:    make([]float64, len(wb1.biasH2)),
		biasO:     make([]float64, len(wb1.biasO)),
	}

	for i := range wb1.weightsIH {
		result.weightsIH[i] = make([]float64, len(wb1.weightsIH[i]))
		for j := range wb1.weightsIH[i] {
			result.weightsIH[i][j] = wb1.weightsIH[i][j] + wb2.weightsIH[i][j]
		}
	}

	for i := range wb1.weightsHH {
		result.weightsHH[i] = make([]float64, len(wb1.weightsHH[i]))
		for j := range wb1.weightsHH[i] {
			result.weightsHH[i][j] = wb1.weightsHH[i][j] + wb2.weightsHH[i][j]
		}
	}

	for i := range wb1.weightsHO {
		result.weightsHO[i] = make([]float64, len(wb1.weightsHO[i]))
		for j := range wb1.weightsHO[i] {
			result.weightsHO[i][j] = wb1.weightsHO[i][j] + wb2.weightsHO[i][j]
		}
	}

	for i := range wb1.biasH {
		result.biasH[i] = wb1.biasH[i] + wb2.biasH[i]
	}

	for i := range wb1.biasH2 {
		result.biasH2[i] = wb1.biasH2[i] + wb2.biasH2[i]
	}

	for i := range wb1.biasO {
		result.biasO[i] = wb1.biasO[i] + wb2.biasO[i]
	}

	return result
}

func (a *AveragerActor) AverageFirstN() (WeightsBiases, bool) {
	n := 0
	if len(a.queueTrainersWB) <  len(a.queueInterfacesWB) {
		n = len(a.queueTrainersWB)
	}else{
		n = len(a.queueInterfacesWB)
	}
	if(n ==0){
		return WeightsBiases{}, false
	}

	sum := WeightsBiases{
		weightsIH: make([][]float64, len(a.queueTrainersWB[0].weightsIH)),
		weightsHH: make([][]float64, len(a.queueTrainersWB[0].weightsHH)),
		weightsHO: make([][]float64, len(a.queueTrainersWB[0].weightsHO)),
		biasH:     make([]float64, len(a.queueTrainersWB[0].biasH)),
		biasH2:    make([]float64, len(a.queueTrainersWB[0].biasH2)),
		biasO:     make([]float64, len(a.queueTrainersWB[0].biasO)),
	}

	for i := 0; i < n; i++ {
		sum = addWeightsBiases(sum, a.queueTrainersWB[i])
		sum = addWeightsBiases(sum, a.queueInterfacesWB[i])
	}

	average := WeightsBiases{
		weightsIH: make([][]float64, len(sum.weightsIH)),
		weightsHH: make([][]float64, len(sum.weightsHH)),
		weightsHO: make([][]float64, len(sum.weightsHO)),
		biasH:     make([]float64, len(sum.biasH)),
		biasH2:    make([]float64, len(sum.biasH2)),
		biasO:     make([]float64, len(sum.biasO)),
	}

	for i := range sum.weightsIH {
		average.weightsIH[i] = make([]float64, len(sum.weightsIH[i]))
		for j := range sum.weightsIH[i] {
			average.weightsIH[i][j] = sum.weightsIH[i][j] / float64(2*n)
		}
	}

	for i := range sum.weightsHH {
		average.weightsHH[i] = make([]float64, len(sum.weightsHH[i]))
		for j := range sum.weightsHH[i] {
			average.weightsHH[i][j] = sum.weightsHH[i][j] / float64(2*n)
		}
	}

	for i := range sum.weightsHO {
		average.weightsHO[i] = make([]float64, len(sum.weightsHO[i]))
		for j := range sum.weightsHO[i] {
			average.weightsHO[i][j] = sum.weightsHO[i][j] / float64(2*n)
		}
	}

	for i := range sum.biasH {
		average.biasH[i] = sum.biasH[i] / float64(2*n)
	}

	for i := range sum.biasH2 {
		average.biasH2[i] = sum.biasH2[i] / float64(2*n)
	}

	for i := range sum.biasO {
		average.biasO[i] = sum.biasO[i] / float64(2*n)
	}

	a.queueTrainersWB = a.queueTrainersWB[n:]
	a.queueInterfacesWB = a.queueInterfacesWB[n:]
	return average, true
}