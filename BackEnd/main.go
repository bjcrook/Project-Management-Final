package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	openConnection()
	createTable()
	http.HandleFunc("/departments/", departmentsHandler)
	http.HandleFunc("/employees/", employeesHandler)
	http.HandleFunc("/newhire/", newhireHandler)
	http.HandleFunc("/newhire/post/", newhirePostHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/update/", updateHandler)
	http.HandleFunc("/delete/", deleteHandler)
	http.HandleFunc("/payroll/", payrollHandler)
	http.HandleFunc("/", defaultHandler)
	fmt.Println("Starting server on port 8080")
	http.ListenAndServe(":8080", nil)
	log.Fatal(closeConnection())
}
