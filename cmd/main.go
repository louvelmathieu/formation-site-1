package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

type NewTable struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Value string `json:"value" gorm:"type:text"`
}

func Connect() (*gorm.DB, error) {
	var c = mysql.Open(os.Getenv("RDS_USERNAME") + ":" + os.Getenv("RDS_PASSWORD") + "@tcp(" + os.Getenv("RDS_HOSTNAME") + ":" + os.Getenv("RDS_PORT") + ")/" + os.Getenv("RDS_DB_NAME") + "?charset=utf8&parseTime=True&loc=UTC")

	db, err := gorm.Open(c, &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	db = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8 auto_increment=1")
	return db.Session(&gorm.Session{}), nil
}

func init() {

	if os.Getenv("ENV") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

func main() {
	var router = mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "application/json; charset=utf-8")
		w.Write([]byte(`{"result":"ok"}`))
	}).Methods("GET")

	router.HandleFunc("/test_db", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "application/json; charset=utf-8")
		db, _ := Connect()
		var t NewTable
		err := db.Unscoped().Last(&t).Error
		if err != nil {
			fmt.Println(err)
		}
		resp, _ := json.Marshal(t)
		w.Write(resp)
	}).Methods("GET")

	e := http.ListenAndServe(":5000", router)
	if e != nil {
		log.Fatal("ListenAndServe: ", e)
	}
}
