package repository

import (
	"context"
	"time"

	"github.com/diana-gemini/drevmass/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type CompanyRepository struct {
	db *pgxpool.Pool
}

func NewCompanyRepository(db *pgxpool.Pool) models.CompanyRepository {
	return &CompanyRepository{db: db}
}

func (h *CompanyRepository) CreateCompanyInformation(c context.Context, company models.Company) (int, error) {
	var companyID int
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	userQuery := `INSERT INTO info_company(
		name, description, image, created_at)
		VALUES ($1, $2, $3, $4) returning id;`
	err := h.db.QueryRow(c, userQuery, company.Name, company.Description, company.Image, currentTime).Scan(&companyID)
	if err != nil {
		return 0, err
	}
	return companyID, nil
}

func (h *CompanyRepository) UpdateCompanyInformation(c context.Context, company models.Company, updateCompany models.Company) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	userQuery := `UPDATE info_company SET name=$1, description=$2, image=$3, updated_at=$4 WHERE id=$5;`

	_, err := h.db.Exec(c, userQuery, updateCompany.Name, updateCompany.Description, updateCompany.Image, currentTime, company.ID)
	if err != nil {
		return err
	}
	return nil
}

func (h *CompanyRepository) GetCompanyInformation(c context.Context) (models.Company, error) {
	company := models.Company{}

	query := `SELECT id, name, description, image FROM info_company`
	row := h.db.QueryRow(c, query)
	err := row.Scan(&company.ID, &company.Name, &company.Description, &company.Image)

	if err != nil {
		return company, err
	}

	return company, nil
}

func (h *CompanyRepository) DeleteCompanyInformation(c context.Context, company models.Company) error {
	userQuery := `DELETE FROM info_company WHERE id=$1;`
	_, err := h.db.Exec(c, userQuery, company.ID)
	if err != nil {
		return err
	}
	return nil
}
