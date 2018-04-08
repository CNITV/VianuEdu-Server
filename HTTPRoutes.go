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
)

// A Route is a variable that can represents a route to be digested by the mux router declared in HTTPRouter.go
// It contains all the parameters required for such a route to be declared.
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// A Routes variable is merely a slice of Routes. That's it. Too lazy to create a slice literally ONE row below.
// Why, you may ask? I have no idea. Sometimes I over-compartimentalize all of my code and, therefore, end up with
// stupid declarations such as this one. Maybe this is a life lesson for myself. Maybe I will learn from this experience
// and realize that I should take Object-Oriented Programming down a peg, and not completely break everything down in
// little pieces. Also, it might be a good idea to stop writing documentation for a type only used ONCE in the entire
// project, since this has actually become the most documented part of the entire project. Oh well.
type Routes []Route

// This variable contains all of the routes used by VianuEdu-Server that are not static site routes. You can use this
// segment of code as an API documentation of sorts, since it contains all the possible entry points to the project.
// In fact, this will probably remain the official API documentation for VianuEdu-Server.
//
// A guide for the API documentation would be the following:
//		-The name of each route represents the name of a request
// 		-The method of each route represents the HTTP method accepted by the router for the specific entry point.
//		-The pattern of each route represents the URL required to access the API. i.e. http://www.example.com/[PATTERN]
//		-The HandlerFunc of each route points to the HandlerFunc to which the router will take the request to. All
//		 HandlerFuncs are found in API*.go files.
//
// Reviewing this part of the source code allows for you to easily access all of the ways that you can query this
// server for.
var routes = Routes{

	Route{
		"GetStudent",
		"GET",
		"/api/getStudent/{id}",
		getStudent,
	},
	Route{
		"FindStudentID",
		"POST",
		"/api/findStudentID",
		findStudentID,
	},
	Route{
		"RegisterStudent",
		"POST",
		"/api/registerStudent",
		registerStudent,
	},
	Route{
		"GetTeacher",
		"GET",
		"/api/getTeacher/{id}",
		getTeacher,
	},
	Route{
		"FindTeacherID",
		"POST",
		"/api/findTeacherID",
		findTeacherID,
	},
	Route{
		"RegisterTeacher",
		"POST",
		"/api/registerTeacher",
		registerTeacher,
	},
	Route{
		"GetAnswerSheet",
		"GET",
		"/api/getAnswerSheet/{studentID}/{testID}",
		getAnswerSheet,
	},
	Route{
		"SubmitAnswerSheet",
		"POST",
		"/api/submitAnswerSheet/{testID}",
		submitAnswerSheet,
	},
	Route{
		"GetGrade",
		"GET",
		"/api/getGrade/{studentID}/{testID}",
		getGrade,
	},
	Route{
		"SubmitGrade",
		"POST",
		"/api/submitGrade/{testID}",
		submitGrade,
	},
	Route{
		"ListLessons",
		"GET",
		"/api/listLessons/{subject}/{grade}",
		listLessons,
	},
	Route{
		"UploadLesson",
		"POST",
		"/api/uploadLesson/{subject}/{grade}",
		uploadLesson,
	},
	Route{
		"AdminDownloadLogs",
		"GET",
		"/api/downloadLogs",
		downloadLogs,
	},
}
