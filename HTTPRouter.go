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
	"github.com/gorilla/mux"
	"net/http"
	"github.com/sirupsen/logrus"
	"os"
	"github.com/gorilla/handlers"
)

func CreateRouter() http.Handler {

	router := mux.NewRouter()

	HTTPLogger.WithFields(logrus.Fields{}).Info("[BOOT] Configuring route handling for API...")
	for _, route := range routes {

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)

		HTTPLogger.WithFields(logrus.Fields{
			"method":  route.Method,
			"path":    route.Pattern,
			"name":    route.Name,
			"handler": route.HandlerFunc,
		}).Info("[BOOT] Configured route for " + route.Name)
	}

	router.NotFoundHandler = http.HandlerFunc(dispense404)

	HTTPLogger.WithFields(logrus.Fields{}).Info("[BOOT] Configuring static site mapping...")
	router.Methods("GET").
		PathPrefix("/").
		Name("StaticSiteRoute").
		Handler(http.StripPrefix("/", http.FileServer(http.Dir("static/"))))

	requestLog, err := os.OpenFile("log/HTTPRequests.log",  os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		HTTPLogger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error opening HTTP requests log!")
	}

	loggedRouter := handlers.LoggingHandler(requestLog, router)

	return loggedRouter

}
