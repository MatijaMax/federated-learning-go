package actors

import (
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"project/messages"
	"strconv"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type TrainerActor struct {
	count               int
	message             string
	spawnedInterfacePID *actor.PID
	spawnedAveragerPID  *actor.PID
	startState          bool
}

type NeuralNetwork struct {
	inputNodes   int
	hiddenNodes  int
	hiddenNodes2 int
	outputNodes  int
	weightsIH    [][]float64 // Weights between input and hidden layer
	weightsHH    [][]float64 // Weights between input and hidden layer
	weightsHO    [][]float64 // Weights between hidden and output layer
	biasH        []float64   // Bias for the hidden layer
	biasH2       []float64   // Bias for the hidden layer
	biasO        []float64   // Bias for the output layer
}

func NewNeuralNetwork(inputNodes, hiddenNodes, hiddenNodes2, outputNodes int) *NeuralNetwork {
	nn := &NeuralNetwork{
		inputNodes:   inputNodes,
		hiddenNodes:  hiddenNodes,
		hiddenNodes2: hiddenNodes2,
		outputNodes:  outputNodes,
		weightsIH:    make([][]float64, inputNodes),
		weightsHH:    make([][]float64, hiddenNodes),
		weightsHO:    make([][]float64, hiddenNodes2),
		biasH:        make([]float64, hiddenNodes),
		biasH2:       make([]float64, hiddenNodes2),
		biasO:        make([]float64, outputNodes),
	}

	for i := range nn.weightsIH {
		nn.weightsIH[i] = make([]float64, hiddenNodes)
		for j := range nn.weightsIH[i] {
			nn.weightsIH[i][j] = randRange(-1, 1)
		}
	}
	for i := range nn.weightsHH {
		nn.weightsHH[i] = make([]float64, hiddenNodes2)
		for j := range nn.weightsHH[i] {
			nn.weightsHH[i][j] = randRange(-1, 1)
		}
	}
	for i := range nn.weightsHO {
		nn.weightsHO[i] = make([]float64, outputNodes)
		for j := range nn.weightsHO[i] {
			nn.weightsHO[i][j] = randRange(-1, 1)
		}
	}

	for i := range nn.biasH {
		nn.biasH[i] = randRange(-1, 1)
	}
	for i := range nn.biasH2 {
		nn.biasH2[i] = randRange(-1, 1)
	}
	for i := range nn.biasO {
		nn.biasO[i] = randRange(-1, 1)
	}

	return nn
}

func NewNeuralNetworkWithWeights(inputNodes, hiddenNodes, hiddenNodes2, outputNodes int, weightsIH, weightsHH, weightsHO [][]float64, biasH, biasH2, biasO []float64) *NeuralNetwork {
	nn := &NeuralNetwork{
		inputNodes:   inputNodes,
		hiddenNodes:  hiddenNodes,
		hiddenNodes2: hiddenNodes2,
		outputNodes:  outputNodes,
		weightsIH:    make([][]float64, inputNodes),
		weightsHH:    make([][]float64, hiddenNodes),
		weightsHO:    make([][]float64, hiddenNodes2),
		biasH:        make([]float64, hiddenNodes),
		biasH2:       make([]float64, hiddenNodes2),
		biasO:        make([]float64, outputNodes),
	}

    fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA1")
    fmt.Print(weightsIH)
	for i := range nn.weightsIH {
		nn.weightsIH[i] = make([]float64, hiddenNodes)
		for j := range nn.weightsIH[i] {
			nn.weightsIH[i][j] = randRange(-1, 1)
		}
	}
    fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA2")
	for i := range nn.weightsHH {
		nn.weightsHH[i] = make([]float64, hiddenNodes2)
		for j := range nn.weightsHH[i] {
			nn.weightsHH[i][j] = weightsHH[i][j]
		}
	}
    fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA3")
	for i := range nn.weightsHO {
		nn.weightsHO[i] = make([]float64, outputNodes)
		for j := range nn.weightsHO[i] {
			nn.weightsHO[i][j] = weightsHO[i][j]
		}
	}
    fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA4")
	for i := range nn.biasH {
		nn.biasH[i] = biasH[i]
	}
	for i := range nn.biasH2 {
		nn.biasH2[i] = biasH2[i]
	}
	for i := range nn.biasO {
		nn.biasO[i] = biasO[i]
	}
fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	return nn
}

func sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

func randRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func (nn *NeuralNetwork) TrainNN(inputData [][]float64, targetData [][]float64, epochs int, learningRate float64) {
	for epoch := 0; epoch < epochs; epoch++ {
		for i := range inputData {

			inputs := inputData[i]
			targets := targetData[i]
			_, hiddenOutputs, _, finalOutputs, _, hiddenOutputs2 := nn.forwardPass(inputs)

			outputErrors := make([]float64, nn.outputNodes)
			for j := 0; j < nn.outputNodes; j++ {
				outputErrors[j] = targets[j] - finalOutputs[j]
			}
			hiddenErrors2 := make([]float64, nn.hiddenNodes2)
			for j := 0; j < nn.hiddenNodes2; j++ {
				errorSum := 0.0
				for k := 0; k < nn.outputNodes; k++ {
					errorSum += outputErrors[k] * nn.weightsHO[j][k]
				}
				hiddenErrors2[j] = errorSum * hiddenOutputs2[j] * (1 - hiddenOutputs2[j])
			}
			hiddenErrors := make([]float64, nn.hiddenNodes)
			for j := 0; j < nn.hiddenNodes; j++ {
				errorSum := 0.0
				for k := range hiddenErrors2 {
					errorSum += hiddenErrors2[k] * nn.weightsHH[j][k]
				}
				hiddenErrors[j] = errorSum * hiddenOutputs[j] * (1 - hiddenOutputs[j])
			}

			for j := 0; j < nn.hiddenNodes2; j++ {
				for k := 0; k < nn.outputNodes; k++ {
					nn.weightsHO[j][k] += learningRate * outputErrors[k] * hiddenOutputs2[j]
				}
			}
			for j := 0; j < nn.hiddenNodes; j++ {
				for k := 0; k < nn.hiddenNodes2; k++ {
					nn.weightsHH[j][k] += learningRate * hiddenErrors2[k] * hiddenOutputs[j]
				}
			}
			for j := 0; j < nn.inputNodes; j++ {
				for k := 0; k < nn.hiddenNodes; k++ {
					nn.weightsIH[j][k] += learningRate * hiddenErrors[k] * inputs[j]
				}
			}

			for j := 0; j < nn.outputNodes; j++ {
				nn.biasO[j] += learningRate * outputErrors[j]
			}
			for j := 0; j < nn.hiddenNodes2; j++ {
				nn.biasH2[j] += learningRate * hiddenErrors2[j]
			}
			for j := 0; j < nn.hiddenNodes; j++ {
				nn.biasH[j] += learningRate * hiddenErrors[j]
			}
		}
	}
}

func (nn *NeuralNetwork) forwardPass(inputs []float64) ([]float64, []float64, []float64, []float64, []float64, []float64) {
	hiddenInputs := make([]float64, nn.hiddenNodes)
	hiddenOutputs := make([]float64, nn.hiddenNodes)
	hiddenInputs2 := make([]float64, nn.hiddenNodes2)
	hiddenOutputs2 := make([]float64, nn.hiddenNodes2)
	finalInputs := make([]float64, nn.outputNodes)
	finalOutputs := make([]float64, nn.outputNodes)

	for i := 0; i < nn.hiddenNodes; i++ {
		sum := 0.0
		for j := 0; j < nn.inputNodes; j++ {
			sum += inputs[j] * nn.weightsIH[j][i]
		}
		hiddenInputs[i] = sum + nn.biasH[i]
		hiddenOutputs[i] = sigmoid(hiddenInputs[i])
	}

	for i := 0; i < nn.hiddenNodes2; i++ {
		sum := 0.0
		for j := 0; j < nn.hiddenNodes; j++ {
			sum += hiddenOutputs[j] * nn.weightsHH[j][i]
		}
		hiddenInputs2[i] = sum + nn.biasH2[i]
		hiddenOutputs2[i] = sigmoid(hiddenInputs2[i])
	}

	for i := 0; i < nn.outputNodes; i++ {
		sum := 0.0
		for j := 0; j < nn.hiddenNodes2; j++ {
			sum += hiddenOutputs2[j] * nn.weightsHO[j][i]
		}
		finalInputs[i] = sum + nn.biasO[i]
		finalOutputs[i] = sigmoid(finalInputs[i])
	}

	return hiddenInputs, hiddenOutputs, finalInputs, finalOutputs, hiddenInputs2, hiddenOutputs2
}

func Train(context actor.Context, state *TrainerActor) []float64 {
	time.Sleep(time.Second * 2)

	featuresIn, labelsIn, err := ReadDataset("../dataset/Diabetes.csv")
	if err != nil {
		fmt.Println("Error reading dataset:", err)
		// return
	}

	// featuresIn = dropColumn(featuresIn, 3)
	// featuresIn = dropColumn(featuresIn, 3)
	//ucitano
	// fmt.Println("Features:", featuresIn)
	// fmt.Println("Labels:", labelsIn)

	inputNodes := 8
	hiddenNodes := 8
	hiddenNodes2 := 4
	outputNodes := 1

	nn := NewNeuralNetwork(inputNodes, hiddenNodes, hiddenNodes2, outputNodes)


	trainingData := featuresIn[:len(featuresIn)-20]
	targetData := labelsIn[:len(featuresIn)-20]

	nn.TrainNN(trainingData, targetData, 10, 0.04)


	validationData := featuresIn[len(featuresIn)-20:]
	validationLabels := labelsIn[len(labelsIn)-20:]

	recall := nn.EvaluateRecall(validationData, validationLabels, 1.0)
	fmt.Printf("Validation Recall nakon 10 epoha: %f\n", recall)


	var twoDArrayProtoIH []*messages.FloatArray
	for _, row := range nn.weightsIH {
		floatArray := &messages.FloatArray{}
		for _, value := range row {
			floatArray.Column = append(floatArray.Column, value)
		}
		twoDArrayProtoIH = append(twoDArrayProtoIH, floatArray)
	}

	var twoDArrayProtoHH []*messages.FloatArray
	for _, row := range nn.weightsHH {
		floatArray := &messages.FloatArray{}
		for _, value := range row {
			floatArray.Column = append(floatArray.Column, value)
		}
		twoDArrayProtoHH = append(twoDArrayProtoHH, floatArray)
	}

	var twoDArrayProtoHO []*messages.FloatArray
	for _, row := range nn.weightsHO {
		floatArray := &messages.FloatArray{}
		for _, value := range row {
			floatArray.Column = append(floatArray.Column, value)
		}
		twoDArrayProtoHO = append(twoDArrayProtoHO, floatArray)
	}

	myMessage := &messages.TrainerWeightsMessage{
		NizFloatova: "Saljem ti tezine",
		WeightsIH:   twoDArrayProtoIH,
		WeightsHH:   twoDArrayProtoHH,
		WeightsHO:   twoDArrayProtoHO,
		BiasH:       nn.biasH,
		BiasH2:      nn.biasH2,
		BiasO:       nn.biasO,
	}

	context.Send(state.spawnedInterfacePID, myMessage)
	context.Send(state.spawnedAveragerPID, myMessage)

	return nil
}

func TrainAgain(context actor.Context, state *TrainerActor, weightsIH, weightsHH, weightsHO [][]float64, biasH, biasH2, biasO []float64) []float64 {
	time.Sleep(time.Second * 2)

	fmt.Println("Opet treniram")
	featuresIn, labelsIn, err := ReadDataset("../dataset/Diabetes.csv")
	if err != nil {
		fmt.Println("Error reading dataset:", err)
	}

	inputNodes := 8
	hiddenNodes := 8
	hiddenNodes2 := 4
	outputNodes := 1

	nn := NewNeuralNetworkWithWeights(inputNodes, hiddenNodes, hiddenNodes2, outputNodes, weightsIH, weightsHH, weightsHO, biasH, biasH2, biasO)

	trainingData := featuresIn[:len(featuresIn)-20]
	targetData := labelsIn[:len(featuresIn)-20]

	nn.TrainNN(trainingData, targetData, 10, 0.04)

	validationData := featuresIn[len(featuresIn)-20:]
	validationLabels := labelsIn[len(labelsIn)-20:]

	recall := nn.EvaluateRecall(validationData, validationLabels, 1.0)
	fmt.Printf("Validation Recall nakon 10 epoha: %f\n", recall)

	var twoDArrayProtoIH []*messages.FloatArray
	for _, row := range nn.weightsIH {
		floatArray := &messages.FloatArray{}
		for _, value := range row {
			floatArray.Column = append(floatArray.Column, value)
		}
		twoDArrayProtoIH = append(twoDArrayProtoIH, floatArray)
	}

	var twoDArrayProtoHH []*messages.FloatArray
	for _, row := range nn.weightsHH {
		floatArray := &messages.FloatArray{}
		for _, value := range row {
			floatArray.Column = append(floatArray.Column, value)
		}
		twoDArrayProtoHH = append(twoDArrayProtoHH, floatArray)
	}

	var twoDArrayProtoHO []*messages.FloatArray
	for _, row := range nn.weightsHO {
		floatArray := &messages.FloatArray{}
		for _, value := range row {
			floatArray.Column = append(floatArray.Column, value)
		}
		twoDArrayProtoHO = append(twoDArrayProtoHO, floatArray)
	}

	myMessage := &messages.TrainerWeightsMessage{
		NizFloatova: "Saljem ti tezine",
		WeightsIH:   twoDArrayProtoIH,
		WeightsHH:   twoDArrayProtoHH,
		WeightsHO:   twoDArrayProtoHO,
		BiasH:       nn.biasH,
		BiasH2:      nn.biasH2,
		BiasO:       nn.biasO,
	}

	context.Send(state.spawnedInterfacePID, myMessage)
	context.Send(state.spawnedAveragerPID, myMessage)

	return nil
}

func (nn *NeuralNetwork) EvaluateRecall(inputData [][]float64, targetData [][]float64, positiveLabel float64) float64 {
	truePositives := 0
	allPositives := 0

	for i := range inputData {
		inputs := inputData[i]
		targets := targetData[i]

		// Forward pass
		_, _, _, finalOutputs, _, _ := nn.forwardPass(inputs)

		// fmt.Printf("%v\n%v\n%v\n",inputs, targets, finalOutputs)

		predictedLabel := 0.0
		if finalOutputs[0] > 0.5 {
			predictedLabel = 1.0
		}

		if predictedLabel == positiveLabel {
			truePositives++
			// fmt.Println("Dobar bas")
		}
		if predictedLabel == targets[0] {
			allPositives++
			// fmt.Println("Dobar")
		}
		fmt.Println(predictedLabel, targets[0], finalOutputs[0])
	}

	recall := 0.0
	if allPositives > 0 {
		recall = float64(truePositives) / float64(allPositives)
		print(truePositives, allPositives)
	}
	return recall
}

func normalizeTo01(data [][]float64) [][]float64 {
	numFeatures := len(data[0])
	minValues := make([]float64, numFeatures)
	maxValues := make([]float64, numFeatures)

	for i := range data {
		for j := range data[i] {
			if i == 0 {
				minValues[j] = data[i][j]
				maxValues[j] = data[i][j]
			} else {
				if data[i][j] < minValues[j] {
					minValues[j] = data[i][j]
				}
				if data[i][j] > maxValues[j] {
					maxValues[j] = data[i][j]
				}
			}
		}
	}

	for i := range data {
		for j := range data[i] {
			data[i][j] = (data[i][j] - minValues[j]) / (maxValues[j] - minValues[j])
		}
	}

	return data
}

func dropColumn(data [][]float64, columnIndex int) [][]float64 {
	var newData [][]float64

	for i := range data {
		var newRow []float64
		for j := range data[i] {
			if j != columnIndex {
				newRow = append(newRow, data[i][j])
			}
		}
		newData = append(newData, newRow)
	}

	return newData
}

func ReadDataset(filename string) ([][]float64, [][]float64, error) {
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
	var labels [][]float64

	for _, record := range records[1:] {
		var featureRow []float64
		for _, value := range record[:8] {
			f, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, nil, err
			}
			featureRow = append(featureRow, f)
		}

		label, err := strconv.ParseFloat(record[8], 64)
		if err != nil {
			return nil, nil, err
		}
		var labelRow []float64
		labelRow = append(labelRow, label)
		labels = append(labels, labelRow)

		features = append(features, featureRow)
	}

	normalizedFeatures := normalizeTo01(features)
	// fmt.Print(normalizedFeatures)
	return normalizedFeatures, labels, nil
}

func (state *TrainerActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.Response:
		state.count++
		state.message = "Input" + string(state.count)
		fmt.Println(msg.GetSomeValue()+":", state.count)

	case *messages.Echo:
		fmt.Printf(msg.GetMessage() + "\n")

	case *messages.SpawnedAveragerPID:
		fmt.Println("TRENER dobavio PID Averagera: ", msg.ThePid)
		state.spawnedAveragerPID = msg.ThePid
		if state.spawnedInterfacePID != nil {
			state.startState = true
			Train(context, state)
		}
		fmt.Printf("Start stanje je: %v \n", state.startState)

	case *messages.SpawnedInterfacePID:
		fmt.Println("TRENER dobavio PID Interfejsa: ", msg.ThePid)
		state.spawnedInterfacePID = msg.ThePid
		if state.spawnedAveragerPID != nil {
			state.startState = true
			Train(context, state)
		}
		fmt.Printf("Start stanje je: %v \n", state.startState)

	case *messages.AveragerWeightsMessage:
		fmt.Println("ADSSSSSSSSSSSSSSS")
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
		TrainAgain(context, state, weightsIH, weightsHH, weightsHO, msg.BiasH, msg.BiasH2, msg.BiasO)
	}

}
