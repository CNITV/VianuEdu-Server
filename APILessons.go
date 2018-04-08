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
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// listLessons lists all of the lessons in the grade and subject provided in the request URL.
//
// Currently, the grades are only between 9 and 12. This is mostly due to the fact that, as the project stands, it will
// be highly unlikely that any 1-8th grade will use this educational software.
//
// This function is most primarily constructed in order to allow for easy listing for all the available lessons that one
// student may be interested in downloading to the client. It is given in such a way that it is very easy to split the
// strings formed by the filenames into something easily readable.
func listLessons(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)
	responseCode := http.StatusOK

	grade, err := strconv.Atoi(requestVars["grade"])
	dirName := "lessons/" + requestVars["subject"] + "/" + strconv.Itoa(grade) + "/"
	lessonsDir, _ := ioutil.ReadDir(dirName)
	fileList := ""

	if err != nil || (grade < 9 || grade > 12) {
		responseCode = http.StatusBadRequest
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "Invalid grade! Must be between 9-12!")
		goto log
	}

	if len(lessonsDir) == 0 {
		responseCode = http.StatusNotFound
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "404 lessons not found")
		goto log
	}

	for _, file := range lessonsDir {
		name := strings.Split(file.Name(), ".")[0]
		fileList = fileList + name + "\n"
	}

	fmt.Fprint(w, fileList)

log:
	APILogger.WithFields(logrus.Fields{
		"host":         r.RemoteAddr,
		"userAgent":    r.UserAgent(),
		"responseCode": responseCode,
	}).Info("listLessons hit")
}

// uploadLesson uploads a lesson to the repository, provided it is given valid teacher credentials.
//
// Currently, it does not differentiate whether a teacher should be able to upload for a specific subject, depending on
// what they teach. This is in development and debate. TODO add check for subject-teacher validation
//
// Should the credentials provided be invalid, the HTTP handler responds with a Unauthorized (401) response code.
// Currently, the function only allows the upload of PDF files. This will soon change, however.
// TODO change from PDF to HTML
func uploadLesson(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)

	username, password, authOK := r.BasicAuth()

	grade, err := strconv.Atoi(requestVars["grade"])
	fileName := "lessons/" + requestVars["subject"] + "/" + strconv.Itoa(grade) + "/"

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

	if r.Header.Get("Content-Type") != "application/pdf" {
		responseCode = http.StatusBadRequest
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "Invalid file! Upload a PDF file!")
		return
	}

	fileName = fileName + r.Header.Get("filename")

	out, err := os.Create(fileName)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Cannot open file for writing!")
	}

	defer out.Close()

	_, err = io.Copy(out, r.Body)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Cannot copy body into file!")
	}

	fmt.Fprint(w, "File uploaded successfully!")

	APILogger.WithFields(logrus.Fields{
		"host":         r.RemoteAddr,
		"userAgent":    r.UserAgent(),
		"responseCode": responseCode,
	}).Info("listLessons hit")
}
