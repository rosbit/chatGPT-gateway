// global conf
// ENV:
//   CONF_FILE      --- 配置文件名
//   TZ             --- 时区名称"Asia/Shanghai"
//
// YAML
// ---
// listen-host: ""
// listen-port: 7080
// api-base-url: https://api.openai.com
// default-timeout: 30s
// apps:
//   - name: app-name
//     api-key: xxxx
//     text:
//       model: text-davinci-003
//       temperature: 0
//       max-tokens: 128
//     chat-completions:
//       role: user
//       model: gpt-3.5-turbo
//       temperature: 0
//       max-tokens: 128
//     image:
//       size: 1024x1024
//       num: 1
//       repsonse-format: url
//   #-
// common-endpoints:
//   health-check: "/health"
//   chat: "/chat"
//   chat-completions: "/chat-completions"
//   context-chat: "/context-chat"
//   image: "/image"
//
// Rosbit Xu

package conf

import (
	"gopkg.in/yaml.v2"
	"fmt"
	"os"
	"time"
)

type ServiceConfT struct {
	ListenHost string `yaml:"listen-host"`
	ListenPort int    `yaml:"listen-port"`
	APIBaseURL string `yaml:"api-base-url"`
	DefaultTimeout string `yaml:"default-timeout"`
	defaultTimeout time.Duration
	Apps []AppParams `yaml:"apps"`
	CommonEndpoints struct {
		HealthCheck string `yaml:"health-check"`
		Chat  string `yaml:"chat"`
		ChatCompletions string `yaml:"chat-completions"`
		Image string `yaml:"image"`
		ContextChat string `yaml:"context-chat"`
	} `yaml:"common-endpoints"`
}

var (
	ServiceConf ServiceConfT
	Loc = time.FixedZone("UTC+8", 8*60*60)
	appConfs = map[string]*AppParams{}
)


func getEnv(name string, result *string, must bool) error {
	s := os.Getenv(name)
	if s == "" {
		if must {
			return fmt.Errorf("env \"%s\" not set", name)
		}
	}
	*result = s
	return nil
}

func CheckGlobalConf() error {
	var p string
	getEnv("TZ", &p, false)
	if p != "" {
		if loc, err := time.LoadLocation(p); err == nil {
			Loc = loc
		}
	}

	var confFile string
	if err := getEnv("CONF_FILE", &confFile, true); err != nil {
		return err
	}

	fp, err := os.Open(confFile)
	if err != nil {
		return err
	}
	defer fp.Close()

	dec := yaml.NewDecoder(fp)
	if err := dec.Decode(&ServiceConf); err != nil {
		return err
	}

	if err = checkMust(confFile); err != nil {
		return err
	}

	return nil
}

func DumpConf() {
	fmt.Printf("conf: %v\n", ServiceConf)
	fmt.Printf("TZ time location: %v\n", Loc)
}

func checkMust(confFile string) error {
	if ServiceConf.ListenPort <= 0 {
		return fmt.Errorf("listen-port expected in conf")
	}
	if len(ServiceConf.APIBaseURL) == 0 {
		return fmt.Errorf("api-base-url expected in conf")
	}
	if len(ServiceConf.DefaultTimeout) == 0 {
		ServiceConf.defaultTimeout = 30 * time.Second
	} else {
		d, err := time.ParseDuration(ServiceConf.DefaultTimeout)
		if err != nil {
			return fmt.Errorf("failed to parse default-timeout: %v\n", err)
		}
		ServiceConf.defaultTimeout = d
	}

	apps := ServiceConf.Apps
	if len(apps) == 0 {
		return fmt.Errorf("apps expected in conf")
	}
	for i, _ := range apps {
		appConf := &apps[i]
		if err := appConf.checkMust(i); err != nil {
			return err
		}
		appConfs[appConf.Name] = appConf
	}

	ce := &ServiceConf.CommonEndpoints
	if len(ce.HealthCheck) == 0 {
		return fmt.Errorf("common-endpoints/health-check expected in conf")
	}
	if len(ce.Chat) == 0 {
		return fmt.Errorf("common-endpoints/chat expected in conf")
	}
	if len(ce.ChatCompletions) == 0 {
		return fmt.Errorf("common-endpoints/chat-completions expected in conf")
	}
	if len(ce.Image) == 0 {
		return fmt.Errorf("common-endpoints/image expected in conf")
	}
	if len(ce.ContextChat) == 0 {
		return fmt.Errorf("common-endpoints/context-chat expected in conf")
	}

	return nil
}

func GetAppConf(appName string) *AppParams {
	if appConf, ok := appConfs[appName]; ok {
		return appConf
	}
	return nil
}

func GetDefaultTimeout() time.Duration {
	return ServiceConf.defaultTimeout
}
