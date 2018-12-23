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
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
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

	// this is guaranteed to work, no need for error-checking
	result, _ := jsonparser.GetString(data, "_id", "$oid")

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

	// this is guaranteed to work, no need for error-checking
	result, _ := jsonparser.GetString(data, "_id", "$oid")

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
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not unmarshal byte-slice into document!")
	}

	err = studentsAccountsCollection.Insert(document)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not register student!")
	}
}

// RegisterTeacher merely adds a JSON Teacher document on the database in the right collection.
//
// This function validates nothing from the document, so any method that might call this one must be certain the
// inserted document is valid JSON for a Teacher object.
func RegisterTeacher(body string) {
	teachersAccountsCollection := session.DB(dbName).C("Teachers.Accounts")

	var document map[string]interface{}

	err := json.Unmarshal([]byte(body), &document)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not unmarshal byte-slice into document!")
	}

	err = teachersAccountsCollection.Insert(document)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not register teacher!")
	}
}

// ChangeStudentPassword changes the document associated with studentID so that the entry "account.password" contains a
// new string, newPassword.
//
// This only changes documents in the Students.Accounts collection.
func ChangeStudentPassword(studentID, newPassword string) {
	studentsAccountsCollection := session.DB(dbName).C("Students.Accounts")

	err := studentsAccountsCollection.UpdateId(bson.ObjectIdHex(studentID), bson.M{"$set": bson.M{"account.password": newPassword}})
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Cannot change password in database for student!")
	}
}

// ChangeTeacherPassword changes the document associated with teacherID so that the entry "account.password" contains a
// new string, newPassword.
//
// This only changes documents in the Teachers.Accounts collection.
func ChangeTeacherPassword(teacherID, newPassword string) {
	teachersAccountsCollection := session.DB(dbName).C("Teachers.Accounts")

	err := teachersAccountsCollection.UpdateId(bson.ObjectIdHex(teacherID), bson.M{"$set": bson.M{"account.password": newPassword}})
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Cannot change password in database for teacher!")
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
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not unmarshal byte-slice into document!")
	}

	err = submittedAnswersCollection.Insert(document)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not add answer sheet!")
	}
}

// GetAnswerSheetsForTest queries the database for all the submitted answers attached to a test and returns a string
// containing, on each line, the student ID of each student who has submitted an answer for this test.
//
// Will return "notFound" if no answer sheets are found.
func GetAnswerSheetsForTest(testID string) string {
	submittedAnswersCollection := session.DB(dbName).C("Students.SubmittedAnswers")

	var testQuery []bson.M

	err := submittedAnswersCollection.Find(bson.M{"testID": testID}).All(&testQuery)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not find any answer sheet in database!")
	}

	query, err := bson.MarshalJSON(testQuery)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not marshal query into variable!")
	}

	result := string(query)

	if result == "null\n" {
		return "notFound"
	}

	result = ""
	_, err = jsonparser.ArrayEach(query, func(value []byte, dataType jsonparser.ValueType, offset int, err1 error) {
		studentUser, _ := jsonparser.GetString(value, "student", "account", "userName")
		studentPass, _ := jsonparser.GetString(value, "student", "account", "password")

		result = result + FindStudentID(studentUser, studentPass) + "\n"
	})
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Can't iterate JSON array for answer sheets to evaluate!")
	}
	return result
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
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error":  err,
			"testID": testID,
		}).Warn("Could not find graade in database!")
	}

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
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not unmarshal byte-slice into document!")
	}

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

	err := testList.FindId(testID).All(&testQuery)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Cannot find test type in database!")
	}

	query, _ := bson.MarshalJSON(testQuery)

	query = query[1:]
	query = query[:len(query)-2]

	result, _ := jsonparser.GetString(query, "course")

	return result
}

