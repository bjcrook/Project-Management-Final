package main

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
	Pw       string
	Database string
}
