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
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.SetOutput(&lumberjack.Logger{
			Filename:   "log/APIRequests.log",
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		})

		log.WithFields(log.Fields{
			"method": r.Method,
			"requestURI": r.RequestURI,
			"name": name,
			"time": time.Since(start),
		}).Info()
	})
}
