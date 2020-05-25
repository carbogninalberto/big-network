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
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"
)

type bigNet []person

// Hyperparameters configuration of Simulation
const (
	//nTheoryNodes     = 4905854
	//nNodes = 1070340 // 1070340 number of people in Trentino-Alto Adige
	nNodes           = 4905854 // 4905854 number of people in Veneto
	nEdges           = 150     //Dunbar number 150
	cpus             = 1
	bedPlaces        = 450 //https://www.aulss2.veneto.it/amministrazione-trasparente/disposizioni-generali/atti-generali/regolamenti?p_p_id=101&p_p_lifecycle=0&p_p_state=maximized&p_p_col_id=column-1&p_p_col_pos=22&p_p_col_count=24&_101_struts_action=%2Fasset_publisher%2Fview_content&_101_assetEntryId=10434368&_101_type=document
	r0               = 1
	medianR0         = 2.28 //https://pubmed.ncbi.nlm.nih.gov/32097725/ 2.06-2.52 95% CI 0,22/1.96 = 0.112
	infectiveEpochs  = 14
	simulationEpochs = 63
	trials           = 1
)

type person struct {
	//Edges     []*person
	//Id        uint32     `json:Id`
	//Edges     []*relation `json:"Edges"`
	Edges        []uint32 `json:"Edges"`
	RelationType []byte   `json:"RelationType"`
	Infective    bool     `json:"Infective"`
	Survived     bool     `json:"Survived"`
	Dead         bool     `json:"Dead"`
	//Age 			bool 	 	`json:"Age`
	InfectiveEpochs uint32 // ottimizza, aumenta tot giorni per terapia intensiva (14+21)
}

// relationType
//	'P': family relationship
//	'A': friends
//	'C': acquaintances
//	'O': others
type relation struct {
	Id           uint32 `json:"Id"`
	RelationType byte   `json:"RelationType"`
}

// spreadingDesease runs a simulation over n epochs on a bigNet ([]person)
func spreadingDesease(networkPointer *bigNet, epochs int, epochsResultsPointer *[simulationEpochs][3]int) error {
	for epoch := 0; epoch < epochs; epoch++ {
		// on epoch 0 choose a random node
		healedCounter := 0

		if epoch == 0 {
			case0 := rand.Intn(nNodes)
			(*networkPointer)[case0].Infective = true
			(*networkPointer)[case0].InfectiveEpochs = infectiveEpochs
			log.Println("CASE 0:", case0)
			for r := 0; r < r0; r++ {
				randomInfect := rand.Intn(len((*networkPointer)[case0].Edges))
				infected := (*networkPointer)[case0].Edges[randomInfect]
				if (*networkPointer)[infected].InfectiveEpochs > 0 {
					(*networkPointer)[infected].Infective = true
				}

			}

			result := reduceInfectiveEpochs(&(*networkPointer)[case0])
			if result {
				healedCounter++
			}
		} else {
			infected := getInfected(networkPointer)

			for _, infectedID := range infected {
				for r := 0; r < r0; r++ {
					randomInfect := rand.Intn(len((*networkPointer)[infectedID].Edges))
					infected := (*networkPointer)[infectedID].Edges[randomInfect]

					if (*networkPointer)[infected].Infective == false &&
						(*networkPointer)[infected].Dead == false &&
						(*networkPointer)[infected].Survived == false &&
						(*networkPointer)[infected].InfectiveEpochs > 0 {
						(*networkPointer)[infected].Infective = true
					}

				}
				result := reduceInfectiveEpochs(&(*networkPointer)[infectedID])
				if result {
					healedCounter++
				}
			}

		}
		infectNumber := countInfected(networkPointer)
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

		// number of people healed
		(*epochsResultsPointer)[epoch][0] = infectNumber
		//log.Println(epoch, "*epochResultsPointer =", (*epochsResultsPointer)[epoch])
		runtime.GC()
	}
	return nil
}

func reduceInfectiveEpochs(personPointer *person) bool {
	//log.Println("reduceInfective", personPointer.Infective, personPointer.Survived, personPointer.InfectiveEpochs)
	if personPointer.InfectiveEpochs > 1 {
		personPointer.InfectiveEpochs--
		//log.Println("personPointer.InfectiveEpochs > 1", personPointer.Infective, personPointer.Survived, personPointer.InfectiveEpochs)
	} else if personPointer.InfectiveEpochs == 1 {
		personPointer.InfectiveEpochs--
		personPointer.Infective = false
		personPointer.Survived = true
		return true
		//log.Println("personPointer.InfectiveEpochs == 1", personPointer.Infective, personPointer.Survived, personPointer.InfectiveEpochs)
	} else {
		log.Panicln("ERROR", personPointer.InfectiveEpochs)
	}
	return false
}

func getInfected(networkPointer *bigNet) []int {
	infected := make([]int, 0, 1)
	for node := 0; node < nNodes; node++ {
		if (*networkPointer)[node].Infective == true &&
			(*networkPointer)[node].Survived == false &&
			(*networkPointer)[node].Dead == false {
			infected = append(infected, node)
		}
	}
	return infected
}

