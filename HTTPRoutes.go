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
	"net/http"
	"github.com/stormfirefox1/VianuEdu-Server/API"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"APIHelloWorld",
		"GET",
		"/api",
		API.HelloWorld,
	},
}
