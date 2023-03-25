package rest

import (
	"github.com/rosbit/mgin"
	"chatGPT-gateway/chatgpt"
	"net/http"
)

// POST /chat-completions?app=xxx
// {
//    "role": "user",
//    "prompt": "xxxx",
//    "model": "xxx",
//    "max-tokens": xxx
// }
func ChatCompletions(c *mgin.Context) {
	var params struct {
		App string `query:"app"`
	}
	if code, err := c.ReadParams(&params); err != nil {
		c.Error(code, err.Error())
		return
	}
	var promptParams struct {
		Role string `json:"role"`
		Prompt string `json:"prompt"`
		Model  string `json:"model"`
		MaxTokens uint16 `json:"max-tokens"`
	}
	if code, err := c.ReadJSON(&promptParams, true); err != nil {
		c.Error(code, err.Error())
		return
	}
	res, err := chatgpt.CreateChatComplection(params.App, promptParams.Role, promptParams.Prompt, promptParams.Model, promptParams.MaxTokens)
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
