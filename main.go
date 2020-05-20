/*
IMPORTANT NOTE: 12gb of RAM at least required
*/
package main

import (
	"log"
	"math/rand"
	"runtime"
	"time"
)

const (
	nNodes          = 4905854 //number of people in Veneto
	nEdges          = 150     //Dunbar number
	cpus            = 1
	nTheoryNodes    = 4905854
	bedPlaces       = 450 //https://www.aulss2.veneto.it/amministrazione-trasparente/disposizioni-generali/atti-generali/regolamenti?p_p_id=101&p_p_lifecycle=0&p_p_state=maximized&p_p_col_id=column-1&p_p_col_pos=22&p_p_col_count=24&_101_struts_action=%2Fasset_publisher%2Fview_content&_101_assetEntryId=10434368&_101_type=document
	r0              = 3
	infectiveEpochs = 2
)

type person struct {
	//Edges     []*person
	Edges           []uint32 `json:"Edges"`
	Infective       bool     `json:"Infective"`
	Survived        bool     `json:"Survived"`
	Dead            bool     `json:"Dead"`
	InfectiveEpochs uint32
}
type bigNet []person

// spreadingDesease runs a simulation over n epochs on a bigNet ([]person)
func spreadingDesease(networkPointer *bigNet, epochs int) error {
	for epoch := 0; epoch < epochs; epoch++ {
		// on epoch 0 choose a random node
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

			reduceInfectiveEpochs(&(*networkPointer)[case0])
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
				reduceInfectiveEpochs(&(*networkPointer)[infectedID])
			}

		}
		log.Println("EPOCH\t", epoch, "\tINFECTED:\t", countInfected(networkPointer))
		runtime.GC()
	}
	return nil
}

func reduceInfectiveEpochs(personPointer *person) {
	//log.Println("reduceInfective", personPointer.Infective, personPointer.Survived, personPointer.InfectiveEpochs)
	if personPointer.InfectiveEpochs > 1 {
		personPointer.InfectiveEpochs--
		//log.Println("personPointer.InfectiveEpochs > 1", personPointer.Infective, personPointer.Survived, personPointer.InfectiveEpochs)
	} else if personPointer.InfectiveEpochs == 1 {
		personPointer.InfectiveEpochs--
		personPointer.Infective = false
		personPointer.Survived = true
		//log.Println("personPointer.InfectiveEpochs == 1", personPointer.Infective, personPointer.Survived, personPointer.InfectiveEpochs)
	} else {
		log.Panicln("ERROR", personPointer.InfectiveEpochs)
	}
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

	spreadingDesease(&network, 100)

	log.Println((&network))

	log.Println("Marshaling...")
	//file, _ := json.Marshal(network)

	log.Println("Marshaled.")

	runtime.GC()
	log.Println("Garbage Collector freed.")

	//_ = ioutil.WriteFile("network.json", file, 0644)

}
