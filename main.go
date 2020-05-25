/*
IMPORTANT NOTE: 12gb of RAM at least required
*/
package main

import (
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"
	"encoding/csv"
	"strconv"
)

const (
	nNodes          = 4905854 //number of people in Veneto
	nEdges          = 150      //Dunbar number 150
	cpus            = 1
	nTheoryNodes    = 4905854
	bedPlaces       = 450 //https://www.aulss2.veneto.it/amministrazione-trasparente/disposizioni-generali/atti-generali/regolamenti?p_p_id=101&p_p_lifecycle=0&p_p_state=maximized&p_p_col_id=column-1&p_p_col_pos=22&p_p_col_count=24&_101_struts_action=%2Fasset_publisher%2Fview_content&_101_assetEntryId=10434368&_101_type=document
	r0              = 1
	medianR0        = 2.28 //https://pubmed.ncbi.nlm.nih.gov/32097725/ 2.06-2.52 95% CI 0,22/1.96 = 0.112
	infectiveEpochs = 14
	simulationEpochs = 100
)

type person struct {
	//Edges     []*person
	Edges           []relation	`json:"Edges"`
	Infective       bool     	`json:"Infective"`
	Survived        bool     	`json:"Survived"`
	Dead            bool     	`json:"Dead"`
	//Age 			bool 	 	`json:"Age`
	InfectiveEpochs uint32 // ottimizza, aumenta tot giorni per terapia intensiva (14+21)
}

type relation struct {
	Id				uint32
	relationship	byte

}

type bigNet []person
type resultMatrix [][]uint32

// spreadingDesease runs a simulation over n epochs on a bigNet ([]person)
func spreadingDesease(networkPointer *bigNet, epochs int, epochsResultsPointer *[simulationEpochs][3]string) error {
	for epoch := 0; epoch < epochs; epoch++ {
		// on epoch 0 choose a random node
		healedCounter := 0

		if epoch == 0 {
			case0 := rand.Intn(nTheoryNodes)
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
		infectNumber :=  countInfected(networkPointer)
		log.Println("EPOCH\t", epoch, "\tINFECTED:\t", infectNumber)
		// number of infected today
		(*epochsResultsPointer)[epoch][0] = string(infectNumber)
		// new number of infected today regards yesterday
		if epoch != 0 {
			lastInfected, _ := strconv.ParseInt((*epochsResultsPointer)[epoch-1][0], 10, 32)
			(*epochsResultsPointer)[epoch][1] = string(infectNumber-int(lastInfected))
		} else {
			(*epochsResultsPointer)[epoch][1] = string(infectNumber)
		}
		
		// number of people healed
		(*epochsResultsPointer)[epoch][0] = string(infectNumber)
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
	//c := make(chan bool, cpus)

	// random seed
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	// initialize network
	network := make(bigNet, nNodes)
	var epochsResults [simulationEpochs][3]string
	
	log.Println("Calling Python script...")

	// call python script

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

	log.Println("Creating network...")

	for i := 0; i < nNodes; i++ {

		newNode := person{
			Infective:       false,
			Survived:        false,
			Dead:            false,
			InfectiveEpochs: uint32(rand.Intn(infectiveEpochs)),
		}

		for j := 0; j < nEdges; j++ {
			edgeID := uint32(rand.Intn(nTheoryNodes))
			if edgeID != uint32(i) {
				newNode.Edges = append(newNode.Edges, edgeID)
			}
		}

		network[i] = newNode

	}
	log.Println("Network nodes allocated.")

	runtime.GC()
	log.Println("Garbage Collector freed.")

	// Montecarlo Simulation
	trials := 100
	for i := 0; i < trials; i++ {
		spreadingDesease(&network, simulationEpochs, &epochsResults)
		log.Println(trials, "\t infected...")
		// CI: INFETTI TOTALI
		// CI: MORTI TOTALI
		// CI: GUARITI TOTALI

	}
	

	//log.Println((&network))

	log.Println("Marshaling...")
	//file, _ := json.Marshal(network)

	log.Println("Marshaled.")

	runtime.GC()
	log.Println("Garbage Collector freed.")

	//_ = ioutil.WriteFile("network.json", file, 0644)

	log.Println("Save results on csv")
	csvFile, err := os.Create("simulation_results.csv")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)
 
	for _, epoch := range epochsResults {
		_ = csvwriter.Write(epoch[:])
	}

	csvwriter.Flush()

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
