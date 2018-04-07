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
	"archive/zip"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

// downloadLogs will download all of the logs currently present in the "log" folder.
// It will create a ZIP archive out of the 3 current logs and send them to the client, provided the credentials
// match with the ones saved inside of the HTTPServer.json configuration file.
func downloadLogs(w http.ResponseWriter, r *http.Request) {

	username, password, authOK := r.BasicAuth()
	responseCode := http.StatusOK

	user, pass := GetAdminCreds()

	if !authOK || (username != user || password != pass) {
		responseCode = http.StatusUnauthorized
		w.WriteHeader(responseCode)
		w.Header().Set("WWW-Authenticate", `Basic realm="127.0.0.1"`)
		fmt.Fprint(w, "Invalid authentication scheme!")
		return
	}

	currentTime := time.Now()
	year, month, day := currentTime.Date()

	files := []string{"log/APIRequests.log", "log/HTTPRequests.log", "log/HTTPServer.log"}

	err := ZipFiles("log/VianuEdu_Server_Logs-"+strconv.Itoa(year)+"-"+month.String()+"-"+strconv.Itoa(day)+".zip", files)

	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Cannot zip up files!")
		responseCode = http.StatusInternalServerError
		w.WriteHeader(responseCode)
		fmt.Fprint(w, "Cannot zip up log files!")
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	http.ServeFile(w, r, "log/VianuEdu_Server_Logs-"+strconv.Itoa(year)+"-"+month.String()+"-"+strconv.Itoa(day)+".zip")
	os.Remove("log/VianuEdu_Server_Logs-" + strconv.Itoa(year) + "-" + month.String() + "-" + strconv.Itoa(day) + ".zip")
}

// ZipFiles creates a ZIP archive by receiving the filepath to each of the respective files.
// The first parameter determines the filepath of the ZIP archive, while the second parameter determines the files to be inserted into the archive.
func ZipFiles(filename string, files []string) error {

	newfile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newfile.Close()

	zipWriter := zip.NewWriter(newfile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {

		zipfile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer zipfile.Close()

		// Get the file information
		info, err := zipfile.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Change to deflate to gain better compression
		// see http://golang.org/pkg/archive/zip/#pkg-constants
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, zipfile)
		if err != nil {
			return err
		}
	}
	return nil
}
