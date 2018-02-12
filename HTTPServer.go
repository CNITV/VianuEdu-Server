/*
 *      This file is part of VianuEdu.
 *
 *      VianuEdu is free software: you can redistribute it and/or modify
 *      it under the terms of the GNU General Public License as published by
 *      the Free Software Foundation, either version 3 of the License, or
 *      (at your option) any later version.
 *
 *      VianuEdu is distributed in the hope that it will be useful,
 *      but WITHOUT ANY WARRANTY; without even the implied warranty of
 *      MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *      GNU General Public License for more details.
 *
 *      You should have received a copy of the GNU General Public License
 *      along with VianuEdu.  If not, see <http://www.gnu.org/licenses/>.
 *
 *      Developed by Matei Gardus <matei@gardus.eu>
 */

package main

import (
	"github.com/buger/jsonparser"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {

	log.SetOutput(&lumberjack.Logger{
		Filename:   "log/HTTPServer.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	})

	log.WithFields(log.Fields{
		"time":     time.Now(),
		"function": "main",
	}).Info("VianuEdu-Server v0.1 -- BEGIN LOG")
	log.Print("[BOOT] Reading configuration file...")

	configFile, err := os.Open("config/HTTPServer.json")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Error opening HTTPServer configuration file!")
	}
	defer configFile.Close()

	mainConfig, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Error reading HTTPServer configuration file!")
	}
	listenPortInt, err := jsonparser.GetInt(mainConfig, "listenPort")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Error parsing HTTPServer configuration file! (can't parse listenPort)")
	}

	listenPort := strconv.FormatInt(int64(listenPortInt), 10)
	listenPort = ":" + listenPort

	log.Println("DONE")
	log.Print("[BOOT] Configuring HTTP Server...")

	router := mux.NewRouter()
	log.Println("[BOOT] Configuring static site mapping...")
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("static/"))))

	log.Println("DONE")
	log.Print("Booting HTTP Server... DONE! Listening on port " + listenPort[1:])

	http.ListenAndServe(listenPort, router)

}
