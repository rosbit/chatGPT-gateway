package chatgpt

import (
	"chatGPT-gateway/conf"
	"fmt"
)

type imageResult struct {
	Created int64 `json:"created"`
	Data []struct {
		Url string `json:"url"`
	} `json:"data"`
}
func (r *imageResult) checkFormat() bool {
	return len(r.Data) > 0
}

func CreateImage(appName string, prompt string, size string, n uint8) (created int64, urls []string, err error) {
	appConf := conf.GetAppConf(appName)
	if appConf == nil {
		err = fmt.Errorf("no conf found for %s", appName)
		return
	}
	body := appConf.Image.MakeParams(prompt, n, size)
	var imgRes imageResult
	cgErr, err := callChatGPT(appConf, "/v1/images/generations", body, &imgRes)
	if err != nil {
		return
	}
	if cgErr != nil {
		err = fmt.Errorf("%s", cgErr.Error.Message)
		return
	}
	urls = make([]string, len(imgRes.Data))
	for i, d := range imgRes.Data {
		urls[i] = d.Url
	}
	created = imgRes.Created
	return
}
