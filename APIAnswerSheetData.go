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
	"regexp"
)

func getAnswerSheet(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)

	regexTest, _ := regexp.Compile("T-([0123456789])\\w+")

	responseCode := http.StatusOK

	if regexTest.Match([]byte(requestVars["testID"])) == false {
		responseCode = http.StatusBadRequest
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "Test ID invalid!")
		return
	}

	student := GetStudentObjectByID(requestVars["studentID"])

	if student == "null\n" {
		responseCode = http.StatusBadRequest
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "Student ID invalid!")
		return
	}

	answerSheet := GetAnswerSheet(student, requestVars["testID"])

	w.Header().Set("Content-Type", "application/json")

	if answerSheet == "null\n" {
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
		templateString, err := ioutil.ReadAll(templateFile)

		answerSheetTemplate := gojsonschema.NewStringLoader(string(templateString))

		body, err := ioutil.ReadAll(r.Body)

		answerSheetResponse := gojsonschema.NewStringLoader(string(body))

		validation, err := gojsonschema.Validate(answerSheetTemplate, answerSheetResponse)
		if err != nil {
			APILogger.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Could not validate JSON schema and document for adding answer sheet!")
		}

		if validation.Valid() {
			//don't need errors here because I've already validated the JSON and know that it will work
			username, _ := jsonparser.GetString(body, "account", "userName")
			password, _ := jsonparser.GetString(body, "account", "password")

			checkID := FindStudentID(username, password)

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
