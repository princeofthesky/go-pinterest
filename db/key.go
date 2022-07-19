package db

import (
	"bytes"
	"strconv"
)

var ImageCountKey = "image_count"

var ConfigCrawler = "config"

var ImageMapIdHash = "image_check"

func ImageInfoHash(i int64) string {
	var b bytes.Buffer
	b.WriteString("i_")
	b.WriteString(strconv.FormatInt(i, 10))
	return b.String()
}

func ImageByCategoryZset(category string) string {
	var b bytes.Buffer
	b.WriteString("c_")
	b.WriteString(category)
	return b.String()
}

func CategorySet() string {
	var b bytes.Buffer
	b.WriteString("c")
	return b.String()
}

func HomeImagesZset() string {
	var b bytes.Buffer
	b.WriteString("home")
	return b.String()
}

func AllImageSet() string {
	var b bytes.Buffer
	b.WriteString("all_image")
	return b.String()
}
