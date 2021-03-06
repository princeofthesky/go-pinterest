package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"go-pinterest/db"
	"go-pinterest/model"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	SessionKey    = "TWc9PSY5UzBwT01YdnZsbG55d3phZGVHNzVFWC93WTg3Z0dOS0FkUUtrWkRBRFRFVE94WFB2K011TjE5aC9GZHFuTTFDOTNsMEJoOENPWU5EODI3TGQ0TTlqTGVGNjY1T3lpUWZqRUgrSG1rYmpNZlpEcFBHN1JRV1BGc2VoOXYwS2FWc1BwNitrdjVpdkltRitwbXVlMnFVbVY2KzgwMkFpclBEcTBHOC9nVkYzbTNoWStTeGxGR1V6dWI1MUR5ZnRSWC8vVjhVaXFSTWhkamJlVUsvdmFNZVJBN2dJRCs2TEFiM1hvMy9yTm5sUjVZNEZlOGtFb3k2dTl0VmliUEZzNmROV1k0ODd5ajR3andzcVM1OG9tSDBMSHFCdGJ4clZjK1VpM3h4Rm9jQU9PcnFoNU91a2F4c3JPWnVuaTM2TUxhVTZUcFFNZDcyNDFqdXlyR0g0WlJpSXJncEkrS3B6c3VCTlNvVnFTZ1BzcWdKUlJacFlCMnZ6NjdkSTJqQ0VQd0gyOFYvdWhEM3dZVS9HNlUvRDlMOUdRUk1uWHlEQUlFbi9WVWF6UHd0QjFHVGRoNWRMbnhwaTJCREZFMnhmelBLaktZaTlEWEJBcXI4dWllc0RwbFl2Y3FRalBWcmtIOXhpMVUxcUZiWGhjenBFOGhBOWF3MkVDMUVrNFlZa3R6eGNYb0VEb2YyeHdQZDdrM2N1SkpYTDVVNjV5bkVYZDdEeEtXTzZRcjZqZHpBSjFtVVZxaGdTbkRnUzhwN0FEbWJrZnUrR2JsSWxxYUw2UjFEb2xQR0x5SUpibXN0dUhlU3YvbzFFdlBZQVJvbTFiL0ZJVUxoWm4vVEZVTXRnbm1ZR04zRmpBdS82cENRNFdlOTNJU2t6VjlwRVBVcTZlZm5kMTNmdVlwYTJhQURrNW8ya05WeVlUTjRRMUhzUmpwTWV0c1EzZ3I4UXlQd1YxUnZKWWM5ZGhmcG8vbnUrRWNXOFNJMW9pamxLbmQ5dEduZTJRU09YYmpUZ0E1cmR6Y0k2azJMSUo5Zm1GVGczUTVZdzFHNmZBPT0mMVFjUHNDQUZaZkFRWjduSDl6Zy9hRGJqcEpzPQ=="
	Csrftoken     = "e7525248933e9ce4fe7cd9507627bd03"
	source        = flag.String("source", "/home/tamnb/projects/src/github.com/nguyenbatam/go-pinterest/sample_crawl_source.csv", "source category Pinterest")
	maxRelated    = flag.Int("max_related_pin", 50, "max items per query related")
	maxPageSearch = flag.Int("max_page_search", 4, "max items page  per search , 25 items per page")
)

