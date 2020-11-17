// Manages the first component: collecting the configuration
// there are three different sources:
//   1. conf.d files              (lowest priority)
//   2. Local conf file
//   3. Remote configuration
//   4. CLI arguments             (highest priority)
package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

const confdPath = "./conf.d"
const localConfPath = "./conf.yml"

type Config struct {
	FrameworkSettings *FrameworkSettings     `mapstructure:"framework_settings"`
	ArgumentSets      map[string]interface{} `mapstructure:"checks_arguments"`
}

type FrameworkSettings struct {
	CmtNode  string `mapstructure:"cmt_node"`
	CmtGroup string `mapstructure:"cmt_group"`

	Checks                 []string
	GraylogUDPGelfServers  []*UDPGelfAddress  `mapstructure:"graylog_udp_gelf_servers"`
	GraylogHTTPGelfServers []*HTTPGelfAddress `mapstructure:"graylog_http_gelf_servers"`

	TeamsChannel   []*TeamsAddress `mapstructure:"teams_channel"`
	TeamsRateLimit int             `mapstructure:"teams_rate_limit"`

	DatabaseFile string `mapstructure:"database_file"`
}

type UDPGelfAddress struct {
	Name string
	Host string
	Port int
}

type HTTPGelfAddress struct {
	Name string
	URL  string
}

type TeamsAddress struct {
	Name string
	URL  string
}

func loadConf() Config {
	base := make(map[string]interface{})
	loadConfInPlaceFromConfd(base)
	loadConfInPlaceFromFile(localConfPath, base)
	loadConfInPlaceFromRemote(base)
	loadConfFromArguments(os.Args[1:])

	var adapted = make(map[string]map[string]interface{})
	adapted["framework_settings"] = make(map[string]interface{})
	adapted["checks_arguments"] = make(map[string]interface{})

	t := reflect.TypeOf(FrameworkSettings{})
	frameworkSettingsKeyNames := make([]string, t.NumField())
	for i := 0; i < len(frameworkSettingsKeyNames); i++ {
		frameworkSettingsKeyNames[i] = t.Field(i).Tag.Get("mapstructure")
		if frameworkSettingsKeyNames[i] == "" {
			frameworkSettingsKeyNames[i] = strings.ToLower(t.Field(i).Name)
		}
	}

	for key, value := range base {

		isFrameworkKey := false
		for _, k := range frameworkSettingsKeyNames {
			if k == key {
				isFrameworkKey = true
				break
			}
		}

		if isFrameworkKey {
			adapted["framework_settings"][key] = value
		} else {
			adapted["checks_arguments"][key] = value
		}

	}

	var conf Config
	if err := mapstructure.Decode(adapted, &conf); err != nil {
		// panic because that means we don't have any conf here
		panic(err)
	}

	return conf
}

func loadConfInPlaceFromConfd(base map[string]interface{}) {
	files, err := ioutil.ReadDir(confdPath)
	if err != nil {
		log.Printf("[conf] couldn't list directory conf.d, ignoring: %s", err)
	}

	// sort alphabetically: NOT numerically aware
	// 2-foo.yml comes after 10-bar.yml
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		loadConfInPlaceFromFile(filepath.Join(confdPath, file.Name()), base)
	}
}

func loadConfInPlaceFromFile(filename string, base map[string]interface{}) {
	f, err := os.Open(filename)
	if err != nil {
		log.Printf("[conf] couldn't open %s conf file: %s\n", filename, err)
		return
	}
	defer f.Close()

	conf, err := readConf(f)
	if err != nil {
		log.Printf("[conf] couldn't read %s conf: %s\n", filename, err)
		return
	}
	mergeConfInBase(base, conf)
}

func loadConfInPlaceFromRemote(base map[string]interface{}) {
	log.Println("[conf] TODO: load configuration from remote")
}

func loadConfFromArguments(args []string) {
	log.Println("[conf] TODO: fancy config from argv: dictionary.key.0=10 list.1 = 2")
}

func readConf(r io.Reader) (map[string]interface{}, error) {
	var conf map[string]interface{}
	if err := yaml.NewDecoder(r).Decode(&conf); err != nil {
		return nil, err
	}
	return conf, nil
}

// mergeConfInBase adds the key/value pairs from extend to base (in place). If
// the key already exists in base, and the value is a list in both the base and
// the extend, then it is concatenated. ONLY FIRST LEVEL LIST ARE CONCATENATED
// LIKE THIS. Others are replaced (as per doc)
func mergeConfInBase(base map[string]interface{}, extend map[string]interface{}) {
	for key, value := range extend {
		base_list, base_ok := base[key].([]interface{})
		extend_list, extend_ok := extend[key].([]interface{})
		if base_ok && extend_ok {
			base[key] = append(base_list, extend_list...)
		} else {
			base[key] = value
		}
	}
}
