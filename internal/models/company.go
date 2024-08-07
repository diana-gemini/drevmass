package models

import "context"

type Company struct {
	ID          uint   `json:"id"`
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
	Image       string `json:"image" form:"image"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type CompanyRepository interface {
	CreateCompanyInformation(c context.Context, company Company) (int, error)
	UpdateCompanyInformation(c context.Context, company Company, updateCompany Company) error
	GetCompanyInformation(c context.Context) (Company, error)
	DeleteCompanyInformation(c context.Context, company Company) error
}
