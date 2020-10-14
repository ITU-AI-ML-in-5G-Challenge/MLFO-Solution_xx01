package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

//struct to understand model query response from db
type response struct {
	id           string
	uri          string
	accessType   string
	trainingTime string
	nGPU         int
	resourceReq  string
}

//struct to understand cluster query response from db
type clusterResponse struct {
	id         string
	gpuPresent bool
	nGPU       int
}

var numGPU int

func main() {

	arg := os.Args[1]

	m := make(map[string]interface{})
	accessType := ""
	trainingTime := ""
	resourceReq := ""
	yamlFile, err := ioutil.ReadFile(arg) //Read yaml file
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println("Parsing Intent...")

	//iterate over the yaml structure and decide model requirements based on use case
	for k, v := range m {
		if k == "model" {
			vnew, _ := v.(map[interface{}]interface{})
			switch vnew["usecase"] {
			case "edge":
				accessType = "private"
				trainingTime = "low"
				resourceReq = "low"
				fmt.Println("\nSelecting model for Edge use case...")
				fmt.Println("\nModel selection requirements: accessType = " + accessType + "\t trainingTime = " + trainingTime + "\t resourceRequirements = " + resourceReq)

			case "cloud":
				accessType = "public"
				trainingTime = "high"
				resourceReq = "high"
				fmt.Println("\nSelecting model for Cloud use case...")
				fmt.Println("\nModel selection requirements: accessType = " + accessType + "\t trainingTime = " + trainingTime + "\t resourceRequirements = " + resourceReq)

			default:
				log.Printf("Invalid use case. Now exiting")
				return
			}
		}
	}

	//start connection to the db
	db, err := sql.Open("mysql", "root:mlfo1234@tcp(db:3306)/modelrepo")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}
	// defer the close till after the main function has finished
	defer db.Close()

	fmt.Println("\nQuerying model repository for appropriate model...")

	// perform a db Query for model selection
	query := "SELECT * FROM models WHERE accessType='" + accessType + "' AND trainingTime='" + trainingTime + "' AND resourceReq='" + resourceReq + "';"
	result, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	//Scan the received query response using predefined struct
	for result.Next() {
		var r response
		err = result.Scan(&r.id, &r.uri, &r.accessType, &r.trainingTime, &r.nGPU, &r.resourceReq)
		if err != nil {
			panic(err.Error())
		}
		// print received model information
		fmt.Println("\nReceived model...")
		fmt.Printf("\nModel id: %s\n", r.id)
		fmt.Printf("Model URI: %s\n", r.uri)
		fmt.Printf("Model accessType: %s\n", r.accessType)
		fmt.Printf("Model trainingTime: %s\n", r.trainingTime)
		fmt.Printf("Model GPU requirement: %d GPUs\n", r.nGPU)
		numGPU = r.nGPU
		fmt.Printf("Model resourceRequirements: %s\n", r.resourceReq)
	}

	fmt.Println("-----------------------Model selection complete-----------------------")

	fmt.Println("\nSelecting appropriate training cluster based on model GPU requirement.....")

	// perform a db Query for training cluster selection
	tcquery := "SELECT * FROM trainingClusters WHERE nGPU='" + strconv.Itoa(numGPU) + "';"
	fmt.Println("\nQuerying db for available training clusters.......")
	tcresult, err := db.Query(tcquery)
	if err != nil {
		panic(err.Error())
	}
	defer tcresult.Close()

	//Scan the received query response using predefined struct
	for tcresult.Next() {
		var tc clusterResponse
		err = tcresult.Scan(&tc.id, &tc.gpuPresent, &tc.nGPU)
		if err != nil {
			panic(err.Error())
		}
		// print received model information
		fmt.Println("\nSelected Training cluster...")
		fmt.Printf("\nCluster id: %s\n", tc.id)
		fmt.Printf("GPUs present in the cluster: %t\n", tc.gpuPresent)
		fmt.Printf("Number of available GPUs: %d\n", tc.nGPU)
	}

	fmt.Println("-----------------------Training Cluster selection complete-----------------------")

	fmt.Println("\nSelecting appropriate inference cluster based on model GPU requirement.....")
	// perform a db Query for training cluster selection
	icquery := "SELECT * FROM inferenceClusters WHERE nGPU='" + strconv.Itoa(numGPU/5) + "';"
	fmt.Println("\nQuerying db for available inference clusters.......")
	icresult, err := db.Query(icquery)
	if err != nil {
		panic(err.Error())
	}
	defer icresult.Close()

	//Scan the received query response using predefined struct
	for icresult.Next() {
		var ic clusterResponse
		err = icresult.Scan(&ic.id, &ic.gpuPresent, &ic.nGPU)
		if err != nil {
			panic(err.Error())
		}
		// print received model information
		fmt.Println("\nSelected Inference cluster...")
		fmt.Printf("\nCluster id: %s\n", ic.id)
		fmt.Printf("GPUs present in the cluster: %t\n", ic.gpuPresent)
		fmt.Printf("Number of available GPUs: %d\n", ic.nGPU)
	}
	fmt.Println("-----------------------Inference Cluster selection complete-----------------------")
}
