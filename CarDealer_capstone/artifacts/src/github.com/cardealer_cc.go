package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// AutoTraceChaincode example simple Chaincode implementation
type AutoTraceChaincode struct {
}

type vehiclePart struct {
	ObjectType   string `json:"docType"`      //docType is used to distinguish the various types of objects in state database
	SerialNumber string `json:"serialNumber"` //the fieldtags are needed to keep case from bouncing around
	Assembler    string `json:"assembler"`
	AssemblyDate int    `json:"assemblyDate"`
	Name         string `json:"name"`
	Owner        string `json:"owner"`
	Recall       bool   `json:"recall"`
	RecallDate   int    `json:"recallDate"`
}

type vehicle struct {
	ObjectType         string `json:"docType"`       //docType is used to distinguish the various types of objects in state database
	ChassisNumber      string `json:"chassisNumber"` //the fieldtags are needed to keep case from bouncing around
	Manufacturer       string `json:"manufacturer"`
	Model              string `json:"model"`
	AssemblyDate       int    `json:"assemblyDate"`
	AirbagSerialNumber string `json:"airbagSerialNumber"`
	Owner              string `json:"owner"`
	Recall             bool   `json:"recall"`
	RecallDate         int    `json:"recallDate"`
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(AutoTraceChaincode))
	if err != nil {
		fmt.Printf("Error starting Parts Trace chaincode: %s", err)
	}
}

