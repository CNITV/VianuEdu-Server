/*
 * This file is part of VianuEdu.
 *
 *  VianuEdu is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 *  VianuEdu is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with VianuEdu.  If not, see <http://www.gnu.org/licenses/>.
 *
 * Developed by Matei Gardus <matei@gardus.eu>
 */

package main

import (
	"net/http"
	"strconv"
)

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://127.0.0.1:443" + r.RequestURI, http.StatusMovedPermanently)
}

func main() {

	HTTPLogger.Println("VianuEdu-Server v0.2 ########################################################## BEGIN NEW LOG ##########################################################")
	HTTPLogger.Println("[BOOT] Reading configuration file...")

	listenPortInt := GetListenPort()

	listenPort := strconv.FormatInt(listenPortInt, 10)
	listenPort = ":" + listenPort

	HTTPLogger.Println("[BOOT] Done reading configuration file")
	HTTPLogger.Println("[BOOT] Initializing database backend...")

	ConnectToDatabase()

	HTTPLogger.Print("[BOOT] Configuring HTTP Server...")

	router := CreateRouter()

	/* go http.ListenAndServeTLS(":443", "keys/VianuEdu_Server.crt", "keys/VianuEdu_Server.key", router)
	 *
	 * Commented for development purposes, unfortunately you cannot make anything work when you have to add exceptions everywhere.
	 * However, when SSL is enabled, http.HandlerFunc(redirect) must be passed as router to http.ListenAndServe below.
	 */
	http.ListenAndServe(listenPort, router)

	HTTPLogger.Println("[BOOT] Booting HTTP Server... DONE! Listening on port " + listenPort[1:])

}
