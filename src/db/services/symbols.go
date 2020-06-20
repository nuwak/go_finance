package services

import "github.com/nuwak/go_finance/src/db"

type Currency int

const (
	_            = iota
	USD Currency = iota
	RUB Currency = iota
)

type SymbolsItem struct {
	Symbol   string
	Currency Currency
}

type SymbolsStruct struct{}

var symbol *SymbolsStruct

func Symbols() *SymbolsStruct {
	if symbol == nil {
		symbol = &SymbolsStruct{}
	}

	return symbol
}

func (symbols *SymbolsStruct) Get(symbol *string) (*SymbolsItem, error) {
	item := &SymbolsItem{}

	err := db.DB.QueryRow("select symbol, currency from symbols where symbol = ?", *symbol).
		Scan(&item.Symbol, &item.Currency)
	if err != nil {
		return item, err
	}

	return item, nil
}
