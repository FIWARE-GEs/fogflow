package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	mux "github.com/gufranmirza/go-router"
	"net/http"
	"log"
)

//var ctxUpdateBuffer []*ContextObject

var isConfigured = false
var brokerURL = ""
var myReferenceURL = ""

var inputEntityId = ""
var inputEntityType = ""

func readConfig(fileName string) []interface{} {
	config, e := ioutil.ReadFile(fileName)
	if e != nil {
		fmt.Printf("File Error: [%v]\n", e)
		os.Exit(1)
	}

	dec := json.NewDecoder(bytes.NewReader(config))
	var commands []interface{}
	dec.Decode(&commands)

	return commands
}

func startApp() {
	fmt.Println("start to receive input data streams via a listening port")
}

func stopApp() {
	fmt.Println("clean up the app")
}

// handle the commands received from the engine
func handleAdmin(commands []interface{}) {
	fmt.Println("=============configuration commands=============")
	fmt.Println(commands)

	handleCmds(commands)

	isConfigured = true
}

func onNotify(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello fmt")
}

func notify2execution() {
	 //apply the configuration
	adminCfg := os.Getenv("adminCfg")
	fmt.Println("handle the initial admin configuration " + adminCfg)
	var commands []interface{}
	json.Unmarshal([]byte(adminCfg), commands)
	handleCmds(commands)

	// get the listening port number from the environment variables given by the FogFlow edge worker
	myport, err := strconv.Atoi(os.Getenv("myport"))
	if err != nil {
		fmt.Println("myport is not set up properly for receiving notification")
		return
	}

	router := mux.NewRouter()
	router.POST("/notifyContext", onNotify)
	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(myport), router))
	startApp()

	// wait for the signal to stop the main thread
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	<-c

	stopApp()
}

func element2Object(element map[string]interface{}) map[string]interface{}{
	ctxObj := make(map[string]interface{})
	for key,ele := range element {
		ctxObj[key] = ele
	}
	return ctxObj
}

func object2Element(element map[string]interface{}) map[string]interface{} {
	ctxObject := make(map[string]interface{})
	ctxObject["id"] = element["id"]
	ctxObject["type"] = element["type"]
	for key, _:= range element {
		if (key != "id") && (key != "type") && (key != "modifiedAt") && (key != "createdAt") && (key != "observationSpace") && (key != "operationSpace") && (key != "location") && (key != "@context") {
			ctxObject[key] =  element[key]
		}
	}
	return ctxObject
}

func query2execution() map[string]interface{} {
	//client := NGSILdClient{IoTBrokerURL: brokerURL}
	//_, err := 
	ctxObjects, err := queryContext(inputEntityId,brokerURL)
	if err != nil {
		fmt.Println("failed to update context",err)
		return nil
	}
	ctxObj := element2Object(ctxObjects)
	return ctxObj
}

/*func handleCmds(commands []ConfigCommand) {
	for _, cmd := range commands {
		handleCmd(cmd)
	}

	// send the updates in the buffer
	sendUpdateWithinBuffer()
}*/

func handleCmds(cmd []interface{}) {
	for _, cmdEle := range cmd {
		cmdMapEle := cmdEle.(map[string]interface{})
		CommandType := cmdMapEle["command"].(string)
		fmt.Println(CommandType)
		switch  CommandType { 
		case "CONNECT_BROKER":
			connectBroker(cmdMapEle)
		case "SET_INPUTS":
			setInputs(cmdMapEle)
		case "SET_OUTPUTS":
			setOutputs(cmdMapEle)
		case "SET_REFERENCE":
			setReferenceURL(cmdMapEle)
		}
	}
}

func connectBroker(cmd map[string]interface{}) {
	brokerURL = cmd["brokerURL"].(string)
	fmt.Println("set brokerURL = " + brokerURL)
}

func setInputs(cmd map[string]interface{}) {
	inputEntityId = cmd["id"].(string)
	inputEntityType = cmd["type"].(string)
	fmt.Println("input has been set to (Id: " + inputEntityId + ", Type : " + inputEntityType + ")")
}

func setOutputs(cmd map[string]interface{}) {
	fmt.Println("SetOutputs  in config.json")
}

func setReferenceURL(cmd map[string]interface{}) {
	//myReferenceURL = cmd.ReferenceURL
	//fmt.Println("your application can subscribe addtional inputs under the reference URL: " + myReferenceURL)
	fmt.Println("Set ReferenceURL  in config.json")
}

//
// publish context entities:
//

func publish(ctxUpdate map[string]interface{}) {

	if brokerURL == "" {
		fmt.Println("=== broker is not configured for your update")
		return
	}

	ctxUpdateEle := object2Element(ctxUpdate)
	err := UpdateLdContext(ctxUpdateEle, brokerURL)
	if err != nil {
		fmt.Println(err)
	}
}

func runInTestMode(runOnce bool) {
	fmt.Println("=== TEST MODE ====")

	// load the configuration
	commands := readConfig("config.json")
	fmt.Println(commands)
	handleCmds(commands)

	// query the required inputs and trigger the data processing function
	element := query2execution()
	fogfunction(element, publish)
	fmt.Println(element)
}

func runInOperationMode() {
	fmt.Println("=== OPERATION MODE ====")

	syncMode := os.Getenv("sync")
	if syncMode == "yes" {
		// query the required inputs and trigger the data processing function
		query2execution()
	} else {
		// trigger the data processing function to handle the received notification
		notify2execution()
	}
}

func main() {
	if len(os.Args) == 2 && os.Args[1] == "-o" {
		runInOperationMode()
	} else {
		runInTestMode(true)
	}
}