// ===================================================================================
//  init...initializes chaincode
// ===================================================================================
func (t *AutoTraceChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// ===================================================================================
// Invoke - Our entry point for Invocations
// ===================================================================================
func (t *AutoTraceChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "initVehiclePart" { //create a new vehiclePart
		return t.initVehiclePart(stub, args)
	} else if function == "transferVehiclePart" { //change owner of a specific vehicle part
		return t.transferVehiclePart(stub, args)
	} else if function == "deleteVehiclePart" { //delete a vehicle part
		return t.deleteVehiclePart(stub, args)
	} else if function == "readVehiclePart" { //read a vehiclePart
		return t.readVehiclePart(stub, args)
	} else if function == "queryVehiclePartByOwner" { //find vehicle part for owner X (rich query)
		return t.queryVehiclePartByOwner(stub, args)
	} else if function == "queryVehiclePart" { //find vehicle part based on an ad hoc rich query (rich query)
		return t.queryVehiclePart(stub, args)
	} else if function == "getHistoryForRecord" { //get history of values for a record
		return t.getHistoryForRecord(stub, args)
	} else if function == "queryVehiclePartByNameOwner" { //get vehicle part based on part name and owner (rich query)
		return t.queryVehiclePartByNameOwner(stub, args)
	} else if function == "initVehicle" { //create a new vehicle
		return t.initVehicle(stub, args)
	} else if function == "transferVehicle" { //change owner of a specific vehicle
		return t.transferVehicle(stub, args)
	} else if function == "readVehicle" { //read a vehicle
		return t.readVehicle(stub, args)
	} else if function == "deleteVehicle" { //delete a vehicle
		return t.deleteVehicle(stub, args)
	} else if function == "transferPartToVehicle" { // transfer airbag to vehicle
		return t.transferPartToVehicle(stub, args)
	} else if function == "setPartRecallState" { // set recall state of vehicle part
		return t.setPartRecallState(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ============================================================
// initVehiclePart - create a new vehicle part, store into chaincode state
// ============================================================
func (t *AutoTraceChaincode) initVehiclePart(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	// data model with recall fields
	//   0       	1      		2     		3				4						5	  6
	// "ser1234", "tata", "1502688979", "airbag 2020", "aaimler ag / mercedes", "false", "0"

	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	// ==== Input sanitation ====
	fmt.Println("- start init vehicle part")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return shim.Error("5th argument must be a non-empty string")
	}

	serialNumber := args[0]
	assembler := strings.ToLower(args[1])
	assemblyDate, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("3rd argument must be a numeric string")
	}
	name := strings.ToLower(args[3])
	owner := strings.ToLower(args[4])

	recall, err := strconv.ParseBool(args[5])
	if err != nil {
		return shim.Error("6th argument must be a boolean string")
	}
	recallDate, err := strconv.Atoi(args[6])
	if err != nil {
		return shim.Error("7th argument must be a numeric string")
	}

	// ==== Check if vehicle part already exists ====
	vehiclePartAsBytes, err := stub.GetState(serialNumber)
	if err != nil {
		return shim.Error("Failed to get vehicle part: " + err.Error())
	} else if vehiclePartAsBytes != nil {
		fmt.Println("This vehicle part already exists: " + serialNumber)
		return shim.Error("This vehicle part already exists: " + serialNumber)
	}

	// ==== Create vehiclePart object and marshal to JSON ====
	objectType := "vehiclePart"
	vehiclePart := &vehiclePart{objectType, serialNumber, assembler, assemblyDate, name, owner, recall, recallDate}
	vehiclePartJSONasBytes, err := json.Marshal(vehiclePart)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save vehiclePart to state ===
	err = stub.PutState(serialNumber, vehiclePartJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== Vehicle part saved and indexed. Return success ====
	fmt.Println("- end init vehicle part")
	return shim.Success(nil)
}

// ============================================================
// setPartRecallState - sets recall field of a vehicle
// ============================================================
func (t *AutoTraceChaincode) setPartRecallState(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	// dexpects following arguements
	//   	0       		1
	// "serialNumber", "status (boolean)"
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}

	serialNumber := args[0]
	recall, err := strconv.ParseBool(args[1])
	if err != nil {
		return shim.Error("2nd argument must be a boolean string")
	}

	// ==== Check if vehicle part already exists ====
	vehiclePartAsBytes, err := stub.GetState(serialNumber)
	if err != nil {
		return shim.Error("Failed to get vehicle part: " + err.Error())
	} else if vehiclePartAsBytes == nil {
		fmt.Println("This vehicle part does not exist: " + serialNumber)
		return shim.Error("This vehicle part does not exist:: " + serialNumber)
	}

	vehiclePartJSON := vehiclePart{}
	err = json.Unmarshal(vehiclePartAsBytes, &vehiclePartJSON) //unmarshal it aka JSON.parse()
	if err != nil {
		fmt.Println("Unable to unmarshall vehicle part from byte to JSON object: " + serialNumber)
		return shim.Error("Unable to unmarshall vehicle part from byte to JSON object: " + serialNumber)
	}

	// ==== Create vehiclePart object and marshal to JSON ====
	objectType := "vehiclePart"
	vehiclePart := &vehiclePart{objectType, serialNumber, vehiclePartJSON.Assembler, vehiclePartJSON.AssemblyDate, vehiclePartJSON.Name, vehiclePartJSON.Owner, recall, 1502688979}
	vehiclePartJSONasBytes, err := json.Marshal(vehiclePart)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save vehiclePart to state ===
	err = stub.PutState(serialNumber, vehiclePartJSONasBytes)

	// ==== Vehicle part saved. Return success ====
	fmt.Println("- end setPartRecallState")
	return shim.Success(nil)
}

// ============================================================
// initVehicle - create a new vehicle , store into chaincode state
// ============================================================
func (t *AutoTraceChaincode) initVehicle(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	// data model with recall fields
	//   0       		1      		2     		3			   4		5	       6			7
	// "mer1000001", "mercedes", "c class", "1502688979", "ser1234", "mercedes", "false", "1502688979"

	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	// ==== Input sanitation ====
	fmt.Println("- start init vehicle")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return shim.Error("5th argument must be a non-empty string")
	}
	if len(args[5]) <= 0 {
		return shim.Error("6th argument must be a non-empty string")
	}

	chassisNumber := args[0]
	manufacturer := strings.ToLower(args[1])
	model := strings.ToLower(args[2])
	assemblyDate, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("3rd argument must be a numeric string")
	}
	airbagSerialNumber := strings.ToLower(args[4])
	owner := strings.ToLower(args[5])

	recall, err := strconv.ParseBool(args[6])
	if err != nil {
		return shim.Error("7th argument must be a boolean string")
	}
	recallDate, err := strconv.Atoi(args[7])
	if err != nil {
		return shim.Error("8th argument must be a numeric string")
	}

	// ==== Check if vehicle already exists ====
	vehicleAsBytes, err := stub.GetState(chassisNumber)
	if err != nil {
		return shim.Error("Failed to get vehicle: " + err.Error())
	} else if vehicleAsBytes != nil {
		return shim.Error("This vehicle already exists: " + chassisNumber)
	}

	// ==== Create vehicle object and marshal to JSON ====
	objectType := "vehicle"
	vehicle := &vehicle{objectType, chassisNumber, manufacturer, model, assemblyDate, airbagSerialNumber, owner, recall, recallDate}
	vehicleJSONasBytes, err := json.Marshal(vehicle)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save vehicle to state ===
	err = stub.PutState(chassisNumber, vehicleJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== Vehicle part saved and indexed. Return success ====
	fmt.Println("- end init vehicle")
	return shim.Success(nil)
}

// ===============================================
// readVehiclePart - read a vehicle part from chaincode state
// ===============================================
func (t *AutoTraceChaincode) readVehiclePart(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var serialNumber, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting serial number of the vehicle part to query")
	}

	serialNumber = args[0]
	valAsbytes, err := stub.GetState(serialNumber) //get the vehiclePart from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + serialNumber + "\"}"
		fmt.Println(jsonResp)
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Vehicle part does not exist: " + serialNumber + "\"}"
		fmt.Println(jsonResp)
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// ===============================================
// readVehicle - read a vehicle from chaincode state
// ===============================================
func (t *AutoTraceChaincode) readVehicle(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var chassisNumber, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting chassis number of the vehicle to query")
	}

	chassisNumber = args[0]
	valAsbytes, err := stub.GetState(chassisNumber) //get the vehicle from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + chassisNumber + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Vehicle does not exist: " + chassisNumber + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// ==================================================
// deleteVehiclePart - remove a vehiclePart key/value pair from state
// ==================================================
func (t *AutoTraceChaincode) deleteVehiclePart(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var jsonResp string
	var vehiclePartJSON vehiclePart
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	serialNumber := args[0]

	// to maintain the assember~serialNumber index, we need to read the vehiclePart first and get its assembler
	valAsbytes, err := stub.GetState(serialNumber) //get the vehiclePart from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + serialNumber + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"VehiclePart does not exist: " + serialNumber + "\"}"
		return shim.Error(jsonResp)
	}

	err = json.Unmarshal([]byte(valAsbytes), &vehiclePartJSON)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to decode JSON of: " + serialNumber + "\"}"
		return shim.Error(jsonResp)
	}

	err = stub.DelState(serialNumber) //remove the vehiclePart from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}

	return shim.Success(nil)
}