func ReadSourcePinterestFile(source string) (map[string]string, error) {
	mapKeywordCategory := map[string]string{}
	// Open CSV file
	fileContent, err := os.Open(source)
	if err != nil {
		return mapKeywordCategory, err
	}
	defer fileContent.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(fileContent).ReadAll()
	if err != nil {
		return mapKeywordCategory, err
	}

	for i := 1; i < len(lines); i++ {
		categry := lines[i][0]
		keywords := strings.Split(lines[i][1], "\n")
		for j := 0; j < len(keywords); j++ {
			keyword := strings.TrimSpace(keywords[j])
			if len(keyword) > 0 {
				mapKeywordCategory[keyword] = categry
			}
		}

	}
	return mapKeywordCategory, nil
}
func main() {
	flag.Parse()
	mapKeywordCategory, err := ReadSourcePinterestFile(*source)
	fmt.Println(os.Args)
	fmt.Println(err)
	db.Init()
	db.InitDataId()
	LoadConfigCrawler()
	IdCrawled := make(map[string]bool)
	IdsSearched := make(map[string]map[string]bool)
	for keyword, category := range mapKeywordCategory {
		depth:=1
		imageSearchinfos, bookmark, _ := GetPinFromKeyWordSearch(keyword)
		fmt.Println(keyword, "search  ", len(imageSearchinfos))
		if err != nil {
			fmt.Println(keyword)
			fmt.Println("err", err)
		}

		for i := 0; i < len(imageSearchinfos); i++ {
			imageInfo := imageSearchinfos[i]
			id := imageInfo.SourceId
			//fmt.Println("find new image search id =", id)
			if _, exits := IdsSearched[id]; exits {
				IdsSearched[id][keyword] = true
				continue
			}
			IdsSearched[id] = make(map[string]bool)
			IdsSearched[id][keyword] = true
			IdCrawled[id] = true
			if id, _ := db.GetImageId(imageInfo.Link); id > 0 {
				imageInfo.Id = id
				imageInfo.CrawledTime = time.Now().UnixNano()
				db.AddImageToCategory(imageInfo, category)
				db.UpdateKeywordImageInfo(imageInfo.Id, keyword)
				db.UpdateCategoryImageInfo(imageInfo.Id,category)
				db.AddImageToCategoryAndDepth(imageInfo,category,keyword,depth)
				continue
			}
			dataId, _ := db.GetMaxImageId()
			imageInfo.Id = dataId + 1
			imageInfo.CrawledTime = time.Now().UnixNano()
			imageInfo.KeyWords = make([]string, 1)
			imageInfo.KeyWords[0] = keyword
			imageInfo.Category=mapKeywordCategory[keyword]

			db.SetMaxImageId(dataId + 1)
			db.SetImageInfo(imageInfo)
			db.AddImageToCategory(imageInfo, category)
			db.AddImageToCategoryAndDepth(imageInfo,category,keyword,depth)
		}
		for j := 0; j < *maxPageSearch-1; j++ {
			imageNextSearchInfos, nextBookmark, err := GetPinFromNextPageSearch(keyword, bookmark)
			bookmark = nextBookmark
			if err != nil {
				fmt.Println(keyword, " page ", j, "bookmark", bookmark)
				fmt.Println("err", err)
			}
			fmt.Println(keyword, "search next ", len(imageNextSearchInfos))
			for i := 0; i < len(imageNextSearchInfos); i++ {
				imageInfo := imageNextSearchInfos[i]
				id := imageInfo.SourceId
				//fmt.Println("find new image search next page id =", id)
				if _, exits := IdsSearched[id]; exits {
					IdsSearched[id][keyword] = true
					continue
				}
				IdsSearched[id] = make(map[string]bool)
				IdsSearched[id][keyword] = true
				IdCrawled[id] = true
				if id, _ := db.GetImageId(imageInfo.Link); id > 0 {
					imageInfo.Id = id
					imageInfo.CrawledTime = time.Now().UnixNano()
					db.AddImageToCategory(imageInfo, category)
					db.UpdateKeywordImageInfo(imageInfo.Id, keyword)
					db.UpdateCategoryImageInfo(imageInfo.Id,category)
					db.AddImageToCategoryAndDepth(imageInfo,category,keyword,depth)
					continue
				}
				dataId, _ := db.GetMaxImageId()
				imageInfo.Id = dataId + 1
				imageInfo.CrawledTime = time.Now().UnixNano()
				imageInfo.KeyWords = make([]string, 1)
				imageInfo.KeyWords[0] = keyword
				imageInfo.Category=mapKeywordCategory[keyword]

				db.SetMaxImageId(dataId + 1)
				db.SetImageInfo(imageInfo)
				db.AddImageToCategory(imageInfo, category)
				db.AddImageToCategoryAndDepth(imageInfo,category,keyword,depth)
			}
			fmt.Println("new", len(IdCrawled))
		}
		db.AddACategory(category)
		db.AddAKeywordToCategory(category,keyword)
	}
	for pinId, keywords := range IdsSearched {
		depth:=2
		relatedImages, err := GetPinFromRelatedPin(pinId, *maxRelated)
		if err != nil {
			fmt.Println(pinId)
			fmt.Println("err", err)
		}
		for i := 0; i < len(relatedImages); i++ {
			IdCrawled[relatedImages[i].SourceId] = true
			imageInfo := relatedImages[i]
			if id, _ := db.GetImageId(imageInfo.Link); id > 0 {
				imageInfo.Id = id
				imageInfo.CrawledTime = time.Now().UnixNano()
				categories := make(map[string]interface{})
				for keyword, _ := range keywords {
					category:=mapKeywordCategory[keyword]
					categories[category] = true
					db.UpdateKeywordImageInfo(imageInfo.Id, keyword)
					db.AddImageToCategoryAndDepth(imageInfo,category,keyword,depth)
				}
				for category, _ := range categories {
					db.AddImageToCategory(imageInfo, category)
				}
				continue
			}
			dataId, _ := db.GetMaxImageId()
			imageInfo.Id = dataId + 1
			imageInfo.CrawledTime = time.Now().UnixNano()
			imageInfo.KeyWords = make([]string, 0)
			categories := make(map[string]interface{})
			for keyword, _ := range keywords {
				category:=mapKeywordCategory[keyword]
				categories[category] = true
				imageInfo.KeyWords = append(imageInfo.KeyWords, keyword)
				imageInfo.Category=category
			}
			db.SetMaxImageId(dataId + 1)
			db.SetImageInfo(imageInfo)
			for keyword, _ := range keywords {
				category:=mapKeywordCategory[keyword]
				db.AddImageToCategoryAndDepth(imageInfo, category, keyword, depth)
			}
			for category, _ := range categories {
				db.AddImageToCategory(imageInfo, category)
			}
		}
	}
	fmt.Println("total", len(IdCrawled))
}

