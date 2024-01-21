package log

import (
	"encoding/json"
	"regexp"
	"strconv"
)

type Log struct {
	Timestamp   int64
	Context     string
	LogType     LOG_TYPE
	Environment ENVIRONMENT
	Service     string
	Feature     string
	Code        string
	Message     string
	Data        map[string]interface{}
}

type ENVIRONMENT string

const (
	PRODUCTION  ENVIRONMENT = "PRODUCTION"
	DEVELOPMENT ENVIRONMENT = "DEVELOPMENT"
	STAGING     ENVIRONMENT = "STAGING"
	LOCAL       ENVIRONMENT = "LOCAL"
)

type LOG_TYPE string

const (
	ERROR   LOG_TYPE = "ERROR"
	INFO    LOG_TYPE = "INFO"
	SUCCESS LOG_TYPE = "SUCCESS"
)

func CreateLogFromJson(jsonStr string) string {
	var log Log
	if err := json.Unmarshal([]byte(jsonStr), &log); err != nil {
		panic(err)
	}

	dataString, err := json.Marshal(log.Data)
	if err != nil {
		panic(err)
	}
	return "\nLOG " + strconv.FormatInt(log.Timestamp, 10) + " " + string(log.LogType) + " " + string(log.Environment) + " " + log.Context + " " + log.Service + " " + log.Feature + " " + log.Code + ` "` + log.Message + `" "` + string(dataString[:]) + `"`
}

func ParseLogFromString(str string) Log {
	nlre := regexp.MustCompile(`\r?\n`)
	wsre := regexp.MustCompile(`\s+`)

	timestampRegex := regexp.MustCompile(`[0-9]{10}`)
	propertiesRegex := regexp.MustCompile(`\S+`)
	logTypeRegex := regexp.MustCompile(`(ERROR|SUCCESS|INFO)`)
	environmentRegex := regexp.MustCompile(`(PRODUCTION|DEVELOPMENT|STAGING|LOCAL)`)
	messageRegex := regexp.MustCompile(`"(.*?)"`)
	dataRegex := regexp.MustCompile(`"{(.*?)}"`)

	properties := propertiesRegex.FindAllString(str, -1)
	timestamp := timestampRegex.FindString(str)
	logType := logTypeRegex.FindString(str)
	environment := environmentRegex.FindString(str)
	message := messageRegex.FindString(str)
	data := dataRegex.FindString(str)
	context := properties[4]
	service := properties[5]
	feature := properties[6]
	code := properties[7]

	parsedTimestamp, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		panic(err)
	}
	var dataMap map[string]interface{}
	data = nlre.ReplaceAllString(data, "")
	data = wsre.ReplaceAllString(data, "")
	data = data[1:]
	data = data[:len(data)-1]

	if err := json.Unmarshal([]byte(data), &dataMap); err != nil {
		panic(err)
	}

	var log = Log{
		Timestamp:   parsedTimestamp,
		Context:     context,
		LogType:     LOG_TYPE(logType),
		Environment: ENVIRONMENT(environment),
		Service:     service,
		Feature:     feature,
		Code:        code,
		Message:     message,
		Data:        dataMap,
	}

	return log

}
