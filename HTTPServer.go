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
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
	"strconv"
)

func init() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "log/HTTPServer.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	})
}

func main() {

	log.WithFields(log.Fields{}).Info("VianuEdu-Server v0.1 ----------- BEGIN NEW LOG ------------")
	log.Print("[BOOT] Reading configuration file...")

	listenPortInt := GetListenPort()

	listenPort := strconv.FormatInt(listenPortInt, 10)
	listenPort = ":" + listenPort

	log.Println("Done reading configuration file")
	log.Print("[BOOT] Configuring HTTP Server...")

	router := CreateRouter()
	log.Println("Booting HTTP Server... DONE! Listening on port " + listenPort[1:])

	http.ListenAndServe(listenPort, router)

}