func RefreshSessionKey(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")
	req.Header.Set("authority", "www.pinterest.com")
	req.Header.Set("path", "/manifest.json")
	req.Header.Set("scheme", "https")
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("accept-language", "vi-VN,vi;q=0.9,en-US;q=0.8,en;q=0.7,fr-FR;q=0.6,fr;q=0.5")
	req.Header.Set("referer", "https://www.pinterest.com/")
	req.Header.Set("sec-ch-ua", "\".Not/A)Brand\";v=\"99\", \"Google Chrome\";v=\"103\", \"Chromium\";v=\"103\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "Linux")
	req.Header.Set("sec-fetch-des", "manifest")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	fmt.Println(resp.Header.Values("Set-Cookie")[0])
	fmt.Println(resp.Header.Values("Set-Cookie")[1])
	for i := 0; i < len(resp.Header.Values("Set-Cookie")); i++ {
		value := resp.Header.Values("Set-Cookie")[i]
		if strings.Contains(value, "csrftoken") {
			Csrftoken = strings.Split(value, ";")[0]
			Csrftoken = strings.Split(Csrftoken, "=")[1]
			fmt.Println(Csrftoken)
		}
		if strings.Contains(value, "_pinterest_sess") {
			SessionKey = strings.Split(value, ";")[0]
			SessionKey = strings.Split(SessionKey, "sess=")[1]
			fmt.Println(SessionKey)
		}
	}
	return nil
}

func LoadConfigCrawler() error {
	config, err := db.GetConfigCrawler()
	if err != nil {
		return err
	}
	if len(config["pinterest_csrf"]) > 0 {
		Csrftoken = config["pinterest_csrf"]
	}
	if len(config["pinterest_session"]) > 0 {
		SessionKey = config["pinterest_session"]
	}
	return nil
}
func RefreshCsrfToken(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	for i := 0; i < len(resp.Cookies()); i++ {
		if strings.Contains(resp.Cookies()[i].Name, "csrftoken") {
			Csrftoken = resp.Cookies()[i].Value
		}
	}
	return nil
}
func GetPinFromKeyWordSearch(keyword string) ([]model.ImageInfo, string, error) {
	dataQuery := "{\"options\":{\"article\":null,\"applied_filters\":null,\"appliedProductFilters\":\"---\",\"auto_correction_disabled\":false,\"corpus\":null,\"customized_rerank_type\":null,\"filters\":null,\"query\":\"" +
		keyword +
		"\",\"query_pin_sigs\":null,\"redux_normalize_feed\":true,\"rs\":\"rs\",\"scope\":\"pins\",\"source_id\":null,\"no_fetch_context_on_resource\":false},\"context\":{}}"
	uri := "https://www.pinterest.com/resource/BaseSearchResource/get/?source_url=" + url.QueryEscape("/search/pins/?q="+url.QueryEscape(keyword)) + "&data=" + url.QueryEscape(dataQuery)
	req, err := http.NewRequest("GET", uri, nil)
	req.AddCookie(&http.Cookie{Name: "csrftoken", Value: Csrftoken})
	req.AddCookie(&http.Cookie{Name: "_pinterest_sess", Value: SessionKey})
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, "", fmt.Errorf("failed to fetch bodyBytes: %d %s", resp.StatusCode, resp.Status)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	response := map[string]interface{}{}
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, "", err
	}
	resourceResponse := response["resource_response"].(map[string]interface{})
	data := resourceResponse["data"].(map[string]interface{})
	results := data["results"].([]interface{})

	imageInfos := make([]model.ImageInfo, 0)
	for i := 0; i < len(results); i++ {
		info := results[i].(map[string]interface{})
		if info["type"].(string) != "pin" {
			continue
		}
		images := info["images"].(map[string]interface{})
		pinner := info["pinner"].(map[string]interface{})
		imageInfo := model.ImageInfo{}
		imageInfo.SourceId = info["id"].(string)
		if _, err := strconv.ParseUint(imageInfo.SourceId, 10, 64); err != nil {
			continue
		}
		imageInfo.Link = "https://www.pinterest.com/pin/" + imageInfo.SourceId
		if info["title"] != nil && len(info["title"].(string)) > 0 {
			imageInfo.Title = info["title"].(string)
		} else if info["grid_title"] != nil && len(info["grid_title"].(string)) > 0 {
			imageInfo.Title = info["grid_title"].(string)
		} else if info["description"] != nil && len(info["description"].(string)) > 0 {
			imageInfo.Title = info["description"].(string)
		}
		imageInfo.Images = make([]model.ImageSize, 0)
		setUrlImage := make(map[int]bool)
		for _typeImage, _ := range images {
			sizeValue := images[_typeImage].(map[string]interface{})
			if _typeImage == "orig" {
				imageInfo.Image = sizeValue["url"].(string)
			}
			imageSize := model.ImageSize{}
			imageSize.Url = sizeValue["url"].(string)
			imageSize.Width = int(sizeValue["width"].(float64))
			imageSize.Height = int(sizeValue["height"].(float64))
			if setUrlImage[imageSize.Width*imageSize.Height] {
				continue
			}
			setUrlImage[imageSize.Width*imageSize.Height] = true
			imageInfo.Images = append(imageInfo.Images, imageSize)
		}
		imageInfo.Link = "https://www.pinterest.com/pin/" + imageInfo.SourceId
		createdTime := info["created_at"].(string) //Fri, 25 Sep 2020 16:51:58 +0000
		created, err := time.Parse(time.RFC1123, createdTime)
		if err != nil {
			fmt.Println("err when parser time", createdTime)
			continue
		}
		imageInfo.CreatedTime = created.UnixNano()
		imageInfo.OwnerName = pinner["full_name"].(string)
		imageInfo.OwnerUrl = "https://www.pinterest.com/" + pinner["username"].(string)
		if info["board"] != nil {
			board := info["board"].(map[string]interface{})
			imageInfo.BoardName = board["name"].(string)
			imageInfo.BoardUrl = "https://www.pinterest.com" + board["url"].(string)
		} else {
			imageInfo.BoardName = ""
			imageInfo.BoardUrl = ""
		}
		imageInfos = append(imageInfos, imageInfo)

	}
	resource := response["resource"].(map[string]interface{})
	options := resource["options"].(map[string]interface{})
	bookmarks := options["bookmarks"].([]interface{})
	bookmark := ""
	for _, value := range bookmarks {
		bookmark = value.(string)
		break
	}
	return imageInfos, bookmark, nil
}

