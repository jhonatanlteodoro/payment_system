package workers

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"time"
)

type DummyWorker struct {
	data []string
}

func (m *DummyWorker) HasMessageToProcess() bool {
	if len(m.data) > 0 {
		return true
	}

	// mock some random items
	num := rand.Intn(10)
	if num%2 == 0 {
		m.data = append(m.data, strconv.Itoa(num))
		fmt.Println("Added value to process: ", num)
	}
	return false
}

func (m *DummyWorker) Process(serverDown chan os.Signal) {
	for {
		select {
		case <-serverDown:
			log.Println("Quit Signal Received - Worker is shutting down...")
			return
		default:
			if m.HasMessageToProcess() {
				fmt.Println("dummy msg: ", m.data[len(m.data)-1])
				m.data = slices.Delete(m.data, len(m.data)-1, len(m.data))
				continue
			}

			fmt.Println("No message to process")
			time.Sleep(2 * time.Second)
		}
	}
}

func StartWorker(serverDown chan os.Signal) {
	w := &DummyWorker{}
	w.Process(serverDown)
}
