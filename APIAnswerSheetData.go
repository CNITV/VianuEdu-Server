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
	"regexp"
)

// getAnswerSheet gets an AnswerSheet object from the database and send it to the client.
// It will query the database for the AnswerSheet object attached to the student ID and test ID received.
//
// It will fail if the test ID is invalid, if the student ID is invalid and send back a Bad Request (400) response code.
// It will also respond with a Resource Not Found response code (404) if there is no answer sheet attached to the
// student ID and test ID combination.
func getAnswerSheet(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)

	regexTest, _ := regexp.Compile(`T-([0123456789])\w+`)

	responseCode := http.StatusOK

	if !regexTest.Match([]byte(requestVars["testID"])) {
		responseCode = http.StatusBadRequest
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "Test ID invalid!")
		return
	}

	student := GetStudentObjectByID(requestVars["studentID"])

	if student == "notFound" {
		responseCode = http.StatusBadRequest
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "Student ID invalid!")
		return
	}

	answerSheet := GetAnswerSheet(student, requestVars["testID"])

	w.Header().Set("Content-Type", "application/json")

	if answerSheet == "notFound" {
		responseCode = http.StatusNotFound
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "404 answer sheet not found")
	} else {
		fmt.Fprint(w, answerSheet)
	}

	APILogger.WithFields(logrus.Fields{
		"host":         r.RemoteAddr,
		"userAgent":    r.UserAgent(),
		"studentID":    requestVars["studentID"],
		"testID":       requestVars["testID"],
		"responseCode": responseCode,
	}).Info("getAnswerSheet hit")

}

// submitAnswerSheet adds an answer sheet to the database, on the provided test ID by the provided student ID.
//
// It will fail if the authentication scheme is invalid, and send back a Unauthorized (401) response code.
// It will fail if the student ID is not found and send back a Resource Not Found (404) response code.
// It will also fail if the submitted answer sheet has an invalid JSON schema and send back Bad Request (400)
// response code.
// Any invalid combination of student ID - test ID will be responded with a Bad Request (400) response code.
func submitAnswerSheet(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)

	//first we strip out the authentication from the header
	username, password, authOK := r.BasicAuth()

	responseCode := http.StatusOK

	studentID := FindStudentID(username, password)

	templateFile, _ := os.Open("templates/AnswerSheetTemplate.json")

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

	//let's go!
	if responseCode == http.StatusOK {

		//validate JSON!
		// we pretty much only care for the final error, since the rest of the stuff here is unlikely to ever fail randomly.
		templateString, _ := ioutil.ReadAll(templateFile)

		answerSheetTemplate := gojsonschema.NewStringLoader(string(templateString))

		body, _ := ioutil.ReadAll(r.Body)

		answerSheetResponse := gojsonschema.NewStringLoader(string(body))

		validation, err := gojsonschema.Validate(answerSheetTemplate, answerSheetResponse)
		if err != nil {
			APILogger.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Could not validate JSON schema and document for adding answer sheet!")
		}

		if validation.Valid() {
			//don't need errors here because I've already validated the JSON and know that it will work
			user, _ := jsonparser.GetString(body, "student", "account", "userName")
			pass, _ := jsonparser.GetString(body, "student", "account", "password")

			checkID := FindStudentID(user, pass)

			if checkID != studentID {
				responseCode = http.StatusUnauthorized
				w.WriteHeader(responseCode)
				fmt.Fprint(w, "Malformed answer sheet! (can't upload answer sheet on someone else's behalf")
				return
			}

			testID, _ := jsonparser.GetString(body, "testID")
			if testID != requestVars["testID"] {
				responseCode = http.StatusBadRequest
				w.WriteHeader(responseCode)
				fmt.Fprint(w, "Cannot submit answer sheet from another test to this one!")
				return
			}

			student, _, _, _ := jsonparser.Get(body, "student")

			if GetAnswerSheet(string(student), testID) != "notFound" {
				responseCode = http.StatusAlreadyReported
				w.WriteHeader(responseCode)
				fmt.Fprint(w, "Cannot submit an answer sheet after it has already been submitted!")
				return
			}

			AddAnswerSheet(string(body))
			fmt.Fprint(w, "Answer sheet added! You can no longer add anything to this test!")
		}
	}

	APILogger.WithFields(logrus.Fields{
		"host":         r.RemoteAddr,
		"userAgent":    r.UserAgent(),
		"studentID":    studentID,
		"responseCode": responseCode,
	}).Info("submitAnswerSheet hit")
}

func getAnswerSheetsForTest(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)
	responseCode := http.StatusOK

	answerSheets := GetAnswerSheetsForTest(requestVars["testID"])

	if answerSheets == "notFound" {
		responseCode = http.StatusNotFound
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "404 no answer sheets found for this test!")
		return
	}
	if answerSheets == "" {
		responseCode = http.StatusInternalServerError
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "Oops! We messed up somewhere! Sorry! Try again")
		return
	}

	fmt.Fprint(w, answerSheets)

	APILogger.WithFields(logrus.Fields{
		"host":         r.RemoteAddr,
		"userAgent":    r.UserAgent(),
		"testID":       requestVars["testID"],
		"responseCode": responseCode,
	}).Info("getAnswerSheetsForTest hit")

}
