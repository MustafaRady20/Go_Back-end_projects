package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-postgres/models"
	"log"
	"net/http"
	"os"
	"strconv"
	_"github.com/lib/pq"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func createConnection() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfuly connected to postgres")
	return db
}

func CreateStock(w http.ResponseWriter, r *http.Request) {
	var stock models.Stock

	err := json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		log.Fatalf("unable to decode request body . %v", err)
	}

	stockId := insertStock(stock)
	json.NewEncoder(w).Encode(stockId)
}



func GetStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("unable to convert string into int , %v", err)
	}

	stock, err := getStock(int64(id))
	if err != nil {
		log.Fatalf("unable to get stock. %v", err)
	}

	json.NewEncoder(w).Encode(stock)

}
func GetAllStocks(w http.ResponseWriter, r *http.Request) {

	stocks, err := getAllStocks()
	if err != nil {
		log.Fatalf("unable to get stocks. %v", err)
	}

	json.NewEncoder(w).Encode(stocks)

}
func UpdateStock(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("unable to convert string into int, %v", err)
	}

	var stock models.Stock

	err = json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		log.Fatalf("unable to decode the request %v", err)
	}

	upadtedRows := updateStock(int64(id), stock)

	msg := fmt.Sprintf("stock updated sccessfuly,Total rows affected %v", upadtedRows)

	json.NewEncoder(w).Encode(msg)
}
func DeleteStocks(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("unable to convert string into int %v", err)
	}

	deletedRows := deletStock(int64(id))

	msg := fmt.Sprintf("stock deleted successfuly , total rows affected %v", deletedRows)

	json.NewEncoder(w).Encode(msg)

}

func insertStock(stock models.Stock) int64 {
	db := createConnection()

	defer db.Close()

	sqlStatement := `INSERT INTO stokes(name,price,company) VALUES($1,$2,$3) RETURNING stokesid`
	var id int64
	err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company).Scan(&id)

	if err != nil {
		log.Fatalf("unable to execute querey %v", err)
	}

	fmt.Printf("Inserted a single row %v", id)

	return id
}

func getStock(id int64) (models.Stock, error) {
	db := createConnection()

	defer db.Close()

	sqlStatement := `SELECT * FROM stocks WHERE stocksid = $1`
	var stock models.Stock
	rows := db.QueryRow(sqlStatement, id)
	err := rows.Scan(&stock.Name, &stock.Price, &stock.Company)

	switch err {
	case sql.ErrNoRows:
		fmt.Printf("No Rows were returned...!")
		return stock, nil
	case nil:
		return stock, nil
	default:
		log.Fatalf("unable to scan the row.! %v", err)
	}

	return stock, err

}

func getAllStocks() ([]models.Stock, error) {

	db := createConnection()
	defer db.Close()

	sqlStatement := `SELECT * FROM stokes`
	var stocks []models.Stock

	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("unable to execute the querey")
	}

	defer rows.Close()

	for rows.Next() {
		var stock models.Stock
		err := rows.Scan(&stock.ID,&stock.Name, &stock.Price, &stock.Company)
		if err != nil {
			log.Fatalf("unable to scan the rows,%v", err)
		}

		stocks = append(stocks, stock)
	}

	return stocks, err
}

func deletStock(id int64) int64 {
	db := createConnection()
	defer db.Close()

	sqlStatment := `DELETE FROM stocks WHRER stocksid=$1`

	res, err := db.Exec(sqlStatment, id)
	if err != nil {
		log.Fatalf("unable to execute the querey %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf(" Error while checking for the affected rows %v", rowsAffected)
	}

	return rowsAffected

}

func updateStock(id int64, stock models.Stock) int64 {
	db := createConnection()
	defer db.Close()

	sqlStatment := `UPDATE stocks SET name=$2,price=$3,company=$4 WHRER stocksid=$1`

	res, err := db.Exec(sqlStatment, id, stock.Name, stock.Price, stock.Company)
	if err != nil {
		log.Fatalf("unable to execute the querey %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf(" Error while checking for the affected rows %v", rowsAffected)
	}

	return rowsAffected
}
