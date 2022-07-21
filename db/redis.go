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
	info.Id = Id
	info.Title = val[TitleField.String()]
	info.Image = val[ImageField.String()]
	info.SourceId = val[SourceIdField.String()]
	info.Link = val[LinkField.String()]
	info.OwnerName = val[OwnerNameField.String()]
	info.OwnerUrl = val[OwnerUrlField.String()]
	info.BoardName = val[BoardNameField.String()]
	info.BoardUrl = val[BoardUrlField.String()]
	info.Images = make([]model.ImageSize, 0)
	json.Unmarshal([]byte(val[ImagesField.String()]), &info.Images)
	info.CreatedTime, _ = strconv.ParseInt(val[CreatedTimeField.String()], 10, 64)
	info.CrawledTime, _ = strconv.ParseInt(val[CrawledTimeField.String()], 10, 64)
	info.Description = val[DescriptionField.String()]
	json.Unmarshal([]byte(val[KeyWordsField.String()]), &info.KeyWords)
	json.Unmarshal([]byte(val[AnnotationsField.String()]), &info.Annotations)
	json.Unmarshal([]byte(val[HashtagsField.String()]), &info.Hashtags)
	if info.Annotations == nil {
		info.Annotations = make([]string, 0)
	}
	if info.Hashtags == nil {
		info.Hashtags = make([]string, 0)
	}
	info.BoardDescription = val[BoardDescriptionField.String()]
	return info, nil
}

func SetImageInfo(info model.ImageInfo) (bool, error) {
	values := make([]interface{}, 28)
	values[0] = TitleField.String()
	values[1] = info.Title
	values[2] = ImageField.String()
	values[3] = info.Image
	values[4] = SourceIdField.String()
	values[5] = info.SourceId
	values[6] = LinkField.String()
	values[7] = info.Link
	values[8] = OwnerNameField.String()
	values[9] = info.OwnerName
	values[10] = OwnerUrlField.String()
	values[11] = info.OwnerUrl
	values[12] = BoardNameField.String()
	values[13] = info.BoardName
	values[14] = BoardUrlField.String()
	values[15] = info.BoardUrl
	values[16] = ImagesField.String()
	sizeDetails, _ := json.Marshal(info.Images)
	values[17] = string(sizeDetails)
	values[18] = CreatedTimeField.String()
	values[19] = strconv.FormatInt(info.CreatedTime, 10)
	values[20] = CrawledTimeField.String()
	values[21] = strconv.FormatInt(info.CrawledTime, 10)
	values[22] = DescriptionField.String()
	values[23] = info.Description
	values[24] = KeyWordsField.String()
	keyWords, _ := json.Marshal(info.KeyWords)
	values[25] = string(keyWords)
	values[26] = CategoryField.String()
	values[27] = info.Category

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

func UpdateImageInfo(imageId int64, title string, description string, ownerName, ownerUrl string, annotations, hashtags []string, boardDescription string) (bool, error) {
	values := make([]interface{}, 14)
	values[0] = TitleField.String()
	values[1] = title
	values[2] = DescriptionField.String()
	values[3] = description
	values[4] = OwnerNameField.String()
	values[5] = ownerName
	values[6] = OwnerUrlField.String()
	values[7] = ownerUrl
	values[8] = AnnotationsField.String()
	annotationsByte, _ := json.Marshal(annotations)
	values[9] = string(annotationsByte)
	values[10] = HashtagsField.String()
	hashtagsByte, _ := json.Marshal(hashtags)
	values[11] = string(hashtagsByte)
	values[12] = BoardDescriptionField.String()
	values[13] = boardDescription
	check, err := Rbd.HMSet(context.Background(), ImageInfoHash(imageId), values...).Result()
	if err != nil {
		return check, err
	}
	return true, nil
}

func UpdateCategoryImageInfo(imageId int64, category string) error {
	_, err := Rbd.HSet(context.Background(), ImageInfoHash(imageId), CategoryField.String(), category).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}
	return nil
}

func UpdateKeywordImageInfo(imageId int64, keywords ...string) error {
	oldKeywords, err := Rbd.HGet(context.Background(), ImageInfoHash(imageId), KeyWordsField.String()).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}
	var exitKeywords []string
	err = json.Unmarshal([]byte(oldKeywords), &exitKeywords)
	if err != nil {
		return err
	}
	mapOldKeys := make(map[string]bool)
	for _, value := range exitKeywords {
		mapOldKeys[value] = true
	}
	needUpdate := false
	for _, value := range keywords {
		if !mapOldKeys[value] {
			exitKeywords = append(exitKeywords, value)
			mapOldKeys[value] = true
			needUpdate = true
		}
	}
	if !needUpdate {
		return nil
	}
	values := make([]interface{}, 2)
	values[0] = KeyWordsField.String()
	newKeywords, _ := json.Marshal(exitKeywords)
	values[1] = string(newKeywords)
	_, err = Rbd.HMSet(context.Background(), ImageInfoHash(imageId), values...).Result()
	if err != nil {
		return err
	}
	return nil
}
func AddImageToCategory(info model.ImageInfo, category string) (int64, error) {
	member := strconv.FormatInt(info.Id, 10)
	score, err := Rbd.ZScore(context.Background(), ImageByCategoryZset(category), member).Result()
	if err != nil {
		if err != redis.Nil {
			return 0, err
		}
	}
	if score > 0 {
		return 0, nil
	}
	return Rbd.ZAdd(context.Background(), ImageByCategoryZset(category), &redis.Z{Member: member, Score: float64(info.CrawledTime)}).Result()
}

func AddImageToCategoryAndDepth(info model.ImageInfo, category, keyword string, depth int) (int64, error) {
	member := strconv.FormatInt(info.Id, 10)
	score, err := Rbd.ZScore(context.Background(), ImageByCategoryAndDepthZset(category, keyword, depth), member).Result()
	if err != nil {
		if err != redis.Nil {
			return 0, err
		}
	}
	if score > 0 {
		return 0, nil
	}
	return Rbd.ZAdd(context.Background(), ImageByCategoryAndDepthZset(category, keyword, depth), &redis.Z{Member: member, Score: float64(info.CrawledTime)}).Result()
}

func GetImageByCategory(category string, offset int64, length int64) (model.ListImageInfo, error) {
	images := make([]model.ImageInfo, 0)
	listImages := model.ListImageInfo{
		images, -1,
	}
	imageIds, err := Rbd.ZRevRangeByScoreWithScores(context.Background(), ImageByCategoryZset(category), &redis.ZRangeBy{Max: strconv.FormatInt(offset, 10), Min: "0", Offset: 0, Count: length}).Result()
	if err != nil {
		if err == redis.Nil {
			return listImages, nil
		}
		return listImages, err
	}
	for i := 0; i < len(imageIds); i++ {
		id, _ := strconv.ParseInt(imageIds[i].Member.(string), 10, 64)
		imageInfo, _ := GetImageInfo(id)
		imageInfo.Category = category
		listImages.Images = append(listImages.Images, *imageInfo)
		listImages.NextOffset = int64(imageIds[i].Score - 1)
	}
	if int64(len(imageIds)) < length {
		listImages.NextOffset = -1
	}
	return listImages, nil
}

func GetAllImageIdByCategory(category string) ([]string, error) {
	imageIds, _ := Rbd.ZRevRange(context.Background(), ImageByCategoryZset(category), 0, -1).Result()
	return imageIds, nil
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
func AddAKeywordToCategory(category, keyword string) error {
	_, err := Rbd.SAdd(context.Background(), AllKeywordInCategorySet(category), keyword).Result()
	if err != nil {
		return err
	}
	return nil
}
