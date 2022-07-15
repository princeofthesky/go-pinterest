package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go-pinterest/db"
	"go-pinterest/model"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var (
	SessionKey = "TWc9PSZFc3ZOOFRnZDl2UWRiOGJvaXZ6WEd6ZjArSDNxVGJ1Ly9NbGZZSjhTdW8yVTdtcHMxc2F0T0JvZW4rVmM1dzJIRkFmelFtUkVHRk9POCs2UEtJdE1JaGdscVBnK2lqSDdidFE3TVlkZmI2a0hmb3lzR3QyZzNKVmNoRFhwN0hhS0tyS1RmKzY2SkdTTWZvL3NLeDZmOVQzSlVyRGY0emJURHBZYzhFNUo0cU9kOExwdmh6NzE4bnp6NnMyTy9CNnR6SlpjNUZyc1RyeWUrSi9zOVBEa3hsVDl4MUlOTHF0QUdieUQwcXlQaTgyNHVNMitUclNjZURVOSt6ajN2RVN4Rks2U2Zkd1JQQldTd2U0WllaUGhnVDVjNzhQYWllOU9rNzBYQmc1SU8rL2NzN1FIVkFpMUszMXpRZmIzRWpoUFg2eW44RmM1RnVLWndna2ZDWjNEUWFJclpiaFJyRWI0Y2NiUHdUZFZHUHQwSk5EdjduaG53VDlVR2FhaFRkcllRT1U2YVZoV0xxb3gzZTYxOTBpUC8wQnlLYS8wL3JQU01sc2FaNmxVbjFNOXZxRWl4OVdFSnhPWDF3aTZNTHh4d00vb08rbHNybkZTSmpDNFhPbUpOeEpoZ2JIeThVMThodUlWUEF4Q0hrcW15K3dad056OExaMHZYU0ltQWJWMUViKzZjd0VnYS9FRkVOdGs2MXB5M0NjdFZNU2JJOGlSWjhDUkhFVEF1M0RheUduZnVXSGFzRi9SQUI2THJsYzlTZW5tdG1UMk1DUWxjK3Jjd0hoYWN0UEM2bndSVzZlbVZNM29KZEIrV0ttMm5Id0Rnb3FNQU16KzRLRDFRVjFtZDBHcDY1Q1l2TS9rSGdBYnZMQTY4M0xIYXlYSnNTYlRoWjVKdUVIZ2Uxa014KzQwZ2tldnhZcTE2R25JYXVqUXhnSWd3dXIwOWZmY2dqVndTUmxmSUc2TFYyQUlmWkgrbkRVcTVXNHp3cUlGOFJQUXZzQVhkYk5Pb0p6THBXWmtUTTRvMVR4Qks5K0FWZVBHejRsQ09nPT0mNjdwREVYRnNHM1JOODNzYUYyMW5BREkvTGc4PQ=="
	Csrftoken  = "cb39a1e889e0f39fc26c75615973fb21"
	Category ="Magazines"
	Keyword="magazines cover poses vogue"
)

func main() {
	db.Init()
	LoadConfigCrawler()
	//imageinfos:=GetPinFromKeyWordSearch("character face cyberpunk")

	imageinfos,err:=GetPinFromNextPageSearch("character face cyberpunk")
	if err !=nil {
		fmt.Println("err",err)
	}
	Ids:=make(map[string]bool)
	for i:=0;i<len(imageinfos);i++ {
		Ids[imageinfos[i].SourceId]=true
	}
	for Key,_ := range Ids{
		GetPinFromRelatedPin(Key, 50)
	}
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
func GetPinFromKeyWordSearch(keyword string) error {
	query := url.QueryEscape("/search/pins/?q=" + url.QueryEscape(keyword))
	resp, err := http.Get("https://www.pinterest.com/resource/BaseSearchResource/get/?source_url=" + query)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to fetch bodyBytes: %d %s", resp.StatusCode, resp.Status)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	response := map[string]interface{}{}
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return err
	}
	resourceResponse := response["resource_response"].(map[string]interface{})
	data := resourceResponse["data"].(map[string]interface{})
	results := data["results"].([]interface{})
	for i := 0; i < len(results); i++ {
		info := results[i].(map[string]interface{})
		//images:=info["images"]
		fmt.Println(i)
		pinner := info["pinner"].(map[string]interface{})
		title := info["title"]
		createdTime := info["created_at"]
		IDPinner := pinner["id"]
		IDUsername := pinner["username"]
		IDFullname := pinner["full_name"]
		fmt.Println(title)
		fmt.Println(createdTime)
		fmt.Println(IDPinner)
		fmt.Println(IDUsername)
		fmt.Println(IDFullname)
	}
	return nil
}

