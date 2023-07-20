package main

import (
	"github.com/gorilla/websocket"
	"github.com/rtsoy/toll-calculator/types"
	"log"
	"math"
	"math/rand"
	"time"
)

const (
	sendInterval = 1 * time.Second
	wsEndpoint   = "ws://127.0.0.1:30000/ws"
)

func generateOBUIDs(n int) []int {
	ids := make([]int, n)

	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(math.MaxInt)
	}

	return ids
}

func generateLatLong() (float64, float64) {
	return generateCoordinate(), generateCoordinate()
}

func generateCoordinate() float64 {
	var (
		n = float64(rand.Intn(100) + 1)
		f = rand.Float64()
	)

	return n + f
}

func main() {
	obuIDs := generateOBUIDs(20)
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	for {
		for i := 0; i < len(obuIDs); i++ {
			lat, long := generateLatLong()

			data := types.OBUData{
				OBUID: obuIDs[i],
				Lat:   lat,
				Long:  long,
			}

			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err)
			}
		}

		time.Sleep(sendInterval)
	}
}

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}
