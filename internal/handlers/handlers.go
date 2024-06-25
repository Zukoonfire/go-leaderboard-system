package handlers

import(
	"encoding/json"
	"leaderboard-system/internal/auth"
	"leaderboard-system/internal/database"
	"leaderboard-system/internal/leaderboard"
	"leaderboard-system/pkg/models"
	"net/http"
	"log"
	"strconv"
	"context"
	"github.com/go-redis/redis/v8"
)

var ctx=context.Background()

func RegisterHandler(w http.ResponseWriter,r *http.Request){
	var user models.User
	if err:=json.NewDecoder(r.Body).Decode(&user);err!=nil{
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	query:=`INSERT INTO users(username,email,password)VALUES($1,$2,$3)`
	_,err:=database.DB.Exec(query,user.Username,user.Email,user.Password)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message":"user registered sucessfully"})
}

func LoginHandler(w http.ResponseWriter,r *http.Request){
	var user models.User
	if err:=json.NewDecoder(r.Body).Decode(&user);err!=nil{
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}

	//Authenticating User
	var username string
	query:=`SELECT username FROM users WHERE email=$1 AND password=$2`

	err:=database.DB.QueryRow(query,user.Email,user.Password).Scan(&username)
	if err!=nil{
		http.Error(w,"Invalid email or password",http.StatusUnauthorized)
		return
	}
	//Generating JWT
	token,err:=auth.GenerateJWT(username)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"tooken":token})
}

func SubmitScoreHandler(w http.ResponseWriter,r *http.Request){
	if r.Method!=http.MethodPost{
		http.Error(w,"Invalid request method",http.StatusMethodNotAllowed)
		return
	}

	tokenString:=r.Header.Get("Authorization")
	if tokenString==""{
		http.Error(w,"Authorization header is missing",http.StatusUnauthorized)
		return
	}

	tokenString=tokenString[len("Bearer "):]
	claims,err:=auth.ValidateJWT(tokenString)
	if err!=nil{
		http.Error(w,"Invalid Token",http.StatusUnauthorized)
		return
	}

	log.Printf("Authenticated user: %s", claims.Username)

	var scoreReq models.ScoreRequest
	if err:=json.NewDecoder(r.Body).Decode(&scoreReq);err!=nil{
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
//Fetch user ID from the database
var userID int
query:=`SELECT id FROM users WHERE username=$1`
err=database.DB.QueryRow(query,claims.Username).Scan(&userID)
if err!=nil{
	log.Printf("Error fetching user ID for %s: %v", claims.Username, err)
	http.Error(w,"User not Found",http.StatusInternalServerError)
	return
}
	//Adding score to Redis 
	err =leaderboard.RDB.ZAdd(ctx,"leaderboard",&redis.Z{
		Score: float64(scoreReq.Score),
		Member:strconv.Itoa(userID),

	}).Err()
	if err!=nil{
		log.Printf("Error adding score to Redis: %v", err)
		http.Error(w,"Failed to add score to leaderboard",http.StatusInternalServerError)
		return
	}
//Saving Score to PostgreSQL
query=`INSERT INTO scores (user_id,score)VALUES($1,$2)`
_,err=database.DB.Exec(query,userID,scoreReq.Score)
if err!=nil{
	log.Printf("Error saving score to PostgreSQL: %v", err)
	http.Error(w,"Failed to save score",http.StatusInternalServerError)
	return
}
w.WriteHeader(http.StatusOK)
json.NewEncoder(w).Encode(map[string]string{"message":"Score submitted sucessfully"})
}

func FetchLeaderboardHandler(w http.ResponseWriter,r *http.Request){
	if r.Method!=http.MethodGet{
		http.Error(w,"Invalid request method",http.StatusMethodNotAllowed)
		return
	}
	entries,err:=leaderboard.RDB.ZRevRangeWithScores(ctx,"leaderboard",0,9).Result()
	if err!=nil{
		http.Error(w,"Failed to fetch leaderboard",http.StatusInternalServerError)
		return
	}

	var leaderboardEntries []models.LeaderboardEntry
	for i,entry:=range entries{
		userID,_:=strconv.Atoi(entry.Member.(string))
		leaderboardEntries=append(leaderboardEntries,models.LeaderboardEntry{
			UserID: userID,
			Score: int(entry.Score),
			Rank:i+1,
		})
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(leaderboardEntries)
}