/*
 *      This file is part of VianuEdu.
 *
 *      VianuEdu is free software: you can redistribute it and/or modify
 *      it under the terms of the GNU General Public License as published by
 *      the Free Software Foundation, either version 3 of the License, or
 *      (at your option) any later version.
 *
 *      VianuEdu is distributed in the hope that it will be useful,
 *      but WITHOUT ANY WARRANTY; without even the implied warranty of
 *      MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *      GNU General Public License for more details.
 *
 *      You should have received a copy of the GNU General Public License
 *      along with VianuEdu.  If not, see <http://www.gnu.org/licenses/>.
 *
 *      Developed by Matei Gardus <matei@gardus.eu>
 */

package main

import (
	"net/http"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"fmt"
	"os"
)

func dispense404(w http.ResponseWriter, r *http.Request) {
	templateFile, err := os.Open("errors/404.html")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Error opening 404 template file!")
	}
	defer templateFile.Close()

	HTMLOutput, err := ioutil.ReadAll(templateFile)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Error reading 404 template variable!")
	}
	fmt.Fprint(w, HTMLOutput)
}