func GetPinFromRelatedPin(pinID string, size int) ([]model.ImageInfo, error) {

	dataQuery := "{\"options\":{\"pin_id\":\"" +
		pinID +
		"\",\"context_pin_ids\":[],\"page_size\":" +
		strconv.Itoa(size) +
		",\"search_query\":\"\",\"source\":\"deep_linking\",\"top_level_source\":\"deep_linking\",\"top_level_source_depth\":1,\"is_pdp\":false,\"no_fetch_context_on_resource\":false},\"context\":{}}"
	query := url.QueryEscape("/pin/"+pinID+"/") + "&data=" + url.QueryEscape(dataQuery)
	uri := "https://www.pinterest.com/resource/RelatedModulesResource/get/?source_url=" + query
	req, err := http.NewRequest("GET", uri, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")
	req.AddCookie(&http.Cookie{Name: "csrftoken", Value: Csrftoken})
	req.AddCookie(&http.Cookie{Name: "_pinterest_sess", Value: SessionKey})

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch bodyBytes: %d %s", resp.StatusCode, resp.Status)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	response := map[string]interface{}{}
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, err
	}
	resourceResponse := response["resource_response"].(map[string]interface{})
	data := resourceResponse["data"].([]interface{})
	imageInfos := make([]model.ImageInfo, 0)
	for i := 0; i < len(data); i++ {
		info := data[i].(map[string]interface{})
		if info["type"].(string) != "pin" {
			continue
		}
		images := info["images"].(map[string]interface{})
		pinner := info["pinner"].(map[string]interface{})
		imageInfo := model.ImageInfo{}
		imageInfo.SourceId = info["id"].(string)
		if _, err := strconv.ParseUint(imageInfo.SourceId, 10, 64); err != nil {
			continue
		}
		imageInfo.Link = "https://www.pinterest.com/pin/" + imageInfo.SourceId
		if info["title"] != nil && len(info["title"].(string)) > 0 {
			imageInfo.Title = info["title"].(string)
		} else if info["grid_title"] != nil && len(info["grid_title"].(string)) > 0 {
			imageInfo.Title = info["grid_title"].(string)
		} else if info["description"] != nil && len(info["description"].(string)) > 0 {
			imageInfo.Title = info["description"].(string)
		}
				imageInfo.Images = make([]model.ImageSize, 0)
		setUrlImage := make(map[int]bool)
		for _typeImage, _ := range images {
			sizeValue := images[_typeImage].(map[string]interface{})
			if _typeImage == "orig" {
				imageInfo.Image = sizeValue["url"].(string)
			}
			imageSize := model.ImageSize{}
			imageSize.Url = sizeValue["url"].(string)
			imageSize.Width = int(sizeValue["width"].(float64))
			imageSize.Height = int(sizeValue["height"].(float64))
			if setUrlImage[imageSize.Width*imageSize.Height] {
				continue
			}
			setUrlImage[imageSize.Width*imageSize.Height] = true
			imageInfo.Images = append(imageInfo.Images, imageSize)
		}
		createdTime := info["created_at"].(string) //Fri, 25 Sep 2020 16:51:58 +0000
		created, err := time.Parse(time.RFC1123, createdTime)
		if err != nil {
			fmt.Println("err when parser time", createdTime)
			continue
		}
		imageInfo.CreatedTime = created.UnixNano()
		imageInfo.OwnerName = pinner["full_name"].(string)
		imageInfo.OwnerUrl = "https://www.pinterest.com/" + pinner["username"].(string)
		if info["board"] != nil {
			board := info["board"].(map[string]interface{})
			imageInfo.BoardName = board["name"].(string)
			imageInfo.BoardUrl = "https://www.pinterest.com" + board["url"].(string)
		} else {
			imageInfo.BoardName = ""
			imageInfo.BoardUrl = ""
		}
		imageInfos = append(imageInfos, imageInfo)
	}
	return imageInfos, nil
}

