package main

import (
  "database/sql"
  "encoding/json"
  "io/ioutil"
  "os"
  "strconv"

  _ "github.com/go-sql-driver/mysql"
)

/*
  This file contains all of the functions needed to interact with the SQL database.

  openConnection()  => this function establishes a connection with the SQL database
  createTable()     => this function creates the needed database tables if they don't already exist
  employeeList()    => this function returns a list of employees based on the passed in parameters
  getDepartments()  => this function returns a list of the various departments
  addEmployee()     => this function is used to add employee to the database
  deleteEmployee()  => this function is used to remove an employee from the database
  updateEmployee()  => this function is used to update an employee record
  payrollList()     => this function is used to generate the payroll data
  closeConnection() => this function destroys the connection with the SQL database
*/

var db *sql.DB
var loggedIn bool

//                                       FUNCTION <openConnection>
// ***********************************************************************************************************
func openConnection(temp loginStruct) bool {
  var pw passwordStruct // creates a struct to hold password data

  jsonFile, err := os.Open("./passwords/password.json") // opens the file containing password information

  // this if statement checks for errors with reading the file
  if err != nil {
    panic(err.Error())
  }

  // reads the .json file into a byte array then closes the file
  byteValue, _ := ioutil.ReadAll(jsonFile)
  jsonFile.Close()

  err = json.Unmarshal(byteValue, &pw) // converts and places the json data into the password struct

  // loops through each user read in from the password file
  for _, user := range pw.Users {
    if temp.Users == user { // true if the user was in the file
      if temp.Pw == pw.Pw { // true if the correct password was entered
        // creates the database connection and creates tables as needed
        db, _ = sql.Open("mysql", "root:"+pw.Pw+"@/test_server")
        createTable()
        return true
      }
    }
  }

  return false
}
// ***********************************************************************************************************

//                                          FUNCTION <createTable>
// ***********************************************************************************************************
func createTable() {
  // the variable command will create the needed SQL table if it doesn't already exist
  command := `CREATE TABLE IF NOT EXISTS Employee_T
  (First_Name CHAR(255), Last_Name CHAR(255),
  Date_Of_Birth DATE, Hire_Date DATE, Department CHAR(255),
  Job_Title CHAR(255), Is_Salaried BOOLEAN, Compensation FLOAT(10,2),
  PRIMARY KEY (First_Name, Last_Name));`

  _, err := db.Exec(command) // runs the SQL statement

  // this if statement checks for errors created by the executed statement
  if err != nil {
    panic(err.Error())
  }
}
// ***********************************************************************************************************

//                                          FUNCTION <employeeList>
// ***********************************************************************************************************
func employeeList(employee ...string) []employeeStuct {
  // creates variables that will be used in this function
  var data []employeeStuct
  var temp employeeStuct
  var params string

  // true when the number of passed in parameters equals two
  if len(employee) == 2 {
    params = " WHERE First_Name = \"" + employee[0] + "\" AND Last_Name = \"" + employee[1] + "\""
  }

  // true when the number of passed in parameters equals one
  if len(employee) == 1 {
    // true when the "default" parameter is the only thing passed in
    if employee[0] != "All" {
      params = " WHERE Department = \"" + employee[0] + "\""
    }
  }

  // executes the query, stores the results in results and checks for errors
  results, err := db.Query("SELECT * FROM Employee_T" + params)
  if err != nil {
    panic(err.Error())
  }

  // loops while there is still data in results
  for results.Next() {
    // stores the results in the temporary employee struct
    err = results.Scan(&temp.FirstName, &temp.LastName,
      &temp.DateOfBirth, &temp.HireDate, &temp.Department,
      &temp.JobTitle, &temp.IsSalaried, &temp.Compensation)

    if err != nil { // true when there is an error
      panic(err.Error())
    } else {
      data = append(data, temp) // adds the employee to the data slice
    }
  }

  // true when the SQL query returns no data
  if len(data) == 0 {
    return nil
  }

  return data
}
// ***********************************************************************************************************

