/*
IMPORTANT NOTE: 12gb of RAM at least required
*/
package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

type bigNet []person

// In Veneto public healthcare beds are 3.3 per 1000 or in total 16.201
// plus 1503 of private healthcare beds or 0.3 per 1000
// http://www.salute.gov.it/imgs/C_17_pubblicazioni_2859_allegato.pdf page 58 of document (71 of pdf)

// Hyperparameters configuration of Simulation
const (
	ISDEBUG             = false
	nNodes              = 4905854 // 4905854 number of people in Veneto
	nEdges              = 150     //Dunbar number 150
	bedIntensiveCare    = 45      //0    //https://www.aulss2.veneto.it/amministrazione-trasparente/disposizioni-generali/atti-generali/regolamenti?p_p_id=101&p_p_lifecycle=0&p_p_state=maximized&p_p_col_id=column-1&p_p_col_pos=22&p_p_col_count=24&_101_struts_action=%2Fasset_publisher%2Fview_content&_101_assetEntryId=10434368&_101_type=document
	bedSubIntensiveCare = 16201   //number of beds
	pIntensiveCare      = 0.02    //probability of requiring intensive Care
	pSubIntensiveCare   = 0.15    //probability of requiring sub intensive care
	hospitalDays        = 7       //the number of day to add to the duration of the disease
	medianR0            = 2.1     //2.28  //https://pubmed.ncbi.nlm.nih.gov/32097725/ 2.06-2.52 95% CI 0,22/1.96 = 0.112
	stdR0               = 0.7     //0.112
	infectiveEpochs     = 3
	simulationEpochs    = 180 //DURATION OF SIMULATION
	deadRate            = 0.025
	//muskEpoch           = 30   //30   //starting epoch of musk set -1 to disable
	//muskProb            = 0.05 //95 //prevention probability
	//socDisEpoch         = -1   //40	//starting epoch of social distacing set -1 to disable
	incubationEpochs = 1 //number of epochs in incubation
)

type person struct {
	Edges           []uint32 `json:"e"`
	RelationType    []string `json:"r"`
	Infective       bool     `json:"-"`
	Survived        bool     `json:"-"`
	Dead            bool     `json:"-"`
	InfectiveEpochs uint32   `json:"-"` // ottimizza per terapia intensiva (14+21)
	InfectiveDays   []int8   `json:"-"`
}

