package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/departments/", http.StatusSeeOther)
}

func departmentsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("We got a request on", r.URL.String())

	var data []employeeStuct

	departmentURL := r.URL.String()[len("/departments/"):]

	if len(departmentURL) > 0 {
		data = employeeList(strings.ReplaceAll(departmentURL, "%20", " "))
	} else {
		data = employeeList("All")
	}

	departments := getDepartments()
	t, _ := template.ParseFiles("./FrontEnd/department_header.html")
	t.Execute(w, departments)
	t, _ = template.ParseFiles("./FrontEnd/departments.html")
	t.Execute(w, data)
}

func employeesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("We got a request on", r.URL.String())
	employeeURL := strings.Split(r.URL.String()[len("/employees/"):], "_")

	if len(employeeURL) == 2 {
		data := employeeList(employeeURL[0], employeeURL[1])
		t, _ := template.ParseFiles("./FrontEnd/employees.html")
		t.Execute(w, data[0])
	} else {
		t, _ := template.ParseFiles("./FrontEnd/employees.html")
		t.Execute(w, nil)
	}
}

func newhireHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("We got a request on", r.URL.String())

	t, _ := template.ParseFiles("./FrontEnd/newhire.html")
	t.Execute(w, nil)
}

func newhirePostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("We got a request on", r.URL.String())
	var isSalaried bool

	checked := r.FormValue("IsSalaried")

	if checked == "on" {
		isSalaried = true
	} else {
		isSalaried = false
	}

	compensation, _ := strconv.ParseFloat(r.FormValue("Compensation"), 32)

	temp := employeeStuct{
		FirstName:    r.FormValue("FirstName"),
		LastName:     r.FormValue("LastName"),
		DateOfBirth:  r.FormValue("DateOfBirth"),
		HireDate:     r.FormValue("HireDate"),
		Department:   r.FormValue("Department"),
		JobTitle:     r.FormValue("JobTitle"),
		IsSalaried:   isSalaried,
		Compensation: float32(compensation),
	}

	success := addEmployee(temp)

	t, _ := template.ParseFiles("./FrontEnd/newhire.html")
	t.Execute(w, nil)
	t, _ = template.ParseFiles("./FrontEnd/newhire_foot.html")
	t.Execute(w, success)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("We got a request on", r.URL.String())
	employeeURL := strings.Split(r.URL.String()[len("/edit/"):], "_")

	if len(employeeURL) == 2 {
		data := employeeList(employeeURL[0], employeeURL[1])
		t, _ := template.ParseFiles("./FrontEnd/edit.html")
		t.Execute(w, data[0])
	} else {
		t, _ := template.ParseFiles("./FrontEnd/employees.html")
		t.Execute(w, nil)
	}
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("We got a request on", r.URL.String())
	employeeURL := strings.Split(r.URL.String()[len("/update/"):], "_")

	employeeURL[1] = strings.ReplaceAll(employeeURL[1], "/", "")
	newURL := "/employees/" + employeeURL[0] + "_" + employeeURL[1]

	var isSalaried bool

	checked := r.FormValue("IsSalaried")

	if checked == "on" {
		isSalaried = true
	} else {
		isSalaried = false
	}

	compensation, _ := strconv.ParseFloat(r.FormValue("Compensation"), 32)

	temp := employeeStuct{
		FirstName:    employeeURL[0],
		LastName:     employeeURL[1],
		Department:   r.FormValue("Department"),
		JobTitle:     r.FormValue("JobTitle"),
		IsSalaried:   isSalaried,
		Compensation: float32(compensation),
	}

	if len(employeeURL) == 2 {
		updateEmployee(employeeURL, temp)
		http.Redirect(w, r, newURL, http.StatusSeeOther)
	} else {
		t, _ := template.ParseFiles("./FrontEnd/employees.html")
		t.Execute(w, nil)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("We got a request on", r.URL.String())

	employeeURL := strings.Split(r.URL.String()[len("/update/"):], "_")

	employeeURL[1] = strings.ReplaceAll(employeeURL[1], "/", "")

	temp := employeeStuct{
		FirstName: employeeURL[0],
		LastName:  employeeURL[1],
	}

	if len(employeeURL) == 2 {
		deleteEmployee(employeeURL)
		t, _ := template.ParseFiles("./FrontEnd/remove.html")
		t.Execute(w, temp)
	} else {
		t, _ := template.ParseFiles("./FrontEnd/employees.html")
		t.Execute(w, nil)
	}
}

func payrollHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("We got a request on", r.URL.String())

	var grandTotal grandTotalStruct
	data := payrollList()

	for i := range data {
		if data[i].IsSalaried {
			data[i].Sum /= 52
			data[i].SumString = fmt.Sprintf("%.2f", data[i].Sum)
			data[i].Name = data[i].Name + " Salaried"
		} else {
			data[i].Sum *= 40
			data[i].SumString = fmt.Sprintf("%.2f", data[i].Sum)
			data[i].Name = data[i].Name + " Hourly"
		}
		grandTotal.Total += data[i].Sum
	}

	grandTotal.TotalString = fmt.Sprintf("%.2f", grandTotal.Total)

	t, _ := template.ParseFiles("./FrontEnd/payroll.html")
	t.Execute(w, data)
	t, _ = template.ParseFiles("./FrontEnd/payroll_foot.html")
	t.Execute(w, grandTotal)
}