func GetPinFromNextPageSearch(keyword string, bookmark string) ([]model.ImageInfo, string, error) {
	bodyPost := []byte("source_url=" + url.QueryEscape("/search/pins/?q="+url.QueryEscape(keyword)) + "&data=" +
		url.QueryEscape("{\"options\":{\"article\":null,\"applied_filters\":null,\"appliedProductFilters\":\"---\",\"auto_correction_disabled\":false,\"corpus\":null,\"customized_rerank_type\":null,\"filters\":null,\"query\":\""+
			keyword+
			"\",\"query_pin_sigs\":null,\"redux_normalize_feed\":true,\"rs\":\"typed\",\"scope\":\"pins\",\"source_id\":null,\"bookmarks\":[\""+
			bookmark+
			"\"],\"no_fetch_context_on_resource\":false},\"context\":{}}"))
	req, err := http.NewRequest("POST", "https://www.pinterest.com/resource/BaseSearchResource/get/", bytes.NewBuffer(bodyPost))
	req.AddCookie(&http.Cookie{Name: "csrftoken", Value: Csrftoken})
	req.AddCookie(&http.Cookie{Name: "_pinterest_sess", Value: SessionKey})
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")
	req.Header.Set("x-csrftoken", Csrftoken)
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, "", fmt.Errorf("failed to fetch bodyBytes: %d %s", resp.StatusCode, resp.Status)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	response := map[string]interface{}{}
	err = json.Unmarshal(bodyBytes, &response)

	if err != nil {
		//fmt.Println(err)
		return nil, "", err
	}
	resourceResponse := response["resource_response"].(map[string]interface{})
	data := resourceResponse["data"].(map[string]interface{})
	results := data["results"].([]interface{})
	imageInfos := make([]model.ImageInfo, 0)
	for i := 0; i < len(results); i++ {
		info := results[i].(map[string]interface{})
		if info["type"].(string) != "pin" {
			continue
		}
		images := info["images"].(map[string]interface{})
		pinner := info["pinner"].(map[string]interface{})
		imageInfo := model.ImageInfo{}
		imageInfo.SourceId = info["id"].(string)
		if _, err := strconv.ParseUint(imageInfo.SourceId, 10, 64); err != nil {
			continue
		}
		imageInfo.Link = "https://www.pinterest.com/pin/" + imageInfo.SourceId
		if info["title"] != nil && len(info["title"].(string)) > 0 {
			imageInfo.Title = info["title"].(string)
		} else if info["grid_title"] != nil && len(info["grid_title"].(string)) > 0 {
			imageInfo.Title = info["grid_title"].(string)
		} else if info["description"] != nil && len(info["description"].(string)) > 0 {
			imageInfo.Title = info["description"].(string)
		}
		imageInfo.Images = make([]model.ImageSize, 0)
		setUrlImage := make(map[int]bool)
		for _typeImage, _ := range images {
			sizeValue := images[_typeImage].(map[string]interface{})
			if _typeImage == "orig" {
				imageInfo.Image = sizeValue["url"].(string)
			}
			imageSize := model.ImageSize{}
			imageSize.Url = sizeValue["url"].(string)
			imageSize.Width = int(sizeValue["width"].(float64))
			imageSize.Height = int(sizeValue["height"].(float64))
			if setUrlImage[imageSize.Width*imageSize.Height] {
				continue
			}
			setUrlImage[imageSize.Width*imageSize.Height] = true
			imageInfo.Images = append(imageInfo.Images, imageSize)
		}
		createdTime := info["created_at"].(string) //Fri, 25 Sep 2020 16:51:58 +0000
		created, err := time.Parse(time.RFC1123, createdTime)
		if err != nil {
			fmt.Println("err when parser time", createdTime)
			continue
		}
		imageInfo.CreatedTime = created.UnixNano()
		imageInfo.OwnerName = pinner["full_name"].(string)
		imageInfo.OwnerUrl = "https://www.pinterest.com/" + pinner["username"].(string)
		if info["board"] != nil {
			board := info["board"].(map[string]interface{})
			imageInfo.BoardName = board["name"].(string)
			imageInfo.BoardUrl = "https://www.pinterest.com" + board["url"].(string)
		} else {
			imageInfo.BoardName = ""
			imageInfo.BoardUrl = ""
		}
		imageInfos = append(imageInfos, imageInfo)
	}

	resource := response["resource"].(map[string]interface{})
	options := resource["options"].(map[string]interface{})
	bookmarks := options["bookmarks"].([]interface{})
	nextBookmark := ""
	for _, value := range bookmarks {
		nextBookmark = value.(string)
		break
	}
	return imageInfos, nextBookmark, nil
}

