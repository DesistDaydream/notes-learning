package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

func main() {
	// 读取文件内容
	file, err := os.Open("/mnt/d/Projects/DesistDaydream/notes-learning/content/zh-cn/.obsidian/plugins/obsidian-bartender/data.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	// 解析 JSON 数据
	var jsonMap map[string]interface{}
	err = json.Unmarshal(data, &jsonMap)
	if err != nil {
		panic(err)
	}

	// 对 "fileExplorerOrder" 排序
	fileExplorerOrder := jsonMap["fileExplorerOrder"].(map[string]interface{})
	sortedKeys := make([]string, 0, len(fileExplorerOrder))
	for key := range fileExplorerOrder {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	sortedFileExplorerOrder := make(map[string]interface{}, len(fileExplorerOrder))
	for _, key := range sortedKeys {
		sortedFileExplorerOrder[key] = fileExplorerOrder[key]
	}
	jsonMap["fileExplorerOrder"] = sortedFileExplorerOrder

	// 将处理后的数据输出到文件
	newData, err := json.MarshalIndent(jsonMap, "", "  ")
	if err != nil {
		panic(err)
	}

	file, err = os.Create("/mnt/d/Projects/DesistDaydream/notes-learning/content/zh-cn/.obsidian/plugins/obsidian-bartender/data_new.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write(newData)
	if err != nil {
		panic(err)
	}

	fmt.Println("JSON 文件处理完成，处理结果输出到 json_file_new.json")
}
