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
	"sort"

	"gopkg.in/yaml.v2"
)

const confd_path = "./conf.d"
const local_conf_path = "./conf.yml"

func loadConf() map[string]interface{} {
	base := make(map[string]interface{})
	loadConfInPlaceFromConfd(base)
	loadConfInPlaceFromFile(local_conf_path, base)
	loadConfInPlaceFromRemote(base)
	loadConfFromArguments(os.Args[1:])
	return base
}

func loadConfInPlaceFromConfd(base map[string]interface{}) {
	files, err := ioutil.ReadDir(confd_path)
	if err != nil {
		log.Printf("couldn't list directory conf.d, ignoring: %s", err)
	}

	// sort alphabetically: NOT numerically aware
	// 2-foo.yml comes after 10-bar.yml
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		loadConfInPlaceFromFile(filepath.Join(confd_path, file.Name()), base)
	}
}

func loadConfInPlaceFromFile(filename string, base map[string]interface{}) {
	f, err := os.Open(filename)
	if err != nil {
		log.Printf("couldn't open %s conf file: %s\n", filename, err)
		return
	}
	defer f.Close()

	conf, err := readConf(f)
	if err != nil {
		log.Printf("couldn't read %s conf: %s\n", filename, err)
		return
	}
	mergeConfInBase(base, conf)
}

func loadConfInPlaceFromRemote(base map[string]interface{}) {
	log.Println("TODO: load configuration from remote")
}

func loadConfFromArguments(args []string) {
	log.Println("TODO: do something fancy like: dictionary.key.0=10 list.1 = 2")
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
