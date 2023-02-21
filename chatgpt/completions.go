package chatgpt

import (
	"chatGPT-gateway/conf"
	"fmt"
	"bytes"
)

type completionsResult struct {
	Id string `json:"id"`
	Object string `json:"object"`
	Created int64 `json:"created"`
	Model string `json:"model"`
	Choices []struct{
		Text string `json:"text"`
		Index int `json:"index"`
		Logprobs []int `json:"logprobs"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
}
func (r *completionsResult) checkFormat() bool {
	return r.Object == "text_completion" || len(r.Object) > 0
}
func (r *completionsResult) String() string {
	if r.Object != "text_completion" {
		return r.Object
	}

	if len(r.Choices) == 0 {
		return "no result"
	}
	sb := &bytes.Buffer{}
	for i, _ := range r.Choices {
		choice := &r.Choices[i]
		fmt.Fprintf(sb, "%s", choice.Text)
	}
	return sb.String()
}

func CreateComplection(appName string, prompt, model string, maxTokens uint16) (res string, err error) {
	appConf := conf.GetAppConf(appName)
	if appConf == nil {
		err = fmt.Errorf("no conf found for %s", appName)
		return
	}
	body := appConf.Text.MakeParams(prompt, model, maxTokens)
	var textRes completionsResult
	cgErr, e := callChatGPT(appConf, "/v1/completions", body, &textRes)
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
