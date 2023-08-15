package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	// "io"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type currency struct {
	UUID             string
	amount           float32
	converted_amount float32
	from_currency    string
	to_currency      string
	conversion_rate  float32
}

type initial_currency struct {
	amount           float32
	conversion_rate  float32
	from_currency    string
	to_currency      string
}

type data_exit struct{
	Valor_convertido float32 `json:"valorConvertido"`
	Simbolo_moeda    string  `json: "simboloMoeda"`
}

func dbConn() (db *sql.DB) {
    dbDriver := "mysql"
    dbUser := "apirestgo"
    dbPass := "t3MzFHDD0Fvd"
    dbName := "apirestgo"
	dbHost := "databak"
    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbHost+")/"+dbName)
    if err != nil {
        panic(err.Error())
    }
    return db
}

func insertDB(currency currency) bool{

	db := dbConn()

	insForm, err := db.Prepare("INSERT INTO currency_exchange_log (UUID, amount, converted_amount, from_currency, to_currency, conversion_rate) VALUES(?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	
	insForm.Exec(
		currency.UUID,
		currency.amount,
		currency.converted_amount,
		currency.from_currency,
		currency.to_currency,
		currency.conversion_rate,
	)

	defer db.Close()

	return true
}


func convertCurrency(currency initial_currency) (float32){

	converted_amount := currency.amount * currency.conversion_rate

	return converted_amount
}

func getRoot(w http.ResponseWriter, r *http.Request){

	var currency_symbol = ""

	q := r.URL.Path

	currencyToConvert := getCurrency(q)

	fmt.Printf("converting from %s to %s \n", currencyToConvert.from_currency, currencyToConvert.to_currency)

	converted_amount := convertCurrency(currencyToConvert)

	currencyConverted := currency{
		uuid.NewString(),
		currencyToConvert.amount,
		converted_amount,
		currencyToConvert.from_currency,
		currencyToConvert.to_currency,
		currencyToConvert.conversion_rate,
	}

	switch currencyToConvert.to_currency {
    case "USD":
        currency_symbol = "$"
    case "BRL":
        currency_symbol = "R$"
    case "EUR":
        currency_symbol = "\u20ac"
	case "BTC":
		currency_symbol = "\u20bf"
    }
	
	de := data_exit{
		converted_amount,
		currency_symbol,
	}

	insertDB(currencyConverted)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	jsonResp, err := json.Marshal(de)
	if err != nil {
		panic(err.Error())
	}

	w.Write(jsonResp)
}

func getCurrency(e string) (initial_currency){

	s := strings.Split(e, "/")

	if(s[1] != "exchange"){
		panic("Wrong URL")
	}

	value, err := strconv.ParseFloat(s[2], 32)
	if err != nil {
		panic(err.Error())
	}
	amount := float32(value)

	value1, err := strconv.ParseFloat(s[5], 32)
	if err != nil {
		panic(err.Error())
	}
	conversion_rate := float32(value1)

	c := initial_currency{
		amount,
		conversion_rate,
		s[3],
		s[4],
	}

	return c
}


func main(){

	fmt.Println("Currency conversion started")

	http.HandleFunc("/", getRoot)

	err := http.ListenAndServe(":7000", nil)

	if err != nil {
		panic(err.Error())
	}
	
}