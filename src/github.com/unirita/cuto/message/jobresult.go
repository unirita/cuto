package message

import (
	"encoding/json"
	"fmt"
)

type JobResult struct {
	Type    string `json:"type"`
	Version string `json:"version"`
	NID     int    `json:"nid"`
	JID     string `json:"jid"`
	RC      int    `json:"rc"`
	Stat    int    `json:"stat"`
	Var     string `json:"var"`
	St      string `json:"st"`
	Et      string `json:"et"`
}

const jobResultMessageType = "jobresult"

// ジョブ正常終了確認JSONメッセージをパースし、JobCheckオブジェクトのメンバをセットする。
//
// param : message 受信メッセージ文字列
func (j *JobResult) ParseJSON(message string) error {
	byteMessage := []byte(message)
	err := json.Unmarshal(byteMessage, j)
	if err != nil {
		return err
	}
	if j.Type != jobResultMessageType {
		return fmt.Errorf("Invalid message type.")
	}
	return nil
}

// JobCheckオブジェクトの値を元に、ジョブ正常終了確認JSONメッセージを生成する
//
// return : JSONメッセージフォーマットの文字列。
func (j JobResult) GenerateJSON() (string, error) {
	j.Type = jobResultMessageType
	j.Version = ServantVersion
	byteMessage, err := json.Marshal(j)
	if err != nil {
		return ``, err
	}
	return string(byteMessage), nil
}
