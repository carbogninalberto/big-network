package main

import (
	"log"
	"math/rand"
	"runtime"
	"time"
)

const (
	nNodes       = 4905854
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

func generateNodesWorker(id, start, end int, network *[nNodes]person, c chan bool) error {
	//log.Println(id, "START ALLOC")
	for i := start; i < end; i++ {
		//log.Println(rand.Intn(nNodes))
		if i%43201 == 0 { //43201 28813
			runtime.GC()
			//log.Println("adding node\t", i, "\tinterval [", start, ",", end, "]\t\t", id)
		}
		newNode := person{
			//Infective: false,
			//Survived:  false,
			//Dead:      false,
		}

		for j := 0; j < nEdges; j++ {
			if j%50 == 0 && j != 0 {
				//runtime.GC()
			}
			edgeID := uint32(rand.Intn(nTheoryNodes))
			if edgeID != uint32(i) {
				newNode.Edges = append(newNode.Edges, edgeID)
			}
		}

		(*network)[i] = newNode

		//log.Println("node", i, "added WOW!")
	}
	network = nil
	c <- true
	//time.Sleep(100 * time.Second)
	//log.Println("[NODE GENERATOR] Service", id, "finished.")
	return nil

}

/*
func generateEdgesWorker(id, start, end int, network *[nNodes]*person, c chan bool) {
	tenth := 0
	for i := start; i < end; i++ {
		//log.Println(rand.Intn(nNodes))
		if i%43201 == 0 { //43201 28813
			runtime.GC()
		}

		if i%((end-start)/10) == 0 {
			//log.Println(id, "\t", tenth*10, "%")
			tenth++
		}

		for j := 0; j < nEdges; j++ {
			if j%50 == 0 {
				runtime.GC()
			}
			edgeID := rand.Intn(nNodes)
			if edgeID != i {
				randomEdge := (*network)[edgeID]
				(*network)[i].edges = append((*network)[i].edges, randomEdge)
			}
		}
	}
	runtime.GC()
	c <- true
	//time.Sleep(100 * time.Second)
	log.Println("[EDGE GENERATOR] Service", id, "finished.")
}

*/

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	runtime.GOMAXPROCS(16)
	//c := make(chan bool, cpus)

	// random seed
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	// initialize network
	network := make([]person, nNodes)
	//networkPointer := &network

	//populate network
	/*
		slices := nNodes / cpus
		for j := 0; j < cpus; j++ {
			go generateNodesWorker(j, j*slices, (j+1)*slices, networkPointer, c)

		} */
	//time.Sleep(5 * time.Second)
	log.Println("Creating network...")

	for i := 0; i < nNodes; i++ {
		/*
			if i%500000 == 0 {
				println(i)
				runtime.GC()
				time.Sleep(1 * time.Second)
			}*/

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

		//log.Println("node", i, "added WOW!")
	}
	log.Println("Network nodes allocated.")
	time.Sleep(5 * time.Second)

	/*
		for i := 0; i < cpus; i++ {
			<-c
			log.Println(i, "END ALLOC")
		}
		log.Println("Network nodes allocated.")
		close(c)
		networkPointer = nil
	*/

	runtime.GC()
	log.Println("Garbage Collector freed.")
	//time.Sleep(2 * time.Second)

	log.Println("Marshaling...")
	//file, _ := json.Marshal(network)

	log.Println("Marshaled.")
	//time.Sleep(2 * time.Second)
	runtime.GC()
	log.Println("Garbage Collector freed.")
	//time.Sleep(2 * time.Second)

	//_ = ioutil.WriteFile("network.json", file, 0644)

	/*
		//populate edges
		for j := 0; j < cpus; j++ {
			go generateEdgesWorker(j, j*slices, (j+1)*slices, networkPointer, c)
		}

		log.Print("Creating edges...")

		for i := 0; i < cpus; i++ {
			<-c
		}

		log.Println("Network edges allocated.")

		log.Println(*network[0])
		/*
			for node := range network {
				if len((*network[node]).edges) != 0 {
					log.Println(*network[node])
				}
			}
	*/

	//

}