// GetTest searches the database for a JSON Test associated with a specific test ID and returns it.
//
// The session initially finds the JSON document with the aforementioned conditions, removes the square brackets
// inherent with the string representation of a []bson.M variable and returns it.
//
// If no such test is found, the method returns "notFound".
func GetTest(testID string) string {

	var testQuery []bson.M

	testType := GetTestType(testID)

	testCollection := session.DB(dbName).C(testType + "Edu.Tests")

	err := testCollection.Find(bson.M{"testID": testID}).All(&testQuery)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not query test in database by testID!")
	}

	test, err := bson.MarshalJSON(testQuery)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not marshal getGrade request in JSON!")
	}

	result := string(test)

	if result == "null\n" {
		return "notFound"
	}

	result = strings.Trim(result, "[")
	result = result[:len(result)-2]

	return result
}

// GetTestQueue searches the database for all tests a specific class might need to take and filters them by specific
// conditions.
//
// Essentially, this takes a JSON array for every single test a class might be able to take and filters them by seeing
// if the time has expired for the test.
//
// If there is no test to be taken, the method returns an empty string.
func GetTestQueue(subject string, grade int64, gradeLetter string) string {

	var testQuery []bson.M

	testCollection := session.DB(dbName).C(subject + "Edu.Tests")

	err := testCollection.Find(bson.M{"grade": grade, "gradeLetter": gradeLetter}).All(&testQuery)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"err": err,
		}).Warn("Cannot find tests in database for this grade!")
	}

	testArray, err := bson.MarshalJSON(testQuery)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"err": err,
		}).Warn("Cannot marshal testArray variable!")
	}

	result := string(testArray)

	if result == "null\n" {
		return "notFound"
	}

	result = ""
	_, err = jsonparser.ArrayEach(testArray, func(value []byte, dataType jsonparser.ValueType, offset int, err1 error) {
		startTime, err2 := jsonparser.GetString(value, "startTime")
		if err2 != nil {
			return
		}
		endTime, err2 := jsonparser.GetString(value, "endTime")
		if err2 != nil {
			return
		}
		testID, err2 := jsonparser.GetString(value, "testID")
		if err2 != nil {
			return
		}

		zone, _ := time.LoadLocation("Europe/Bucharest")

		const layout = "Jan 2, 2006 3:04:05 PM"

		start, _ := time.ParseInLocation(layout, startTime, zone)
		end, _ := time.ParseInLocation(layout, endTime, zone)

		now := time.Now().In(zone)

		if start.Before(now) && end.After(now) {
			result = result + testID + "\n"
		}
	})
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Unable to iterate JSON array!")
	}
	return result
}

// GetNextTestID queries the database for the last test added to it, and returns the next test ID to be used.
//
// i.e If the last test ID taken is T-000001, then the next test ID is T-000002, so it returns the next one.
func GetNextTestID() string {

	var testQuery []bson.M

	testList := session.DB(dbName).C("VianuEdu.TestList")

	err := testList.Find(bson.M{}).Sort("-_id").Limit(1).All(&testQuery)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Cannot find stuff in VianuEdu.TestList!")

	}
	query, err := bson.MarshalJSON(testQuery)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Cannot marshal variable for GetNextTestID!")
	}

	formatter := string(query)
	formatter = strings.Trim(formatter, "[")
	formatter = formatter[:len(formatter)-2]

	query = []byte(formatter)
	lastTestID, _ := jsonparser.GetString(query, "_id")
	testNumber, _ := strconv.Atoi(lastTestID[2:])
	testNumber++
	newTestNumber := fmt.Sprintf("%06d", testNumber)
	newTestID := "T-" + newTestNumber
	return newTestID
}

// AddTest adds a Test JSON document to the database in the right collection.
//
// This function validates nothing from the document, so any method that might call this one must be certain the
// inserted document is valid JSON for an Test object.
func AddTest(subject string, test string, testID string) {
	testList := session.DB(dbName).C("VianuEdu.TestList")

	var testProps = []byte("{\n" +
		"    \"_id\": \"" + testID + "\", \n" +
		"    \"course\": \"" + subject + "\"\n" +
		"}")

	var document map[string]interface{}
	err := json.Unmarshal(testProps, &document)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Cannot unmarshal into document! (How, though? This part is hardcoded)")
	}

	err = testList.Insert(document)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Cannot insert test properties in database!")
	}

	var document2 map[string]interface{}
	err = json.Unmarshal([]byte(test), &document2)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Cannot unmarshal test into document!")
	}

	testCollection := session.DB(dbName).C(subject + "Edu.Tests")
	err = testCollection.Insert(document2)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Cannot insert test in database!")
	}
}