// ==================================================
// deleteVehicle - remove a vehicle key/value pair from state
// ==================================================
func (t *AutoTraceChaincode) deleteVehicle(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var jsonResp string
	var vehicleJSON vehicle
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	chassisNumber := args[0]

	// to maintain the manufacturer~chassisNumber index, we need to read the vehicle first and get its assembler
	valAsbytes, err := stub.GetState(chassisNumber) //get the vehicle from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + chassisNumber + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Vehicle does not exist: " + chassisNumber + "\"}"
		return shim.Error(jsonResp)
	}

	err = json.Unmarshal([]byte(valAsbytes), &vehicleJSON)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to decode JSON of: " + chassisNumber + "\"}"
		return shim.Error(jsonResp)
	}

	err = stub.DelState(chassisNumber) //remove the vehicle from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}

	return shim.Success(nil)
}

// ===========================================================
// transfer a vehicle part by setting a new owner name on the vehiclePart
// ===========================================================
func (t *AutoTraceChaincode) transferVehiclePart(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0       1       3
	// "name", "from", "to"
	if len(args) < 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	serialNumber := args[0]
	currentOwner := strings.ToLower(args[1])
	newOwner := strings.ToLower(args[2])
	fmt.Println("- start transferVehiclePart ", serialNumber, currentOwner, newOwner)

	message, err := t.transferPartHelper(stub, serialNumber, currentOwner, newOwner)
	if err != nil {
		return shim.Error(message + err.Error())
	} else if message != "" {
		return shim.Error(message)
	}

	fmt.Println("- end transferVehiclePart (success)")
	return shim.Success(nil)
}

