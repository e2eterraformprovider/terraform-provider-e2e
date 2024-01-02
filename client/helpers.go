package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

func RemoveExtraKeysLoadBalancer(buf *bytes.Buffer) (bytes.Buffer, error) {
	jsonData := buf.Bytes()
	var data map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		return *buf, err
	}
	log.Println("===========================REMOVE EXTRA KEY FROM LB=====================================")
	enableEosLogger, ok := data["enable_eos_logger"].(map[string]interface{})
	log.Println("===============ENABLE EOS LOGGER================", enableEosLogger)
	if !ok {
		return *buf, nil
	}
	accessKey, ok := enableEosLogger["access_key"].(string)
	log.Println("==========================ACCESS KEY VALUE==================", accessKey)
	log.Println(ok, len(accessKey))
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
	log.Println("======================NEW BUFFER===================", newBuffer.String())
	return *newBuffer, nil
}
