package models

import "context"

type Contact struct {
	ID        uint   `json:"id"`
	Name      string `json:"name" form:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type ContactRepository interface {
	CreateContactInformation(c context.Context, contact Contact) (int, error)
	UpdateContactInformation(c context.Context, contact Contact, updateContact Contact) error
	GetContactInformation(c context.Context) (Contact, error)
	DeleteContactInformation(c context.Context, contact Contact) error
}
