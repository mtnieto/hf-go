package main

import (
	"fmt"

	fab "github.com/hyperledger/fabric-sdk-go/api/apifabclient"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	"github.com/hyperledger/fabric-sdk-go/pkg/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

const (
	channelID = "serieschannel"
	orgName   = "netflix"
	orgAdmin  = "Admin"
	ccID      = "mycc"
)
const ExampleCCInitB = "200"

var orgTestPeer1 fab.Peer
var initArgs = [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte(ExampleCCInitB)}

//ExampleCCInitArgs returns example cc initialization args
func ExampleCCInitArgs() [][]byte {
	return initArgs
}

func main() {

	// Create SDK setup for the integration tests
	sdk, err := fabsdk.New(config.FromFile("config.yaml"))
	if err != nil {
		fmt.Print("Failed to create new SDK: %s", err)
	}
	fmt.Println("Creado el SDK")

	// // Org1 resource management client (Org1 is default org)
	// org1ResMgmt, err := sdk.NewClient(fabsdk.WithUser("Admin")).ResourceMgmt()
	// if err != nil {
	// 	fmt.Print("Failed to create new resource management client: %s", err)
	// }
	// fmt.Println("Creado el cliente")

	// // Create chaincode package for example cc
	// ccPkg, err := packager.NewCCPackage("chaincode_example02", "chaincodes/")
	// if err != nil {
	// 	fmt.Print(err)
	// }
	// fmt.Println("Llega aqui")

	// installCCReq := resmgmt.InstallCCRequest{Name: "mycc", Path: "chaincode_example02", Version: "1", Package: ccPkg}
	// fmt.Println("Creada request")
	// // Install example cc to Org1 peers
	// _, err = org1ResMgmt.InstallCC(installCCReq)
	// if err != nil {
	// 	fmt.Println("Error al instalar", err)
	// }
	// fmt.Println("Instalado")

	// // Set up chaincode policy to 'any of two msps'
	// ccPolicy := cauthdsl.SignedByAnyMember([]string{"hboMSP", "netflixMSP"})

	// instantciateCReq := resmgmt.InstantiateCCRequest{Name: "mycc", Path: "chaincode_example02", Version: "1", Args: integration.ExampleCCInitArgs(), Policy: ccPolicy}
	// // Org1 resource manager will instantiate 'example_cc' on 'orgchannel'
	// err = org1ResMgmt.InstantiateCC("serieschannel", instantciateCReq)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println("Instanciado")

	// Org1 user queries initial value on both peers

	chClientOrg1User, err := sdk.NewClient(fabsdk.WithUser("Admin"), fabsdk.WithOrg("netflix")).Channel("serieschannel")
	if err != nil {
		fmt.Println("Failed to create new channel client for Org1 user: %s", err)
	}
	queryArgs := [][]byte{[]byte("b")}
	initialValue, err := chClientOrg1User.Query(apitxn.QueryRequest{ChaincodeID: "mycc", Fcn: "query", Args: queryArgs})
	if err != nil {
		fmt.Println("Failed to query funds: %s", err)
	}
	fmt.Println("B value: ", string(initialValue))

	// //invoke
	invokeArgs := [][]byte{[]byte("b"), []byte("a"), []byte("100")}
	//no entiendo que es esto :D
	txOpts := apitxn.ExecuteTxOpts{ProposalProcessors: []apitxn.ProposalProcessor{orgTestPeer1}}
	_, _, err = chClientOrg1User.ExecuteTxWithOpts(apitxn.ExecuteTxRequest{ChaincodeID: "exampleCC", Fcn: "invoke", Args: invokeArgs}, txOpts)
	if err == nil {
		fmt.Print("Should have failed to move funds due to cc policy")
	}

	// ialue, err := chClientOrg1User.Query(apitxn.QueryRequest{ChaincodeID: "mycc", Fcn: "query", Args: queryArgs})
	// if err != nil {
	// 	fmt.Println("Failed to query funds: %s", err)
	// }
	// fmt.Println(ialue)
}
