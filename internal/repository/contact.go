package repository

import (
	"context"
	"time"

	"github.com/diana-gemini/drevmass/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ContactRepository struct {
	db *pgxpool.Pool
}

func NewContactRepository(db *pgxpool.Pool) models.ContactRepository {
	return &ContactRepository{db: db}
}

func (h *ContactRepository) CreateContactInformation(c context.Context, contact models.Contact) (int, error) {
	var contactID int
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	userQuery := `INSERT INTO info_contact(
		name, created_at)
		VALUES ($1, $2) returning id;`
	err := h.db.QueryRow(c, userQuery, contact.Name, currentTime).Scan(&contactID)
	if err != nil {
		return 0, err
	}
	return contactID, nil
}

func (h *ContactRepository) UpdateContactInformation(c context.Context, contact models.Contact, updateContact models.Contact) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	userQuery := `UPDATE info_contact SET name=$1, updated_at=$2 WHERE id=$3;`

	_, err := h.db.Exec(c, userQuery, updateContact.Name, currentTime, contact.ID)
	if err != nil {
		return err
	}
	return nil
}

func (h *ContactRepository) GetContactInformation(c context.Context) (models.Contact, error) {
	contact := models.Contact{}

	query := `SELECT id, name FROM info_contact`
	row := h.db.QueryRow(c, query)
	err := row.Scan(&contact.ID, &contact.Name)
	if err != nil {
		return contact, err
	}

	return contact, nil
}

func (h *ContactRepository) DeleteContactInformation(c context.Context, contact models.Contact) error {
	userQuery := `DELETE FROM info_contact WHERE id=$1;`
	_, err := h.db.Exec(c, userQuery, contact.ID)
	if err != nil {
		return err
	}
	return nil
}
