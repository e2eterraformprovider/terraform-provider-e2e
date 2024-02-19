package client

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func RemoveExtraKeysLoadBalancer(buf *bytes.Buffer) (bytes.Buffer, error) {
	jsonData := buf.Bytes()
	var data map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		return *buf, err
	}
	enableEosLogger, ok := data["enable_eos_logger"].(map[string]interface{})
	if !ok {
		return *buf, nil
	}
	accessKey, ok := enableEosLogger["access_key"].(string)
	if ok && (len(accessKey) != 0) {
		return *buf, nil
	}
	delete(data, "enable_eos_logger")
	NewjsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding data:", err)
		return *buf, nil
	}
	newBuffer := bytes.NewBuffer(NewjsonData)
	return *newBuffer, nil
}
