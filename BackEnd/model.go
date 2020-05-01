package main

/*
  This file defines the various structures used throughout the program.
  The structures are used either to store user information, or to help
  make interaction between the program and the SQL database easier.
  Some of the structures also help with the login process.
*/

type employeeStuct struct {
  FirstName    string
  LastName     string
  DateOfBirth  string
  HireDate     string
  Department   string
  JobTitle     string
  IsSalaried   bool
  Compensation float32
}

type summaryStruct struct {
  Name       string
  Count      uint8
  Sum        float64
  SumString  string
  IsSalaried bool
}

type grandTotalStruct struct {
  Total       float64
  TotalString string
}

type departmentsStruct struct {
  Name string
}

type addEmployeeStruct struct {
  Success bool
}

type passwordStruct struct {
  Pw    string
  Users []string
}

type loginStruct struct {
  Pw    string
  Users string
}
