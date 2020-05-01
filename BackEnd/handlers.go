package main

import (
  "fmt"
  "net/http"
  "strconv"
  "strings"
  "text/template"
)

/*
  This function contains the implementation of the handlers defined in the main.go file.
  The functions are invoked by visiting the various URL's on the website.

  defaultHandler()     => this handler is called for unspecified URL entires, redirects the user back to the
                          login screen
  departmentsHandler() => this handler deals with getting and displaying data to the list on the departments
                          page
  employeesHandler()   => this handler deals with displaying the specific employee that was chosen from the
                          departments page
  newhireHandler()     => this handler deals with displaying a form that adds employees to the database
  newhirePostHandler() => this handler deals with the input from the new hire form and actually adds the
                          employee to the database
  editHandler()        => this handler deals with displaying the form that allows for employee edits
  updateHandler()      => this handler deals with the form input that actually modifies employee data
  deleteHandler()      => this handler deals with the deletion of an employee form the database
  payrollHandler()     => this handler deals with the data gathering and display of payroll information
  loginHandler()       => this handler displays the login screen / login form
  loginPostHandler()   => this handler deals with the login forms input and decides if the username/password
                          combo was valid or not
  logoutHandler()      => this handler logs the user out of the website and returns them to the login screen
*/

//                                        FUNCTION <defaultHandler>
// ***********************************************************************************************************
func defaultHandler(w http.ResponseWriter, r *http.Request) {
  // this is the default handler that program defaults to
  http.Redirect(w, r, "/login/", http.StatusSeeOther)
}
// ***********************************************************************************************************

//                                       FUNCTION <departmentsHandler>
// ***********************************************************************************************************
func departmentsHandler(w http.ResponseWriter, r *http.Request) {
  // checks to see if the user is logged in before continuing
  if !loggedIn {
    http.Redirect(w, r, "/login/", http.StatusSeeOther)
    return
  }

  var data []employeeStuct

  // value set based on the users URL
  departmentURL := r.URL.String()[len("/departments/"):]
  urlExtra := strings.Split(r.URL.String()[len("/departments/"):], "/")

  // if the URL contains "departments/hourly" 
  if urlExtra[0] == "hourly" {
    // asks for all of the employees
    data = employeeList("All")
    departments := getDepartments()
    // displays all of the hourly employees
    t, _ := template.ParseFiles("./FrontEnd/department_header.html")
    t.Execute(w, departments)
    t, _ = template.ParseFiles("./FrontEnd/departments_hourly.html")
    t.Execute(w, data)
    return
  } else if urlExtra[0] == "salary" { // if the URL contains "departments/salary" 
    // asks for all of the employees
    data = employeeList("All")
    departments := getDepartments()
    // displays all of the salaried employees
    t, _ := template.ParseFiles("./FrontEnd/department_header.html")
    t.Execute(w, departments)
    t, _ = template.ParseFiles("./FrontEnd/departments_salary.html")
    t.Execute(w, data)
    return
  }

  // true if the department URL has no additional parameters
  if len(departmentURL) > 0 {
    data = employeeList(strings.ReplaceAll(urlExtra[0], "%20", " ")) // removes the spaces in the URL
  } else {
    data = employeeList("All")
  }

  // true when the URL contains a specific department and a specific pay type
  if len(urlExtra) == 2 {
    // true when the pay type is salary
    if urlExtra[1] == "salary" {
      // loads and executes the appropriate html file
      departments := getDepartments()
      t, _ := template.ParseFiles("./FrontEnd/department_header.html")
      t.Execute(w, departments)
      t, _ = template.ParseFiles("./FrontEnd/departments_salary.html")
      t.Execute(w, data)
      return
      // true when the pay type is hourly
    } else if urlExtra[1] == "hourly" {
      // loads and executes the appropriate html file
      departments := getDepartments()
      t, _ := template.ParseFiles("./FrontEnd/department_header.html")
      t.Execute(w, departments)
      t, _ = template.ParseFiles("./FrontEnd/departments_hourly.html")
      t.Execute(w, data)
      return
    }
  }

  // loads and executes the appropriate html file
  departments := getDepartments()
  t, _ := template.ParseFiles("./FrontEnd/department_header.html")
  t.Execute(w, departments)
  t, _ = template.ParseFiles("./FrontEnd/departments.html")
  t.Execute(w, data)
}
// ***********************************************************************************************************

