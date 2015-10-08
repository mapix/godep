package vcs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

var customVCSPaths []*vcsPath
var verboseMirror bool

type customVCSPath struct {
	Prefix string `json:"prefix"`
	Re     string `json:"re"`
	Repo   string `json:"repo"`
	Vcs    string `json:"vcs"`
}

func init() {
	if os.Getenv("GO_GET_VERBOSE_MIRROR") != "" {
		verboseMirror = true
	}
	configFile := os.Getenv("GO_GET_MIRROR_CONFIG_FILE")
	if configFile == "" {
		currentDir, _ := os.Getwd()
		configFile = path.Join(currentDir, "go-get-config.json")
		mirrorPrintf("try load custom go-get-config from: %v", configFile)
	} else {
		mirrorPrintf("try use custom go-get-config location : %v", configFile)
	}
	if _, errStat := os.Stat(configFile); errStat == nil {
		if configContent, errRead := ioutil.ReadFile(configFile); errRead == nil {
			var objmap map[string]*json.RawMessage
			if errUnmarshal := json.Unmarshal(configContent, &objmap); errUnmarshal == nil {
				var tmpVCSPaths []*customVCSPath
				if errMetas := json.Unmarshal(*objmap["Metas"], &tmpVCSPaths); errMetas == nil {
					for _, tmpVCSPath := range tmpVCSPaths {
						customVCSPaths = append(customVCSPaths, &vcsPath{tmpVCSPath.Prefix, tmpVCSPath.Re, tmpVCSPath.Repo, tmpVCSPath.Vcs, nil, false, nil})
					}
				} else {
					mirrorErrorf("failed unmarshal vcs metas from %v : %s", configFile, errMetas)
				}
			} else {
				mirrorErrorf("failed unmarshal go-get-config from %v : %s", configFile, errUnmarshal)
			}
		} else {
			mirrorErrorf("failed read go-get-config from %v : %s", configFile, errRead)
		}
	} else {
		mirrorPrintf("failed load go-get-config from %v : %s", configFile, errStat)
	}
}

func getCustomVCSPaths() []*vcsPath {
	return customVCSPaths
}

func mirrorErrorf(format string, args ...interface{}) {
	if verboseMirror {
		fmt.Fprintf(os.Stderr, "cfg: "+format+"\n", args...)
	}
}

func mirrorPrintf(format string, args ...interface{}) {
	if verboseMirror {
		fmt.Fprintf(os.Stdout, "cfg: "+format+"\n", args...)
	}
}
