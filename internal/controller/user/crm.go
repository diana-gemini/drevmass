package user

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	//"github.com/nurzzaat/group_1/internal/models"
)

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

func SendDeal(c *gin.Context) {
	// Данные сделки
	deals := []Deal{
		{
			Name:  "Название сделки",
			Price: 1000,
			CustomFieldsValues: []CustomFieldValue{
				{
					FieldID: 207367,
					Values: []FieldValue{
						{Value: "бир бале 3 300\nеки бале 4 400\nуш бале 3 300"},
					},
				},
				{
					FieldID: 113851,
					Values: []FieldValue{
						{Value: 300},
					},
				},
				{
					FieldID: 177119,
					Values: []FieldValue{
						{Value: 700},
					},
				},
				{
					FieldID: 207419,
					Values: []FieldValue{
						{Value: "саке баке 7075554433"},
					},
				},
			},
		},
	}

	accessToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImp0aSI6IjNkZTI1M2Y4MTZiY2FkM2FhYjZhN2NlZGEwOGNjMmRkMTg5ODY5NGRiMTRmNTU0ZTUwMTg2NTJlNzhiNzVlNTg0MGU3ZTQwZmE1MDg3Y2U3In0.eyJhdWQiOiI2MDJlODVkNC00N2QzLTRkMjktOGNmMy00ZGI3MWVkNzE1YzMiLCJqdGkiOiIzZGUyNTNmODE2YmNhZDNhYWI2YTdjZWRhMDhjYzJkZDE4OTg2OTRkYjE0ZjU1NGU1MDE4NjUyZTc4Yjc1ZTU4NDBlN2U0MGZhNTA4N2NlNyIsImlhdCI6MTcxOTUxMDAwNCwibmJmIjoxNzE5NTEwMDA0LCJleHAiOjE3MjM2ODAwMDAsInN1YiI6IjExMTkwMzIyIiwiZ3JhbnRfdHlwZSI6IiIsImFjY291bnRfaWQiOjMxODE0MTUwLCJiYXNlX2RvbWFpbiI6ImFtb2NybS5ydSIsInZlcnNpb24iOjIsInNjb3BlcyI6WyJjcm0iLCJmaWxlcyIsImZpbGVzX2RlbGV0ZSIsIm5vdGlmaWNhdGlvbnMiLCJwdXNoX25vdGlmaWNhdGlvbnMiXSwiaGFzaF91dWlkIjoiYTQ5NDY1M2MtMTVlNi00OTYyLThjMDMtMTAyNzZkYTI1MmUyIn0.TNVmn0BRd1MsrZEgqpV6kLSWLPq9H08nSNWRRN-hgKrVvO-mpQZH2ozlTbnr8KZKyKb6hyl6y_bsYbH6DCeP5nEKDR9ikEgQ5eXKvz83saWtcrm27GcVkPMkTp3eOhSa-mEjGKd3LLbb4cmqJ3AFpyFQoPsiknTP4Yyj2X9SNkuerVXVi1-Kn2_vfMIYURRh8fwBhkyEfGHcpB_dMHjmflBmRDGGPdZDpV4BOcOSimJJDWZZVpz7n_bRkPU3qsYdq2vsNzlJ8S86UNsf3QpwV3jpZehS75yzemrzaRtz4iDfzeSC-75hV6Y7KICIKBitEDM4g0y9xO3QRkcKhwmBJw" // Замените на ваш токен доступа
	url := "https://almasdjan05.amocrm.ru/api/v4/leads"                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 // Замените на ваш домен

	jsonData, err := json.Marshal(deals)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error response from server: %v", string(body))
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deal sent", "response": string(body)})
}

func SendDealFunc(products, user string, price int, bonus, overall float64) error {
	// Данные сделки
	deals := []Deal{
		{
			Name:  "Название сделки",
			Price: price,
			CustomFieldsValues: []CustomFieldValue{
				{
					FieldID: 207367,
					Values: []FieldValue{
						{Value: products},
					},
				},
				{
					FieldID: 113851,
					Values: []FieldValue{
						{Value: bonus},
					},
				},
				{
					FieldID: 177119,
					Values: []FieldValue{
						{Value: overall},
					},
				},
				{
					FieldID: 207419,
					Values: []FieldValue{
						{Value: user},
					},
				},
			},
		},
	}

	accessToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImp0aSI6IjNkZTI1M2Y4MTZiY2FkM2FhYjZhN2NlZGEwOGNjMmRkMTg5ODY5NGRiMTRmNTU0ZTUwMTg2NTJlNzhiNzVlNTg0MGU3ZTQwZmE1MDg3Y2U3In0.eyJhdWQiOiI2MDJlODVkNC00N2QzLTRkMjktOGNmMy00ZGI3MWVkNzE1YzMiLCJqdGkiOiIzZGUyNTNmODE2YmNhZDNhYWI2YTdjZWRhMDhjYzJkZDE4OTg2OTRkYjE0ZjU1NGU1MDE4NjUyZTc4Yjc1ZTU4NDBlN2U0MGZhNTA4N2NlNyIsImlhdCI6MTcxOTUxMDAwNCwibmJmIjoxNzE5NTEwMDA0LCJleHAiOjE3MjM2ODAwMDAsInN1YiI6IjExMTkwMzIyIiwiZ3JhbnRfdHlwZSI6IiIsImFjY291bnRfaWQiOjMxODE0MTUwLCJiYXNlX2RvbWFpbiI6ImFtb2NybS5ydSIsInZlcnNpb24iOjIsInNjb3BlcyI6WyJjcm0iLCJmaWxlcyIsImZpbGVzX2RlbGV0ZSIsIm5vdGlmaWNhdGlvbnMiLCJwdXNoX25vdGlmaWNhdGlvbnMiXSwiaGFzaF91dWlkIjoiYTQ5NDY1M2MtMTVlNi00OTYyLThjMDMtMTAyNzZkYTI1MmUyIn0.TNVmn0BRd1MsrZEgqpV6kLSWLPq9H08nSNWRRN-hgKrVvO-mpQZH2ozlTbnr8KZKyKb6hyl6y_bsYbH6DCeP5nEKDR9ikEgQ5eXKvz83saWtcrm27GcVkPMkTp3eOhSa-mEjGKd3LLbb4cmqJ3AFpyFQoPsiknTP4Yyj2X9SNkuerVXVi1-Kn2_vfMIYURRh8fwBhkyEfGHcpB_dMHjmflBmRDGGPdZDpV4BOcOSimJJDWZZVpz7n_bRkPU3qsYdq2vsNzlJ8S86UNsf3QpwV3jpZehS75yzemrzaRtz4iDfzeSC-75hV6Y7KICIKBitEDM4g0y9xO3QRkcKhwmBJw" // Замените на ваш токен доступа
	url := "https://almasdjan05.amocrm.ru/api/v4/leads"                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 // Замените на ваш домен

	jsonData, err := json.Marshal(deals)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
		return err

	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error response from server: %v", string(body))
		return err
	}

	return nil

}
