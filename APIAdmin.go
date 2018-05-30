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
	"archive/zip"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"io/ioutil"
	"github.com/buger/jsonparser"
	"strings"
	"runtime"
	update "github.com/inconshreveable/go-update"
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
		w.Header().Set("WWW-Authenticate", `Basic realm="Access to the admin section"`)
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

// updateServer calls for the server to download the latest binary from the latest release of this server's GitHub
// repository.
//
// The method calls the GitHub API and downloads the latest release whose binary matches this program's name (for OS
// compatibility) and self-updates. After that, a manual restart is required for the server to apply its changes.
func updateServer(w http.ResponseWriter, r *http.Request) {

	username, password, authOK := r.BasicAuth()
	responseCode := http.StatusOK

	user, pass := GetAdminCreds()

	if !authOK || (username != user || password != pass) {
		responseCode = http.StatusUnauthorized
		w.Header().Set("WWW-Authenticate", `Basic realm="Access to the admin section"`)
		http.Error(w, "Invalid authentication scheme!", responseCode)
		return
	}

	resp, err := http.Get("https://api.github.com/repos/CNITV/VianuEdu-Server/releases/latest")
	if err != nil {
		responseCode = http.StatusInternalServerError
		w.WriteHeader(responseCode)
		fmt.Fprintf(w, "Could not get API URL for update! (%s)", err.Error())
		return
	}
	defer resp.Body.Close()

	downloadURLJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		responseCode = http.StatusInternalServerError
		w.WriteHeader(responseCode)
		fmt.Fprintf(w, "Could not extract body from HTTP request! (%s)", err.Error())
		return
	}

	updateURL := ""
	_, err = jsonparser.ArrayEach(downloadURLJSON, func(value []byte, dataType jsonparser.ValueType, offset int, err1 error) {
		asset, err2 := jsonparser.GetString(value, "browser_download_url")
		if err2 != nil {
			return
		}

		if runtime.GOOS == "windows" {
			if strings.Contains(asset, ".exe") {
				updateURL = asset
			}
		} else {
			if !strings.Contains(asset, ".exe") {
				updateURL = asset
			}
		}
	}, "assets")
	if err != nil {
		responseCode = http.StatusInternalServerError
		w.WriteHeader(responseCode)
		fmt.Fprintf(w, "Could not get asset for update from JSON array! (%s)", err.Error())
		return
	}

	updateResp, err := http.Get(updateURL)
	if err != nil {
		responseCode = http.StatusInternalServerError
		w.WriteHeader(responseCode)
		fmt.Fprintf(w, "Could not update from URL! (%s)", err.Error())
		return
	}
	defer updateResp.Body.Close()

	err = update.Apply(updateResp.Body, update.Options{})
	if err != nil {
		responseCode = http.StatusInternalServerError
		w.WriteHeader(responseCode)
		fmt.Fprintf(w, "Could not apply update! (%s)", err.Error())
		return
	}

	fmt.Fprint(w, "Update succesful! Restart the server!")
	HTTPLogger.Warn("[WARN] Server has been updated through updateServer HTTP Handler! Restart!")
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