// EditTest updates a test in the database which has a specific test ID.
//
// This function validates nothing from the document, so any method that might call this one must be certain the
// inserted document is valid JSON for an Test object.
func EditTest(subject string, test string, testID string) {
	var document2 map[string]interface{}
	err := json.Unmarshal([]byte(test), &document2)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Cannot unmarshal test into document!")
	}

	testCollection := session.DB(dbName).C(subject + "Edu.Tests")
	err = testCollection.Update(bson.M{"testID": testID}, document2)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Cannot insert test in database!")
	}
}

// GetPlannedTests searches the database for all tests a specific course has and filters them by specific
// conditions.
//
// Essentially, this takes a JSON array for every single test a course has and filters them by seeing
// if the test hasn't started yet.
//
// If there is no test to be taken, the method returns an empty string.
func GetPlannedTests(subject string) string {
	var testQuery []bson.M

	testCollection := session.DB(dbName).C(subject + "Edu.Tests")

	err := testCollection.Find(bson.M{}).All(&testQuery)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"err": err,
		}).Warn("Cannot find tests in database for this subjects!")
	}

	testArray, err := bson.MarshalJSON(testQuery)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"err": err,
		}).Warn("Cannot marshal testArray variable!")
	}

	result := string(testArray)

	if result == "null\n" {
		return "notFound"
	}

	result = ""
	_, err = jsonparser.ArrayEach(testArray, func(value []byte, dataType jsonparser.ValueType, offset int, err1 error) {
		startTime, err2 := jsonparser.GetString(value, "startTime")
		if err2 != nil {
			return
		}
		testID, err2 := jsonparser.GetString(value, "testID")
		if err2 != nil {
			return
		}

		class, err2 := jsonparser.GetInt(value, "grade")
		if err2 != nil {
			return
		}
		letter, err2 := jsonparser.GetString(value, "gradeLetter")
		if err2 != nil {
			return
		}
		grade := strconv.Itoa(int(class))
		grade = grade + letter

		zone, _ := time.LoadLocation("Europe/Bucharest")

		const layout = "Jan 2, 2006 3:04:05 PM"

		start, _ := time.ParseInLocation(layout, startTime, zone)

		now := time.Now().In(zone)

		if start.After(now) {
			result = result + testID + " // " + grade + "\n"
		}
	})
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Unable to iterate JSON array!")
	}
	return result
}

// GetGradesForTest queries the database for all the grades attached to the provided username and password for the past
// 150 days.
//
// This searches the database for all the grades added to a specific student and checks their IDs for the timestamp in
// order to run the check for time elapsed on grade submission. Then it extracts the test IDs from each grade and returns
// them in a string.
func GetGradesForTest(studentUser, studentPass, subject string) string {

	var gradeQuery []bson.M

	gradeCollection := session.DB(dbName).C(subject + "Edu.Grades")

	zone, _ := time.LoadLocation("Europe/Bucharest")

	checkTime := time.Now().Add(-150 * 24 * time.Hour).In(zone)

	checkerID := bson.NewObjectIdWithTime(checkTime)

	err := gradeCollection.Find(bson.M{"_id": bson.M{"$gte": checkerID}, "studentAnswerSheet.student.account.userName": studentUser, "studentAnswerSheet.student.account.password": studentPass}).All(&gradeQuery)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"err": err,
		}).Warn("Cannot find grades in database for this subject!")
	}

	gradeArray, err := bson.MarshalJSON(gradeQuery)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"err": err,
		}).Warn("Cannot marshal gradeArray variable!")
	}

	result := string(gradeArray)

	if result == "null\n" {
		return "notFound"
	}

	result = ""
	_, err = jsonparser.ArrayEach(gradeArray, func(value []byte, dataType jsonparser.ValueType, offset int, err1 error) {
		testID, err2 := jsonparser.GetString(value, "answerKey", "testID")
		if err2 != nil {
			return
		}
		result = result + testID + "\n"
	})
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Unable to iterate JSON array!")
	}
	return result
}