//                                        FUNCTION <employeesHandler>
// ***********************************************************************************************************
func employeesHandler(w http.ResponseWriter, r *http.Request) {
  // checks to see if the user is logged in before continuing
  if !loggedIn {
    http.Redirect(w, r, "/login/", http.StatusSeeOther)
    return
  }

  // value set based on the current URL
  employeeURL := strings.Split(r.URL.String()[len("/employees/"):], "_")

  // loads in the employees html file
  t, _ := template.ParseFiles("./FrontEnd/employees.html")

  // true when there are two URL parameters
  if len(employeeURL) == 2 {
    data := employeeList(employeeURL[0], employeeURL[1])
    // true if no employees where returned (most likely a bad SQL call)
    if data != nil {
      t.Execute(w, data[0])
      return
    }
  }

  // executes the previously loaded html file
  t.Execute(w, nil)
}
// ***********************************************************************************************************

//                                         FUNCTION <newhireHandler>
// ***********************************************************************************************************
func newhireHandler(w http.ResponseWriter, r *http.Request) {
  // checks to see if the user is logged in before continuing
  if !loggedIn {
    http.Redirect(w, r, "/login/", http.StatusSeeOther)
    return
  }

  // loads in and executes the newhire html file
  t, _ := template.ParseFiles("./FrontEnd/newhire.html")
  t.Execute(w, nil)
}
// ***********************************************************************************************************

//                                       FUNCTION <newhirePostHandler>
// ***********************************************************************************************************
func newhirePostHandler(w http.ResponseWriter, r *http.Request) {
  // checks to see if the user is logged in before continuing
  if !loggedIn {
    http.Redirect(w, r, "/login/", http.StatusSeeOther)
    return
  }

  var isSalaried bool

  // determines if the salaried box is checked or not
  checked := r.FormValue("IsSalaried")

  // true when the box is checked
  if checked == "on" {
    isSalaried = true
  } else {
    isSalaried = false
  }

  // converts the float value into a string
  compensation, _ := strconv.ParseFloat(r.FormValue("Compensation"), 32)

  // reads in the data from the form and places it in the temporary employee structure
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

  // attempts to add the employee to the database
  success := addEmployee(temp)

  // loads and executes the html file
  t, _ := template.ParseFiles("./FrontEnd/newhire.html")
  t.Execute(w, nil)
  t, _ = template.ParseFiles("./FrontEnd/newhire_foot.html")
  t.Execute(w, success)
}
// ***********************************************************************************************************

//                                           FUNCTION <editHandler>
// ***********************************************************************************************************
func editHandler(w http.ResponseWriter, r *http.Request) {
  // checks to see if the user is logged in before continuing
  if !loggedIn {
    http.Redirect(w, r, "/login/", http.StatusSeeOther)
    return
  }

  // sets the value based on the current URL then loads the html file
  employeeURL := strings.Split(r.URL.String()[len("/edit/"):], "_")
  t, _ := template.ParseFiles("./FrontEnd/edit.html")

  // true when there are two URL parameters
  if len(employeeURL) == 2 {
    data := employeeList(employeeURL[0], employeeURL[1])
    // true if the employeeList call returned a employee
    if data != nil {
      // executes the html page
      t.Execute(w, data[0])
      return
    }
  }

  // executes the html page
  t.Execute(w, nil)
}
// ***********************************************************************************************************

//                                          FUNCTION <updateHandler>
// ***********************************************************************************************************
func updateHandler(w http.ResponseWriter, r *http.Request) {
  // checks to see if the user is logged in before continuing
  if !loggedIn {
    http.Redirect(w, r, "/login/", http.StatusSeeOther)
    return
  }

  // sets the value based on the current URL
  employeeURL := strings.Split(r.URL.String()[len("/update/"):], "_")

  // replaces the "/" in the read in URL then modifies the URL's output
  employeeURL[1] = strings.ReplaceAll(employeeURL[1], "/", "")
  newURL := "/employees/" + employeeURL[0] + "_" + employeeURL[1]

  var isSalaried bool

  // checks the form to see if the salaried box is checked or not
  checked := r.FormValue("IsSalaried")

  // true when the salaried box was checked
  if checked == "on" {
    isSalaried = true
  } else {
    isSalaried = false
  }

  // converts the float to a string
  compensation, _ := strconv.ParseFloat(r.FormValue("Compensation"), 32)

  // reads in all of the remaining values from the update form
  temp := employeeStuct{
    FirstName:    employeeURL[0],
    LastName:     employeeURL[1],
    Department:   r.FormValue("Department"),
    JobTitle:     r.FormValue("JobTitle"),
    IsSalaried:   isSalaried,
    Compensation: float32(compensation),
  }

  // true if the employee URL had two parameters
  if len(employeeURL) == 2 {
    // updates the employee record then redirects the user
    updateEmployee(employeeURL, temp)
    http.Redirect(w, r, newURL, http.StatusSeeOther)
  } else {
    // parses the html file before executing it
    t, _ := template.ParseFiles("./FrontEnd/employees.html")
    t.Execute(w, nil)
  }
}
// ***********************************************************************************************************

