package handler

import "fmt"

func errParamIsRequired(name, typ string) error {
	return fmt.Errorf("param: %s (type: %s) is required", name, typ)
}

// CreateOpening
type CreateOpeningRequest struct {
	Role     string `json:"role" example:"Backend Developer"`
	Company  string `json:"company" example:"Google"`
	Location string `json:"location" example:"Remote"`
	Remote   *bool  `json:"remote" example:"true"`
	Link     string `json:"link" example:"https://google.com/careers/123"`
	Salary   int64  `json:"salary" example:"120000"`
}

func (r *CreateOpeningRequest) Validate() error {

	if r.Role == "" && r.Company == "" && r.Location == "" && r.Link == "" && r.Remote == nil && r.Salary <= 0 {
		return fmt.Errorf("request body is empty or invalid")
	}
	if r.Role == "" {
		return errParamIsRequired("role", "string")
	}
	if r.Company == "" {
		return errParamIsRequired("company", "string")
	}
	if r.Location == "" {
		return errParamIsRequired("location", "string")
	}
	if r.Link == "" {
		return errParamIsRequired("link", "string")
	}
	if r.Remote == nil {
		return errParamIsRequired("remote", "bool")
	}
	if r.Salary <= 0 {
		return errParamIsRequired("salary", "int64")
	}

	return nil
}

// UpdateOpening
type UpdateOpeningRequest struct {
	Role     string `json:"role" example:"Senior Backend Developer"`
	Company  string `json:"company" example:"Google"`
	Location string `json:"location" example:"Lisbon, Portugal"`
	Remote   *bool  `json:"remote" example:"true"`
	Link     string `json:"link" example:"https://google.com/careers/456"`
	Salary   int64  `json:"salary" example:"150000"`
}

func (r *UpdateOpeningRequest) Validate() error {
	// If any field is provided, it is valid

	if r.Role != "" || r.Company != "" || r.Location != "" || r.Link != "" || r.Remote != nil || r.Salary > 0 {
		return nil
	}

	// Otherwise, return an error
	return fmt.Errorf("at least one field must be provided to update")
}