// ===========================================================
// transferParts : helper method for transferVehiclePart
// ===========================================================
func (t *AutoTraceChaincode) transferPartHelper(stub shim.ChaincodeStubInterface, serialNumber string, currentOwner string, newOwner string) (string, error) {
	// attempt to get the current vehiclePart object by serial number.
	// if sucessful, returns us a byte array we can then us JSON.parse to unmarshal
	fmt.Println("Transfering part with serial number: " + serialNumber + " To: " + newOwner)
	vehiclePartAsBytes, err := stub.GetState(serialNumber)
	if err != nil {
		return "Failed to get vehicle part: " + serialNumber, err
	} else if vehiclePartAsBytes == nil {
		return "Vehicle part does not exist: " + serialNumber, nil
	}

	vehiclePartToTransfer := vehiclePart{}
	err = json.Unmarshal(vehiclePartAsBytes, &vehiclePartToTransfer) //unmarshal it aka JSON.parse()
	if err != nil {
		return "", err
	}

	// if currentOwner != vehiclePartToTransfer.Owner {
	// 	return "This asset is currently owned by another entity.", err
	// }

	vehiclePartToTransfer.Owner = newOwner //change the owner

	vehiclePartJSONBytes, _ := json.Marshal(vehiclePartToTransfer)
	err = stub.PutState(serialNumber, vehiclePartJSONBytes) //rewrite the vehiclePart
	if err != nil {
		return "", err
	}

	return "", nil
}

// ===========================================================
// transfer a vehicle part by setting a new owner name on the vehiclePart
// ===========================================================
func (t *AutoTraceChaincode) transferPartToVehicle(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- start transferPartToVehicle")
	//   	0      			 1
	// "serialNumber", "chassisNumber"
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	serialNumber := args[0]
	chassisNumber := args[1]

	message, err := t.transferPartToVehicleHelper(stub, serialNumber, chassisNumber)
	if err != nil {
		return shim.Error(message)
	}

	fmt.Println("- end transferPartToVehicle (success)")
	return shim.Success(nil)
}

// ===========================================================
// transferPartToVehicleHelper : helper for transferPartToVehicle
// ===========================================================
func (t *AutoTraceChaincode) transferPartToVehicleHelper(stub shim.ChaincodeStubInterface, serialNumber string, chassisNumber string) (string, error) {
	vehiclePartAsBytes, err := stub.GetState(serialNumber)
	if err != nil {
		return "Failed to get vehicle part: " + serialNumber, err
	} else if vehiclePartAsBytes == nil {
		return "Vehicle part does not exist: " + serialNumber, nil
	}

	vehicleAsBytes, err := stub.GetState(chassisNumber)
	if err != nil {
		return "Failed to get vehicle: " + chassisNumber, err
	} else if vehicleAsBytes == nil {
		return "Vehicle does not exist: " + chassisNumber, err
	}

	part := vehiclePart{}
	err = json.Unmarshal(vehiclePartAsBytes, &part) //unmarshal it aka JSON.parse()
	if err != nil {
		return "", err
	}

	car := vehicle{}
	err = json.Unmarshal(vehicleAsBytes, &car) //unmarshal it aka JSON.parse()
	if err != nil {
		return "", err
	}

	if car.Owner != part.Owner {
		return "Illegal Transfer.", err
	}

	vehicleToModify := vehicle{}
	err = json.Unmarshal(vehicleAsBytes, &vehicleToModify) //unmarshal it aka JSON.parse()
	if err != nil {
		return "", err
	}
	vehicleToModify.AirbagSerialNumber = serialNumber //change the serialnumber of the vehicle

	vehicleJSONBytes, _ := json.Marshal(vehicleToModify)
	err = stub.PutState(chassisNumber, vehicleJSONBytes) //rewrite the vehicle
	if err != nil {
		return "", err
	}
	return "", nil
}

