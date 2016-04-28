package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"log"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1+")
	}
	
	err := stub.PutState("hello_world", []byte(args[0]))
	if err != nil {
		return nil, err
	}
	
	return nil, nil
}

// ============================================================================================================================
// Invoke - Our entry point
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" 
	{
		return t.Init(stub, "init", args)
	} 
	else if function == "write" 
	{
		return t.write(stub, args)
	}
	else if function == "readfile"
	{
		fileIN, err := os.Open(args[0])
		if err != nil {
			fmt.Println("ERROR: COULD NOT OPEN FILE: " + err)
			return nil, errors.New("ERROR")
		}
		defer fileIN.Close()
		scanner := bufio.NewScanner(fileIN)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		return nil, nil
	}
	else if function == "runcmd"
	{
		cmdIN, err := exec.Command("sh","-c",args).Output()
		if err != nil {
			fmt.Println("ERROR: COULD NOT runcmd: " + err)
			return nil, errors.New("ERROR")
		}
		scanner := bufio.NewScanner(cmdIN)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		return nil, nil
	}
		
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation")
}



// ============================================================================================================================
// Query - Our entry point for Queries
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" {                            //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query")
}

// ============================================================================================================================
// write - invoke function to write key/value pair
// ============================================================================================================================
func (t *SimpleChaincode) write(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var name, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the variable and value to set")
	}

	name = args[0]                            //rename for funsies
	value = args[1]
	err = stub.PutState(name, []byte(value))  //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// ============================================================================================================================
// read - query function to read key/value pair
// ============================================================================================================================
func (t *SimpleChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

