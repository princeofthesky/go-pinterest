package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-pinterest/db"
	"strconv"
	"time"
)

var (
	maxHomePage = flag.Int("max", 200, "max items in home page ")
)

func main() {

	flag.Parse()
	db.Init()
	for true {
		images, err := db.Rbd.SRandMemberN(context.Background(), db.AllImageSet(), int64(*maxHomePage)).Result()
		if err != nil {
			fmt.Println("err", err)
		}
		for _, imageId := range images {
			id, _ := strconv.ParseInt(imageId, 10, 64)
			info, _ := db.GetImageInfo(id)
			db.Rbd.ZAdd(context.Background(), db.HomeImagesZset(), &redis.Z{Member: imageId, Score: float64(info.CrawledTime)})
		}
		fmt.Println("Updated home page with size ", len(images))
		time.Sleep(30 * time.Minute)
	}
}