func main() {
	runtime.GOMAXPROCS(8)

	// flags
	loadNetwork := flag.Bool("loadnet", false, "default value is false, if true it load a network from a file called Network.json, to change the loading file name check flag namenet")
	saveNetwork := flag.Bool("savenet", false, "default value is false, if true saves network on timestamp/Network.json")
	fileNetwork := flag.String("namenet", "Network.json", "default value is Network.json, it's the name of the network file")
	mctrials := flag.Int("mctrials", 1, "default value is 1, you can choose how many trials run on the Montecarlo Simulation")
	computeCI := flag.Bool("computeCI", false, "default value is false, set to true when use flag -mctrials > 1 to get Confidence Intervals of metrics")
	computeSSN := flag.Bool("computeSSN", false, "default value is false, set to true to get information about national healthcare system")
	runPyScript := flag.Bool("runpyscript", false, "default valuse is false, set to true if you want to print graphs of simulation with matplotlib")
	folderFlag := flag.String("folder", "", "default value is '', set the name of the folder to generate")
	muskEpoch := flag.Int("muskEpoch", -1, "default value is -1, number of epochs before applying measure")
	muskProb := flag.Float64("muskProb", 0.2, "default value is 0.2, musk policy efficency")
	socDisEpoch := flag.Int("socDisEpoch", -1, "default value is -1, number of epochs before applying measure")
	flag.Parse()

	// random seed
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	// initialize network
	network := make(bigNet, nNodes)
	var epochsResults [simulationEpochs][5]int
	var trialsResults = make([][3]int, *mctrials)
	var ssnEpochsResults [simulationEpochs][2]int
	// creating run folder
	var folderName string
	if *folderFlag != "" {
		folderName = *folderFlag
	} else {
		folderName = strconv.Itoa(int(time.Now().UnixNano()))
	}

	os.MkdirAll(folderName, os.ModePerm)
	os.MkdirAll(folderName+"/trials", os.ModePerm)

	if !*loadNetwork {
		log.Println("Creating network...")

		// init counter for relations assignment
		//	'P': family relationship
		//	'A': friends
		//	'C': acquaintances
		//	'O': others
		counters := map[string]int{
			"P": 10,
			"A": 10,
			"C": 30,
			"O": 100,
		}
		// array of relationship types
		relTypes := [...]string{"P", "A", "C", "O"}
		// edge map
		//edjeMap := make(map[uint32]bool, nNodes)

		for i := 0; i < nNodes; i++ {

			rand.Seed(time.Now().UnixNano())

			// Days where the node will infect others
			// custom Normal Distribution
			tmpR0 := int(math.Round(rand.NormFloat64()*stdR0 + medianR0))
			if tmpR0 < 0 {
				tmpR0 = 0
			}

			infectiveDays := make([]int8, tmpR0)

			for r := 0; r < tmpR0; r++ {
				infectiveDays[r] = int8(rand.Intn(infectiveEpochs-incubationEpochs) + incubationEpochs)
			}

			newNode := person{
				Infective:       false,
				Survived:        false,
				Dead:            false,
				InfectiveEpochs: infectiveEpochs,
				InfectiveDays:   infectiveDays,
			}

			// this index is used to access relTypes byte array
			currentIndex := 0

			// Initialize Relationships
			for j := 0; j < nEdges; j++ {
				// generate a random ID
				edgeID := uint32(rand.Intn(nNodes))
				// check that the random ID is not equal to the vertex we are considering

				contained := false
				for _, a := range newNode.Edges {
					if a == edgeID {
						contained = true
					}
				}

				if edgeID != uint32(i) && !contained {
					// initialize the relation struct with the random ID
					newNode.Edges = append(newNode.Edges, edgeID)
					newNode.RelationType = append(newNode.RelationType, relTypes[currentIndex])
				}
				// check if it is the last element of index to add
				if counters[relTypes[currentIndex]] > 1 {
					// decreasing counter in every case if it's not the last element
					counters[relTypes[currentIndex]]--
				} else {
					// check if increasing current index is possible
					if currentIndex+1 < len(relTypes) {
						currentIndex++
					}
				}

			}
			// copy the node into the network array
			network[i] = newNode

		}
		log.Println("Network nodes allocated.")

		runtime.GC()
		log.Println("Garbage Collector freed.")

		// Save the network
		if *saveNetwork {
			log.Println("Saving network on file..\nMarshaling...")

			file, _ := json.Marshal(network)

			log.Println("Marshaled.")
			log.Println("Writing on json file...")

			ioutil.WriteFile(folderName+"/"+*fileNetwork, file, 0644)

			log.Println("Written on json file.")

			runtime.GC()
			log.Println("Garbage Collector freed.")
		} else {
			log.Println("<Skip saving network>")
		}
	} else {
		log.Println("Loading network from file {}...", *fileNetwork)
		file, _ := ioutil.ReadFile(*fileNetwork)

		_ = json.Unmarshal([]byte(file), &network)

		// reset network to default
		resetNetwork(&network)

		log.Println("Network loaded.")
	}

	// Allocating Policy Measure
	muskPointer := &muskMeasure{
		Active:    true,
		FromEpoch: *muskEpoch,
		Psucc:     *muskProb,
	}
	socialDistancingPointer := &socialDistancingMeasure{
		Active:    true,
		FromEpoch: *socDisEpoch,
		AllowContacts: map[string]bool{
			"P": true,
			"A": true,
			"C": false,
			"O": false,
		},
	}
	// Allocating SSN (National Healthcare System)
	ssnPointer := &nationalHealthcareSystem{
		intensiveCare:    bedIntensiveCare,
		subIntensiveCare: bedSubIntensiveCare,
	}

	// Montecarlo Simulation
	for i := 0; i < *mctrials; i++ {
		log.Println("TRIAL:\t", i, "______________________________")
		spreadingDesease(&network, simulationEpochs, &epochsResults, muskPointer, socialDistancingPointer, ssnPointer, &trialsResults, &ssnEpochsResults, i)
		log.Println("clear graph network...")
		resetNetwork(&network)
		// reset national healthcare system
		ssnPointer = &nationalHealthcareSystem{
			intensiveCare:    bedIntensiveCare,
			subIntensiveCare: bedSubIntensiveCare,
		}

		if *computeCI {
			log.Println("Save for compute CI on csv...")

			csvFile, err := os.Create(folderName + "/trials/" + strconv.Itoa(i) + "_trial_results.csv")

			if err != nil {
				log.Fatalf("failed creating file: %s", err)
			}

			csvwriter := csv.NewWriter(csvFile)

			for _, epoch := range epochsResults {
				// convert row of []int into []string
				var row [len(epoch)]string
				for i := 0; i < len(epoch); i++ {
					row[i] = strconv.Itoa(epoch[i])
				}
				err = csvwriter.Write(row[:])

				if err != nil {
					log.Println("ERROR ON CREATING CSV:", err)
				}
			}

			csvwriter.Flush()
			csvFile.Close()

			log.Println("Saved for compute CI and closed csv.")
		}

	}

	if *computeCI {
		log.Println("Save for compute CI on csv...")

		csvFile, err := os.Create(folderName + "/simulation_trials_results.csv")

		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}

		csvwriter := csv.NewWriter(csvFile)

		for _, trial := range trialsResults {
			// convert row of []int into []string
			var row [len(trial)]string
			for i := 0; i < len(trial); i++ {
				row[i] = strconv.Itoa(trial[i])
			}
			err = csvwriter.Write(row[:])

			if err != nil {
				log.Println("ERROR ON CREATING CSV:", err)
			}
		}

		csvwriter.Flush()
		csvFile.Close()

		log.Println("Saved for compute CI and closed csv.")
	}

	if *computeSSN {
		log.Println("Save for compute SSN on csv...")

		csvFile, err := os.Create(folderName + "/simulation_ssn_results.csv")

		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}

		csvwriter := csv.NewWriter(csvFile)

		for _, epoch := range ssnEpochsResults {
			// convert row of []int into []string
			var row [len(epoch)]string
			for i := 0; i < len(epoch); i++ {
				row[i] = strconv.Itoa(epoch[i])
			}
			err = csvwriter.Write(row[:])

			if err != nil {
				log.Println("ERROR ON CREATING CSV:", err)
			}
		}

		csvwriter.Flush()
		csvFile.Close()

		log.Println("Saved for compute SSN and closed csv.")
	}

	log.Println("Save results on csv...")

	csvFile, err := os.Create(folderName + "/simulation.csv")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(csvFile)

	for _, epoch := range epochsResults {
		// convert row of []int into []string
		var row [len(epoch)]string
		for i := 0; i < len(epoch); i++ {
			row[i] = strconv.Itoa(epoch[i])
		}
		err = csvwriter.Write(row[:])

		if err != nil {
			log.Println("ERROR ON CREATING CSV:", err)
		}
	}

	csvwriter.Flush()
	csvFile.Close()

	log.Println("Saved and closed csv.")

	// call python script *working*
	if *runPyScript {
		if *computeCI && *computeSSN {
			log.Println("Calling python CI script...")

			pathOfData := "--trialsFile=simulation_trials_results.csv"
			pathOfFolder := "--folder=" + folderName + "/"
			plotSSN := "--ssnFile=simulation_ssn_results.csv"
			nTrials := "--trials=" + strconv.Itoa(*mctrials)

			log.Println("python3 ./Scripts/plotgraphs.py", pathOfData, pathOfFolder, plotSSN, nTrials)

			out, err := exec.Command("python", "./Scripts/plotgraphs.py", pathOfData, pathOfFolder, plotSSN, nTrials).Output()

			if err != nil {
				log.Println(string(out))
				log.Panicln("ERROR ON EXECUTING PYTHON SCRIPT", err)
			}

			log.Println("Output:\n---\n\n", string(out), "\n----")

			errOutputFile := ioutil.WriteFile(folderName+"/output.txt", out, 0644)
			if errOutputFile != nil {
				log.Panicln("ERROR ON WRITING OUTPUT OF PYTHON SCRIPT", err)
			}
		}

	} else {
		log.Println("<Skip calling python script>")
	}

}
