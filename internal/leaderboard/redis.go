package leaderboard

import
(
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

var RDB *redis.Client
var ctx =context.Background()

func InitRedis(){
	RDB =redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password:"",
		DB:0,
	})

	_,err:=RDB.Ping(ctx).Result()
	if err!=nil{
		log.Fatal(err)
	}
	log.Println("Redis connected")
}