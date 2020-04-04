package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	papers = iota
	tobacco
	matches
)

var resourcesMap = map[int]string{
	papers:  "Papers",
	tobacco: "Tobacco",
	matches: "Matches",
}

var keepersMap = map[int]string{
	papers:  "PapersKeeper",
	tobacco: "TobaccoKeeper",
	matches: "MatchesKeeper",
}

type Table struct {
	papers  chan int
	tobacco chan int
	matches chan int
}

func main() {
	var wg *sync.WaitGroup = new(sync.WaitGroup)
	var table = new(Table)

	table.papers = make(chan int, 1)
	table.tobacco = make(chan int, 1)
	table.matches = make(chan int, 1)

	var signals [3]chan int

	for i := 0; i < 3; i++ {
		s := make(chan int, 1)
		signals[i] = s
		go smoker(table, i, s, wg)
	}

	bartender(table, signals, wg)
}

func bartender(table *Table, smokers [3]chan int, wg *sync.WaitGroup) {
	fmt.Println("bartender started\n")
	for {
		time.Sleep(time.Second * 1)
		selected := rand.Intn(3) // select random smoker to make cigarette

		// send data to channels (put resources on table)
		switch selected {
		case papers:
			table.tobacco <- 1
			table.matches <- 1

		case tobacco:
			table.papers <- 1
			table.matches <- 1

		case matches:
			table.papers <- 1
			table.tobacco <- 1
		}

		// send selected to smoker channel
		for _, smoker := range smokers {
			smoker <- selected
		}
		// add running to wait group and then wait
		wg.Add(1)
		wg.Wait()
	}
}

func smoker(t *Table, selected int, signal chan int, wg *sync.WaitGroup) {
	for {
		// check if bartender selected current smoker to make cig
		if selected != <-signal {
			continue
		}

		// print current resources
		fmt.Printf("Papers: %d | tobacco: %d | matches: %d\n", len(t.papers), len(t.tobacco), len(t.matches))

		// 2 select blocks to get 2 resources from channels
		// because each takes one only
		select {
		case <-t.papers:
		case <-t.tobacco:
		case <-t.matches:
		}
		time.Sleep(10 * time.Millisecond)
		select {
		case <-t.papers:
		case <-t.tobacco:
		case <-t.matches:
		}
		fmt.Printf("%v smokes a cigarette\n\n", keepersMap[selected])
		time.Sleep(time.Millisecond * 500)
		// decrement wait group counter
		wg.Done()
		time.Sleep(time.Millisecond * 100)
	}
}
