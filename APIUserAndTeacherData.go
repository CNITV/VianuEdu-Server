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
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Gets a student from the database based on the student ID presented.
// Will return application/json content type unless text/plain is requested.
func getStudent(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)

	id := requestVars["id"]

	student := GetStudentObjectByID(id)

	responseCode := http.StatusOK

	if student == "null\n" {
		w.WriteHeader(http.StatusNotFound)
		responseCode = http.StatusNotFound
	}

	APILogger.WithFields(logrus.Fields{
		"host":         r.RemoteAddr,
		"userAgent":    r.UserAgent(),
		"id":           id,
		"responseCode": responseCode,
	}).Info("getStudent hit")

	if r.Header.Get("Accept") == "text/plain" {
		w.Header().Set("Content-Type", "text/plain")
	} else {
		w.Header().Set("Content-Type", "application/json")
	}

	if responseCode == http.StatusNotFound {
		fmt.Fprint(w, "404 student not found")
	} else {
		fmt.Fprint(w, student)
	}

}

// Gets the student ID by checking the database for the user with the provided username and password
// Will return ID in text/plain form.
func findStudentID(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	studentID := FindStudentID(username, password)

	responseCode := http.StatusOK

	if studentID == "" {
		w.WriteHeader(http.StatusNotFound)
		responseCode = http.StatusNotFound
	}

	APILogger.WithFields(logrus.Fields{
		"host":         r.RemoteAddr,
		"userAgent":    r.UserAgent(),
		"teacherID":    studentID,
		"responseCode": responseCode,
	}).Info("findStudentID hit")

	w.Header().Set("Content-Type", "text/plain")

	if responseCode == http.StatusNotFound {
		fmt.Fprint(w, "404 student not found")
	} else {
		fmt.Fprint(w, studentID)
	}
}

// Gets a teacher from the database based on the teacher ID presented.
// Will return application/json content type unless text/plain is requested.
func getTeacher(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)

	id := requestVars["id"]

	teacher := GetTeacherObjectByID(id)

	responseCode := http.StatusOK

	if teacher == "null\n" {
		w.WriteHeader(http.StatusNotFound)
		responseCode = http.StatusNotFound
	}

	APILogger.WithFields(logrus.Fields{
		"host":         r.RemoteAddr,
		"userAgent":    r.UserAgent(),
		"id":           id,
		"responseCode": responseCode,
	}).Info("getTeacher hit")

	if r.Header.Get("Accept") == "text/plain" {
		w.Header().Set("Content-Type", "text/plain")
	} else {
		w.Header().Set("Content-Type", "application/json")
	}

	if responseCode == http.StatusNotFound {
		fmt.Fprint(w, "404 teacher not found")
	} else {
		fmt.Fprint(w, teacher)
	}
}

func findTeacherID(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")

	teacherID := FindTeacherID(username, password)

	responseCode := http.StatusOK

	if teacherID == "" {
		w.WriteHeader(http.StatusNotFound)
		responseCode = http.StatusNotFound
	}

	APILogger.WithFields(logrus.Fields{
		"host":         r.RemoteAddr,
		"userAgent":    r.UserAgent(),
		"teacherID":    teacherID,
		"responseCode": responseCode,
	}).Info("findTeacherID hit")

	w.Header().Set("Content-Type", "text/plain")

	if responseCode == http.StatusNotFound {
		fmt.Fprint(w, "404 teacher not found")
	} else {
		fmt.Fprint(w, teacherID)
	}
}
