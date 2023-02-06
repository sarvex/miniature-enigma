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
	 Weight 	int			`json:"weight"`
	 Grower		string		`json:"grower"`
	 Readings	[]Reading	`json:"readings"`
 }

 type Reading struct {
	 Type	int		`json:"type"`
	 Value 	string	`json:"value"`
	 Member	string	`json:"member"`
 }
 
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

	 
 }