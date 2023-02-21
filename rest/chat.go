package rest

import (
	"github.com/rosbit/mgin"
	"chatGPT-gateway/chatgpt"
	"net/http"
)

// POST /chat?app=xxx
// {
//    "prompt": "xxxx",
//    "model": "xxx",
//    "max-tokens": xxx
// }
func Chat(c *mgin.Context) {
	var params struct {
		App string `query:"app"`
	}
	if code, err := c.ReadParams(&params); err != nil {
		c.Error(code, err.Error())
		return
	}
	var promptParams struct {
		Prompt string `json:"prompt"`
		Model  string `json:"model"`
		MaxTokens uint16 `json:"max-tokens"`
	}
	if code, err := c.ReadJSON(&promptParams); err != nil {
		c.Error(code, err.Error())
		return
	}
	res, err := chatgpt.CreateComplection(params.App, promptParams.Prompt, promptParams.Model, promptParams.MaxTokens)
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
		"result": res,
	})
}
