package converter

import (
	"encoding/json"
	"gin-shop-api/internal/helpers/logging"
	"gin-shop-api/internal/helpers/types"
	"log"
	"strconv"
)

func ConvertJsonToDict(jsonString string) types.Dict {
	defer logging.HandlePanic()

	var result types.Dict
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		msg := "JSON was not converted"
		log.Panicf("%s: %s", msg, err)
	}

	return result
}

func ConvertObjectToJson(data any) string {
	defer logging.HandlePanic()

	r, err := json.Marshal(data)
	if err != nil {
		msg := "JSON was not converted"
		log.Panicf("%s: %s", msg, err)
	}

	return string(r)
}

func ConvertByteToObject(data []byte) types.Dict {
	defer logging.HandlePanic()

	var result types.Dict
	err := json.Unmarshal(data, &result)
	if err != nil {
		msg := "JSON was not converted"
		log.Panicf("%s: %s", msg, err)
	}

	return result
}

func ConvertDictToByte(data types.Dict) []byte {
	defer logging.HandlePanic()

	result, err := json.Marshal(data)
	if err != nil {
		msg := "JSON was not converted"
		log.Panicf("%s: %s", msg, err)
	}

	return result
}

func CastStringToInt(stringValue string) int {
	defer logging.HandlePanic()

	stringToIntValue, err := strconv.Atoi(stringValue)
	if err != nil {
		msg := "Failed to convert string value to integer"
		log.Panicf("%s: %s", msg, err)
	}
	return stringToIntValue
}
