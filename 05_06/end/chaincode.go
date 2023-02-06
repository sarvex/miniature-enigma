/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type Batch struct {
	Id			string		`json:"id"`
	Weight 		int			`json:"weight"`
	Grower		string		`json:"grower"`
	Readings	[]Reading 	`json:"readings"`
	Grade		string		`json:"grade"`
	Producer	string		`json:"producer"`
}

type Reading struct {
	Type	int		`json:"type"`
	Value	string	`json:"value"`
	Member 	string	`json:"member"`
}

type ReadingType int

const (
	Soil_Quality ReadingType=iota
	Pesticide
	Acidity
)

// Chaincode is the definition of the chaincode structure.
type BatchManager struct {
}

// Init is called when the chaincode is instantiated by the blockchain network.
func (cc *BatchManager) Init(stub shim.ChaincodeStubInterface) sc.Response {
	fcn, params := stub.GetFunctionAndParameters()
	fmt.Println("Init()", fcn, params)
	return shim.Success(nil)
}

// Invoke is called as a result of an application request to run the chaincode.
func (cc *BatchManager) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	fcnName, _ := stub.GetFunctionAndParameters()

	if callFunction, ok := functions[fcnName]; ok{
		fmt.Printf("Invoking %s\n", fcnName)
		return callFunction(stub)
	}

	return shim.Error(fmt.Sprintf("Unknown function %s", fcnName))
}

var functions = map[string]func(stub shim.ChaincodeStubInterface) sc.Response{
	"addReadings":		addReadings,
	"createBatch":		createBatch,
	"determineGrade":	determineGrade,
	"queryBatch":		queryBatch,
}

/*
* type - soil quality, pesticide, acidity
* batchId - the individual collecting the readings
* weight - the organization responsible for t
* grower - the batch on the which the reading is being conducted
*/
func createBatch(stub shim.ChaincodeStubInterface) sc.Response{
	_, parameters := stub.GetFunctionAndParameters()
	if len(parameters) != 3 {
		return shim.Error("Wrong number of arguments supplied. Expected: <batchId> <weight> <grower>")
	}

	//if batch exists throw an error
	batchId := parameters[0]
	batchJSON, err := stub.GetState(batchId)
	if batchJSON != nil {
		return shim.Error(fmt.Sprintf("Batch already exists for id: %s", string(parameters[0])))
	}

	weight, err := strconv.Atoi(parameters[1])
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to parse weight: %s", string(parameters[1])))
	}

	grower := string(parameters[2])
	newBatch := Batch{Id: batchId, Weight: weight, Grower: grower}
	
	batchasBytes, err := json.Marshal(newBatch)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(batchId, batchasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("batch created"))
}

/**
* add readings transaction
* readingType - Integer representing soil quality, pesticide, acidity
* value - reading value
* member - who took the reading
* batchId - the batch on the which the reading is being conducted
*/
func addReadings(stub shim.ChaincodeStubInterface) sc.Response {
	_, parameters := stub.GetFunctionAndParameters()
	if len(parameters) != 4 {
		return shim.Error("Wrong number of arguments supplied. Expected: <batchID> <readingType> <value> <member>")
	}

	batchId := parameters[0]
	batchJSON, err := stub.GetState(batchId)
	if err != nil {
		return shim.Error(err.Error())
	}
	if batchJSON == nil {
		return shim.Error(fmt.Sprintf("Unable to find batch %s", parameters[0]))
	}

	var selectedBatch Batch
	err = json.Unmarshal(batchJSON, &selectedBatch)
	if err != nil {
		return shim.Error(err.Error())
	}

	readingType, err := strconv.Atoi(parameters[1])
	if err != nil {
		return shim.Error(err.Error())
	}

	value := string(parameters[2])
	member := string(parameters[3])
	newReading := Reading{Type: readingType, Value: value, Member: member} 

	readingSet := &selectedBatch.Readings
	*readingSet = append(*readingSet, newReading)

	fmt.Printf("selectedBatch output: %v\n", selectedBatch)

	selectedBatchJSONasBytes, err := json.Marshal(selectedBatch)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(batchId, selectedBatchJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("reading added"))
}

/*
* type - soil quality, pesticide, acidity
* batchId - the individual collecting the readings
* producer - the oil producer
*/
func determineGrade(stub shim.ChaincodeStubInterface) sc.Response{
	_, parameters := stub.GetFunctionAndParameters()

	if len(parameters) != 2 {
		return shim.Error("Wrong number of arguments supplied. Expected: <batchId> <producer>")
	}

	batchId := parameters[0]
	batchJSON, err := stub.GetState(batchId)
	if err != nil {
		return shim.Error(err.Error())
	}
	if batchJSON == nil {
		return shim.Error(fmt.Sprintf("Unable to find batch %s", parameters[0]))
	}

	producer := parameters[1]

	var selectedBatch Batch
	err = json.Unmarshal(batchJSON, &selectedBatch)
	if err != nil {
		return shim.Error(err.Error())
	}

	var oilTypeVal string
	readingSet := selectedBatch.Readings
	
	for _, v := range readingSet {
		if (v.Type == int(Soil_Quality) && v.Value != "NONE") {
			oilTypeVal = "REJECTED"
			fmt.Println("rejected batch - soil quality bad")
			break
		} else if (v.Type == int(Pesticide) && v.Value != "CHEMICALS") {
			oilTypeVal = "ORGANIC_";
			fmt.Println("batch is organic")
		} else if v.Type == int(Acidity) {
			fmt.Println("evaluating acidity: ", v.Value)
			acidityVal, err := strconv.Atoi(v.Value)
			if err != nil {
				return shim.Error(fmt.Sprintf("Unable to get acidity reading: %s", v.Value))
			}

			if(acidityVal <= 1) {
				oilTypeVal += "DELICATE_EXTRA_VIRGIN"
			} else if (acidityVal <= 2 && acidityVal > 1) {
				oilTypeVal += "EXTRA_VIRGIN"
			} else if (acidityVal <= 4 && acidityVal > 2) {
				oilTypeVal += "VIRGIN"
			} else {
				oilTypeVal = "BEAUTY"
			}
		}
	}

	fmt.Printf("Determined oil type %s for production %s", oilTypeVal, batchId)

	selectedBatch.Grade = oilTypeVal
	selectedBatch.Producer = producer

	batchasBytes, err := json.Marshal(selectedBatch)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(batchId, batchasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("graded"))
}

func queryBatch(stub shim.ChaincodeStubInterface) sc.Response {
	_, parameters := stub.GetFunctionAndParameters()

	if len(parameters) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	batchAsBytes, _ := stub.GetState(parameters[0])
	return shim.Success(batchAsBytes)
}

func main() {
	err := shim.Start(new(BatchManager))
	if err != nil {
		fmt.Printf("Error starting BatchManager chaincode: %s", err)
	}
}