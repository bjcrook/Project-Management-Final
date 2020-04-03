package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func openConnection() {
	var pw passwordStruct

	jsonFile, err := os.Open("./passwords/password.json")

	if err != nil {
		panic(err.Error())
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &pw)

	db, _ = sql.Open("mysql", "root:"+pw.Pw+"@/"+pw.Database)

	jsonFile.Close()
}

func createTable() {
	command := `CREATE TABLE IF NOT EXISTS Employee_T
	(First_Name CHAR(255), Last_Name CHAR(255),
	Date_Of_Birth DATE, Hire_Date DATE, Department CHAR(255),
	Job_Title CHAR(255), Is_Salaried BOOLEAN, Compensation FLOAT(10,2),
	PRIMARY KEY (First_Name, Last_Name));`

	_, err := db.Exec(command)
	if err != nil {
		panic(err.Error())
	}
}

func employeeList(employee ...string) []employeeStuct {
	var data []employeeStuct
	var temp employeeStuct
	var params string

	if len(employee) == 2 {
		params = " WHERE First_Name = \"" + employee[0] + "\" AND Last_Name = \"" + employee[1] + "\""
	}

	if len(employee) == 1 {
		if employee[0] != "All" {
			params = " WHERE Department = \"" + employee[0] + "\""
		}
	}

	results, err := db.Query("SELECT * FROM Employee_T" + params)
	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		err = results.Scan(&temp.FirstName, &temp.LastName,
			&temp.DateOfBirth, &temp.HireDate, &temp.Department,
			&temp.JobTitle, &temp.IsSalaried, &temp.Compensation)

		if err != nil {
			panic(err.Error())
		} else {
			data = append(data, temp)
		}
	}

	if len(data) == 0 {
		return nil
	}

	return data
}

func getDepartments() []departmentsStruct {
	departments := []departmentsStruct{{"All"}}
	var temp departmentsStruct

	results, err := db.Query("SELECT DISTINCT Department from Employee_T")
	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		err = results.Scan(&temp.Name)

		if err != nil {
			panic(err.Error())
		} else {
			departments = append(departments, temp)
		}
	}

	return departments
}

func addEmployee(employee employeeStuct) addEmployeeStruct {
	var success addEmployeeStruct

	command := "INSERT INTO Employee_T VALUES(\""
	command += employee.FirstName + "\", \""
	command += employee.LastName + "\", \""
	command += employee.DateOfBirth + "\", \""
	command += employee.HireDate + "\", \""
	command += employee.Department + "\", \""
	command += employee.JobTitle + "\", "
	command += strconv.FormatBool(employee.IsSalaried) + ", "
	command += strconv.FormatFloat(float64(employee.Compensation), 'f', 2, 32) + ");"

	_, err := db.Exec(command)
	if err != nil {
		return success
	}

	success.Success = true
	return success
}

func deleteEmployee(data []string) {
	_, err := db.Query("DELETE FROM Employee_T WHERE First_Name=\"" + data[0] + "\" AND Last_Name=\"" + data[1] + "\"")
	if err != nil {
		panic(err.Error())
	}
}

func updateEmployee(name []string, data employeeStuct) {
	results, err := db.Query("SELECT Date_Of_Birth, Hire_Date FROM Employee_T WHERE First_Name=\"" + name[0] + "\" AND Last_Name=\"" + name[1] + "\"")
	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		err = results.Scan(&data.DateOfBirth, &data.HireDate)

		if err != nil {
			panic(err.Error())
		}
	}

	deleteEmployee(name)
	addEmployee(data)
}

func payrollList() []summaryStruct {
	var data []summaryStruct
	var temp summaryStruct

	results, err := db.Query("select Department, COUNT(Department), SUM(Compensation), Is_Salaried from Employee_T group by Department, Is_Salaried")
	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		err = results.Scan(&temp.Name, &temp.Count, &temp.Sum, &temp.IsSalaried)

		if err != nil {
			panic(err.Error())
		} else {
			data = append(data, temp)
		}
	}

	return data
}

func closeConnection() int {
	db.Close()
	return 0
}
