package main

import (
	"encoding/json"
	"github.com/lenistwo/conf"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

const (
	ConfigFileName  = "config.json"
	DiscordCacheDir = "\\discord\\Cache"
	BackSlash       = "\\"
	Space           = " "
	Dot             = "."
	EmptyString     = ""
	Underscore      = "_"
	AppDataEnvVar   = "APPDATA"
	AllInDir        = 0
	DateTimePattern = "2006-04-02 15_04_05"
)

var (
	appdataPath string
	config      conf.Config
	wg          sync.WaitGroup
)

func init() {
	loadConfig()
	createOutputDirIfNotExist()
	if !isFileExtensionValid() {
		panic("File Extension Should Start With '.'")
	}
	appdataPath = os.Getenv(AppDataEnvVar)
}

func main() {
	open, err := os.Open(appdataPath + DiscordCacheDir)
	checkError(err, "Cannot Open Discord Cache Folder")
	defer open.Close()
	cacheFiles, err := open.Readdir(AllInDir)
	wg.Add(len(cacheFiles))
	for _, fileInfo := range cacheFiles {
		go writeFile(fileInfo)
	}
	wg.Wait()
}

func writeFile(fileInfo os.FileInfo) {
	baseFile, err := os.Open(appdataPath + DiscordCacheDir + BackSlash + fileInfo.Name())
	checkError(err, "Cannot Open Discord File In Directory")
	destFile, err := os.Create(config.OutputPath + BackSlash + prepareFileName(fileInfo) + config.FileExtension)
	checkError(err, "Cannot Create File In Directory")
	_, _ = io.Copy(destFile, baseFile)
	wg.Done()
}

func prepareFileName(info os.FileInfo) string {
	if !config.WithModificationTime {
		return info.Name()
	}
	modificationTime := info.ModTime().Format(DateTimePattern)
	validFileName := replaceSpaces(info.Name()) + Underscore + modificationTime
	return validFileName
}

func checkError(err error, message string) {
	if err != nil {
		panic(message)
	}
}

func loadConfig() {
	file, err := ioutil.ReadFile(ConfigFileName)
	checkError(err, "Error During Loading Config File")
	err = json.Unmarshal(file, &config)
	checkError(err, "Error During Unmarshaling Config File")
}

func createOutputDirIfNotExist() {
	if _, err := os.Stat(config.OutputPath); os.IsNotExist(err) {
		err = os.Mkdir(config.OutputPath, os.ModeDir)
		checkError(err, "Error During Creating Output Directory")
	}
}

func isFileExtensionValid() bool {
	return strings.HasPrefix(config.FileExtension, Dot)
}

func replaceSpaces(s string) string {
	return strings.ReplaceAll(s, Space, EmptyString)
}
