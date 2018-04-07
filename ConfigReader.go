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
	"github.com/buger/jsonparser"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

// GetListenPort reads the configuration file HTTPServer.json for the HTTP server listening port.
//
// It extracts the config file JSON into memory and parses it, looking for the listenPort entry.
// The configuration file must follow the template provided with the source code and release distribution,
// otherwise the server exits immediately.
func GetListenPort() int64 {

	configFile, err := os.Open("config/HTTPServer.json")
	if err != nil {
		HTTPLogger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error opening HTTPServer configuration file!")
	}
	defer configFile.Close()

	mainConfig, err := ioutil.ReadAll(configFile)
	if err != nil {
		HTTPLogger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error reading HTTPServer configuration variable!")
	}
	HTTPLogger.Println("[BOOT] Reading listen port...")
	listenPort, err := jsonparser.GetInt(mainConfig, "listenPort")
	if err != nil {
		HTTPLogger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error parsing HTTPServer configuration file! (can't parse listenPort)")
	}
	return listenPort
}

// GetDBConnectionURL reads the configuration file DatabaseSettings.json for the connection URL required for the mgo
// driver to successfully connect to the MongoDB instance.
//
// It extracts the config file JSON into memory and parses it, looking for all entries, and compiles them into a valid
// URL that can be used by the driver.
// The configuration file must follow the template provided with the source code and release distribution,
// otherwise the server exits immediately.
func GetDBConnectionURL() string {
	configFile, err := os.Open("config/DatabaseSettings.json")
	if err != nil {
		HTTPLogger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error opening DatabaseSettings configuration file!")
	}
	defer configFile.Close()
	mainConfig, err := ioutil.ReadAll(configFile)
	if err != nil {
		HTTPLogger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error reading DatabaseSettings configuration variable!")
	}
	HTTPLogger.Println("[BOOT] Reading connection information...")
	serverIP, err := jsonparser.GetString(mainConfig, "serverIP")
	if err != nil {
		HTTPLogger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error parsing DatabaseSettings configuration file! (can't parse serverIP)")
	}
	serverPort, err := jsonparser.GetString(mainConfig, "serverPort")
	if err != nil {
		HTTPLogger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error parsing DatabaseSettings configuration file! (can't parse serverPort)")
	}
	userName, err := jsonparser.GetString(mainConfig, "userName")
	if err != nil {
		HTTPLogger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error parsing DatabaseSettings configuration file! (can't parse userName)")
	}
	userPass, err := jsonparser.GetString(mainConfig, "userPass")
	if err != nil {
		HTTPLogger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error parsing DatabaseSettings configuration file! (can't parse userPass)")
	}
	dbName, err := jsonparser.GetString(mainConfig, "databaseName")
	if err != nil {
		HTTPLogger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error parsing DatabaseSettings configuration file! (can't parse databaseName)")
	}

	connectionURL := "mongodb://" + userName + ":" + userPass + "@" + serverIP + ":" + serverPort + "/" + dbName

	HTTPLogger.Println("[BOOT] Done reading configuration...")
	return connectionURL
}

// getDBName reads the configuration file DatabaseSettings.json for the name of the database on the MongoDB instance
// which contains the predefined VianuEdu schema. The schema is defined in README.md
//
// This method extracts the config file JSON into memory and parses it, looking for the "databaseName" entry
// The configuration file must follow the template provided with the source code and release distribution,
// otherwise the server exits immediately.
func getDBName() string {
	configFile, err := os.Open("config/DatabaseSettings.json")
	if err != nil {
		HTTPLogger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error opening DatabaseSettings configuration file!")
	}
	defer configFile.Close()

	mainConfig, err := ioutil.ReadAll(configFile)
	if err != nil {
		HTTPLogger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error reading DatabaseSettings configuration variable!")
	}

	dbName, err := jsonparser.GetString(mainConfig, "databaseName")
	if err != nil {
		HTTPLogger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error parsing DatabaseSettings configuration file! (can't parse databaseName)")
	}

	return dbName
}

// getAdminCreds reads the configuration file HTTPServer.json for the credentials required for all administrative
// functions found in APIAdmin.go
//
// This method extracts the config file JSON into memory and parses it, looking for the "adminUser" and "adminPass"
// entry.
// The configuration file must follow the template provided with the source code and release distribution,
// otherwise the server exits immediately.
func GetAdminCreds() (string, string) {

	configFile, err := os.Open("config/HTTPServer.json")
	if err != nil {
		HTTPLogger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error opening HTTPServer configuration file!")
	}
	defer configFile.Close()

	mainConfig, err := ioutil.ReadAll(configFile)
	if err != nil {
		HTTPLogger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error reading HTTPServer configuration variable!")
	}
	HTTPLogger.Println("[BOOT] Reading admin credentials...")
	adminUser, err := jsonparser.GetString(mainConfig, "adminUser")
	if err != nil {
		HTTPLogger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error parsing HTTPServer configuration file! (can't parse adminUser)")
	}
	adminPass, err := jsonparser.GetString(mainConfig, "adminPass")
	if err != nil {
		HTTPLogger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error parsing HTTPServer configuration file! (can't parse adminPass)")
	}

	return adminUser, adminPass
}

// ShouldIGoSecure reads the configuration file HTTPServer.json and checks whether the current server instance should
// restrict all communication to TLS-only (HTTPS). The server will automatically take care of all redirects and
// everything in between.
// The configuration file must follow the template provided with the source code and release distribution,
// otherwise the server exits immediately.
func ShouldIGoSecure() (bool, string, string) {
	configFile, err := os.Open("config/HTTPServer.json")
	if err != nil {
		HTTPLogger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error opening HTTPServer configuration file!")
	}
	defer configFile.Close()

	mainConfig, err := ioutil.ReadAll(configFile)
	if err != nil {
		HTTPLogger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error reading HTTPServer configuration variable!")
	}
	HTTPLogger.Println("[BOOT] Reading enableTLS value...")
	value, err := jsonparser.GetBoolean(mainConfig, "enableTLS")
	if err != nil {
		HTTPLogger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error parsing HTTPServer configuration file! (can't parse enableTLS)")
	}

	if value {
		HTTPLogger.Println("[BOOT] TLS enabled! Reading TLS key files...")
		keyFile, err := jsonparser.GetString(mainConfig, "keyFile")
		if err != nil {
			HTTPLogger.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("Error parsing HTTPServer configuration file! (can't parse keyFile)")
		}
		certFile, err := jsonparser.GetString(mainConfig, "certFile")
		if err != nil {
			HTTPLogger.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("Error parsing HTTPServer configuration file! (can't parse certFile)")
		}

		return true, certFile, keyFile
	}
	HTTPLogger.Warn("[BOOT][WARN] TLS disabled! Recheck configuration if this is non-intentional!")
	return false, "", ""
}
