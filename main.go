package main

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"

	resmgmt "github.com/hyperledger/fabric-sdk-go/api/apitxn/resmgmtclient"

	"github.com/hyperledger/fabric-sdk-go/pkg/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

const (
	channelID = "serieschannel"
	orgName   = "netflix"
	orgAdmin  = "Admin"
	ccID      = "mycc"
)
const ExampleCCInitB = "200"

var initArgs = [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte(ExampleCCInitB)}

func main() {

	// Creación del SDK
	sdk, err := fabsdk.New(config.FromFile("config.yaml"))
	if err != nil {
		fmt.Print("Failed to create new SDK: %s", err)
	}
	fmt.Println("Creado el SDK")

	// Creación del usuario para interactuar.
	org1ResMgmt, err := sdk.NewClient(fabsdk.WithUser("Admin")).ResourceMgmt()
	if err != nil {
		fmt.Print("Failed to create new resource management client: %s", err)
	}
	fmt.Println("Creado el cliente")

	// Crear el package para el chaincode
	ccPkg, err := packager.NewCCPackage("chaincode_example02", "chaincodes/")
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("Creado el paquete")

	installCCReq := resmgmt.InstallCCRequest{Name: "mycc", Path: "chaincode_example02", Version: "1", Package: ccPkg}
	fmt.Println("Creada request")
	// Install example cc to Org1 peers
	_, err = org1ResMgmt.InstallCC(installCCReq)
	if err != nil {
		fmt.Println("Error al instalar", err)
	}
	fmt.Println("Instalado chaincode en peer")

	// Set up chaincode policy to 'any of two msps'
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"hboMSP", "netflixMSP"})

	instantciateCReq := resmgmt.InstantiateCCRequest{Name: "mycc", Path: "chaincode_example02", Version: "1", Args: initArgs, Policy: ccPolicy}
	// Instanciación del chaincode
	err = org1ResMgmt.InstantiateCC("serieschannel", instantciateCReq)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Instanciado")

	//Query del valor B

	chClientOrg1User, err := sdk.NewClient(fabsdk.WithUser("Admin"), fabsdk.WithOrg("netflix")).Channel("serieschannel")
	if err != nil {
		fmt.Println("Failed to create new channel client for Org1 user: %s", err)
	}
	queryArgs := [][]byte{[]byte("b")}
	initialValue, err := chClientOrg1User.Query(apitxn.Request{ChaincodeID: "mycc", Fcn: "query", Args: queryArgs})
	if err != nil {
		fmt.Println("Failed to query funds: %s", err)
	}
	fmt.Println("B value: ", string(initialValue))

	// //invoke
	invokeArgs := [][]byte{[]byte("b"), []byte("a"), []byte("1")}

	// Invoke chaincode
	_, _, err = chClientOrg1User.Execute(apitxn.Request{ChaincodeID: "mycc", Fcn: "invoke", Args: invokeArgs})
	if err != nil {
		fmt.Print("Failed to move funds: %s", err)
	}
	//QUERY de comprobación de invoke OK
	value, err := chClientOrg1User.Query(apitxn.Request{ChaincodeID: "mycc", Fcn: "query", Args: queryArgs})
	if err != nil {
		fmt.Println("Failed to query funds: %s", err)
	}
	fmt.Println("B value: ", string(value))

	//Upgradeando el contrato
	installCCReq2 := resmgmt.InstallCCRequest{Name: "mycc", Path: "chaincode_example02", Version: "2", Package: ccPkg}
	fmt.Println("Creada request")
	// Install example cc to Org1 peers
	_, err = org1ResMgmt.InstallCC(installCCReq2)
	if err != nil {
		fmt.Println("Error al instalar", err)
	}
	fmt.Println("Instalado")

	// Set up chaincode policy to 'any of two msps'

	upgradeRq := resmgmt.UpgradeCCRequest{Name: "mycc", Path: "chaincode_example2", Version: "2", Args: initArgs, Policy: ccPolicy}

	// Org1 resource manager will instantiate 'example_cc' version 1 on 'orgchannel'
	err = org1ResMgmt.UpgradeCC("serieschannel", upgradeRq)
	if err != nil {
		fmt.Println(err)
	}

}
