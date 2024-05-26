package actors

import (
	"encoding/csv"
	"fmt"
	"os"
	"project/messages"
	"strconv"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type SpawnedInterfacePID struct{ PID *actor.PID }

type TrainerActor struct {
	count               int
	message             string
	spawnedInterfacePID *actor.PID
	spawnedAveragerPID  *actor.PID
	startState          bool
}

func Train(features [][]float64, labels []bool, context actor.Context, state *TrainerActor) []float64 {
	//TODO
	time.Sleep(time.Second * 3)
	context.Send(state.spawnedInterfacePID, &messages.TrainerWeightsMessage{NizFloatova: "Saljem ti tezine interfejsu moj"})
	context.Send(state.spawnedAveragerPID, &messages.TrainerWeightsMessage{NizFloatova: "Saljem ti tezine prosijaku moj"})
	return nil
}
func ReadDataset(filename string) ([][]float64, []bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	var features [][]float64
	var labels []bool

	for _, record := range records[1:] {
		var featureRow []float64
		for _, value := range record[:8] {
			f, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, nil, err
			}
			featureRow = append(featureRow, f)
		}

		label, err := strconv.Atoi(record[8])
		if err != nil {
			return nil, nil, err
		}
		// true,false konverzija, videcu hoce li trebati
		labels = append(labels, label == 1)

		features = append(features, featureRow)
	}

	return features, labels, nil
}

func (state *TrainerActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.Response:
		state.count++
		state.message = "Input" + string(state.count)
		fmt.Println(msg.GetSomeValue()+":", state.count)

	case *messages.Echo:
		fmt.Printf(msg.GetMessage() + "\n")

	case SpawnedAveragerPID:
		fmt.Println("TRENER dobavio PID Averagera: ", msg.PID)
		state.spawnedAveragerPID = msg.PID
		if state.spawnedInterfacePID != nil {
			state.startState = true
			Train(nil, nil, context, state)
		}
		fmt.Printf("Start stanje je: %v \n", state.startState)
		/*features, labels, err := ReadDataset("../dataset/Diabetes.csv")
		if err != nil {
			fmt.Println("Error reading dataset:", err)
			return
		}
		//ucitano
		fmt.Println("Features:", features)
		fmt.Println("Labels:", labels)
		*/
	case SpawnedInterfacePID:
		fmt.Println("TRENER dobavio PID Interfejsa: ", msg.PID)
		state.spawnedInterfacePID = msg.PID
		if state.spawnedAveragerPID != nil {
			state.startState = true
			Train(nil, nil, context, state)
		}
		fmt.Printf("Start stanje je: %v \n", state.startState)

	}

}