//                                          FUNCTION <getDepartments>
// ***********************************************************************************************************
func getDepartments() []departmentsStruct {
  // creates variables needed for the function and initializes them where needed
  departments := []departmentsStruct{{"All"}}
  var temp departmentsStruct

  // executes the query and checks for errors
  results, err := db.Query("SELECT DISTINCT Department from Employee_T")
  if err != nil {
    panic(err.Error())
  }

  // loops while there is still data in results
  for results.Next() {
    err = results.Scan(&temp.Name) // stores the data in the temporary department struct

    // true when there is an error
    if err != nil {
      panic(err.Error())
    } else {
      departments = append(departments, temp) // adds the temporary department to the department slice
    }
  }

  return departments
}
// ***********************************************************************************************************

//                                           FUNCTION <addEmployee>
// ***********************************************************************************************************
func addEmployee(employee employeeStuct) addEmployeeStruct {
  var success addEmployeeStruct

  // the command variable holds the SQL command responsible for inserting employee data into the table
  command := "INSERT INTO Employee_T VALUES(\""
  command += employee.FirstName + "\", \""
  command += employee.LastName + "\", \""
  command += employee.DateOfBirth + "\", \""
  command += employee.HireDate + "\", \""
  command += employee.Department + "\", \""
  command += employee.JobTitle + "\", "
  command += strconv.FormatBool(employee.IsSalaried) + ", "
  command += strconv.FormatFloat(float64(employee.Compensation), 'f', 2, 32) + ");"

  // executes the command and then checks for any errors
  _, err := db.Exec(command)
  if err != nil {
    return success
  }

  success.Success = true
  return success
}
// ***********************************************************************************************************

//                                         FUNCTION <deleteEmployee>
// ***********************************************************************************************************
func deleteEmployee(data []string) bool {
  // true if data contains two elements
  if len(data) == 2 {
    temp := employeeList(data[0], data[1])

    if temp == nil {
      return false
    }
  } else {
    return false
  }

  // executes the query that deletes an employee from the database based on the passed in primary key value
  _, err := db.Query("DELETE FROM Employee_T WHERE First_Name=\"" + data[0] + "\" AND Last_Name=\"" + data[1] + "\"")
  if err != nil {
    panic(err.Error())
  }

  return true
}
// ***********************************************************************************************************

//                                         FUNCTION <updateEmployee>
// ***********************************************************************************************************
func updateEmployee(name []string, data employeeStuct) {
  // executes the query that will return the employee that we want to update and checks for errors
  results, err := db.Query("SELECT Date_Of_Birth, Hire_Date FROM Employee_T WHERE First_Name=\"" + name[0] + "\" AND Last_Name=\"" + name[1] + "\"")
  if err != nil {
    panic(err.Error())
  }

  // loops while there is still data in results (in this case it will loop at most once)
  for results.Next() {
    // stores the results and checks for errors
    err = results.Scan(&data.DateOfBirth, &data.HireDate)

    if err != nil {
      panic(err.Error())
    }
  }

  deleteEmployee(name) // calls the deleteEmployee function to first remove the employee from the database
  addEmployee(data) // calls the addEmployee function to then re-add the employee with the updated data
}
// ***********************************************************************************************************

//                                           FUNCTION <payrollList>
// ***********************************************************************************************************
func payrollList() []summaryStruct {
  // creates the needed variables for this function
  var data []summaryStruct
  var temp summaryStruct

  // executes and stores the results of the query, then checks for errors
  results, err := db.Query("select Department, COUNT(Department), SUM(Compensation), Is_Salaried from Employee_T group by Department, Is_Salaried")
  if err != nil {
    panic(err.Error())
  }

  // loops while there is still data in results
  for results.Next() {
    // stores the results in the temp variable
    err = results.Scan(&temp.Name, &temp.Count, &temp.Sum, &temp.IsSalaried)

    // true when there is an error
    if err != nil {
      panic(err.Error())
    } else {
      data = append(data, temp) // adds the temporary data to the data slice
    }
  }

  return data
}
// ***********************************************************************************************************

//                                         FUNCTION <closeConnection>
// ***********************************************************************************************************
func closeConnection() int {
  // closes the database connection and returns 0 to indicate that the operation was successful
  db.Close()
  loggedIn = false
  return 0
}
// ***********************************************************************************************************

