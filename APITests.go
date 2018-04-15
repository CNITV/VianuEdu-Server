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
	"strings"
)

func getTest(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)

	test := GetTest(requestVars["testID"])
	responseCode := http.StatusOK

	if test == "notFound" {
		responseCode = http.StatusNotFound
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "404 test not found!")
		return
	}

	fmt.Fprint(w, test)

	APILogger.WithFields(logrus.Fields{
		"host":         r.RemoteAddr,
		"userAgent":    r.UserAgent(),
		"testID":       requestVars["testID"],
		"responseCode": responseCode,
	}).Info("getTest hit")
}

func getTestQueue(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)
	responseCode := http.StatusOK

	student := GetStudentObjectByID(requestVars["studentID"])

	if student == "notFound" {
		responseCode = http.StatusNotFound
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "404 student not found!")
		return
	}

	if !strings.Contains("GeoPhiInfoMath", requestVars["subject"]) {
		responseCode = http.StatusNotFound
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "404 course not found")
		return
	}

	grade, _ := jsonparser.GetInt([]byte(student), "grade")
	gradeLetter, _ := jsonparser.GetString([]byte(student), "gradeLetter")

	tests := GetTestQueue(requestVars["subject"], grade, gradeLetter)

	if tests == "notFound" {
		responseCode := http.StatusNotFound
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "404 tests not found!")
		return
	}

	fmt.Fprint(w, tests)

	APILogger.WithFields(logrus.Fields{
		"host":         r.RemoteAddr,
		"userAgent":    r.UserAgent(),
		"studentID":    requestVars["studentID"],
		"subject":      requestVars["subject"],
		"responseCode": responseCode,
	}).Info("getTestQueue hit")
}

func getNextTestID(w http.ResponseWriter, r *http.Request) {
	testID := GetNextTestID()

	fmt.Fprint(w, testID)

	APILogger.WithFields(logrus.Fields{
		"host":         r.RemoteAddr,
		"userAgent":    r.UserAgent(),
		"testID":       testID,
		"responseCode": http.StatusOK,
	}).Info("getNextTestID hit")
}

func createTest(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)

	//first we strip out the authentication from the header
	username, password, authOK := r.BasicAuth()

	responseCode := http.StatusOK

	teacherID := FindTeacherID(username, password)

	templateFile, _ := os.Open("templates/TestTemplate.json")

	testID := ""

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

		testTemplate := gojsonschema.NewStringLoader(string(templateString))

		body, _ := ioutil.ReadAll(r.Body)

		testResponse := gojsonschema.NewStringLoader(string(body))

		validation, err := gojsonschema.Validate(testTemplate, testResponse)
		if err != nil {
			APILogger.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Could not validate JSON schema and document for adding test!")
		}

		if validation.Valid() {
			testID = GetNextTestID()

			submittedTestID, _ := jsonparser.GetString(body, "testID")

			if !(testID == submittedTestID) {
				responseCode := http.StatusBadRequest
				w.WriteHeader(responseCode)
				fmt.Fprint(w, "Invalid test ID! Test ID must be acquired from server from GetNextTestID endpoint!")
				return
			}

			AddTest(requestVars["subject"], string(body), testID)

			fmt.Fprint(w, "Test created! New test ID is "+testID)
		}
	}

	APILogger.WithFields(logrus.Fields{
		"host":         r.RemoteAddr,
		"userAgent":    r.UserAgent(),
		"teacherID":    teacherID,
		"testID":       testID,
		"responseCode": responseCode,
	}).Info("submitGrade hit")
}