// ===========================================================
// transferVehicleHelper: transfer a vehicle  by setting a new owner name on the vehicle
// ===========================================================
func (t *AutoTraceChaincode) transferVehicle(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0       1       3
	// "name", "from", "to"
	if len(args) < 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	chassisNumber := args[0]
	currentOnwer := strings.ToLower(args[1])
	newOwner := strings.ToLower(args[2])
	fmt.Println("- start transferVehicle ", chassisNumber, currentOnwer, newOwner)

	// attempt to get the current vehicle object by serial number.
	// if sucessful, returns us a byte array we can then us JSON.parse to unmarshal
	message, err := t.trannsferVehicleHelper(stub, chassisNumber, currentOnwer, newOwner)
	if err != nil {
		return shim.Error(message + err.Error())
	} else if message != "" {
		return shim.Error(message)
	}

	fmt.Println("- end transferVehicle (success)")
	return shim.Success(nil)
}

// ===========================================================
// trannsferVehicleHelper : helper method for transferVehicle
// ===========================================================
func (t *AutoTraceChaincode) trannsferVehicleHelper(stub shim.ChaincodeStubInterface, chassisNumber string, currentOwner string, newOwner string) (string, error) {
	// attempt to get the current vehicle object by serial number.
	// if sucessful, returns us a byte array we can then us JSON.parse to unmarshal
	fmt.Println("Transfering vehicle with chassis number: " + chassisNumber + " To: " + newOwner)
	vehicleAsBytes, err := stub.GetState(chassisNumber)
	if err != nil {
		return "Failed to get vehicle:", err
	} else if vehicleAsBytes == nil {
		return "Vehicle does not exist", err
	}

	vehicleToTransfer := vehicle{}
	err = json.Unmarshal(vehicleAsBytes, &vehicleToTransfer) //unmarshal it aka JSON.parse()
	if err != nil {
		return "", err
	}

	// if currentOwner != vehicleToTransfer.Owner {
	// 	return "This asset is currently owned by another entity.", err
	// }

	vehicleToTransfer.Owner = newOwner //change the owner

	vehicleJSONBytes, _ := json.Marshal(vehicleToTransfer)
	err = stub.PutState(chassisNumber, vehicleJSONBytes) //rewrite the vehicle
	if err != nil {
		return "", err
	}

	return "", nil
}

// ===== Example: Parameterized rich query =================================================
// queryVehiclePartByNameOwner queries for vehicle part based on a passed in name and owner.
// This is an example of a parameterized query where the query logic is baked into the chaincode,
// and accepting a single query parameter (owner).
// =========================================================================================
func (t *AutoTraceChaincode) queryVehiclePartByNameOwner(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	name := strings.ToLower(args[0])
	owner := strings.ToLower(args[1])

	queryString := fmt.Sprintf("SELECT valueJson FROM <STATE> WHERE json_extract(valueJson, '$.docType', '$.name', '$.owner') = '[\"vehiclePart\",\"%s\",\"%s\"]'", name, owner)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// ===== Example: Parameterized rich query =================================================
// queryVehiclePartByOwner queries for vehicle part based on a passed in owner.
// This is an example of a parameterized query where the query logic is baked into the chaincode,
// and accepting a single query parameter (owner).
// =========================================================================================
func (t *AutoTraceChaincode) queryVehiclePartByOwner(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	owner := strings.ToLower(args[0])

	queryString := fmt.Sprintf("SELECT valueJson FROM <STATE> WHERE json_extract(valueJson, '$.docType', '$.owner') = '[\"vehiclePart\",\"%s\"]'", owner)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// ===== Example: Ad hoc rich query ========================================================
// queryVehiclePart uses a query string to perform a query for vehiclePart.
// Query string matching state database syntax is passed in and executed as is.
// Supports ad hoc queries that can be defined at runtime by the client.
// =========================================================================================
func (t *AutoTraceChaincode) queryVehiclePart(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// "queryString"
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	queryString := args[0]

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

// ===========================================================================================
// getHistoryForRecord returns the histotical state transitions for a given key of a record
// ===========================================================================================
func (t *AutoTraceChaincode) getHistoryForRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	recordKey := args[0]

	fmt.Printf("- start getHistoryForRecord: %s\n", recordKey)

	resultsIterator, err := stub.GetHistoryForKey(recordKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the key/value pair
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON vehiclePart)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForRecord returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
