package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/dhaliwal-h/go-postgres/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func createConnection() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error Loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to postgresql")
	return db
}

func CreateStock(w http.ResponseWriter, r *http.Request) {
	var stock models.Stock
	err := json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		fmt.Println("Unable to decode the request body")
		return
	}
	insertId, err := insertStock(stock)
	if err != nil {
		fmt.Println("Unable to insert stock")
		return
	}
	res := response{
		ID:      insertId,
		Message: "stock created successfully",
	}

	json.NewEncoder(w).Encode(res)
}
func GetStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	stockId := params["id"]
	id, err := strconv.Atoi(stockId)
	if err != nil {
		fmt.Println("Unable to parse stock id")
		return
	}
	stock, err := getStock(int64(id))
	if err != nil {
		fmt.Println("Unable to get Stock with id")
		return
	}
	json.NewEncoder(w).Encode(stock)

}
func GetAllStock(w http.ResponseWriter, r *http.Request) {
	stocks, err := getAllStocks()
	if err != nil {
		fmt.Println("Could not laod all stocks")
		return
	}
	json.NewEncoder(w).Encode(stocks)
}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	stockId := params["id"]
	id, err := strconv.Atoi(stockId)
	if err != nil {
		fmt.Println("Unable to parse stock id")
		return
	}
	stock, err := getStock(int64(id))
	if err != nil {
		fmt.Println("Unable to get Stock with id")
		return
	}

	var newStock models.Stock
	if err := json.NewDecoder(r.Body).Decode(&newStock); err != nil {
		fmt.Println("Unable to parse request body into a stock struct")
		return
	}

	if newStock.Company != "" {
		stock.Company = newStock.Company
	}
	if newStock.Name != "" {
		stock.Name = newStock.Name
	}
	if newStock.Price > 0 {
		stock.Price = newStock.Price
	}

	err = updateStock(stock)
	if err != nil {
		fmt.Println("Unable to udpate the give stock")
		return
	}

	json.NewEncoder(w).Encode(stock)
}
func DeleteStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	stockId := params["id"]
	id, err := strconv.Atoi(stockId)
	if err != nil {
		fmt.Println("Unable to parse stock id")
		return
	}
	stock, err := deleteStock(int64(id))
	if err != nil {
		fmt.Println("Unable to get Stock with id")
		return
	}
	json.NewEncoder(w).Encode(stock)

}

func insertStock(s models.Stock) (int64, error) {
	db := createConnection()
	defer db.Close()
	sqlStatement := `INSERT INTO stocks(name, price, company) VALUES($1, $2,$3) RETURNING stockid`
	var id int64
	err := db.QueryRow(sqlStatement, s.Name, s.Price, s.Company).Scan(&id)

	if err != nil {
		fmt.Println("Unalbe to run insert query")
		return id, err
	}
	return id, nil
}

func getStock(id int64) (models.Stock, error) {
	var stock models.Stock
	db := createConnection()
	defer db.Close()
	sqlStatement := `SELECT * FROM stocks WHERE stockid =$1`
	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
	if err != nil {
		fmt.Println(err)
		return stock, err
	}
	return stock, nil
}

func getAllStocks() ([]models.Stock, error) {
	var allStocks []models.Stock
	db := createConnection()
	defer db.Close()
	sqlStatement := `SELECT * FROM stocks`
	rows, err := db.Query(sqlStatement)

	if err != nil {
		fmt.Println(err)
	}

	defer rows.Close()

	for rows.Next() {
		var stock models.Stock
		err := rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
		if err != nil {
			fmt.Println(err)
		}
		allStocks = append(allStocks, stock)
	}
	rows.Scan(&allStocks)
	fmt.Printf("%v", allStocks)
	return allStocks, nil
}

func updateStock(s models.Stock) error {
	db := createConnection()
	defer db.Close()

	sql := `UPDATE stocks SET name = $1, price=$2, company=$3 WHERE stockid=$4`
	res, err := db.Exec(sql, s.Name, s.Price, s.Company, s.StockID)
	if err != nil {
		fmt.Println(err)
	}
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Rows affected %v", rowsAffected)
	return nil
}

func deleteStock(id int64) (models.Stock, error) {
	var stock models.Stock
	db := createConnection()
	defer db.Close()

	sql := `DELETE FROM stocks WHERE stockid=$1`
	res, err := db.Exec(sql, id)
	if err != nil {
		fmt.Println(err)
	}
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Rows affected %v", rowsAffected)
	return stock, nil
}
