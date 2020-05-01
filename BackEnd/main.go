package main

import (
  "fmt"
  "log"
  "net/http"
)

/*
  This file contains the main function. The main function is automatically called
  when the program is ran. The main function creates http handlers that handle the
  various requests made to the website. This function also creates the server and
  destroys the server upon the ending of the program.
*/

func main() {
  http.HandleFunc("/departments/", departmentsHandler)
  http.HandleFunc("/employees/", employeesHandler)
  http.HandleFunc("/newhire/", newhireHandler)
  http.HandleFunc("/newhire/post/", newhirePostHandler)
  http.HandleFunc("/edit/", editHandler)
  http.HandleFunc("/update/", updateHandler)
  http.HandleFunc("/delete/", deleteHandler)
  http.HandleFunc("/payroll/", payrollHandler)
  http.HandleFunc("/login/", loginHandler)
  http.HandleFunc("/login/post/", loginPostHandler)
  http.HandleFunc("/logout/", logoutHandler)
  http.HandleFunc("/", defaultHandler)

  fmt.Println("Starting server on port 8080")
  http.ListenAndServe(":8080", nil)
  log.Fatal(closeConnection())
}
