package conf

import (
	"fmt"
)

type AppParams struct {
	Name   string `yaml:"name"`
	APIKey string `yaml:"api-key"`
	Text TextParams `yaml:"text"`
	ChatCompletions ChatCompletionsParams `yaml:"chat-completions"`
	Image ImageParams `yaml:"image"`
}
func (appConf *AppParams) checkMust(i int) (error) {
	if len(appConf.Name) == 0 {
		return fmt.Errorf("apps[%d]/name expected in conf", i)
	}
	if len(appConf.APIKey) == 0 {
		return fmt.Errorf("apps[%d]/api-key expected in conf", i)
	}
	if err := appConf.Text.checkMust(i); err != nil {
		return err
	}
	if err := appConf.ChatCompletions.checkMust(i); err != nil {
		return err
	}
	if err := appConf.Image.checkMust(i); err != nil {
		return err
	}
	return nil
}

type TextParams struct {
	Model string `yaml:"model"`
	Temperature float32 `yaml:"temperature"`
	MaxTokens uint16 `yaml:"max-tokens"`
}
func (p *TextParams) checkMust(i int) (error) {
	if len(p.Model) == 0 {
		return fmt.Errorf("apps[%d]/text/model expected in conf", i)
	}
	if p.Temperature < 0.0 || p.Temperature > 2.0 {
		return fmt.Errorf("apps[%d]/text/temperature must between 0 and 2 in conf", i)
	}
	if p.MaxTokens == 0 {
		p.MaxTokens = 16
	} else if p.MaxTokens > 2048 {
		return fmt.Errorf("apps[%d]/text/max-tokens must be less than 2048", i)
	}
	return nil
}
func (p *TextParams) MakeParams(prompt, model string, maxTokens uint16) map[string]interface{} {
	if len(model) == 0 {
		model = p.Model
	}
	if maxTokens == 0 {
		maxTokens = p.MaxTokens
	}
	return map[string]interface{}{
		"model": model,
		"prompt": prompt,
		"temperature": p.Temperature,
		"max_tokens": maxTokens,
	}
}

type ChatCompletionsParams struct {
	Model string `yaml:"model"`
	Role string `yaml:"role"`
	Temperature float32 `yaml:"temperature"`
	MaxTokens uint16 `yaml:"max-tokens"`
}
func (p *ChatCompletionsParams) checkMust(i int) (error) {
	if len(p.Model) == 0 {
		return fmt.Errorf("apps[%d]/chat-completions/model expected in conf", i)
	}
	if len(p.Role) == 0 {
		p.Role = "user"
	}
	if p.Temperature < 0.0 || p.Temperature > 2.0 {
		return fmt.Errorf("apps[%d]/chat-completions/temperature must between 0 and 2 in conf", i)
	}
	if p.MaxTokens == 0 {
		p.MaxTokens = 16
	} else if p.MaxTokens > 8096 {
		return fmt.Errorf("apps[%d]/chat-completions/max-tokens must be less than 8096", i)
	}
	return nil
}
func (p *ChatCompletionsParams) MakeParams(role, prompt, model string, maxTokens uint16) map[string]interface{} {
	if len(role) == 0 {
		role = p.Role
	}
	if len(model) == 0 {
		model = p.Model
	}
	if maxTokens == 0 {
		maxTokens = p.MaxTokens
	}
	return map[string]interface{}{
		"model": model,
		"messages": []interface{}{
			map[string]interface{}{
				"role": role,
				"content": prompt,
			},
		},
		"temperature": p.Temperature,
		"max_tokens": maxTokens,
	}
}

type ImageParams struct {
	Size string `yaml:"size"`
	Num  uint8 `yaml:"num"`
	ResponseFormat string `yaml:"response-format"`
}
func (p *ImageParams) checkMust(i int) (error) {
	if len(p.Size) == 0 {
		p.Size = "1024x1024"
	} else {
		switch p.Size {
		case "256x256", "512x512", "1024x1024":
		default:
			return fmt.Errorf("bad format or unsupportable size of apps[%d]/image/size", i)
		}
	}
	if p.Num == 0 {
		p.Num = 1
	}
	if p.Num > 10 {
		p.Num = 1
	}
	if len(p.ResponseFormat) == 0 {
		p.ResponseFormat = "url"
	} else {
		switch p.ResponseFormat{
		case "url", "b64_json":
		default:
			return fmt.Errorf("unknown format for apps[%d]/image/response-format", i)
		}
	}
	return nil
}
func (p *ImageParams) MakeParams(prompt string, n uint8, size string) map[string]interface{} {
	if n == 0 {
		n = p.Num
	}
	if len(size) == 0 {
		size = p.Size
	}

	return map[string]interface{}{
		"prompt": prompt,
		"n": n,
		"size": size,
		"response_format": p.ResponseFormat,
	}
}
