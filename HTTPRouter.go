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
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	_ "net/http/pprof"
	"os"
)

// CreateRouter is... a mess.
//
// It uses gorilla/mux to create a router, and it logs absolutely every route constructed. It reads the variables
// created in HTTPRoutes.go and declares the routes for each of the routes declared in that file. It logs every
// declaration, and, after that, adds a middleware to the router in order to log every HTTP request that pings the
// router.
//
// For more details about the specific routes, check HTTPRoutes.go
func CreateRouter() http.Handler {

	router := mux.NewRouter().StrictSlash(true)

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

	HTTPLogger.WithFields(logrus.Fields{}).Info("[BOOT] Configuring debug/pprof routes...")
	router.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)

	HTTPLogger.WithFields(logrus.Fields{}).Info("[BOOT] Configuring lessons download mapping...")
	router.Methods("GET").
		PathPrefix("/lessons/").
		Name("LessonsRoute").
		Handler(http.StripPrefix("/lessons/", http.FileServer(http.Dir("lessons/"))))

	HTTPLogger.WithFields(logrus.Fields{}).Info("[BOOT] Configuring static site mapping...")
	router.Methods("GET").
		PathPrefix("/").
		Name("StaticSiteRoute").
		Handler(http.StripPrefix("/", http.FileServer(http.Dir("static/"))))

	requestLog, err := os.OpenFile("log/HTTPRequests.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		HTTPLogger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error opening HTTP requests log!")
	}

	loggedRouter := handlers.LoggingHandler(requestLog, router)

	return loggedRouter

}
