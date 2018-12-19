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
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

// listLessons lists all of the lessons in the grade and subject provided in the request URL by returning the IDS
// of the lessons in the database.
//
// Currently, the grades are only between 9 and 12. This is mostly due to the fact that, as the project stands, it will
// be highly unlikely that any 1-8th grade will use this educational software.
func listLessons(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)
	responseCode := http.StatusOK

	grade, err := strconv.Atoi(requestVars["grade"])
	lessonList := ListLessons(requestVars["subject"], grade)

	if err != nil || (grade < 9 || grade > 12) {
		responseCode = http.StatusBadRequest
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "Invalid grade! Must be between 9-12!")
		goto log
	}

	if lessonList == "notFound" {
		responseCode = http.StatusNotFound
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "404 lessons not found")
		goto log
	}

	fmt.Fprint(w, lessonList)

log:
	APILogger.WithFields(logrus.Fields{
		"host":         r.RemoteAddr,
		"userAgent":    r.UserAgent(),
		"responseCode": responseCode,
	}).Info("listLessons hit")
}

// getLesson simply downloads a lesson from the VianuEdu server to the user.
//
// If it is not found, the server returns a 404 Not Found status code.
func getLesson(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)

	lesson := GetLesson(requestVars["course"], requestVars["lessonID"])
	responseCode := http.StatusOK

	if lesson == "notFound" {
		responseCode = http.StatusNotFound
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "404 lesson not found!")
		return
	}

	fmt.Fprint(w, lesson)

	APILogger.WithFields(logrus.Fields{
		"host":         r.RemoteAddr,
		"userAgent":    r.UserAgent(),
		"lessonID":       requestVars["lessonID"],
		"responseCode": responseCode,
	}).Info("getLesson hit")
}

// uploadLesson uploads a lesson to the repository, provided it is given valid teacher credentials.
//
// Currently, it does not differentiate whether a teacher should be able to upload for a specific subject, depending on
// what they teach. This is in development and debate. TODO add check for subject-teacher validation
//
// Should the credentials provided be invalid, the HTTP handler responds with a Unauthorized (401) response code.
// Currently, the function only allows the upload of PNG files.
func uploadLesson(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)

	username, password, authOK := r.BasicAuth()

	grade, err := strconv.Atoi(requestVars["grade"])
	templateFile, _ := os.Open("templates/LessonTemplate.json")

	responseCode := http.StatusOK

	if err != nil || (grade < 9 || grade > 12) {
		responseCode = http.StatusBadRequest
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "Invalid grade! Must be between 9-12!")
		return
	}

	if !authOK {
		responseCode = http.StatusUnauthorized
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "Malformed authentication scheme!")
		return
	}

	teacherID := FindTeacherID(username, password)

	if teacherID == "notFound" {
		responseCode = http.StatusUnauthorized
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "Invalid username and password!")
		return
	}

	//let's go!
	if responseCode == http.StatusOK {

		//validate JSON!
		// we pretty much only care for the final error, since the rest of the stuff here is unlikely to ever fail randomly.
		templateString, _ := ioutil.ReadAll(templateFile)

		lessonTemplate := gojsonschema.NewStringLoader(string(templateString))

		body, _ := ioutil.ReadAll(r.Body)

		lessonResponse := gojsonschema.NewStringLoader(string(body))

		validation, err := gojsonschema.Validate(lessonTemplate, lessonResponse)
		if err != nil {
			APILogger.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Could not validate JSON schema and document for adding lesson!")
			responseCode := http.StatusBadRequest
			w.WriteHeader(responseCode)
			fmt.Fprint(w, "Invalid Lesson object!")
			return
		}

		if validation.Valid() {
			AddLesson(requestVars["course"], grade, string(body))

			fmt.Fprint(w, "Lesson uploaded!")
		}
	}

	APILogger.WithFields(logrus.Fields{
		"host":         r.RemoteAddr,
		"userAgent":    r.UserAgent(),
		"teacherID":    teacherID,
		"responseCode": responseCode,
	}).Info("uploadLesson hit")
}
