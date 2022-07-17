package db

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go-pinterest/model"
	"strconv"
	"time"
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
func GetDataId() (uint64, error) {
	val, err := Rbd.Get(context.Background(), ImageCountKey).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	total, err := strconv.ParseUint(val, 10, 64)

	if err != nil {
		return 0, err
	}
	return total, nil
}

func SetDataId(dataId uint64) error {
	_, err := Rbd.SetXX(context.Background(), ImageCountKey, strconv.FormatUint(dataId, 10), redis.KeepTTL).Result()
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
	//val, err := Rbd.HGetAll(context.Background(), ImageInfoHash(Id)).Result()
	//if err != nil {
	//	if err == redis.Nil {
	//		return nil, nil
	//	}
	//	return nil, err
	//}
	info := &model.ImageInfo{}
	//info.Owner = common.HexToAddress(val["Owner"])
	//info.Image = val["Image"]
	//info.TokenId, err = strconv.ParseInt(val["TokenId"], 10, 64)
	//if err != nil {
	//	return nil, err
	//}
	//info.DataId, err = strconv.ParseInt(val["DataId"], 10, 64)
	//if err != nil {
	//	return nil, err
	//}
	//info.MinedBlock, err = strconv.ParseInt(val["MinedBlock"], 10, 64)
	//if err != nil {
	//	return nil, err
	//}
	//info.SellStatus = false
	//if status, exit := val["SellStatus"]; exit {
	//	info.SellStatus, err = strconv.ParseBool(status)
	//	if err != nil {
	//		return nil, errors.New(fmt.Sprintf("Can not parse sell status with dataId :%d  , %s", dataId, status))
	//	}
	//}
	//
	//check := false
	//if price, exit := val["Price"]; exit {
	//	info.Price, check = new(big.Int).SetString(price, 10)
	//	if !check {
	//		return nil, errors.New(fmt.Sprintf("Can not parse price with dataId :%d  , %s", dataId, price))
	//	}
	//}
	//info.PriceType = val["PriceType"]
	return info, nil
}

func SetImageInfo(info model.ImageInfo) (bool, error) {
	//values := make([]interface{}, 16)
	//values[0] = "Owner"
	//values[1] = info.Owner.Hex()
	//values[2] = "Image"
	//values[3] = info.Image
	//values[4] = "TokenId"
	//values[5] = strconv.FormatInt(info.TokenId, 10)
	//values[6] = "DataId"
	//values[7] = strconv.FormatInt(info.DataId, 10)
	//values[8] = "MinedBlock"
	//values[9] = strconv.FormatInt(info.MinedBlock, 10)
	//values[10] = "SellStatus"
	//values[11] = strconv.FormatBool(info.SellStatus)
	//values[12] = "Price"
	//values[13] = info.Price.String()
	//values[14] = "PriceType"
	//values[15] = info.PriceType
	//check, err := Rbd.HMSet(context.Background(), ImageInfoHash(info.DataId), values...).Result()
	//if err != nil {
	//	return check, err
	//}
	//_, err = Rbd.ZAdd(context.Background(), NFTByAddressZset(info.Owner), &redis.Z{Member: info.DataId, Score: float64(time.Now().UnixNano())}).Result()
	//if err != nil {
	//	return false, err
	//}
	return true, nil
}

func AddImageToCategory(info model.ImageInfo, category string) (int64, error) {
	timestamp := time.Now().UnixNano()
	scoreText := strconv.FormatInt(info.CreatedTime, 10) + "." + strconv.FormatInt(timestamp, 10)
	score, err := strconv.ParseFloat(scoreText, 64)
	if err != nil {
		return 0, err
	}
	return Rbd.ZAdd(context.Background(), NFTByCategoryZset(category), &redis.Z{Member: info.Id, Score: score}).Result()
}

func GetAllCategory() ([]string, error) {
	val, err := Rbd.ZRange(context.Background(), CategoryListZset(), 0, -1).Result()
	if err != nil {
		if err == redis.Nil {
			return []string{}, nil
		}
		return []string{}, err
	}
	return val, nil
}