//
//
//
//source_url=%2Fsearch%2Fpins%2F%3Fq%3Dharry%2Bpotter&
//	data=%7B%22options%22%3A%7B%22article%22%3Anull%2C%22applied_filters%22%3Anull%2C%22appliedProductFilters%22%3A%22---%22%2C%22auto_correction_disabled%22%3Afalse%2C%22corpus%22%3Anull%2C%22customized_rerank_type%22%3Anull%2C%22filters%22%3Anull%2C%22query%22%3A%22the+oscars+winners+portraits%22%2C%22query_pin_sigs%22%3Anull%2C%22redux_normalize_feed%22%3Atrue%2C%22rs%22%3A%22typed%22%2C%22scope%22%3A%22pins%22%2C%22source_id%22%3Anull%2C%22bookmarks%22%3A%5B%22Y2JVSG81V2sxcmNHRlpWM1J5VFVaU1YxWllhRlJXTVVreVZsZDRRMVV4U1hsVlZFWlhVbnBXTTFWNlNrZFdNa3BIWVVaYVdGSXlhRkZXUm1Rd1ZtMVJlRlZ1VWs1V1ZGWlBXVmh3UjFac1draE5WRkpWVFd0YWVWUnNhRXRYUmxsNlVXMUdWVlpzY0hsYVZscExWbFpXZEZKc1RsTmlXRkV4Vm10U1IxVXhUbkpPVlZwUVZsWmFXVmxzYUZOVU1YQllaRVYwYWxKc1NubFhhMVpoWWtaYVZWWnVhRmRpUmtwVVYxWmFTbVF3TVZWV2JGWk9VbXR3U1ZkWGVGWmxSVFYwVW10b2FGSnJTbGhWYlhoYVRXeFplVTFZWkZOaGVrWjVWR3hhVjFadFNsVlNibEpXWWtaS1dGVnFSbUZqVmxKeFZHeEdWbFpFUVRWYWExcFhVMWRLTmxWdGVGZE5XRUpLVm10amVHSXhiRmRUV0dob1RUSlNXVmxVUmt0VU1WSlhWbGhvYWxZd2NFbFpWVlUxWVVkRmVGZFljRmRTTTFKeVZqSXhWMk15VGtkV2JFcFlVakpvYUZadGRHRlpWMDVYVlc1S1ZtSnJOVzlVVm1RMFpVWldjMVZzVGxWaVJYQkpXWHBPYzFaV1dsZFRia1poVmpOTmVGVnNXbE5YVjBwR1QxWmtVMDFFUWpOV2FrbDRaREZrZEZac1pHbFRSa3BXVm10Vk1XRkdiRmhOVjNCc1lrWktlbGRyV21GVWJFcDFVVzVvVmxac1NraFdSRVpMVWpKS1JWZHNhRmRpUlhCRVYyeGFWbVZIVWtkVGJrWm9VbXhhYjFSV1duZFhiR1IwWkVWYVVGWnJTbE5WUmxGNFQwVXhObFJVVWs1bGJHdDVWRmN4U2sxR2NFVlhiV3hoWWxWcmQxUnJVbkpsVm14WVVtMXNXbFpGVmpOWFYzQnVUVEZyZVZKWWNFOVdSVEF3VjJ4U1FrMVdiRFpWVkU1aFlsWktjRmRyWkU5aVJteHhWbGh3V2sxc1NtOVVWbVJHWkRBeFNGTnRNVTVoYkhCdlZGWmtWazFzYTNwbFJUbFRWbTFSTkdaSFZteFBSRWt3VG1wRk1FNVVVVE5OVkUxNFRWZEplRTVFVm0xYWJVWnFXWHBWTkU5WFdtbFpWRTAwVFVkR2FVOUVVbTFOVkdoc1RtcEZlRTFxUlRKT01rbDVXbXBzYUZscVNYaE9SRlUwVG5wYWFFMTZaRGhVYTFaWVprRTlQUT09fFVIbzVUMkl5Tld4bVJFNXJUMVJOTUZwVVFUSlBSR3MxVFRKUmVFNHlWVFZhUkVVelQwZE9iRmxVU1RSWlYxcHRUbFJuTlU1NlNYaGFhbHBvVG5wc2JFNUVRWGhhVkVGNVdrZFpNazR5UlRSTmVtaHFXVlJhYTA5WFVYaGFSRm80Vkd0V1dHWkJQVDA9fGYwYTc0N2I5MjU2MmQyZDA2YmQ2NGM3Y2MxMjVjOTA5YWY4MjAwMmJkODgyNzlhYTFhMjliNTJiYmU4ZGZmN2J8TkVXfA%3D%3D%22%5D%2C%22no_fetch_context_on_resource%22%3Afalse%7D%2C%22context%22%3A%7B%7D%7D
//source_url=%2Fsearch%2Fpins%2F%3Fq%3Dthe%2520oscars%2520winners%2520portraits%26rs%3Dtyped%26term_meta%5B%5D%3Dthe%2520oscars%2520winners%2520portraits%257Ctyped&
//	data=%7B%22options%22%3A%7B%22article%22%3Anull%2C%22applied_filters%22%3Anull%2C%22appliedProductFilters%22%3A%22---%22%2C%22auto_correction_disabled%22%3Afalse%2C%22corpus%22%3Anull%2C%22customized_rerank_type%22%3Anull%2C%22filters%22%3Anull%2C%22query%22%3A%22the%20oscars%20winners%20portraits%22%2C%22query_pin_sigs%22%3Anull%2C%22redux_normalize_feed%22%3Atrue%2C%22rs%22%3A%22typed%22%2C%22scope%22%3A%22pins%22%2C%22source_id%22%3Anull%2C%22bookmarks%22%3A%5B%22Y2JVSG81V2sxcmNHRlpWM1J5VFVaU1YxWllhRlJXTVVreVZsZDRRMVV4U1hsVlZFWlhVbnBXTTFWNlNrZFdNa3BIWVVaYVdGSXlhRkZXUm1Rd1ZtMVJlRlZ1VWs1V1ZGWlBXVmh3UjFac1draE5WRkpWVFd0YWVWUnNhRXRYUmxsNlVXMUdWVlpzY0hsYVZscExWbFpXZEZKc1RsTmlXRkV4Vm10U1IxVXhUbkpPVlZwUVZsWmFXVmxzYUZOVU1YQllaRVYwYWxKc1NubFhhMVpoWWtaYVZWWnVhRmRpUmtwVVYxWmFTbVF3TVZWV2JGWk9VbXR3U1ZkWGVGWmxSVFYwVW10b2FGSnJTbGhWYlhoYVRXeFplVTFZWkZOaGVrWjVWR3hhVjFadFNsVlNibEpXWWtaS1dGVnFSbUZqVmxKeFZHeEdWbFpFUVRWYWExcFhVMWRLTmxWdGVGZE5XRUpLVm10amVHSXhiRmRUV0dob1RUSlNXVmxVUmt0VU1WSlhWbGhvYWxZd2NFbFpWVlUxWVVkRmVGZFljRmRTTTFKeVZqSXhWMk15VGtkV2JFcFlVakpvYUZadGRHRlpWMDVYVlc1S1ZtSnJOVzlVVm1RMFpVWldjMVZzVGxWaVJYQkpXWHBPYzFaV1dsZFRia1poVmpOTmVGVnNXbE5YVjBwR1QxWmtVMDFFUWpOV2FrbDRaREZrZEZac1pHbFRSa3BXVm10Vk1XRkdiRmhOVjNCc1lrWktlbGRyV21GVWJFcDFVVzVvVmxac1NraFdSRVpMVWpKS1JWZHNhRmRpUlhCRVYyeGFWbVZIVWtkVGJrWm9VbXhhYjFSV1duZFhiR1IwWkVWYVVGWnJTbE5WUmxGNFQwVXhObFJVVWs1bGJHdDVWRmN4U2sxR2NFVlhiV3hoWWxWcmQxUnJVbkpsVm14WVVtMXNXbFpGVmpOWFYzQnVUVEZyZVZKWWNFOVdSVEF3VjJ4U1FrMVdiRFpWVkU1aFlsWktjRmRyWkU5aVJteHhWbGh3V2sxc1NtOVVWbVJHWkRBeFNGTnRNVTVoYkhCdlZGWmtWazFzYTNwbFJUbFRWbTFSTkdaSFZteFBSRWt3VG1wRk1FNVVVVE5OVkUxNFRWZEplRTVFVm0xYWJVWnFXWHBWTkU5WFdtbFpWRTAwVFVkR2FVOUVVbTFOVkdoc1RtcEZlRTFxUlRKT01rbDVXbXBzYUZscVNYaE9SRlUwVG5wYWFFMTZaRGhVYTFaWVprRTlQUT09fFVIbzVUMkl5Tld4bVJFNXJUMVJOTUZwVVFUSlBSR3MxVFRKUmVFNHlWVFZhUkVVelQwZE9iRmxVU1RSWlYxcHRUbFJuTlU1NlNYaGFhbHBvVG5wc2JFNUVRWGhhVkVGNVdrZFpNazR5UlRSTmVtaHFXVlJhYTA5WFVYaGFSRm80Vkd0V1dHWkJQVDA9fGYwYTc0N2I5MjU2MmQyZDA2YmQ2NGM3Y2MxMjVjOTA5YWY4MjAwMmJkODgyNzlhYTFhMjliNTJiYmU4ZGZmN2J8TkVXfA%3D%3D%22%5D%2C%22no_fetch_context_on_resource%22%3Afalse%7D%2C%22context%22%3A%7B%7D%7D' \
