package chatgpt

import (
	"github.com/rosbit/gnet"
	"chatGPT-gateway/conf"
	"net/url"
	"os"
	"fmt"
	"encoding/json"
)

type chatGPTResult interface {
	checkFormat() bool
}

type chatgptError struct {
	Error struct {
		Message string `json:"message"`
		Type string `json:"type"`
	} `json:"error"`
}

func callChatGPT(appConf *conf.AppParams, uri string, body interface{}, res chatGPTResult) (cgerr *chatgptError, err error) {
	url := fmt.Sprintf("%s%s", conf.ServiceConf.APIBaseURL, uri)
	_, resBody, _, e := gnet.JSON(url, gnet.Headers(map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", appConf.APIKey),
		}),
		gnet.Params(body),
		gnet.BodyLogger(os.Stderr),
		gnet.WithTimeoutDuration(conf.GetDefaultTimeout()),
	)
	if e != nil {
		err = e
		return
	}
	if err = json.Unmarshal(resBody, res); err != nil {
		return
	}
	if res.checkFormat() {
		return
	}

	var gptErr chatgptError
	if err = json.Unmarshal(resBody, &gptErr); err == nil {
		cgerr = &gptErr
	}
	return
}

func IsTimeoutError(err error) bool {
	if e, ok := err.(*url.Error); ok && e.Timeout() {
		return true
	}
	return false
}

