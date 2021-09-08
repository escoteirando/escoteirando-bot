package mappa_api

import (
	"fmt"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"github.com/guionardo/escoteirando-bot/src/utils"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

var (
	httpClient    = &http.Client{}
	bodyFilesPath = ""
	saveBodyFunc  func(string, []byte)
)

func init() {
	saveBodyFunc = func(url string, body []byte) {
		return
	}
	if utils.GetEnvironmentSetup().BotDebug {
		wd, _ := os.Getwd()
		filePath := path.Join(wd, ".requests")

		if fi, err := os.Stat(filePath); !fi.IsDir() {
			err = os.Mkdir(filePath, 0666)
			if err != nil {
				log.Printf("Failed to create path %s %v", filePath, err)
				return
			}
		}
		bodyFilesPath = filePath
		saveBodyFunc = saveBody
	}
}
func fetchMappa(url string, ctx domain.Context) ([]byte, int, error) {
	env := utils.GetEnvironmentSetup()
	mappaUrl := fmt.Sprintf("%s/%s", env.MappaProxyUrl, url)
	req, err := http.NewRequest("GET", mappaUrl, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("request error: %s - %v", url, err)
	}
	req.Header.Set("Authorization", ctx.AuthKey)
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request error: %s - %v", url, err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("request error: %s - %v", url, err)
	}
	saveBodyFunc(url, body)
	return body, res.StatusCode, nil
}

func saveBody(url string, body []byte) {
	fileUrl := strings.SplitN(strings.ReplaceAll(url, "/", "_"), "?", 1)[0]
	if len(fileUrl) > 20 {
		fileUrl = fileUrl[0:20]
	}

	fileName := path.Join(bodyFilesPath, time.Now().Format("20060102_150405_")+fileUrl)
	err := os.WriteFile(fileName, body, 0666)
	if err != nil {
		log.Printf("Failed to write response body %s %v", url, err)
	} else {
		log.Printf("Response body saved %s %s", url, fileName)
	}
}
