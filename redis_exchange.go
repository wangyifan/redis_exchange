package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"time"
)

var ctx = context.Background()

func test() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	id := 103
	tableSell := fmt.Sprintf("redis_test_sell_%d", id)
	tableBuy := fmt.Sprintf("redis_test_buy_%d", id)
	loops := 10000
	scoreMax := 1000
	fmt.Println(tableSell)
	fmt.Println(tableBuy)
	fmt.Println(fmt.Sprintf("loops: %d", loops))
	for i := 0; i < loops; i++ {
		score := rand.Intn(scoreMax)
		err := rdb.ZAdd(ctx, tableSell, &redis.Z{
			Score:  float64(score),
			Member: fmt.Sprintf("order %d with score %d", i, score)}).Err()
		if err != nil {
			fmt.Printf("err: %d:=%s", i, err)
		}
	}

	for i := 0; i < loops; i++ {
		score := rand.Intn(scoreMax)
		err := rdb.ZAdd(ctx, tableBuy, &redis.Z{
			Score:  float64(score),
			Member: fmt.Sprintf("order %d with score %d", i, score)}).Err()
		if err != nil {
			fmt.Printf("err: %d:=%s", i, err)
		}
	}

	dt := time.Now()
	fmt.Println(dt.Format("2006-02-01 15:04:05"))
	for i := 0; i < loops; i++ {
		sellVals, err := rdb.ZRangeWithScores(ctx, tableSell, 0, 20).Result()
		if err == nil && i == 0 {
			for index := len(sellVals) - 1; index >= 0; index-- {
				v := sellVals[index]
				fmt.Printf("%d) %s\n", index, v.Member)
			}
		} else {
			fmt.Errorf("err: %s", err)
		}

		// insert new item to sell
		score := rand.Intn(scoreMax)
		rdb.ZAdd(ctx, tableSell, &redis.Z{
			Score:  float64(score),
			Member: fmt.Sprintf("order %d with score %d", i, score)}).Err()

		if i == 0 {
			fmt.Println("\n---------------------\n")
		}

		buyVals, err := rdb.ZRangeWithScores(ctx, tableBuy, 0, 20).Result()
		if err == nil && i == 0 {
			for index, v := range buyVals {
				fmt.Printf("%d) %s\n", index, v.Member)
			}
		} else {
			fmt.Errorf("err: %s", err)
		}

		// insert new item to buy
		score = rand.Intn(scoreMax)
		rdb.ZAdd(ctx, tableBuy, &redis.Z{
			Score:  float64(score),
			Member: fmt.Sprintf("order %d with score %d", i, score)}).Err()

	}
	dt = time.Now()
	fmt.Println(dt.Format("2006-02-01 15:04:05"))
}

func main() {
	test()
}
