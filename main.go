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

// Hyperparameters configuration of Simulation
const (
	nNodes           = 49058 //5 //4 // 4905854 number of people in Veneto
	nEdges           = 150   //Dunbar number 150
	bedPlaces        = 450   //https://www.aulss2.veneto.it/amministrazione-trasparente/disposizioni-generali/atti-generali/regolamenti?p_p_id=101&p_p_lifecycle=0&p_p_state=maximized&p_p_col_id=column-1&p_p_col_pos=22&p_p_col_count=24&_101_struts_action=%2Fasset_publisher%2Fview_content&_101_assetEntryId=10434368&_101_type=document
	medianR0         = 2.28  //https://pubmed.ncbi.nlm.nih.gov/32097725/ 2.06-2.52 95% CI 0,22/1.96 = 0.112
	stdR0            = 0.112 //0.112
	infectiveEpochs  = 14
	simulationEpochs = 63
	trials           = 1
	deadRate         = 0.054
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

// spreadingDesease runs a simulation over n epochs on a bigNet ([]person)
func spreadingDesease(networkPointer *bigNet, epochs int, epochsResultsPointer *[simulationEpochs][5]int) error {
	for epoch := 0; epoch < epochs; epoch++ {

		if epoch == 0 {
			// pick a random infect over the graph
			case0 := rand.Intn(nNodes)
			(*networkPointer)[case0].Infective = true
			log.Println("CASE 0:", case0)
			infectiveDaysLen := len((*networkPointer)[case0].InfectiveDays)

			for day := 0; day < infectiveDaysLen; day++ {
				if (*networkPointer)[case0].InfectiveDays[day] == 0 {

					// Add support to different measures
					randomInfect := rand.Intn(len((*networkPointer)[case0].Edges))
					infected := (*networkPointer)[case0].Edges[randomInfect]
					if (*networkPointer)[infected].InfectiveEpochs > 0 {
						(*networkPointer)[infected].Infective = true
					}

					// I set to -1 in order to not consider it anymore
					(*networkPointer)[case0].InfectiveDays[day] = -1

				} else if (*networkPointer)[case0].InfectiveDays[day] > 0 {

					// Add support to different measures
					randomInfect := rand.Intn(len((*networkPointer)[case0].Edges))
					infected := (*networkPointer)[case0].Edges[randomInfect]
					if (*networkPointer)[infected].InfectiveEpochs > 0 {
						(*networkPointer)[infected].Infective = true
					}
					(*networkPointer)[case0].InfectiveDays[day]--
				}
			}
			// make time pass and reduce the remaining infective days
			_ = reduceInfectiveEpochs(&(*networkPointer)[case0])

		} else {
			infected := getInfected(networkPointer)

			for _, infectedID := range infected {

				infectiveDaysLen := len((*networkPointer)[infectedID].InfectiveDays)

				for day := 0; day < infectiveDaysLen; day++ {
					if (*networkPointer)[infectedID].InfectiveDays[day] == 0 {

						// Add support to different measures
						randomInfect := rand.Intn(len((*networkPointer)[infectedID].Edges))
						infected := (*networkPointer)[infectedID].Edges[randomInfect]

						if (*networkPointer)[infected].Infective == false &&
							(*networkPointer)[infected].Dead == false &&
							(*networkPointer)[infected].Survived == false &&
							(*networkPointer)[infected].InfectiveEpochs > 0 {
							(*networkPointer)[infected].Infective = true
						}

						// I set to -1 in order to not consider it anymore
						(*networkPointer)[infectedID].InfectiveDays[day] = -1
					} else if (*networkPointer)[infectedID].InfectiveDays[day] > 0 {
						// Add support to different measures
						randomInfect := rand.Intn(len((*networkPointer)[infectedID].Edges))
						infected := (*networkPointer)[infectedID].Edges[randomInfect]

						if (*networkPointer)[infected].Infective == false &&
							(*networkPointer)[infected].Dead == false &&
							(*networkPointer)[infected].Survived == false &&
							(*networkPointer)[infected].InfectiveEpochs > 0 {
							(*networkPointer)[infected].Infective = true
						}
						(*networkPointer)[infectedID].InfectiveDays[day]--
					}
				}

				// make time pass and reduce the remaining infective days
				_ = reduceInfectiveEpochs(&(*networkPointer)[infectedID])
			}
		}

		infectNumber := countInfected(networkPointer, true, false, false)
		log.Println("EPOCH\t", epoch, "\tINFECTED:\t", infectNumber)

		// number of infected today
		(*epochsResultsPointer)[epoch][0] = infectNumber
		// new number of infected today regards yesterday
		if epoch != 0 {
			lastInfected := (*epochsResultsPointer)[epoch-1][0]
			(*epochsResultsPointer)[epoch][1] = infectNumber - int(lastInfected)
		} else {
			(*epochsResultsPointer)[epoch][1] = infectNumber
		}

		// number of total infected
		(*epochsResultsPointer)[epoch][2] = countTotalInfected(networkPointer)
		// number of total survived
		(*epochsResultsPointer)[epoch][3] = countInfected(networkPointer, false, true, false)
		// number of total dead
		(*epochsResultsPointer)[epoch][4] = countInfected(networkPointer, false, false, true)

		runtime.GC()
	}
	return nil
}

func main() {
	runtime.GOMAXPROCS(8)

	// flags
	loadNetwork := flag.Bool("loadnet", false, "default value is false, if true it load a network from a file called Network.json, to change the loading file name check flag namenet")
	saveNetwork := flag.Bool("savenet", false, "default value is false, if true saves network on timestamp/Network.json")
	fileNetwork := flag.String("namenet", "Network.json", "default value is Network.json, it's the name of the network file")
	mctrials := flag.Int("mctrials", 1, "default value is 1, you can choose how many trials run on the Montecarlo Simulation")
	computeCI := flag.Bool("computeCI", false, "default value is false, set to true when use flag -mctrials > 1 to get Confidence Intervals of metrics")
	runPyScript := flag.Bool("runpyscript", false, "default valuse is false, set to true if you want to print graphs of simulation with matplotlib")
	flag.Parse()

	// random seed
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	// initialize network
	network := make(bigNet, nNodes)
	var epochsResults [simulationEpochs][5]int
	// creating run folder
	folderName := strconv.Itoa(int(time.Now().UnixNano()))
	os.MkdirAll(folderName, os.ModePerm)

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
				infectiveDays[r] = int8(rand.Intn(infectiveEpochs))
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

	// Montecarlo Simulation
	for i := 0; i < *mctrials; i++ {
		log.Println("TRIAL:\t", i, "______________________________")
		spreadingDesease(&network, simulationEpochs, &epochsResults)
		log.Println("clear graph network...")
		resetNetwork(&network)
		if *computeCI {

		}
		// compute CI ecc
		// CI: INFETTI TOTALI
		// CI: MORTI TOTALI
		// CI: GUARITI TOTALI

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
		log.Println("Calling Python script...")

		out, err := exec.Command("python", "./plotgraphs.py").Output()

		if err != nil {
			log.Panicln("ERROR ON EXECUTING PYTHON SCRIPT", err)
		}

		log.Println("Output:\n---\n\n", string(out), "\n----")
	} else {
		log.Println("<Skip calling python script>")
	}

}
