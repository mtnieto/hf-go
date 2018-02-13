package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("Management Chaincode")

// SimpleChaincode representing a class of chaincode
type ManagementChaincode struct{}
type Code struct {
	Name    string `json:"Name"`
	Version string `json:"Version"`
	Source  string `json:"Source"`
	Target  string
}

// Init to initiate the SimpleChaincode class
func (t *ManagementChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Debug("[Management Chaincode][Init]Instanciating chaincode...")
	return shim.Success([]byte("Init called"))
}

// Invoke a method specified in the SimpleChaincode class
func (t *ManagementChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Debug("[Management Chaincode][Invoke]Invoking chaincode...")
	function, args := stub.GetFunctionAndParameters()
	if function == "storeCode" {
		return t.storeCode(stub, args)
	}
	return shim.Success([]byte("Invoke"))
}

func (t *ManagementChaincode) storeCode(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	logger.Debug("[Management Chaincode][StoreCode]Calling storeCode...")
	//Get the transaction ID
	txID := stub.GetTxID()
	logger.Debug("[Management Chaincode][StoreCode]Transaction ID", txID)
	caller, err := stub.GetCreator()

	if err != nil {
		return shim.Error(err.Error())
	}

	rawIn := json.RawMessage(args[1])
	bytes, err := rawIn.MarshalJSON()
	if err != nil {
		panic(err)
	}
	var something Code
	err = json.Unmarshal(bytes, &something)
	fmt.Println(something)
	callerStr := string(caller[:])
	chaincodeName := args[0]
	something.Target = callerStr
	fmt.Println("The caller", callerStr)
	var data []byte

	data = []byte(args[1])
	stub.PutState(chaincodeName, data)
	return shim.Success([]byte(caller))
}

func main() {
	err := shim.Start(new(ManagementChaincode))
	if err != nil {
		logger.Debugf("Error: %s", err)
	}
}
