package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const DELTA_TIME = 10

func putHandler(writer http.ResponseWriter, request *http.Request) {
	buf, _ := ioutil.ReadAll(request.Body)

	reader := csv.NewReader(bytes.NewBuffer(buf))
	req, _ := reader.ReadAll()

	now := time.Now()

	data := make([]float32, len(req[0])-1)
	timestamps := make([]time.Time, len(req[0])-1)

	for i := 0; i < len(data); i++ {
		f, _ := strconv.ParseFloat(req[0][i], 32)
		data[i] = float32(f)
		dur := time.Duration((len(data)-1-i)*DELTA_TIME) * time.Second
		timestamps[i] = now.Add(-dur)
	}

	fmt.Println(timestamps)
	fmt.Println(data)
}

func main() {

	http.HandleFunc("/put", putHandler)

	//Setup DataBase 1st Time

	//Init Database

	err := http.ListenAndServe(":8008", nil)
	if err != nil {
		log.Fatal("HTTP Server Error", err)
	}
}
