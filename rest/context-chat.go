package rest

import (
	"github.com/rosbit/mgin"
	"chatGPT-gateway/chatgpt"
	"net/http"
)

// POST /context-chat?app=xxx
// {
//    "model": "xxx",
//    "max-tokens": xxx,
//    "system-role": "optional",
//    "messages": ["q1", "a1", "q2", "a2", ..., "qn"]
// }
func ContextChat(c *mgin.Context) {
	var params struct {
		App string `query:"app"`
	}
	if code, err := c.ReadParams(&params); err != nil {
		c.Error(code, err.Error())
		return
	}
	var promptParams struct {
		Model  string `json:"model"`
		MaxTokens uint16 `json:"max-tokens"`
		SystemRole string `json:"system-role"`
		Messages []string `json:"messages"`
	}
	if code, err := c.ReadJSON(&promptParams, true); err != nil {
		c.Error(code, err.Error())
		return
	}
	res, err := chatgpt.ContextChat(params.App, promptParams.Model, promptParams.MaxTokens, promptParams. SystemRole, promptParams.Messages)
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
