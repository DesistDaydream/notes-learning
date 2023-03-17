package main

import (
	"encoding/json"
	"io/ioutil"
	"sort"
)

type Config struct {
	StatusBarOrder   []interface{}       `json:"statusBarOrder"`
	RibbonBarOrder   []interface{}       `json:"ribbonBarOrder"`
	FileExplorerData map[string][]string `json:"fileExplorerOrder"`
	ActionBarOrder   struct{}            `json:"actionBarOrder"`
	AutoHide         bool                `json:"autoHide"`
	AutoHideDelay    int                 `json:"autoHideDelay"`
	DragDelay        int                 `json:"dragDelay"`
}

func main() {
	srcFile := "/mnt/d/Projects/DesistDaydream/notes-learning/content/zh-cn/.obsidian/plugins/obsidian-bartender/data.json"
	// dstFile := "/mnt/d/Projects/DesistDaydream/notes-learning/content/zh-cn/.obsidian/plugins/obsidian-bartender/data_new.json"

	fileByte, err := ioutil.ReadFile(srcFile)
	if err != nil {
		panic(err)
	}

	var config Config
	err = json.Unmarshal(fileByte, &config)
	if err != nil {
		panic(err)
	}

	sortedKeys := make([]string, 0, len(config.FileExplorerData))
	for key := range config.FileExplorerData {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	sortedFileExplorerData := make(map[string][]string, len(config.FileExplorerData))
	for _, key := range sortedKeys {
		sortedFileExplorerData[key] = config.FileExplorerData[key]
	}
	config.FileExplorerData = sortedFileExplorerData

	newData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(srcFile, newData, 0666)
	if err != nil {
		panic(err)
	}
}
