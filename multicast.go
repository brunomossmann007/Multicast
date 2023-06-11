package main

import (
	"fmt"
	"time"
)

type Message struct {
	ID        int
	Content   string
	Timestamp int
}

type Process struct {
	ID              int
	LogicalClock    int
	LocalTimestamps map[int]int
	FinalTimestamps map[int]int
	Acked           map[int]bool
	Delivered       map[int]bool
}

func (p *Process) Multicast(message Message, destinations []*Process) {
	p.LogicalClock++
	p.LocalTimestamps[message.ID] = p.LogicalClock
	fmt.Printf("Process %d multicast message ID %d\n", p.ID, message.ID)

	for _, dest := range destinations {
		dest.LocalTimestamps[message.ID] = p.LocalTimestamps[message.ID]
		fmt.Printf("Process %d sends local timestamp for message ID %d to process %d\n", p.ID, message.ID, dest.ID)
	}

	p.Acked[message.ID] = true
	fmt.Printf("Process %d acks message ID %d\n", p.ID, message.ID)

	p.TryDeliver()
}

func (p *Process) TryDeliver() {
	for messageID := range p.FinalTimestamps {
		if !p.Delivered[messageID] && p.Acked[messageID] {
			canDeliver := true
			for _, dest := range p.LocalTimestamps {
				if (p.FinalTimestamps[messageID] < p.FinalTimestamps[dest] && dest != messageID) || (p.FinalTimestamps[messageID] < p.LocalTimestamps[dest]) {
					canDeliver = false
					break
				}
			}
			if canDeliver {
				p.Delivered[messageID] = true
				fmt.Printf("Process %d delivers message ID %d\n", p.ID, messageID)
			}
		}
	}
}

func main() {
	process1 := &Process{
		ID:              1,
		LogicalClock:    0,
		LocalTimestamps: make(map[int]int),
		FinalTimestamps: make(map[int]int),
		Acked:           make(map[int]bool),
		Delivered:       make(map[int]bool),
	}
	process2 := &Process{
		ID:              2,
		LogicalClock:    0,
		LocalTimestamps: make(map[int]int),
		FinalTimestamps: make(map[int]int),
		Acked:           make(map[int]bool),
		Delivered:       make(map[int]bool),
	}
	process3 := &Process{
		ID:              3,
		LogicalClock:    0,
		LocalTimestamps: make(map[int]int),
		FinalTimestamps: make(map[int]int),
		Acked:           make(map[int]bool),
		Delivered:       make(map[int]bool),
	}

	message := Message{
		ID:        1,
		Content:   "Hello, world!",
		Timestamp: 0,
	}
	process1.Multicast(message, []*Process{process2, process3})

	// Aguarda um pouco antes de finalizar o programa
	time.Sleep(1 * time.Second)
}
