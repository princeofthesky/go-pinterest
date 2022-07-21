package db

import (
	"bytes"
	"strconv"
)

type ImageFieldData int

const (
	ImageField ImageFieldData = iota
	IdField
	TitleField
	DescriptionField
	SourceIdField
	LinkField
	OwnerNameField
	OwnerUrlField
	BoardNameField
	BoardUrlField
	BoardDescriptionField
	ImagesField
	KeyWordsField
	AnnotationsField
	HashtagsField
	CreatedTimeField
	CrawledTimeField
	CategoryField
)

func (d ImageFieldData) String() string {
	switch d {
	case ImageField:
		return "Image"
	case IdField:
		return "Id"
	case TitleField:
		return "Title"
	case DescriptionField:
		return "Description"
	case SourceIdField:
		return "SourceId"
	case LinkField:
		return "Link"
	case OwnerNameField:
		return "OwnerName"
	case OwnerUrlField:
		return "OwnerUrl"
	case BoardNameField:
		return "BoardName"
	case BoardUrlField:
		return "BoardUrl"
	case BoardDescriptionField:
		return "BoardDescription"
	case ImagesField:
		return "Images"
	case KeyWordsField:
		return "KeyWords"
	case AnnotationsField:
		return "Annotations"
	case HashtagsField:
		return "Hashtags"
	case CreatedTimeField:
		return "CreatedTime"
	case CrawledTimeField:
		return "CrawledTime"
	case CategoryField:
		return "Category"
	default: break
	}
	return "unknown"
}

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

func ImageByCategoryAndDepthZset(category string, keyword string, depth int) string {
	var b bytes.Buffer
	b.WriteString("c_")
	b.WriteString(category)
	b.WriteString("_k_")
	b.WriteString(keyword)
	b.WriteString("_d_")
	b.WriteString(strconv.Itoa(depth))
	return b.String()
}


func AllKeywordInCategorySet(category string) string {
	var b bytes.Buffer
	b.WriteString("c_")
	b.WriteString(category)
	b.WriteString("_k")
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
