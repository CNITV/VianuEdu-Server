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
	"encoding/json"
	"github.com/buger/jsonparser"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/sirupsen/logrus"
	"strings"
)

var session *mgo.Session
var dbName string

func ConnectToDatabase() {
	HTTPLogger.Println("[BOOT] Connecting to database...")
	sessionTemp, err := mgo.Dial(GetDBConnectionURL())
	if err != nil {
		HTTPLogger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error connecting to database!")
	}
	session = sessionTemp

	dbName = getDBName()

}

func GetStudentObjectByID(id string) string {
	var queryMap []bson.M

	studentAccountsCollection := session.DB(dbName).C("Students.Accounts")

	err := studentAccountsCollection.FindId(bson.ObjectIdHex(id)).All(&queryMap)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not complete query for getStudent! ID might not exist!")
	}

	jsonString, err := bson.MarshalJSON(queryMap)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not marshal getStudent request in JSON!")
	}

	result := string(jsonString)

	result = strings.Trim(result, "[")
	result = result[:len(result)-2]

	return result
}

func GetTeacherObjectByID(id string) string {
	var queryMap []bson.M

	teachersAccountsCollection := session.DB(dbName).C("Teachers.Accounts")

	err := teachersAccountsCollection.FindId(bson.ObjectIdHex(id)).All(&queryMap)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not complete query for getTeacher! ID might not exist!")
	}

	jsonString, err := bson.MarshalJSON(queryMap)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not marshal getTeacher request in JSON!")
	}

	result := string(jsonString)

	result = strings.Trim(result, "[")
	result = result[:len(result)-2]

	return result
}

func FindTeacherID(user string, password string) string {
	var queryMap []bson.M

	teachersAccountsCollection := session.DB(dbName).C("Teachers.Accounts")

	err := teachersAccountsCollection.Find(bson.M{"account.userName": user, "account.password": password}).All(&queryMap)

	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not find teacher ID in database with username" + user + "and password " + password)

	}

	teacher, err := bson.MarshalJSON(queryMap)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not marshal findTeacherID request in JSON!")
	}

	jsonString := string(teacher)

	jsonString = strings.Trim(jsonString, "[")
	jsonString = jsonString[:len(jsonString)-2]

	data := []byte(jsonString)

	result, err := jsonparser.GetString(data, "_id", "$oid")

	return result
}

func FindStudentID(user string, password string) string {
	var queryMap []bson.M

	studentsAccountsCollection := session.DB(dbName).C("Students.Accounts")

	err := studentsAccountsCollection.Find(bson.M{"account.userName": user, "account.password": password}).All(&queryMap)

	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not find student ID in database with username" + user + "and password " + password)
		return "notFound"
	}

	student, err := bson.MarshalJSON(queryMap)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not marshal findStudentID request in JSON!")
	}

	jsonString := string(student)

	jsonString = strings.Trim(jsonString, "[")
	jsonString = jsonString[:len(jsonString)-2]

	data := []byte(jsonString)

	result, err := jsonparser.GetString(data, "_id", "$oid")

	return result
}

func RegisterStudent(body string) {
	studentsAccountsCollection := session.DB(dbName).C("Students.Accounts")

	var document map[string]interface{}

	err := json.Unmarshal([]byte(body), &document)

	err = studentsAccountsCollection.Insert(document)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not register student!")
	}
}

func RegisterTeacher(body string) {
	teachersAccountsCollection := session.DB(dbName).C("Teachers.Accounts")

	var document map[string]interface{}

	err := json.Unmarshal([]byte(body), &document)

	err = teachersAccountsCollection.Insert(document)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not register teacher!")
	}
}

func GetAnswerSheet(student string, testID string) string {
	var queryMap []bson.M

	submittedAnswersCollection := session.DB(dbName).C("Students.SubmittedAnswers")

	studentJSON := []byte(student)

	username, _ := jsonparser.GetString(studentJSON, "account", "userName")
	password, _ := jsonparser.GetString(studentJSON, "account", "password")

	err := submittedAnswersCollection.Find(bson.M{"testID": testID, "student.account.userName": username, "student.account.password": password}).All(&queryMap)

	studentID, _ := jsonparser.GetString([]byte(student), "_id", "$oid")

	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error":     err,
			"studentID": studentID,
			"testID":    testID,
		}).Warn("Could not find answer sheet in database!")

	}

	answerSheet, err := bson.MarshalJSON(queryMap)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not marshal getAnswerSheet request in JSON!")
	}

	result := string(answerSheet)

	result = strings.Trim(result, "[")
	result = result[:len(result)-2]

	return result
}

func AddAnswerSheet(answerSheet string) {
	submittedAnswersCollection := session.DB(dbName).C("Students.SubmittedAnswers")

	var document map[string]interface{}

	err := json.Unmarshal([]byte(answerSheet), &document)

	err = submittedAnswersCollection.Insert(document)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not add answer sheet!")
	}
}