func GetPinFromRelatedPin(pinID string, size int) error {

	dataQuery := "{\"options\":{\"pin_id\":\"" +
		pinID +
		"\",\"context_pin_ids\":[],\"page_size\":" +
		strconv.Itoa(size) +
		",\"search_query\":\"\",\"source\":\"deep_linking\",\"top_level_source\":\"deep_linking\",\"top_level_source_depth\":1,\"is_pdp\":false,\"no_fetch_context_on_resource\":false},\"context\":{}}"
	query := url.QueryEscape("/pin/"+pinID+"/") + "&data=" + url.QueryEscape(dataQuery)
	uri := "https://www.pinterest.com/resource/RelatedModulesResource/get/?source_url=" + query
	fmt.Println(uri)
	req, err := http.NewRequest("GET", uri, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")
	req.Header.Set("cookie", " _pinterest_sess="+SessionKey+"")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to fetch bodyBytes: %d %s", resp.StatusCode, resp.Status)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	response := map[string]interface{}{}
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return err
	}
	fmt.Println(len(bodyBytes))
	resourceResponse := response["resource_response"].(map[string]interface{})
	data := resourceResponse["data"].([]interface{})
	fmt.Println(len(data))
	for i := 1; i < len(data); i++ {
		info := data[i].(map[string]interface{})
		images := info["images"]
		fmt.Println(i)
		pinner := info["pinner"].(map[string]interface{})
		title := info["title"]
		createdTime := info["created_at"]
		IDPinner := pinner["id"]
		IDUsername := pinner["username"]
		IDFullname := pinner["full_name"]
		fmt.Println(title)
		fmt.Println(createdTime)
		fmt.Println(IDPinner)
		fmt.Println(IDUsername)
		fmt.Println(IDFullname)
	}
	return nil
}

func GetPinFromNextPageSearch(keyword string) ([]model.ImageInfo, error) {
	bodyPost := []byte("source_url=%2Fsearch%2Fpins%2F%3Fq%3Dcharacter%2520face%2520cyberpunk")
	req, err := http.NewRequest("POST", "https://www.pinterest.com/resource/BaseSearchResource/get/", bytes.NewBuffer(bodyPost))
	req.Header.Set("cookie", "csrftoken="+Csrftoken)
	req.Header.Set("x-csrftoken", SessionKey)

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
		fmt.Println(err)
		return nil, err
	}
	resourceResponse := response["resource_response"].(map[string]interface{})
	data := resourceResponse["data"].(map[string]interface{})
	results := data["results"].([]interface{})
	imageInfos := make([]model.ImageInfo, 0)
	for i := 0; i < len(results); i++ {
		info := results[i].(map[string]interface{})
		//images:=info["images"]
		fmt.Println(i)
		pinner := info["pinner"].(map[string]interface{})
		imageInfo := model.ImageInfo{}
		imageInfo.Title = info["title"].(string)
		imageInfo.SourceId = info["id"].(string)
		imageInfo.SourceLink = "" + imageInfo.SourceId
		//createdTime := info["created_at"]

		imageInfo.Source = model.Pinterest
		imageInfo.SourceId = pinner["id"].(string)
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
	fmt.Println(info["images"])
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