//                                          FUNCTION <deleteHandler>
// ***********************************************************************************************************
func deleteHandler(w http.ResponseWriter, r *http.Request) {
  // checks to see if the user is logged in before continuing
  if !loggedIn {
    http.Redirect(w, r, "/login/", http.StatusSeeOther)
    return
  }

  // value is set based on the current URL
  employeeURL := strings.Split(r.URL.String()[len("/update/"):], "_")

  // true when there are two URL parameters
  if len(employeeURL) == 2 {
    employeeURL[1] = strings.ReplaceAll(employeeURL[1], "/", "")
  } else {
    // loads and then executes the selected html file
    t, _ := template.ParseFiles("./FrontEnd/employees.html")
    t.Execute(w, nil)
    return
  }

  // creates a temporary employee structure
  temp := employeeStuct{
    FirstName: employeeURL[0],
    LastName:  employeeURL[1],
  }

  // true when there are two URL parameters
  if len(employeeURL) == 2 {
    data := deleteEmployee(employeeURL)

    // true when the employee was deleted successfully
    if data {
      // loads and executes the html file
      t, _ := template.ParseFiles("./FrontEnd/remove.html")
      t.Execute(w, temp)
      return
    }
  }

  // loads and executes the html file
  t, _ := template.ParseFiles("./FrontEnd/employees.html")
  t.Execute(w, nil)
}
// ***********************************************************************************************************

//                                         FUNCTION <payrollHandler>
// ***********************************************************************************************************
func payrollHandler(w http.ResponseWriter, r *http.Request) {
  // checks to see if the user is logged in before continuing
  if !loggedIn {
    http.Redirect(w, r, "/login/", http.StatusSeeOther)
    return
  }

  // creates need variables and calls for all the payroll list information
  var grandTotal grandTotalStruct
  data := payrollList()

  // loops through all of the payroll list information
  for i := range data {
    // true when the employee is a salaried employee
    if data[i].IsSalaried {
      // salaried payment calculations (per week)
      data[i].Sum /= 52
      data[i].SumString = fmt.Sprintf("%.2f", data[i].Sum)
    } else {
      // hourly payment calculations (per week)
      data[i].Sum *= 40
      data[i].SumString = fmt.Sprintf("%.2f", data[i].Sum)
    }
    grandTotal.Total += data[i].Sum // total payment calculations counter
  }

  // formats and sets the grand total
  grandTotal.TotalString = fmt.Sprintf("%.2f", grandTotal.Total)

  // loads and executes the html files
  t, _ := template.ParseFiles("./FrontEnd/payroll.html")
  t.Execute(w, data)
  t, _ = template.ParseFiles("./FrontEnd/payroll_foot.html")
  t.Execute(w, grandTotal)
}
// ***********************************************************************************************************

//                                          FUNCTION <loginHandler>
// ***********************************************************************************************************
func loginHandler(w http.ResponseWriter, r *http.Request) {
  // checks to see if the user is logged in before continuing
  if loggedIn {
    http.Redirect(w, r, "/departments/", http.StatusSeeOther)
    return
  }

  // loads and executes the html file
  t, _ := template.ParseFiles("./FrontEnd/login.html")
  t.Execute(w, nil)
}
// ***********************************************************************************************************

//                                        FUNCTION <loginPostHandler>
// ***********************************************************************************************************
func loginPostHandler(w http.ResponseWriter, r *http.Request) {
  // loads the login information in a temporary login structure
  temp := loginStruct{
    Pw:    r.FormValue("Pw"),
    Users: r.FormValue("Database"),
  }

  // attempts to login to the website, true or false is returned to loggedIn
  loggedIn = openConnection(temp)

  // redirects the user back to the main login screen
  http.Redirect(w, r, "/login/", http.StatusSeeOther)
}
// ***********************************************************************************************************

//                                          FUNCTION <logoutHandler>
// ***********************************************************************************************************
func logoutHandler(w http.ResponseWriter, r *http.Request) {
  // sets the logged in value to false
  loggedIn = false

  // redirects the user back to the main login screen
  http.Redirect(w, r, "/login/", http.StatusSeeOther)
}
// ***********************************************************************************************************