// countInfected counts the total number of infected people
func countInfected(networkPointer *bigNet) int {
	counter := 0
	for node := 0; node < nNodes; node++ {
		/*
			if (*networkPointer)[node].Infective == true ||
				(*networkPointer)[node].Survived == true ||
				(*networkPointer)[node].Dead == true {
				counter++
			}
		*/
		if (*networkPointer)[node].Infective == true {
			counter++
		}
	}
	//log.Println("INFECTED PEOPLE:", counter)
	return counter
}

func main() {
	runtime.GOMAXPROCS(8)

	// flags
	loadNetwork := flag.Bool("loadnet", false, "default value is false, if true it load a network from a file called Network.json, to change the loading file name check flag namenet")
	saveNetwork := flag.Bool("savenet", false, "default value is false, if true saves network on timestamp/Network.json")
	fileNetwork := flag.String("namenet", "Network.json", "default value is Network.json, it's the name of the network file")
	mctrials := flag.Int("mctrials", 1, "default value is 1, you can choose how many trials run on the Montecarlo Simulation")
	flag.Parse()

	// random seed
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	// initialize network
	network := make(bigNet, nNodes)
	var epochsResults [simulationEpochs][3]int
	// creating run folder
	folderName := strconv.Itoa(int(time.Now().UnixNano()))
	os.MkdirAll(folderName, os.ModePerm)

	// call python script *working*
	/*
		log.Println("Calling Python script...")
		cmd := exec.Command("python", "test.py")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		log.Println(cmd.Run())

		//err := cmd.Run()
		out, err := exec.Command("python", "./test.py").Output()

		if err != nil {
			log.Panicln("ERROR ON EXECUTING PYTHON SCRIPT", err)
		}

		log.Println("Output:\n---\n", string(out), "\n----")
	*/

	if !*loadNetwork {
		log.Println("Creating network...")
		// init counter for relations assignment
		//	'P': family relationship
		//	'A': friends
		//	'C': acquaintances
		//	'O': others
		counters := map[byte]int{
			'P': 10,
			'A': 10,
			'C': 30,
			'O': 100,
		}
		// array of relationship types
		relTypes := [...]byte{'P', 'A', 'C', 'O'}

		for i := 0; i < nNodes; i++ {
			/*
				if i%100000 == 0 && i != 0 {
					runtime.GC()
					log.Println(i)
				} */

			newNode := person{
				Infective:       false,
				Survived:        false,
				Dead:            false,
				InfectiveEpochs: uint32(rand.Intn(infectiveEpochs)),
			}

			// this index is used to access relTypes byte array
			currentIndex := 0

			// Initialize Relationships
			for j := 0; j < nEdges; j++ {
				// generate a random ID
				edgeID := uint32(rand.Intn(nNodes))
				// check that the random ID is not equal to the vertex we are considering
				if edgeID != uint32(i) {
					// initialize the relation struct with the random ID
					newNode.Edges = append(newNode.Edges, edgeID)
					newNode.RelationType = append(newNode.RelationType, relTypes[currentIndex])
					/*
						newNode.Edges = append(newNode.Edges, &relation{
							Id:           edgeID,
							RelationType: relTypes[currentIndex], // current Type of releation in generation
						}) */
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
			//ioutil.WriteFile("Network.json", file, 0644)

			/* init new code
			csvFile, err := os.Create(folderName + "/" + *fileNetwork)

			if err != nil {
				log.Fatalf("failed creating file: %s", err)
			}

			csvwriter := csv.NewWriter(csvFile)

			for _, node := range network {
				// convert row of []int into []string
				var row [...]string{}
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
			*/
			// end new code

			runtime.GC()
			log.Println("Garbage Collector freed.")
		} else {
			log.Println("Skip saving network")
		}
	} else {
		log.Println("Loading network from file {}...", *fileNetwork)
		file, _ := ioutil.ReadFile(*fileNetwork)

		_ = json.Unmarshal([]byte(file), &network)
		log.Println("Network loaded.")
	}

	// Montecarlo Simulation

	for i := 0; i < *mctrials; i++ {
		spreadingDesease(&network, simulationEpochs, &epochsResults)
		log.Println(i, "\t statistics...")
		// CI: INFETTI TOTALI
		// CI: MORTI TOTALI
		// CI: GUARITI TOTALI

	}

	//log.Println((&network))

	//_ = ioutil.WriteFile("network.json", file, 0644)

	log.Println("Save results on csv")

	csvFile, err := os.Create("simulation_results.csv")

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

}

/*

package main

import "fmt"
import "os/exec"

func main() {
    cmd := exec.Command("python",  "-c", "import pythonfile; print pythonfile.cat_strings('foo', 'bar')")
    fmt.Println(cmd.Args)
    out, err := cmd.CombinedOutput()
    if err != nil { fmt.Println(err); }
    fmt.Println(string(out))
}
exec.Command("script.py").Run()
*/