// GetUncorrectedTests queries the database for all the tests that currently have an AnswerSheet attached to them in the
// Students.SubmittedAnswers collection.
//
// This functions reads all of the distinct values of testID in that collection, sees which one are for which course and
// returns them, should they match with the provided subject parameter.
func GetUncorrectedTests(subject string) string {
	var testQuery []string

	submittedAnswersCollection := session.DB(dbName).C("Students.SubmittedAnswers")

	err := submittedAnswersCollection.Find(bson.M{}).Distinct("testID", &testQuery)

	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Unable to query submitted answers for testID!")
	}

	result := ""
	for _, testID := range testQuery {
		testSubject := GetTestType(testID)

		if testSubject == subject {
			result = result + testID + "\n"
		}
	}

	if result == "" {
		return "notFound"
	}

	return result
}

// ListClassbook lists all the studentID's matched to a specific grade.
//
// It queries the database for all students from the provided grade and extracts their ID's, returning them.
func ListClassbook(grade, gradeLetter string) string {
	var queryMap []bson.M

	teachersAccountsCollection := session.DB(dbName).C("Students.Accounts")

	gradeInt, _ := strconv.Atoi(grade)

	err := teachersAccountsCollection.Find(bson.M{"grade": gradeInt, "gradeLetter": gradeLetter}).All(&queryMap)

	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not find classbook for grade" + grade + gradeLetter + "!")

	}

	catalog, err := bson.MarshalJSON(queryMap)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not marshal listCatalog request in JSON!")
	}

	jsonString := string(catalog)

	if jsonString == "null\n" {
		return "notFound"
	}

	result := ""
	_, err = jsonparser.ArrayEach(catalog, func(value []byte, dataType jsonparser.ValueType, offset int, err1 error) {
		id, _ := jsonparser.GetString(value, "_id", "$oid")

		result = result + id + "\n"
	})
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Unable to iterate JSON array!")
	}

	return result
}

func ListLessons(course string, grade int) string {
	var queryMap []bson.M

	lessonsCollection := session.DB(dbName).C(course + "Edu.Lessons")

	err := lessonsCollection.Find(bson.M{"grade": grade, "course": course}).All(&queryMap)

	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn(fmt.Sprintf("Could not find lessons for %sEdu course and grade %v!", course, grade))
	}

	lessonArray, err := bson.MarshalJSON(queryMap)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not marshal listLessons request in JSON!")
	}

	if string(lessonArray) == "null\n" {
		return "notFound"
	}

	result := ""
	_, err = jsonparser.ArrayEach(lessonArray, func(value []byte, dataType jsonparser.ValueType, offset int, err1 error) {
		id, _ := jsonparser.GetString(value, "_id", "$oid")
		result = result + id + "\n"
	})
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Unable to iterate JSON array!")
	}

	return result
}

func AddLesson(course string, lesson string) {
	lessonsCollection := session.DB(dbName).C(course + "Edu.Lessons")

	var document map[string]interface{}
	err := json.Unmarshal([]byte(lesson), &document)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Cannot unmarshal lesson into document!")
	}

	err = lessonsCollection.Insert(document)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Cannot insert lesson in database!")
	}
}

func GetLesson(course, id string) string {
	lessonsCollection := session.DB(dbName).C(course + "Edu.Lessons")

	var lessonQuery []bson.M

	err := lessonsCollection.FindId(bson.ObjectIdHex(id)).All(&lessonQuery)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error":  err,
			"lessonID": id,
		}).Warn("Could not find lesson in database!")
	}

	lesson, err := bson.MarshalJSON(lessonQuery)
	if err != nil {
		APILogger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Could not marshal getLesson request in JSON!")
	}

	result := string(lesson)

	if result == "null\n" {
		return "notFound"
	}

	result = strings.Trim(result, "[")
	result = result[:len(result)-2]

	return result
}