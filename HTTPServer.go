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

func main() {

	HTTPLogger.Println("VianuEdu-Server v0.2 ########################################################## BEGIN NEW LOG ##########################################################")
	HTTPLogger.Println("[BOOT] Reading configuration file...")

	listenPortInt := GetListenPort()

	listenPort := strconv.FormatInt(listenPortInt, 10)
	listenPort = ":" + listenPort

	HTTPLogger.Println("[BOOT] Done reading configuration file")
	HTTPLogger.Print("[BOOT] Configuring HTTP Server...")

	router := CreateRouter()

	HTTPLogger.Println("[BOOT] Booting HTTP Server... DONE! Listening on port " + listenPort[1:])

	http.ListenAndServe(listenPort, router)

}
