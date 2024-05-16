package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/souvik150/supplychain-hackathon-proj/types"
)

func main(){
	recv, err := NewDataReciever()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/ws", recv.handleWS)

	http.ListenAndServe(":30000", nil)
}

type DataReciever struct {
	msgCh chan types.OBUData
	conn *websocket.Conn
	
	prod DataProducer

}

func NewDataReciever() (*DataReciever, error) {
	var (
		p DataProducer
		err error
		kafkaTopic = "obudata"
	)
	
	p, err = NewKafkaProducer(kafkaTopic)
	if err != nil {
		return nil, err
	}

	p = NewLogMiddleware(p)

	return &DataReciever{
		msgCh: make(chan types.OBUData, 128),
		prod: p,
	}, nil
}

func (dr *DataReciever) produceData(data types.OBUData) error {
	return dr.prod.ProduceData(data)
}

func(dr *DataReciever) handleWS(w http.ResponseWriter, r *http.Request){
	u := websocket.Upgrader{
		ReadBufferSize: 1028,
		WriteBufferSize: 1028,
	}
	conn, err := u.Upgrade(w, r, nil)
	if(err != nil){
		log.Fatal(err)
	}
	dr.conn = conn

	go dr.wsRecieveLoop()
}

func(dr *DataReciever) wsRecieveLoop(){
	fmt.Println("New OBU Client Connected!")
	for {
		data := types.OBUData{}
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("Error reading json: ", err)
			continue
		}

		// fmt.Printf("Received OBU data from [%d] :: <lat %.2f, long %.2f> \n", data.ODUID, data.Lat, data.Long)
		// dr.msgCh <- data

		err := dr.produceData(data); 
		if err != nil {
			log.Println("Error producing data to kafka: ", err)
		}
		// dr.prod.Flush(15 * 1000)
	}
}