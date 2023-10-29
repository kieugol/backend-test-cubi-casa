package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

type ValidationError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

func ShouldBindHeader(c *gin.Context) bool {
	platform := c.Request.Header.Get("X-PLATFORM")
	deviceType := c.Request.Header.Get("X-DEVICE-TYPE")
	deviceId := c.Request.Header.Get("X-DEVICE-ID")
	lang := c.Request.Header.Get("X-LANG")
	channel := c.Request.Header.Get("X-CHANNEL")

	if platform == "" || deviceType == "" || deviceId == "" || lang == "" || channel == "" {
		return false
	}

	return true
}

func LogPrint(jsonData interface{}) {
	prettyJSON, _ := json.MarshalIndent(jsonData, "", "")
	fmt.Printf("%s\n", strings.ReplaceAll(string(prettyJSON), "\n", ""))
}

func GetStructName(x any) string {
	return reflect.TypeOf(x).Name()
}

func IsEmptyStruct(x, y any) bool {
	return reflect.DeepEqual(x, y)
}

func ParseJSON(data []byte, target interface{}) error {
	if err := json.Unmarshal(data, target); err != nil {
		return err
	}
	return nil
}

func ReadFile(path string) []byte {
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	value, _ := ioutil.ReadAll(jsonFile)

	return value
}
