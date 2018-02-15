package main

import (
	"encoding/json"

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

var targetList []string

// Init to initiate the SimpleChaincode class
func (t *ManagementChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Debug("[Management Chaincode][Init]Instanciating chaincode...")
	return shim.Success([]byte("Init called"))
}

// Invoke a method specified in the SimpleChaincode class
func (t *ManagementChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Critical("[Management Chaincode][Invoke]Invoking chaincode...")
	function, args := stub.GetFunctionAndParameters()
	if function == "storeCode" {
		return t.storeCode(stub, args)
	}
	if function == "registrar" {
		return t.registrar(stub, args)
	}
	if function == "getAlias" {
		return t.getAlias(stub, args)
	}
	return shim.Success([]byte("Invoke"))
}
func (t *ManagementChaincode) registrar(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	logger.Critical("[Management Chaincode][Registrar]Calling registrar...")
	//Get the transaction ID
	txID := stub.GetTxID()
	logger.Debug("[Management Chaincode][StoreCode]Transaction ID", txID)
	// Get certs from transaction sender
	caller, err := stub.GetCreator()
	if err != nil {
		return shim.Error(err.Error())
	}
	//Get alias to store
	alias := args[0]
	logger.Critical("[Management Chaincode][StoreCode]Storing cert, alias", args[0])
	err = stub.PutState(string(caller[:]), []byte(alias))
	logger.Critical("[Management Chaincode][StoreCode]Caller", string(caller[:]))
	if err != nil {
		logger.Error("[Management Chaincode][addTarget]Problem adding new target..", err)
		return shim.Error(err.Error())
	}
	logger.Critical("[Management Chaincode][StoreCode]Stored successful", args[0])

	return shim.Success([]byte(caller))
}

func (t *ManagementChaincode) getAlias(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	logger.Critical("[Management Chaincode][Registrar]Calling registrar...")
	//Get the transaction ID
	txID := stub.GetTxID()
	logger.Debug("[Management Chaincode][StoreCode]Transaction ID", txID)
	// Get certs from transaction sender
	caller, err := stub.GetCreator()
	if err != nil {
		return shim.Error(err.Error())
	}

	logger.Critical("[Management Chaincode][StoreCode]Getting alias from caller---", string(caller[:]))
	state, err := stub.GetState(string(caller[:]))
	if err != nil {
		logger.Error("[Management Chaincode][addTarget]Problem adding new target..", err)
		return shim.Error(err.Error())
	}
	logger.Critical("[Management Chaincode][StoreCode]State", string(state[:]))
	return shim.Success([]byte(""))
}

func (t *ManagementChaincode) storeCode(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	logger.Critical("[Management Chaincode][StoreCode]Calling storeCode...")
	//Get the transaction ID
	txID := stub.GetTxID()
	logger.Critical("[Management Chaincode][StoreCode]Transaction ID", txID)
	caller, err := stub.GetCreator()

	if err != nil {
		return shim.Error(err.Error())
	}

	// rawIn := json.RawMessage(args[1])
	// bytes, err := rawIn.MarshalJSON()
	// if err != nil {
	// 	panic(err)
	// }
	// var something Code
	// err = json.Unmarshal(bytes, &something)
	// callerStr := string(caller[:])
	// chaincodeName := args[0]
	// something.Target = callerStr
	// logger.Critical("The caller", callerStr)
	// var data []byte

	// data = []byte(args[1])
	// stub.PutState(chaincodeName, data)
	return shim.Success([]byte(caller))
}
func addTarget(stub shim.ChaincodeStubInterface, newTarget string) bool {
	state, err := stub.GetState("targetList")
	if err != nil {
		logger.Error("[Management Chaincode][addTarget]Problem adding new target..", err)
		return false
	}
	json.Unmarshal(state, &targetList)
	slice := append(targetList, newTarget)
	logger.Critical("[Management Chaincode][addTarget] Actual state of the slice ..", slice)
	toStore, err := json.Marshal(slice)
	if err != nil {
		return false
	}
	logger.Critical("[Management Chaincode][addTarget] Updating the state...")
	stub.PutState("targetList", toStore)
	logger.Critical("[Management Chaincode][addTarget] Updated the state")

	return true
}
func main() {
	err := shim.Start(new(ManagementChaincode))
	if err != nil {
		logger.Debugf("Error: %s", err)
	}
}
