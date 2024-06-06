package models

type Currency struct {
	ID    int     `db:"ID"`
	Title string  `db:"TITLE"`
	Code  string  `db:"CODE"`
	Value float64 `db:"VALUE"`
	ADate string  `db:"A_DATE"`
}
