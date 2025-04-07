package models

type Column struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	DataType   string `json:"data_type"`
	Confidence uint   `json:"conficence"`
}

type Table struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	RowCount    uint   `json:"row_count"`
	Columns     Column `json:"columns"`
}

type DataBase struct {
	ID               uint   `json:"id"`
	Name             string `json:"name"`
	ConnectionString string `json:"conection_string"`
	Tables           Table  `json:"tables"`
}
