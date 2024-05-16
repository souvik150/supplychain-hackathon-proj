package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"

	types "github.com/souvik150/supplychain-hackathon-proj/types"
)

const wsEndpoint = "ws://127.0.0.1:30000/ws"
var sendInterval = 1 * time.Second



func genLatLong() (float64, float64) {
	return genCoord(), genCoord()
}

func genCoord() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()
	return n + f
}

func main(){
	obuIDS := generateOBUIDS(20)
	conn, _ , err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		for i := 0; i<len(obuIDS); i++ {
			lat, long := genLatLong()
			data := types.OBUData{
				ODUID: obuIDS[i],
				Lat : lat,
				Long: long,
			}
			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%+v\n", data)
		}
		time.Sleep(sendInterval)
	}
}

func generateOBUIDS(n int) []int {
	ids := make([]int, n)
	for i := 0; i<n; i++ {
		ids[i] = rand.Intn(math.MaxInt)
	}
	return ids
}

func init(){
	rand.Seed(time.Now().UnixNano())
}
