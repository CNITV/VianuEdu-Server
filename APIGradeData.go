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
	"github.com/buger/jsonparser"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"net/http"
	"os"
)

// getGrade obtains a grade from the database by querying for student ID and test ID.
//
// It will send back a Resource Not Found (404) response code if there is no grade found.
func getGrade(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)

	student := GetStudentObjectByID(requestVars["studentID"])

	user, _ := jsonparser.GetString([]byte(student), "account", "userName")

	grade := GetGrade(user, requestVars["testID"])

	responseCode := http.StatusOK

	if grade == "notFound" {
		responseCode = http.StatusNotFound
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "404 grade not found")
		return
	}

	fmt.Fprint(w, grade)

	APILogger.WithFields(logrus.Fields{
		"host":         r.RemoteAddr,
		"userAgent":    r.UserAgent(),
		"studentID":    requestVars["studentID"],
		"testID":       requestVars["testID"],
		"responseCode": responseCode,
	}).Info("getGrade hit")
}

// submitGrade adds a grade to the database, to the provided student and test ID.
//
// Only teachers can submit grades to the database! Anyone else attempting to do so will be responded with a
// Unauthorized (401) response code.
//
// Every single validation conducted within this HTTP handler function is directly equivalent in some way, shape, or
// form to the submitAnswerSheet documentation. Refer there for details.
func submitGrade(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)

	//first we strip out the authentication from the header
	username, password, authOK := r.BasicAuth()

	responseCode := http.StatusOK

	teacherID := FindTeacherID(username, password)

	templateFile, _ := os.Open("templates/GradeTemplate.json")

	var studentID = ""

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

	//let's go!
	if responseCode == http.StatusOK {

		//validate JSON!
		// we pretty much only care for the final error, since the rest of the stuff here is unlikely to ever fail randomly.
		templateString, _ := ioutil.ReadAll(templateFile)

		gradeTemplate := gojsonschema.NewStringLoader(string(templateString))

		body, _ := ioutil.ReadAll(r.Body)

		gradeResponse := gojsonschema.NewStringLoader(string(body))

		validation, err := gojsonschema.Validate(gradeTemplate, gradeResponse)
		if err != nil {
			APILogger.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Could not validate JSON schema and document for adding grade!")
		}

		if validation.Valid() {
			//don't need errors here because I've already validated the JSON and know that it will work
			username, _ := jsonparser.GetString(body, "teacher", "account", "userName")
			password, _ := jsonparser.GetString(body, "teacher", "account", "password")

			studentUser, _ := jsonparser.GetString(body, "studentAnswerSheet", "student", "account", "userName")
			studentPass, _ := jsonparser.GetString(body, "studentAnswerSheet", "student", "account", "password")

			studentID = FindStudentID(studentUser, studentPass)

			checkID := FindTeacherID(username, password)

			if checkID != teacherID {
				responseCode = http.StatusUnauthorized
				w.WriteHeader(responseCode)
				fmt.Fprint(w, "Malformed grade! (can't upload grade on someone else's behalf")
				return
			}

			testID, _ := jsonparser.GetString(body, "studentAnswerSheet", "testID")
			keyID, _ := jsonparser.GetString(body, "answerKey", "testID")

			if testID != requestVars["testID"] {
				responseCode = http.StatusBadRequest
				w.WriteHeader(responseCode)
				fmt.Fprint(w, "Cannot submit grade from another test to this one!")
				return
			}

			if testID != keyID {
				responseCode = http.StatusBadRequest
				w.WriteHeader(responseCode)
				fmt.Fprint(w, "Malformed grade! (Cannot have answer sheets from different tests!")
				return
			}

			if GetGrade(studentUser, testID) != "notFound" {
				responseCode = http.StatusAlreadyReported
				w.WriteHeader(responseCode)
				fmt.Fprint(w, "Cannot submit a grade after it has already been submitted!")
				return
			}

			AddGrade(string(body), testID)
			fmt.Fprint(w, "Grade added! You can no longer add anything to this test!")
		}
	}

	APILogger.WithFields(logrus.Fields{
		"host":         r.RemoteAddr,
		"userAgent":    r.UserAgent(),
		"teacherID":    teacherID,
		"studentID":    studentID,
		"responseCode": responseCode,
	}).Info("submitGrade hit")
}
