package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go-pinterest/db"
	"net/http"
	"strconv"
	"time"
)

var (
	update = flag.Int("update", 0, "update info again from pinterest or not : = 1 (force) , =0 : only new info")
)

func main() {
	flag.Parse()
	db.Init()
	for true {
		categories, _ := db.GetAllCategory()
		for _, category := range categories {
			imageIds, _ := db.GetAllImageIdByCategory(category)
			for _, imageId := range imageIds {
				id, err := strconv.ParseInt(imageId, 10, 64)
				if err == nil && id > 0 {
					imageInfo, _ := db.GetImageInfo(id)
					if *update == 1 || (len(imageInfo.Image) > 0 && len(imageInfo.Annotations) == 0) {
						err = UpdatePinInfo(imageInfo.Id, imageInfo.SourceId)
						if err != nil {
							fmt.Println("err when update pin info", imageInfo.Id, imageInfo.SourceId)
							fmt.Println("err", err)
						}
					}
				}
			}
		}
		time.Sleep(30 * time.Minute)
	}
}

func UpdatePinInfo(imageId int64, pinID string) error {
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
	fmt.Println(imageId, pinID)
	//ioutil.WriteFile("1.txt", []byte(data), fs.ModePerm)
	response := map[string]interface{}{}
	err = json.Unmarshal([]byte(data), &response)
	if err != nil {
		return err
	}
	props := response["props"].(map[string]interface{})
	initialReduxState := props["initialReduxState"].(map[string]interface{})
	pins := initialReduxState["pins"].(map[string]interface{})
	if pins[pinID] == nil {
		return fmt.Errorf("can request data with Pin ID :%v", pinID)
	}

	info := pins[pinID].(map[string]interface{})
	title := info["title"].(string)
	grid_title := info["grid_title"].(string)
	if len(title) == 0 && len(grid_title) > 0 {
		title = grid_title
	}
	description := info["description"].(string)
	closeupUnifiedDescription := info["closeup_unified_description"].(string)
	if len(title) == 0 {
		if len(description) > 0 {
			title = description
		} else if len(closeupUnifiedDescription) > 0 {
			title = closeupUnifiedDescription
		} else {
			title = ""
		}
	}
	if len(description) == 0 {
		if len(info["description_html"].(string)) > 0 {
			description = info["description_html"].(string)
		} else {
			description = ""
		}
	}

	var closeupAttribution map[string]interface{}

	if info["closeup_attribution"] != nil {
		closeupAttribution = info["closeup_attribution"].(map[string]interface{})
	} else {
		closeupAttribution = info["pinner"].(map[string]interface{})
	}

	ownerName := closeupAttribution["full_name"].(string)
	ownerUrl := "https://www.pinterest.com/" + closeupAttribution["username"].(string)

	pinJoin := info["pin_join"].(map[string]interface{})
	annotationsWithLinks := pinJoin["annotations_with_links"].(map[string]interface{})
	annotations := make([]string, 0)
	for key, _ := range annotationsWithLinks {
		annotations = append(annotations, key)
	}

	hashtagsInfo := info["hashtags"].([]interface{})
	hashtags := make([]string, 0)
	for _, key := range hashtagsInfo {
		hashtags = append(hashtags, key.(string))
	}
	boards := initialReduxState["boards"].(map[string]interface{})
	boardDescription := ""
	for _, v := range boards {
		board := v.(map[string]interface{})
		boardDescription = board["description"].(string)
		break
	}
	db.UpdateImageInfo(imageId, title, description, ownerName, ownerUrl, annotations, hashtags, boardDescription)
	return nil
}
