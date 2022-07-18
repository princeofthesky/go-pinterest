package db

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"go-pinterest/model"
	"strconv"
)

var (
	Rbd *redis.Client
)

func Init() {
	Rbd = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
		PoolSize: 10,
	})
}

func Close() {
	Rbd.Close()
}
func GetMaxImageId() (int64, error) {
	val, err := Rbd.Get(context.Background(), ImageCountKey).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	total, err := strconv.ParseInt(val, 10, 64)

	if err != nil {
		return 0, err
	}
	return total, nil
}

func SetMaxImageId(dataId int64) error {
	_, err := Rbd.SetXX(context.Background(), ImageCountKey, strconv.FormatInt(dataId, 10), 0).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}
	return nil
}

func InitDataId() (bool, error) {
	exit, err := Rbd.Exists(context.Background(), ImageCountKey).Result()
	if err != nil {
		return false, err
	}
	if exit == 1 {
		return true, nil
	}
	val, err := Rbd.SetNX(context.Background(), ImageCountKey, "0", 0).Result()
	if err != nil {
		return false, err
	}
	return val, nil
}

func SaveNFTSupply(supply int64) error {
	supplyText := strconv.FormatInt(supply, 10)
	_, err := Rbd.SetXX(context.Background(), ImageCountKey, supplyText, 0).Result()
	if err != nil {
		return err
	}
	return nil
}
func GetConfigCrawler() (map[string]string, error) {
	val, err := Rbd.HGetAll(context.Background(), ConfigCrawler).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return val, nil
}

func GetImageId(sourceLink string) (int64, error) {
	val, err := Rbd.HGet(context.Background(), ImageMapIdHash, sourceLink).Result()
	if err != nil {
		if err == redis.Nil {
			return -1, nil
		}
		return -1, err
	}
	id, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func GetImageInfo(Id int64) (*model.ImageInfo, error) {
	val, err := Rbd.HGetAll(context.Background(), ImageInfoHash(Id)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	info := &model.ImageInfo{}
	info.Title = val["Title"]
	info.Image = val["Image"]
	info.SourceId = val["SourceId"]
	info.Link = val["Link"]
	info.OwnerName = val["OwnerName"]
	info.OwnerUrl = val["OwnerUrl"]
	info.BoardName = val["BoardName"]
	info.BoardUrl = val["BoardUrl"]
	info.Images = make([]model.ImageSize, 0)
	json.Unmarshal([]byte(val["Images"]), &info.Images)
	info.CreatedTime, _ = strconv.ParseInt(val["CreatedTime"], 10, 64)
	info.CrawledTime, _ = strconv.ParseInt(val["CrawledTime"], 10, 64)
	return info, nil
}

func SetImageInfo(info model.ImageInfo) (bool, error) {
	values := make([]interface{}, 22)
	values[0] = "Title"
	values[1] = info.Title
	values[2] = "Image"
	values[3] = info.Image
	values[4] = "SourceId"
	values[5] = info.SourceId
	values[6] = "Link"
	values[7] = info.Link
	values[8] = "OwnerName"
	values[9] = info.OwnerName
	values[10] = "OwnerUrl"
	values[11] = info.OwnerUrl
	values[12] = "BoardName"
	values[13] = info.BoardName
	values[14] = "BoardUrl"
	values[15] = info.BoardUrl
	values[16] = "Images"
	sizeDetails, _ := json.Marshal(info.Images)
	values[17] = string(sizeDetails)
	values[18] = "CreatedTime"
	values[19] = strconv.FormatInt(info.CreatedTime, 10)
	values[20] = "CrawledTime"
	values[21] = strconv.FormatInt(info.CrawledTime, 10)
	check, err := Rbd.HMSet(context.Background(), ImageInfoHash(info.Id), values...).Result()
	if err != nil {
		return check, err
	}
	_, err = Rbd.SAdd(context.Background(), AllImageSet(), strconv.FormatInt(info.Id, 10)).Result()
	if err != nil {
		return false, err
	}
	_, err = Rbd.HSet(context.Background(), ImageMapIdHash, info.Link, strconv.FormatInt(info.Id, 10)).Result()
	if err != nil {
		return false, err
	}
	return true, nil
}

func AddImageToCategory(info model.ImageInfo, category string) (int64, error) {
	member := strconv.FormatInt(info.Id, 10)
	score, err := Rbd.ZScore(context.Background(), NFTByCategoryZset(category), member).Result()
	if err != nil {
		if err != redis.Nil {
			return 0, err
		}
	}
	if score > 0 {
		return 0, nil
	}
	return Rbd.ZAdd(context.Background(), NFTByCategoryZset(category), &redis.Z{Member: member, Score: float64(info.CrawledTime)}).Result()
}

func GetImageByCategory(category string, offset int64, length int64) (model.ListImageInfo, error) {
	images := make([]model.ImageInfo, 0)
	listImages := model.ListImageInfo{
		images, -1,
	}
	imageIds, err := Rbd.ZRevRangeByScoreWithScores(context.Background(), NFTByCategoryZset(category), &redis.ZRangeBy{Max: strconv.FormatInt(offset, 10), Min: "0", Offset: 0, Count: length}).Result()
	if err != nil {
		if err == redis.Nil {
			return listImages, nil
		}
		return listImages, err
	}
	for i := 0; i < len(imageIds); i++ {
		id, _ := strconv.ParseInt(imageIds[i].Member.(string), 10, 64)
		imageInfo, _ := GetImageInfo(id)
		listImages.Images = append(listImages.Images, *imageInfo)
		listImages.NextOffset = int64(imageIds[i].Score - 1)
	}
	if int64(len(imageIds)) < length {
		listImages.NextOffset = -1
	}
	return listImages, nil
}

func GetHomeImages(offset int64, length int64) (model.ListImageInfo, error) {
	images := make([]model.ImageInfo, 0)
	listImages := model.ListImageInfo{
		images, -1,
	}
	imageIds, err := Rbd.ZRevRangeByScoreWithScores(context.Background(), HomeImagesZset(), &redis.ZRangeBy{Max: strconv.FormatInt(offset, 10), Min: "0", Offset: 0, Count: length}).Result()
	if err != nil {
		if err == redis.Nil {
			return listImages, nil
		}
		return listImages, err
	}
	for i := 0; i < len(imageIds); i++ {
		id, _ := strconv.ParseInt(imageIds[i].Member.(string), 10, 64)
		imageInfo, _ := GetImageInfo(id)
		listImages.Images = append(listImages.Images, *imageInfo)
		listImages.NextOffset = int64(imageIds[i].Score - 1)
	}
	if int64(len(imageIds)) < length {
		listImages.NextOffset = -1
	}
	return listImages, nil
}

func GetAllCategory() ([]string, error) {
	val, err := Rbd.SMembers(context.Background(), CategorySet()).Result()
	if err != nil {
		if err == redis.Nil {
			return []string{}, nil
		}
		return []string{}, err
	}
	return val, nil
}

func AddACategory(category string) error {
	_, err := Rbd.SAdd(context.Background(), CategorySet(), category).Result()
	if err != nil {
		return err
	}
	return nil
}
