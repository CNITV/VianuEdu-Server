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
	"fmt"
	"github.com/gorilla/mux"
)

// Gets a student from the database based on the student ID presented.
// Will return application/json content type unless text/plain is requested.
func getStudent(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)

	id := requestVars["id"]

	student := GetStudentObjectByID(id)

	if r.Header.Get("Accept") == "text/plain" {
		w.Header().Set("Content-Type", "text/plain")
	} else {
		w.Header().Set("Content-Type", "application/json")
	}

	fmt.Fprint(w, student)
}

func findStudentID(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Not implemented yet!")
}

// Gets a student from the database based on the teacher ID presented.
// Will return application/json content type unless text/plain is requested.
func getTeacher(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)

	id := requestVars["id"]

	student := GetTeacherObjectByID(id)

	if r.Header.Get("Accept") == "text/plain" {
		w.Header().Set("Content-Type", "text/plain")
	} else {
		w.Header().Set("Content-Type", "application/json")
	}

	fmt.Fprint(w, student)
}

func findTeacherID(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Not implemented yet!")
}