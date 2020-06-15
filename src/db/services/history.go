package services

import (
	"github.com/nuwak/go_finance/src/db"
	"time"
)

type HistoryStruct struct {
	QValue       string
	QInsertValue string
}

var history *HistoryStruct

func History() *HistoryStruct {
	if history == nil {
		history = &HistoryStruct{
			QValue:       "SELECT value FROM history where symbol = ? and created_d = ? order by created_d desc limit 1",
			QInsertValue: "INSERT INTO history (symbol, value, created_d) VALUES (?, ?, ?)",
		}
	}

	return history
}

func (history *HistoryStruct) GetValue(symbol *string, isToday bool) float64 {
	var value float64
	var date string

	if isToday {
		date = time.Now().Format("2006-01-02")
	} else {
		date = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	}

	err := db.DB.QueryRow(history.QValue, *symbol, date).Scan(&value)
	if err != nil {
		return 0
	}

	return value
}

func (history *HistoryStruct) AddValue(symbol *string, value *float64) {
	stmt, _ := db.DB.Prepare("INSERT INTO history (symbol, value) VALUES (?, ?)")
	_, err := stmt.Exec(*symbol, *value)
	if err != nil {
		panic(err)
	}
}
