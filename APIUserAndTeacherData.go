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
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"net/http"
	"os"
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
// If body is invalid JSON, the HTTP Handler returns a Bad Request (400) response code.
func findStudentID(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Cannot parse body from request!")
	}

	responseCode := http.StatusOK

	username, err := jsonparser.GetString(body, "userName")
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Cannot parse JSON body for `userName` entry!")
		responseCode = http.StatusBadRequest
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "Invalid body!")
		return
	}
	password, err := jsonparser.GetString(body, "password")
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Cannot parse JSON body for `password` entry!")
		responseCode = http.StatusBadRequest
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "Invalid body!")
		return
	}

	studentID := FindStudentID(username, password)

	if studentID == "notFound" {
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

// Gets the teacher ID by checking the database for the user with the provided username and password
// Will return ID in text/plain form.
// If body is invalid JSON, the HTTP Handler returns a Bad Request (400) response code.
func findTeacherID(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Cannot parse body from request!")
	}

	responseCode := http.StatusOK

	username, err := jsonparser.GetString(body, "userName")
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Cannot parse JSON body for `userName` entry!")
		responseCode = http.StatusBadRequest
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "Invalid body!")
		return
	}
	password, err := jsonparser.GetString(body, "password")
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Cannot parse JSON body for `password` entry!")
		responseCode = http.StatusBadRequest
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "Invalid body!")
		return
	}

	teacherID := FindTeacherID(username, password)

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

// changeStudentPassword changes the password of an already added Student in the database.
//
// It queries for the ID that is found and changes the password with one provided in the body.
// If the student isn't found, the handler returns a 401 Unauthorized error.
func changeStudentPassword(w http.ResponseWriter, r *http.Request) {

	//first we strip out the authentication from the header
	username, password, authOK := r.BasicAuth()

	responseCode := http.StatusOK

	studentID := FindStudentID(username, password)

	//then we check to see if authOK
	if !authOK {
		responseCode = http.StatusUnauthorized
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "Invalid authentication scheme!")
	}

	//see if student exists
	if studentID == "notFound" {
		responseCode = http.StatusUnauthorized
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "Invalid username and password combination!")
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseCode = http.StatusBadRequest
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "Cannot read body! Try again!")
	}

	ChangeStudentPassword(studentID, string(body))

	fmt.Fprint(w, "Password changed!")
}

// changeTeacherPassword changes the password of an already added Teacher in the database.
//
// It queries for the ID that is found and changes the password with one provided in the body.
// If the teacher isn't found, the handler returns a 401 Unauthorized error.
func changeTeacherPassword(w http.ResponseWriter, r *http.Request) {

	//first we strip out the authentication from the header
	username, password, authOK := r.BasicAuth()

	responseCode := http.StatusOK

	teacherID := FindTeacherID(username, password)

	//then we check to see if authOK
	if !authOK {
		responseCode = http.StatusUnauthorized
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "Invalid authentication scheme!")
	}

	//see if teacher exists
	if teacherID == "notFound" {
		responseCode = http.StatusUnauthorized
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "Invalid username and password combination!")
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseCode = http.StatusBadRequest
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "Cannot read body! Try again!")
	}

	ChangeTeacherPassword(teacherID, string(body))

	fmt.Fprint(w, "Password changed!")
}

// registerStudent adds the provided Student object to the database, provided the body contains valid JSON for a Student
// object.
//
// If it isn't valid, then the HTTP handler returns a Bad Request (400) response code.
// If the student if successfully registered, then the handler returns the ID for the brand-new created student.
func registerStudent(w http.ResponseWriter, r *http.Request) {
	templateFile, err := os.Open("templates/StudentTemplate.json")
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not open StudentTemplate file!")
	}

	// we pretty much only care for the final error, since the rest of the stuff here is unlikely to ever fail randomly.
	templateString, _ := ioutil.ReadAll(templateFile)

	studentTemplate := gojsonschema.NewStringLoader(string(templateString))

	body, _ := ioutil.ReadAll(r.Body)

	studentResponse := gojsonschema.NewStringLoader(string(body))

	validation, err := gojsonschema.Validate(studentTemplate, studentResponse)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not validate JSON schema and document for registering Student")
	}

	responseCode := http.StatusOK

	if validation.Valid() {
		RegisterStudent(string(body))

		username, _ := jsonparser.GetString(body, "account", "userName")
		password, _ := jsonparser.GetString(body, "account", "password")

		id := FindStudentID(username, password)

		fmt.Fprint(w, id)
	} else {
		responseCode = http.StatusBadRequest
		fmt.Fprint(w, "Sent student JSON not valid! Reevaluate")
	}

	APILogger.WithFields(logrus.Fields{
		"host":         r.RemoteAddr,
		"userAgent":    r.UserAgent(),
		"responseCode": responseCode,
	}).Info("registerStudent hit")
}

// registerTeacher adds the provided Teacher object to the database, provided the body contains valid JSON for a Teacher
// object.
//
// If it isn't valid, then the HTTP handler returns a Bad Request (400) response code.
// If the teacher if successfully registered, then the handler returns the ID for the brand-new created teacher.
func registerTeacher(w http.ResponseWriter, r *http.Request) {
	templateFile, err := os.Open("templates/TeacherTemplate.json")
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not open TeacherTemplate file!")
	}

	// we pretty much only care for the final error, since the rest of the stuff here is unlikely to ever fail randomly.
	templateString, _ := ioutil.ReadAll(templateFile)

	teacherTemplate := gojsonschema.NewStringLoader(string(templateString))

	body, _ := ioutil.ReadAll(r.Body)

	teacherResponse := gojsonschema.NewStringLoader(string(body))

	validation, err := gojsonschema.Validate(teacherTemplate, teacherResponse)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not validate JSON schema and document for registering Teacher")
	}

	responseCode := http.StatusOK

	if validation.Valid() {
		RegisterTeacher(string(body))

		username, _ := jsonparser.GetString(body, "account", "userName")
		password, _ := jsonparser.GetString(body, "account", "password")

		id := FindTeacherID(username, password)

		fmt.Fprint(w, id)
	} else {
		responseCode = http.StatusBadRequest
		fmt.Fprint(w, "Sent teacher JSON not valid! Reevaluate")
	}

	APILogger.WithFields(logrus.Fields{
		"host":         r.RemoteAddr,
		"userAgent":    r.UserAgent(),
		"responseCode": responseCode,
	}).Info("registerTeacher hit")
}

func listClassbook(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)
	responseCode := http.StatusOK

	catalog := ListClassbook(requestVars["grade"], requestVars["gradeLetter"])

	if catalog == "notFound" {
		responseCode = http.StatusNotFound
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "404 classbook not found")
		return
	}

	fmt.Fprint(w, catalog)

	APILogger.WithFields(logrus.Fields{
		"host":      r.RemoteAddr,
		"userAgent": r.UserAgent(),
		"grade":     requestVars["grade"] + requestVars["gradeLetter"],
	}).Info("listClassbook hit")
}
