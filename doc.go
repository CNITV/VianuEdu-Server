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

/*
VianuEdu-Server is the server-side component for the VianuEdu educational software

At its core, it is a simple REST API which internally uses various methods and mgo driver calls in order to interact
with a MongoDB database with a predefined schema, as follows:
	[dbName]
	│
	├───[COURSE]Edu.Grades
	│   ├───{ ... }
	│   └───{ ... }
	├───[COURSE]Edu.Tests
	│   ├───{ ... }
	│   └───{ ... }
	├───Students.Accounts
	│   ├───{ ... }
	│   └───{ ... }
	├───Students.SubmittedAnswers
	│   ├───{ ... }
	│   └───{ ... }
	├───Teachers.Accounts
	│   ├───{ ... }
	│   └───{ ... }
	└───[dbName].TestList
	│   ├───{ ... }
	│   └───{ ... }
	└─── { ... }
The configuration files contain all variables marked between square brackets.
*/
package main
