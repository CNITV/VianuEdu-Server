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

// ConnectToDatabase links the session variable above to the connection URL read from the configuration files. This
// method obviously requires an Internet connection (but, come on, this is a server). It also reads the database name.
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

// GetStudentObjectByID searches the database for a student by ID and extracts one JSON document which has that ID.
//
// The original method (the mgo call) returns a JSON array. As such, it is imperative that this method removes the
// square brackets created by the string representation of the []bson.M variable.
//
// If no student is found by that ID, the method returns "notFound".
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

	if result == "null\n" {
		return "notFound"
	}

	result = strings.Trim(result, "[")
	result = result[:len(result)-2]

	return result
}

// GetTeacherObjectByID searches the database for a teacher by ID and extracts one JSON document which has that ID.
//
// The original method (the mgo call) returns a JSON array. As such, it is imperative that this method removes the
// square brackets created by the string representation of the []bson.M variable.
//
// If no teacher is found by that ID, the method returns "notFound".
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

	if result == "null\n" {
		return "notFound"
	}

	result = strings.Trim(result, "[")
	result = result[:len(result)-2]

	return result
}

// FindTeacherID searches the database for any teacher with the username and password provided and returns their ID.
//
// The session initially finds the JSON document for the teacher, removes the square brackets encountered in by the
// string representation of a []bson.M variable, and then parses the JSON document obtained for the ID and returns it.
//
// If no teacher is found by that username and password, then the method returns "notFound".
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

	if jsonString == "null\n" {
		return "notFound"
	}

	jsonString = strings.Trim(jsonString, "[")
	jsonString = jsonString[:len(jsonString)-2]

	data := []byte(jsonString)

	result, err := jsonparser.GetString(data, "_id", "$oid")

	return result
}

// FindStudentID searches the database for any student with the username and password provided and returns their ID.
//
// The session initially finds the JSON document for the student, removes the square brackets encountered in by the
// string representation of a []bson.M variable, and then parses the JSON document obtained for the ID and returns it.
//
// If no student is found by that username and password, then the method returns "notFound".
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

	if jsonString == "null\n" {
		return "notFound"
	}

	jsonString = strings.Trim(jsonString, "[")
	jsonString = jsonString[:len(jsonString)-2]

	data := []byte(jsonString)

	result, err := jsonparser.GetString(data, "_id", "$oid")

	return result
}

// RegisterStudent merely adds a JSON Student document on the database in the right collection.
//
// This function validates nothing from the document, so any method that might call this one must be certain the
// inserted document is valid JSON for a Student object.
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

// RegisterStudent merely adds a JSON Teacher document on the database in the right collection.
//
// This function validates nothing from the document, so any method that might call this one must be certain the
// inserted document is valid JSON for a Teacher object.
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

// GetAnswerSheet searches the database for a JSON Answer Sheet associated with a specific student on a specific test ID
// and returns it.
//
// The session initially finds the JSON document with the aforementioned conditions, removes the square brackets
// inherent with the string representation of a []bson.M variable and returns it.
//
// If no such answer sheet is found, the method returns "notFound".
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

	if result == "null\n" {
		return "notFound"
	}

	result = strings.Trim(result, "[")
	result = result[:len(result)-2]

	return result
}

// AddAnswerSheet adds an Answer Sheet JSON document to the database in the right collection.
//
// This function validates nothing from the document, so any method that might call this one must be certain the
// inserted document is valid JSON for an AnswerSheet object.
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

// GetGrade searches the database for a JSON Grade associated with a specific student on a specific test ID
// and returns it.
//
// The session initially finds the JSON document with the aforementioned conditions, removes the square brackets
// inherent with the string representation of a []bson.M variable and returns it.
//
// If no such grade is found, the method returns "notFound".
func GetGrade(studentUser string, testID string) string {

	testType := GetTestType(testID)

	gradesCollection := session.DB(dbName).C(testType + "Edu.Grades")

	var gradeQuery []bson.M

	err := gradesCollection.Find(bson.M{"studentAnswerSheet.testID": testID, "studentAnswerSheet.student.account.userName": studentUser}).All(&gradeQuery)

	grade, err := bson.MarshalJSON(gradeQuery)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not marshal getGrade request in JSON!")
	}

	result := string(grade)

	if result == "null\n" {
		return "notFound"
	}

	result = strings.Trim(result, "[")
	result = result[:len(result)-2]

	return result

}

// AddGrade adds a Grade JSON document to the database in the right collection.
//
// The method checks which subject the grade is for by calling GetTestType(testID) and adds the grade. If it is
// successful, the method deletes the answer sheet JSON document from which the grade was constructed from.
//
// This function validates nothing from the document, so any method that might call this one must be certain the
// inserted document is valid JSON for an Grade object.
func AddGrade(grade string, testID string) {
	testType := GetTestType(testID)

	gradesCollection := session.DB(dbName).C(testType + "Edu.Grades")

	var document map[string]interface{}

	err := json.Unmarshal([]byte(grade), &document)

	err = gradesCollection.Insert(document)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not add grade!")
	} else {
		submittedAnswersCollection := session.DB(dbName).C("Students.SubmittedAnswers")

		username, _ := jsonparser.GetString([]byte(grade), "studentAnswerSheet", "student", "account", "userName")
		password, _ := jsonparser.GetString([]byte(grade), "studentAnswerSheet", "student", "account", "password")

		err := submittedAnswersCollection.Remove(bson.M{"testID": testID, "student.account.userName": username, "student.account.password": password})
		if err != nil {
			APILogger.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Cannot remove answer sheet from database!")
		}
	}
}

// GetTestType checks the "VianuEdu.TestList" collection for the course the test ID provided is for.
func GetTestType(testID string) string {

	var testQuery []bson.M

	testList := session.DB(dbName).C("VianuEdu.TestList")

	testList.FindId(testID).All(&testQuery)

	query, _ := bson.MarshalJSON(testQuery)

	query = query[1:]
	query = query[:len(query)-2]

	result, _ := jsonparser.GetString(query, "course")

	return result
}
