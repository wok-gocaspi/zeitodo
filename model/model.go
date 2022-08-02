package model

type Employee struct {
	ID          string `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
}

type Payload struct {
	Employees []Employee `json:"employees"`
}

type DbConfig struct {
	URL      string
	Database string
}