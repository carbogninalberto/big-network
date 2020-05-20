package main

import (
	"log"
	"math/rand"
	"runtime"
	"time"
)

const (
	nNodes       = 4905854 //number of people in Veneto
	nEdges       = 150
	cpus         = 1
	nTheoryNodes = 4905854
	bedPlaces    = 450 //https://www.aulss2.veneto.it/amministrazione-trasparente/disposizioni-generali/atti-generali/regolamenti?p_p_id=101&p_p_lifecycle=0&p_p_state=maximized&p_p_col_id=column-1&p_p_col_pos=22&p_p_col_count=24&_101_struts_action=%2Fasset_publisher%2Fview_content&_101_assetEntryId=10434368&_101_type=document
)

type person struct {
	//Edges     []*person
	Edges     []uint32 `json:"Edges"`
	Infective bool     `json:"Infective"`
	Survived  bool     `json:"Survived"`
	Dead      bool     `json:"Dead"`
}

func main() {
	runtime.GOMAXPROCS(16)
	//c := make(chan bool, cpus)

	// random seed
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	// initialize network
	network := make([]person, nNodes)

	log.Println("Creating network...")

	for i := 0; i < nNodes; i++ {

		newNode := person{
			Infective: false,
			Survived:  false,
			Dead:      false,
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

	log.Println("Marshaling...")
	//file, _ := json.Marshal(network)

	log.Println("Marshaled.")

	runtime.GC()
	log.Println("Garbage Collector freed.")

	//_ = ioutil.WriteFile("network.json", file, 0644)

}
