package main

import (
	"bytes"
	"encoding/csv"
	"github.com/RobinUS2/golang-moving-average"
	"github.com/wcharczuk/go-chart"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const DELTA_TIME = 60

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

	go AddValuesToDB(data, timestamps)

	writer.WriteHeader(200)
	writer.Write([]byte("OK"))
}

func graphHandler(writer http.ResponseWriter, r *http.Request) {

	uri, err := url.ParseRequestURI(r.RequestURI)
	query := uri.Query()
	tim, ok2 := query["t"]
	ti, err := strconv.Atoi(tim[0])

	values, times := GetValuesFromDB(time.Now().Add(-8 * time.Hour))
	if ok2 && err == nil {
		values, times = GetValuesFromDB(time.Now().Add(-time.Duration(ti) * time.Hour))
	}
	/*
		p, err := plot.New()
		if err != nil {
			log.Fatal(err)
		}

		p.Title.Text = "Plotutil example"
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"

		pts := make(plotter.XYs, len(values))
		for i := range pts {
			pts[i].X = float64(times[i].Day()*24*60 +times[i].Hour()*60 + times[i].Minute() + times[i].Second()/60)
			pts[i].Y = float64(values[i])
		}

		err = plotutil.AddLinePoints(p, "First", pts)
		if err != nil {
			log.Fatal(err)
		}

		w,err := p.WriterTo(4*vg.Inch, 4*vg.Inch, "svg" )
		w.WriteTo(writer)
	*/

	X := make([]float64, len(values))
	Y := make([]float64, len(values))
	e := movingaverage.New(5) 
	for i := range X {
		X[i] = float64(times[i].Day()*24*60 + times[i].Hour()*60 + times[i].Minute() + times[i].Second()/60)
		e.Add(float64(values[i]))
		Y[i] = e.Avg()
	}

	graph := chart.Chart{
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: X,
				YValues: Y,
			},
		},
	}

	//buffer := bytes.NewBuffer([]byte{})
	err = graph.Render(chart.PNG, writer)
	if err != nil {
		log.Fatal("HTTP Server Error", err)
	}

	writer.WriteHeader(200)
	//writer.Write([]byte("OK"))
}

func main() {

	http.HandleFunc("/put", putHandler)
	http.HandleFunc("/graph", graphHandler)
	//Setup DataBase 1st Time
	InitDB()
	defer CloseDB()
	//Init Database
	if !CheckBaseDB() {
		BuildBaseDB()
	}

	err := http.ListenAndServe(":8008", nil)
	if err != nil {
		log.Fatal("HTTP Server Error", err)
	}
	
}
