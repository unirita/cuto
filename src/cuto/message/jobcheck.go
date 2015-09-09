package message

import (
	"encoding/json"
	"fmt"
)

// ジョブ正常終了確認メッセージ。
type JobCheck struct {
	Type    string `json:"type"`
	Version string `json:"version"`
	NID     int    `json:"nid"`
	JID     string `json:"jid"`
}

const jobCheckMessageType = "jobcheck"

// ジョブ正常終了確認JSONメッセージをパースし、JobCheckオブジェクトのメンバをセットする。
//
// param : message 受信メッセージ文字列
func (j *JobCheck) ParseJSON(message string) error {
	byteMessage := []byte(message)
	err := json.Unmarshal(byteMessage, j)
	if err != nil {
		return err
	}
	if j.Type != jobCheckMessageType {
		return fmt.Errorf("Invalid message type.")
	}
	return nil
}

// JobCheckオブジェクトの値を元に、ジョブ正常終了確認JSONメッセージを生成する
//
// return : JSONメッセージフォーマットの文字列。
func (j JobCheck) GenerateJSON() (string, error) {
	j.Type = jobCheckMessageType
	j.Version = MasterVersion
	byteMessage, err := json.Marshal(j)
	if err != nil {
		return ``, err
	}
	return string(byteMessage), nil
}
