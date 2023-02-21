package rest

import (
	"github.com/rosbit/mgin"
	"chatGPT-gateway/chatgpt"
	"net/http"
)

// POST /image?app=xxx
// {
//    "prompt": "xxxx"
//    "size": "optional size",
//    "n": option-count
// }
func CreateImage(c *mgin.Context) {
	var params struct {
		App string `query:"app"`
	}
	if code, err := c.ReadParams(&params); err != nil {
		c.Error(code, err.Error())
		return
	}
	var promptParams struct {
		Prompt string `json:"prompt"`
		Size string `json:"size"`
		N    uint8  `json:"n"`
	}
	if code, err := c.ReadJSON(&promptParams); err != nil {
		c.Error(code, err.Error())
		return
	}
	created, urls, err := chatgpt.CreateImage(params.App, promptParams.Prompt, promptParams.Size, promptParams.N)
	if err != nil {
		var status int
		if chatgpt.IsTimeoutError(err) {
			status = http.StatusGatewayTimeout
		} else {
			status = http.StatusInternalServerError
		}
		c.Error(status, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg": "OK",
		"result": urls,
		"created": created,
	})
}
