// global conf
// ENV:
//   CONF_FILE      --- configuration file name
//   TZ             --- Time Zone name, e.g. "Asia/Shanghai"
//
// YAML
// ---
// listen-host: ""
// http-listen-port: 7080
// jsonl-rpc-listen-prot: 7081
// q-len: 5
// base32-chars: "abcd2efgh3ijkl4mnop5qrst6uvwx7yz"
//
// Rosbit Xu
package conf

import (
	"gopkg.in/yaml.v2"
	"fmt"
	"os"
	"time"
	"path"
)

type DBProxyConf struct {
	ListenHost      string `yaml:"listen-host"`
	HttpListenPort     int `yaml:"http-listen-port"`
	JSONLRpcListenPort int `yaml:"jsonl-rpc-listen-prot"`
	QLen               int `yaml:"q-len"`
	Base32Chars     string `yaml:"base32-chars"`
}

var (
	ServiceConf DBProxyConf
	Loc = time.FixedZone("UTC+8", 8*60*60)
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
	if err = dec.Decode(&ServiceConf); err != nil {
		return err
	}

	if err = checkMust(confFile); err != nil {
		return err
	}

	return nil
}

func checkMust(confFile string) error {
	// confRoot := path.Dir(confFile)

	if ServiceConf.HttpListenPort <= 0 {
		return fmt.Errorf("http-listen-port expected in conf")
	}
	if ServiceConf.JSONLRpcListenPort <= 0 {
		return fmt.Errorf("jsonl-rpc-listen-port expected in conf")
	}
	if ServiceConf.HttpListenPort == ServiceConf.JSONLRpcListenPort {
		return fmt.Errorf("http-listen-port and jsonl-rpc-listen-port must be different")
	}

	if ServiceConf.QLen <= 1 {
		// 1会导致死锁
		ServiceConf.QLen = 5
	}

	return nil
}

func DumpConf() {
	fmt.Printf("conf: %#v\n", ServiceConf)
	fmt.Printf("TZ time location: %v\n", Loc)
}

func checkDir(path, prompt string) error {
	if fi, err := os.Stat(path); err != nil {
		return err
	} else if !fi.IsDir() {
		return fmt.Errorf("%s %s is not a directory", prompt, path)
	}
	return nil
}

func toAbsPath(absRoot, filePath string) string {
	if path.IsAbs(filePath) {
		return filePath
	}
	return path.Join(absRoot, filePath)
}
