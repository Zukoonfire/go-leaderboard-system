package main

import
(
	"leaderboard-system/internal/database"
	"leaderboard-system/internal/leaderboard"
	"leaderboard-system/internal/handlers"
	"log"
	"net/http"
)

func main(){
	database.InitPostgres()
	leaderboard.InitRedis()

	http.HandleFunc("/register",handlers.RegisterHandler)
	http.HandleFunc("/login",handlers.LoginHandler)
	http.HandleFunc("/submit-score",handlers.SubmitScoreHandler)
	http.HandleFunc("/leaderboard",handlers.FetchLeaderboardHandler)

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080",nil))
}

