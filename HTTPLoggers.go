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
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// HTTPLogger is a logger linked to HTTPServer.log file. It uses the logrus framework.
var HTTPLogger *logrus.Logger

// APILogger is a logger linked to APIRequests.log file. It uses the logrus framework.
var APILogger *logrus.Logger

// StartLoggers assigns a lumberjack rotation io.WriteCloser to the loggers.
//
// It also adds a TextFormatter to the loggers in order to allow for easy readability of the log files.
func StartLoggers() {

	HTTPLogger = logrus.New()
	APILogger = logrus.New()

	HTTPLogger.Out = &lumberjack.Logger{
		Filename:   "log/HTTPServer.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}
	APILogger.Out = &lumberjack.Logger{
		Filename:   "log/APIRequests.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}

	HTTPLogger.Formatter = &logrus.TextFormatter{}
	APILogger.Formatter = &logrus.TextFormatter{}
}
