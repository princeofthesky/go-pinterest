package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go-pinterest/db"
	"go-pinterest/model"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var (
	SessionKey = "TWc9PSY5UzBwT01YdnZsbG55d3phZGVHNzVFWC93WTg3Z0dOS0FkUUtrWkRBRFRFVE94WFB2K011TjE5aC9GZHFuTTFDOTNsMEJoOENPWU5EODI3TGQ0TTlqTGVGNjY1T3lpUWZqRUgrSG1rYmpNZlpEcFBHN1JRV1BGc2VoOXYwS2FWc1BwNitrdjVpdkltRitwbXVlMnFVbVY2KzgwMkFpclBEcTBHOC9nVkYzbTNoWStTeGxGR1V6dWI1MUR5ZnRSWC8vVjhVaXFSTWhkamJlVUsvdmFNZVJBN2dJRCs2TEFiM1hvMy9yTm5sUjVZNEZlOGtFb3k2dTl0VmliUEZzNmROV1k0ODd5ajR3andzcVM1OG9tSDBMSHFCdGJ4clZjK1VpM3h4Rm9jQU9PcnFoNU91a2F4c3JPWnVuaTM2TUxhVTZUcFFNZDcyNDFqdXlyR0g0WlJpSXJncEkrS3B6c3VCTlNvVnFTZ1BzcWdKUlJacFlCMnZ6NjdkSTJqQ0VQd0gyOFYvdWhEM3dZVS9HNlUvRDlMOUdRUk1uWHlEQUlFbi9WVWF6UHd0QjFHVGRoNWRMbnhwaTJCREZFMnhmelBLaktZaTlEWEJBcXI4dWllc0RwbFl2Y3FRalBWcmtIOXhpMVUxcUZiWGhjenBFOGhBOWF3MkVDMUVrNFlZa3R6eGNYb0VEb2YyeHdQZDdrM2N1SkpYTDVVNjV5bkVYZDdEeEtXTzZRcjZqZHpBSjFtVVZxaGdTbkRnUzhwN0FEbWJrZnUrR2JsSWxxYUw2UjFEb2xQR0x5SUpibXN0dUhlU3YvbzFFdlBZQVJvbTFiL0ZJVUxoWm4vVEZVTXRnbm1ZR04zRmpBdS82cENRNFdlOTNJU2t6VjlwRVBVcTZlZm5kMTNmdVlwYTJhQURrNW8ya05WeVlUTjRRMUhzUmpwTWV0c1EzZ3I4UXlQd1YxUnZKWWM5ZGhmcG8vbnUrRWNXOFNJMW9pamxLbmQ5dEduZTJRU09YYmpUZ0E1cmR6Y0k2azJMSUo5Zm1GVGczUTVZdzFHNmZBPT0mMVFjUHNDQUZaZkFRWjduSDl6Zy9hRGJqcEpzPQ=="
	Csrftoken  = "e7525248933e9ce4fe7cd9507627bd03"
	Category   = "Magazines"
	Keyword    = "magazines cover poses vogue"
	source     = flag.String("source", "/home/tamnb/projects/src/github.com/nguyenbatam90/go-pinterest/sample_crawl_source.csv", "source category Pinterest")
	size       = flag.Int("size", 50, "max items per query related,search")
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
	//LoadConfigCrawler()

	IdCrawled := make(map[string]bool)
	IdRelatedCrawl := make(map[string]bool)
	for keyword, _ := range mapKeywordCategory {

		imageSearchinfos, _ := GetPinFromKeyWordSearch(keyword)
		imageNextSearchInfos, err := GetPinFromNextPageSearch(keyword)
		if err != nil {
			fmt.Println(keyword)
			fmt.Println("err", err)
		}
		fmt.Println(keyword , "search  " , len(imageSearchinfos))
		fmt.Println(keyword , "search next " , len(imageNextSearchInfos))
		newIds := make(map[string]bool)
		for i := 0; i < len(imageSearchinfos); i++ {
			id := imageSearchinfos[i].SourceId
			//if !IdRelatedCrawl[id] {
				IdRelatedCrawl[id]=true
				IdCrawled[id] = true
				newIds[id] = true
			//}
		}
		for i := 0; i < len(imageNextSearchInfos); i++ {
			id := imageNextSearchInfos[i].SourceId
			//if !IdRelatedCrawl[id] {
				IdRelatedCrawl[id]=true
				IdCrawled[id] = true
				newIds[id] = true
			//}
		}

		for pinId, _ := range newIds {
			IdRelatedCrawl[pinId]=true
			relatedImages, err := GetPinFromRelatedPin(pinId, 50)
			if err != nil {
				fmt.Println(pinId)
				fmt.Println("err", err)
			}
			for i := 0; i < len(relatedImages); i++ {
				IdCrawled[relatedImages[i].SourceId]=true
			}
		}
		fmt.Println("new",len(newIds), " total",len(IdCrawled))
	}
	fmt.Println("total",len(IdCrawled))
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
	Csrftoken = config["pinterest_csrftoken"]
	SessionKey = config["pinterest_SessionKey"]
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
func GetPinFromKeyWordSearch(keyword string) ([]model.ImageInfo, error) {

	uri := "https://www.pinterest.com/resource/BaseSearchResource/get/?source_url=" +url.QueryEscape("/search/pins/?q=" + url.QueryEscape(keyword))
	resp, err := http.Get(uri)
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
	data := resourceResponse["data"].(map[string]interface{})
	results := data["results"].([]interface{})

	imageInfos := make([]model.ImageInfo, 0)
	for i := 0; i < len(results); i++ {
		info := results[i].(map[string]interface{})
		//images:=info["images"]
		pinner := info["pinner"].(map[string]interface{})

		//createdTime := info["created_at"]

		imageInfo := model.ImageInfo{}
		imageInfo.Title = info["title"].(string)
		imageInfo.SourceId = info["id"].(string)
		imageInfo.SourceLink = "" + imageInfo.SourceId
		//createdTime := info["created_at"]

		imageInfo.Source = model.Pinterest
		imageInfo.SourceOwnerName = pinner["full_name"].(string)
		imageInfo.SourceOwnerLink = pinner["username"].(string)
		imageInfos = append(imageInfos, imageInfo)

	}
	return imageInfos, nil
}

func GetPinFromRelatedPin(pinID string, size int) ([]model.ImageInfo, error) {

	dataQuery := "{\"options\":{\"pin_id\":\"" +
		pinID +
		"\",\"context_pin_ids\":[],\"page_size\":" +
		strconv.Itoa(size) +
		",\"search_query\":\"\",\"source\":\"deep_linking\",\"top_level_source\":\"deep_linking\",\"top_level_source_depth\":1,\"is_pdp\":false,\"no_fetch_context_on_resource\":false},\"context\":{}}"
	query := url.QueryEscape("/pin/"+pinID+"/") + "&data=" + url.QueryEscape(dataQuery)
	uri := "https://www.pinterest.com/resource/RelatedModulesResource/get/?source_url=" + query
	//fmt.Println(uri)
	req, err := http.NewRequest("GET", uri, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")
	req.AddCookie(&http.Cookie{Name: "csrftoken",Value: Csrftoken})
	req.AddCookie(&http.Cookie{Name: "_pinterest_sess",Value: SessionKey})

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
	//fmt.Println(len(bodyBytes))
	resourceResponse := response["resource_response"].(map[string]interface{})
	data := resourceResponse["data"].([]interface{})
	//fmt.Println(resourceResponse["http_status"],len(data))
	imageInfos := make([]model.ImageInfo, 0)
	for i := 1; i < len(data); i++ {
		info := data[i].(map[string]interface{})
		if info["type"].(string) != "pin" {
			continue
		}
		//images := info["images"]
		pinner := info["pinner"].(map[string]interface{})
		imageInfo := model.ImageInfo{}
		imageInfo.Title = info["title"].(string)
		imageInfo.SourceId = info["id"].(string)
		imageInfo.SourceLink = "" + imageInfo.SourceId
		//createdTime := info["created_at"]

		imageInfo.Source = model.Pinterest
		imageInfo.SourceOwnerName = pinner["full_name"].(string)
		imageInfo.SourceOwnerLink = pinner["username"].(string)
		imageInfos=append(imageInfos,imageInfo)
	}
	return imageInfos, nil
}

func GetPinFromNextPageSearch(keyword string) ([]model.ImageInfo, error) {

	bodyPost := []byte("source_url=" +url.QueryEscape("/search/pins/?q="+url.QueryEscape(keyword)))

	req, err := http.NewRequest("POST", "https://www.pinterest.com/resource/BaseSearchResource/get/", bytes.NewBuffer(bodyPost))
	req.AddCookie(&http.Cookie{Name: "csrftoken",Value: Csrftoken})
	req.AddCookie(&http.Cookie{Name: "_pinterest_sess",Value: SessionKey})
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")
	req.Header.Set("x-csrftoken",Csrftoken)

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
		//fmt.Println(err)
		return nil, err
	}
	resourceResponse := response["resource_response"].(map[string]interface{})
	data := resourceResponse["data"].(map[string]interface{})
	results := data["results"].([]interface{})
	imageInfos := make([]model.ImageInfo, 0)
	for i := 0; i < len(results); i++ {
		info := results[i].(map[string]interface{})
		//images:=info["images"]
		//fmt.Println(i)
		pinner := info["pinner"].(map[string]interface{})
		imageInfo := model.ImageInfo{}
		imageInfo.Title = info["title"].(string)
		imageInfo.SourceId = info["id"].(string)
		imageInfo.SourceLink = "" + imageInfo.SourceId
		//createdTime := info["created_at"]

		imageInfo.Source = model.Pinterest
		imageInfo.SourceOwnerName = pinner["full_name"].(string)
		imageInfo.SourceOwnerLink = pinner["username"].(string)
		imageInfos = append(imageInfos, imageInfo)
	}
	return imageInfos, nil
}
func GetPinInfo(pinID string) error {
	resp, err := http.Get("https://www.pinterest.com/pin/" + pinID)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to fetch data: %d %s", resp.StatusCode, resp.Status)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}
	data := doc.Find("#__PWS_DATA__").Text()
	response := map[string]interface{}{}
	err = json.Unmarshal([]byte(data), &response)
	if err != nil {
		return err
	}
	props := response["props"].(map[string]interface{})
	initialReduxState := props["initialReduxState"].(map[string]interface{})
	pins := initialReduxState["pins"].(map[string]interface{})
	info := pins[pinID].(map[string]interface{})
	//fmt.Println(info["images"])
	//["initialReduxState"]["pins"]

	title := info["title"]
	createdTime := info["created_at"]

	pinner := info["pinner"].(map[string]interface{})
	IDpinner := pinner["id"]
	IDusername := pinner["username"]
	fmt.Println(title)
	fmt.Println(createdTime)
	fmt.Println(IDpinner)
	fmt.Println(IDusername)
	return nil
}
