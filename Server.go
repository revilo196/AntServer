package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func main() {

	http.HandleFunc("/put", func(writer http.ResponseWriter, request *http.Request) {
		buf,_ := ioutil.ReadAll(request.Body)

		reader := csv.NewReader(bytes.NewBuffer(buf))
		req, _ := reader.ReadAll()


		data := make([]float32,len(req[0])-1)

		for i:=0; i<len(data) ; i++ {
			f,_ := strconv.ParseFloat(req[0][i],32)
			data[i] = float32(f)
		}

		fmt.Println(data)

	})

	http.ListenAndServe(":8008",nil)
}
