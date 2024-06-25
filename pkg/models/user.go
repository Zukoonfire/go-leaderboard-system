package models 

type User struct{
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
}
type ScoreRequest struct {
    Score int `json:"score"`
}

type LeaderboardEntry struct {
    UserID int    `json:"user_id"`
    Score  int    `json:"score"`
    Rank   int    `json:"rank"`
}