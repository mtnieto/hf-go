package main

import (
	"fmt"

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
	ccID      = "mycontract3"
)
const ExampleCCInitB = "200"

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

	// Org1 resource management client (Org1 is default org)
	org1ResMgmt, err := sdk.NewClient(fabsdk.WithUser("Admin")).ResourceMgmt()
	if err != nil {
		fmt.Print("Failed to create new resource management client: %s", err)
	}
	fmt.Println("Creado el cliente")

	// Create chaincode package for example cc
	ccPkg, err := packager.NewCCPackage("chaincode_example02", "chaincodes/")
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("Llega aqui")

	installCCReq := resmgmt.InstallCCRequest{Name: "mycontract3", Path: "chaincode_example02", Version: "2", Package: ccPkg}
	fmt.Println("Creada request")
	// Install example cc to Org1 peers
	_, err = org1ResMgmt.InstallCC(installCCReq)
	if err != nil {
		fmt.Println("Error al instalar", err)
	}
	fmt.Println("Instalado")

	// Set up chaincode policy to 'any of two msps'
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"hboMSP", "netflixMSP"})

	upgradeRq := resmgmt.UpgradeCCRequest{Name: "mycontract3", Path: "chaincode_example2", Version: "2", Args: initArgs, Policy: ccPolicy}

	// Org1 resource manager will instantiate 'example_cc' version 1 on 'orgchannel'
	err = org1ResMgmt.UpgradeCC("serieschannel", upgradeRq)
	if err != nil {
		fmt.Println(err)
	}

}
