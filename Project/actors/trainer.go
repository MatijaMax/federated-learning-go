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

type SpawnedInterfacePID struct{ PID *actor.PID }

type TrainerActor struct {
	count               int
	message             string
	spawnedInterfacePID *actor.PID
	spawnedAveragerPID  *actor.PID
	startState          bool
}


type NeuralNetwork struct {
    inputNodes  int
    hiddenNodes int
    outputNodes int
    weightsIH   [][]float64 // Weights between input and hidden layer
    weightsHO   [][]float64 // Weights between hidden and output layer
    biasH       []float64   // Bias for the hidden layer
    biasO       []float64   // Bias for the output layer
}


func NewNeuralNetwork(inputNodes, hiddenNodes, outputNodes int) *NeuralNetwork {
    nn := &NeuralNetwork{
        inputNodes:  inputNodes,
        hiddenNodes: hiddenNodes,
        outputNodes: outputNodes,
        weightsIH:   make([][]float64, inputNodes),
        weightsHO:   make([][]float64, hiddenNodes),
        biasH:       make([]float64, hiddenNodes),
        biasO:       make([]float64, outputNodes),
    }


    for i := range nn.weightsIH {
        nn.weightsIH[i] = make([]float64, hiddenNodes)
        for j := range nn.weightsIH[i] {
            nn.weightsIH[i][j] = randRange(-1, 1)
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
    for i := range nn.biasO {
        nn.biasO[i] = randRange(-1, 1)
    }

    return nn
}


func sigmoid(x float64) float64 {
    return 1 / (1 + math.Exp(-x))
}


func (nn *NeuralNetwork) FeedForward(inputData []float64) []float64 {

    hiddenInputs := make([]float64, nn.hiddenNodes)
    for i := 0; i < nn.hiddenNodes; i++ {
        sum := 0.0
        for j := 0; j < nn.inputNodes; j++ {
            sum += inputData[j] * nn.weightsIH[j][i]
        }
        hiddenInputs[i] = sum + nn.biasH[i] // Add bias
    }


    hiddenOutputs := make([]float64, nn.hiddenNodes)
    for i := 0; i < nn.hiddenNodes; i++ {
        hiddenOutputs[i] = sigmoid(hiddenInputs[i])
    }


    finalInputs := make([]float64, nn.outputNodes)
    for i := 0; i < nn.outputNodes; i++ {
        sum := 0.0
        for j := 0; j < nn.hiddenNodes; j++ {
            sum += hiddenOutputs[j] * nn.weightsHO[j][i]
        }
        finalInputs[i] = sum + nn.biasO[i] // Add bias
    }


    finalOutputs := make([]float64, nn.outputNodes)
    for i := 0; i < nn.outputNodes; i++ {
        finalOutputs[i] = sigmoid(finalInputs[i])
    }

    return finalOutputs
}


func randRange(min, max float64) float64 {
    return min + rand.Float64()*(max-min)
}


func (nn *NeuralNetwork) TrainNN(inputData [][]float64, targetData [][]float64, epochs int, learningRate float64) {
    for epoch := 0; epoch < epochs; epoch++ {
        for i := range inputData {

            inputs := inputData[i]
            targets := targetData[i]
            _, hiddenOutputs, _, finalOutputs := nn.forwardPass(inputs)


            outputErrors := make([]float64, nn.outputNodes)
            for j := 0; j < nn.outputNodes; j++ {
                outputErrors[j] = targets[j] - finalOutputs[j]
            }
            hiddenErrors := make([]float64, nn.hiddenNodes)
            for j := 0; j < nn.hiddenNodes; j++ {
                errorSum := 0.0
                for k := 0; k < nn.outputNodes; k++ {
                    errorSum += outputErrors[k] * nn.weightsHO[j][k]
                }
                hiddenErrors[j] = errorSum * hiddenOutputs[j] * (1 - hiddenOutputs[j])
            }


            for j := 0; j < nn.hiddenNodes; j++ {
                for k := 0; k < nn.outputNodes; k++ {
                    nn.weightsHO[j][k] += learningRate * outputErrors[k] * hiddenOutputs[j]
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
            for j := 0; j < nn.hiddenNodes; j++ {
                nn.biasH[j] += learningRate * hiddenErrors[j]
            }
        }
    }
}


func (nn *NeuralNetwork) forwardPass(inputs []float64) ([]float64, []float64, []float64, []float64) {
    hiddenInputs := make([]float64, nn.hiddenNodes)
    hiddenOutputs := make([]float64, nn.hiddenNodes)
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


    for i := 0; i < nn.outputNodes; i++ {
        sum := 0.0
        for j := 0; j < nn.hiddenNodes; j++ {
            sum += hiddenOutputs[j] * nn.weightsHO[j][i]
        }
        finalInputs[i] = sum + nn.biasO[i]
        finalOutputs[i] = sigmoid(finalInputs[i])
    }

    return hiddenInputs, hiddenOutputs, finalInputs, finalOutputs
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
	

    inputNodes := 6
    hiddenNodes := 4
    outputNodes := 1

    // Create a new neural network
    nn := NewNeuralNetwork(inputNodes, hiddenNodes, outputNodes)

    // Test the feedforward function with some input data
    input := featuresIn[9]
    fmt.Println("Input to the neural network:", input)
    output := nn.FeedForward(input)
    fmt.Println("Output from the neural network:", output)

    trainingData := featuresIn[:len(featuresIn)-20]
    targetData := labelsIn[:len(featuresIn)-20]

    
    nn.TrainNN(trainingData, targetData, 3, 0.1)

    // output = nn.FeedForward(input)
    // fmt.Println("Output from the neural network:", output)


    validationData := featuresIn[len(featuresIn)-20:]
    validationLabels := labelsIn[len(labelsIn)-20:]

    recall := nn.EvaluateRecall(validationData, validationLabels, 1.0)
    fmt.Printf("Validation Recall: %f\n", recall)

    
    nn.TrainNN(trainingData, targetData, 100, 0.1)

    // output = nn.FeedForward(input)
    // fmt.Println("Output from the neural network:", output)


    recall = nn.EvaluateRecall(validationData, validationLabels, 1.0)
    fmt.Printf("Validation Recall: %f\n", recall)

    // fmt.Println(len(validationData), len(validationLabels))
    // fmt.Println(len(featuresIn), len(labelsIn))


	context.Send(state.spawnedInterfacePID, &messages.TrainerWeightsMessage{NizFloatova: "Saljem ti tezine interfejsu moj"})
	context.Send(state.spawnedAveragerPID, &messages.TrainerWeightsMessage{NizFloatova: "Saljem ti tezine prosijaku moj"})

	return nil
}


func (nn *NeuralNetwork) EvaluateRecall(inputData [][]float64, targetData [][]float64, positiveLabel float64) float64 {
    truePositives := 0
    allPositives := 0

    for i := range inputData {
        inputs := inputData[i]
        targets := targetData[i]


        // Forward pass
        _, _, _, finalOutputs := nn.forwardPass(inputs)
        
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
        fmt.Println(predictedLabel, targets[0])
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

	case SpawnedAveragerPID:
		fmt.Println("TRENER dobavio PID Averagera: ", msg.PID)
		state.spawnedAveragerPID = msg.PID
		if state.spawnedInterfacePID != nil {
			state.startState = true
			Train(context, state)
		}
		fmt.Printf("Start stanje je: %v \n", state.startState)
		
	case SpawnedInterfacePID:
		fmt.Println("TRENER dobavio PID Interfejsa: ", msg.PID)
		state.spawnedInterfacePID = msg.PID
		if state.spawnedAveragerPID != nil {
			state.startState = true
			Train(context, state)
		}
		fmt.Printf("Start stanje je: %v \n", state.startState)

	}

}
