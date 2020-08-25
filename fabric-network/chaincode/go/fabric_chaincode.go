package main

import (

	"fmt"
	"strconv"

	"encoding/json"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
)

// Define the Smart Contract structure
type SmartContract struct {

}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract.
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "store" {
		return s.store(APIstub, args)
	} else if function == "Init" {
		return s.Init(APIstub)
	} else if function == "retrieve" {
		return s.retrieve(APIstub, args)
	} else {fmt.Println("Invalid Smart Contract function name.")
	return shim.Error("Invalid Smart Contract function name.")}
}

type unit struct {
	Value string `json:"value"`
}

func (s *SmartContract) retrieve(APIstub shim.ChaincodeStubInterface, args []string ) sc.Response {
	var iterator,err=APIstub.GetStateByPartialCompositeKey("key~timestamp",  []string{args[0]})
	if err != nil {
		fmt.Printf("error,error")
		return shim.Error(err.Error())
	}
	defer iterator.Close()
	total_result:="\"{"

	// Iterate through result set 
	var i int
	var unitJSON unit
	for i = 0; iterator.HasNext(); i++ {

		responseRange, err := iterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		_, compositeKeyParts, err := APIstub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}

		returnedTimestamp := compositeKeyParts[1]  
		err = json.Unmarshal(responseRange.Value, &unitJSON)
		if err != nil {
			jsonResp := "{\"Error\":\"Failed to decode JSON. \"}"
			return shim.Error(jsonResp)
		}
	
		total_result=total_result+"\\\"Object"+strconv.Itoa(i)+"\\\":{\\\"key\\\":\\\""+args[0]+"\\\",\\\"time\\\":\\\""+returnedTimestamp+"\\\",\\\"value\\\":\\\""+unitJSON.Value+"\\\"}"
		if (iterator.HasNext()){
			total_result=total_result+","
		} else {
			total_result=total_result+",\\\"length\\\":\\\""+strconv.Itoa(i+1)+"\\\""
		}
	}

	total_result=total_result+"}\""

	return shim.Success([]byte(total_result))
}

func (s *SmartContract) store(APIstub shim.ChaincodeStubInterface, args []string ) sc.Response {

	var compositeKey, err=APIstub.CreateCompositeKey("key~timestamp",[]string{args[0],args[1]})

	if err != nil {
		return shim.Error(err.Error())
	}

	APIstub.PutState(compositeKey, []byte(args[2])); 

	return shim.Success(nil)
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface ) sc.Response {
	return shim.Success(nil)
}

func main() {
	//Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil { 
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}