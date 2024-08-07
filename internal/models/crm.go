package models

// Структуры данных для сделки
type FieldValue struct {
	Value interface{} `json:"value"`
}

type CustomFieldValue struct {
	FieldID int          `json:"field_id"`
	Values  []FieldValue `json:"values"`
}

type Deal struct {
	Name               string             `json:"name"`
	Price              int                `json:"price"`
	CustomFieldsValues []CustomFieldValue `json:"custom_fields_values"`
}
