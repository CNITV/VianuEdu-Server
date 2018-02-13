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
	"github.com/gorilla/mux"
	"net/http"
	log "github.com/sirupsen/logrus"
)

func CreateRouter() *mux.Router {

	router := mux.NewRouter()

	log.WithFields(log.Fields{}).Info("Configuring route handling for API...")
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)

		log.WithFields(log.Fields{
			"method":  route.Method,
			"path":    route.Pattern,
			"name":    route.Name,
			"handler": route.HandlerFunc,
		}).Info("Configured route for " + route.Name)
	}

	log.WithFields(log.Fields{}).Info("Configuring static site mapping...")
	router.Methods("GET").
		PathPrefix("/").
		Name("StaticSiteRoute").
		Handler(http.StripPrefix("/", http.FileServer(http.Dir("static/"))))

	return router

}
