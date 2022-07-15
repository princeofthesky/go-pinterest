package db

import (
	"bytes"
	"strconv"
)

var ImageCountKey = "image_count"

var HomeImageSet = "home"

var ConfigCrawler = "config"

var ImageMapIdHash = "image_check"

func ImageInfoHash(i int64) string {
	var b bytes.Buffer
	b.WriteString("img_")
	b.WriteString(strconv.FormatInt(i, 10))
	return b.String()
}

func NFTByCategoryZset(category string) string {
	var b bytes.Buffer
	b.WriteString("c_")
	b.WriteString(category)
	return b.String()
}

func CategoryListZset() string {
	var b bytes.Buffer
	b.WriteString("c")
	return b.String()
}

func AllImageSet() string {
	var b bytes.Buffer
	b.WriteString("nft_sell")
	return b.String()
}
