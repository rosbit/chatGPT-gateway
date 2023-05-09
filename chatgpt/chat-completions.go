package chatgpt

import (
	"chatGPT-gateway/conf"
	"fmt"
	"bytes"
)

type chatCompletionsResult struct {
	Id string `json:"id"`
	Object string `json:"object"`
	Created int64 `json:"created"`
	Choices []struct{
		Index int `json:"index"`
		Message struct {
			Role string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
}
func (r *chatCompletionsResult) checkFormat() bool {
	return r.Object == "chat.completion" || len(r.Object) > 0
}
func (r *chatCompletionsResult) String() string {
	if r.Object != "chat.completion" {
		return r.Object
	}

	if len(r.Choices) == 0 {
		return "no result"
	}
	sb := &bytes.Buffer{}
	for i, _ := range r.Choices {
		choice := &r.Choices[i]
		fmt.Fprintf(sb, "%s", choice.Message.Content)
	}
	return sb.String()
}

func CreateChatComplection(appName string, role, prompt, model string, maxTokens uint16) (res string, err error) {
	appConf := conf.GetAppConf(appName)
	if appConf == nil {
		err = fmt.Errorf("no conf found for %s", appName)
		return
	}
	body := appConf.ChatCompletions.MakeParams(role, prompt, model, maxTokens)
	var textRes chatCompletionsResult
	cgErr, e := callChatGPT(appConf, "/v1/chat/completions", body, &textRes)
	if e != nil {
		res = e.Error()
		err = e
		return
	}
	if cgErr != nil {
		res = cgErr.Error.Message
		return
	}
	res = textRes.String()
	return
}

func ContextChat(appName string, model string, maxTokens uint16, systemRole string, messages []string) (res string, err error) {
	appConf := conf.GetAppConf(appName)
	if appConf == nil {
		err = fmt.Errorf("no conf found for %s", appName)
		return
	}
	body := appConf.ChatCompletions.MakeContextParams(model, maxTokens, systemRole, messages)
	var textRes chatCompletionsResult
	cgErr, e := callChatGPT(appConf, "/v1/chat/completions", body, &textRes)
	if e != nil {
		res = e.Error()
		err = e
		return
	}
	if cgErr != nil {
		res = cgErr.Error.Message
		return
	}
	res = textRes.String()
	return
}
