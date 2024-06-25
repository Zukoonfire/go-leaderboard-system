package database 

import
(
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

var DB *sql.DB

func InitPostgres(){
	var err error
	connStr:="user=prashant password=password dbname=leaderboard sslmode=disable"
	DB,err=sql.Open("postgres",connStr)
	if err!=nil{
		log.Fatal(err)
	}

	err=DB.Ping()
	if err!=nil{
		log.Fatal(err)
	}
	log.Println("Postgres connected")
}