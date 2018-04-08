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

package vianueduserver

import (
	"net/http"
	"strconv"
)

// redirect redirects (duh) all requests from HTTP to HTTPS, if used, and sends a Moved Permanently (301) response code
// to the client.
func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://127.0.0.1:443"+r.RequestURI, http.StatusMovedPermanently)
}

// StartServer starts the server. It reads the configurations, creates the router, and boots the HTTP servers
// according to the configuration files.
func StartServer() {

	StartLoggers()

	HTTPLogger.Println("VianuEdu-Server v0.5-BETA ########################################################## BEGIN NEW LOG ##########################################################")
	HTTPLogger.Println("[BOOT] Reading configuration file...")

	listenPortInt := GetListenPort()

	listenPort := strconv.FormatInt(listenPortInt, 10)
	listenPort = ":" + listenPort

	HTTPLogger.Println("[BOOT] Done reading configuration file")
	HTTPLogger.Println("[BOOT] Initializing database backend...")

	ConnectToDatabase()

	HTTPLogger.Print("[BOOT] Configuring HTTP Server...")

	router := CreateRouter()

	isSSL, certFile, keyFile := ShouldIGoSecure()

	if isSSL {
		go http.ListenAndServeTLS(":443", certFile, keyFile, router)
		HTTPLogger.Println("[BOOT] Booting HTTP Server... DONE! Listening on port " + listenPort[1:] + "! TLS enabled on port 443!")
		http.ListenAndServe(listenPort, http.HandlerFunc(redirect))
	} else {
		HTTPLogger.Println("[BOOT] Booting HTTP Server... DONE! Listening on port " + listenPort[1:] + "! TLS disabled!")
		http.ListenAndServe(listenPort, router)
	}
}
