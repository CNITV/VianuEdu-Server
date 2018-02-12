package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"time"
	"os"
	"io/ioutil"
	"github.com/buger/jsonparser"
	"strconv"
)

func main() {

	log.Printf("VianuEdu-Server v0.1 -- LOG -- %s", time.Now())
	log.Print("Reading configuration...")

	configFile, err := os.Open("config/HTTPServer.json")
	if err != nil {
		log.Fatal(err)
	}
	mainConfig, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Fatal(err)
	}
	listenPortInt, err := jsonparser.GetInt(mainConfig, "listenPort")
	if err != nil {
		log.Fatal(err)
	}

	listenPort := strconv.FormatInt(int64(listenPortInt), 10)
	listenPort = ":" + listenPort

	log.Println("DONE")
	log.Print("Configuring HTTP Server...")

	router := mux.NewRouter()
	log.Println("Configuring static site mapping...")
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("static/"))))

	log.Println("DONE")
	log.Print("Booting HTTP Server... DONE! Listening on port " + listenPort[1:])

	http.ListenAndServe(listenPort, router)

}
