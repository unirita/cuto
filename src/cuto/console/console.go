// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package console

import (
	"fmt"
	"os"

	"cuto/log"
)

// USAGE表示用の定義メッセージ
const USAGE = `Usage :
    master.exe [-v] [-n Jobnetwork] [-s] [-c ConfigFile]

Option :
    -v            :   Print master version.
    -n bpmn name  :   Designate a bpmn file name.(Without extensions.)
    -s            :   Execute Jobnetwork.
    -c ConfigFile :   Designate config file path.
                      If it is omitted, '<Current Directory>/master.ini' will be used.

Copyright 2015 unirita Inc.
`

var stack_msg = []string{"CTM019E", "CTS019E", "CTU003E"}

// showユーティリティのUSAGE表示用の定義メッセージ
const USAGE_SHOW = `Usage :
    show.exe [-v] [-flow="bpmn file name"] [-From="From date"] [-to="To date"] [-status="normal" | "abnormal" | "running"] [-format="json" | "csv"]

Option :
    -v                 :   Print master version.
    -flow=flow name    :   Designate a bpmn name.
    -from=yyyymmdd     :   From date is designated.
    -to=yyyymmdd       :   To date is designated.
    -status=normal     :   Status indicates only something of NORMAL-END.
    -status=abnormal   :   Status indicates only something of ABNORMAL-END.
    -status=running    :   Status indicates only something of RUNNING.
    -format=json       :   It outputs by the form of JSON.
    -format=csv        :   It outputs by the form of CSV.
    
When omitting [-from] and [-to], only Jobnetwork begun today is indicated.
	
Copyright 2015 unirita Inc.
`

// コンソールメッセージ一覧
var msgs = map[string]string{
	"CTM001I": "CUTO MASTER STARTED. PID [%d] VERSION [%v]",
	"CTM002I": "CUTO MASTER ENDED. RC [%d]",
	"CTM003E": "INVALID ARGUMENT.",
	"CTM004E": "FAILED TO READ EXPAND JOB CONFIG FILE [%s].",
	"CTM005E": "CONFIG PARM IS NOT EXACT FORMAT. REASON [%s]",
	"CTM006W": "INVALID VALUE [%s / %s]. USE DEFAULT VALUE [%s]",
	"CTM007I": "BIND ADDRESS[%v] USED PORT[%v].",
	"CTM008I": "REMINDER: EVALUATION PERIOD WILL EXPIRE IN [%d] DAYS.",
	"CTM009I": "LICENSE VALID TILL %04d/%02d/%02d.",
	"CTM010E": "FAILED TO READ BPMN FILE [%s].",
	"CTM011E": "[%s] IS NOT EXACT FORMAT. REASON [%s]",
	"CTM012I": "[%s] STARTED. INSTANCE [%d]",
	"CTM013I": "[%s] ENDED. INSTANCE [%d] STATUS [%s]",
	"CTM014I": "SENDED REQUEST [%s] START MESSAGE. ID [%s] NODE [%s] PORT [%d]",
	"CTM015I": "RECEIVED RESPONSE [%s] END MESSAGE. ID [%s] STATUS [%s]",
	"CTM016W": "UNABLE TO CONNECT SERVANT. NODE[%s] PORT[%d] EC [%d] MSG [%s]",
	"CTM017W": "UNABLE TO SEND MESSAGE. NODE[%s] PORT[%d] EC [%d] MSG [%s]",
	"CTM018W": "TIME OUT WAITING JOB EXECUTE.",
	"CTM019E": "EXCEPTION OCCURED - %s",
	"CTM020I": "BPMN FILE [%s] IS VALID.",
	"CTM021E": "COULD NOT INITIALIZE LOGGER. REASON[%s]",
	"CTM022I": "JOB [%s] IS RUNNING FOR %d MINUTES.",
	"CTM023I": "JOB [%s] STARTED. INSTANCE [%d] JOBID [%s].",
	"CTM024I": "JOB [%s] ENDED. INSTANCE [%d] JOBID [%s] STATUS [%d].",
	"CTM025W": "JOB [%s] ABNORMAL ENDED. INSTANCE [%d] JOBID [%s] STATUS [%d] DETAIL [%s].",
	"1":       "",
	"CTS001I": "CUTO SERVANT STARTED. PID [%v] VERSION [%s]",
	"CTS002I": "CUTO SERVANT ENDED. RC [%d].",
	"CTS003E": "INVALID ARGUMENT.",
	"CTS004W": "FAILED TO READ CONFIG FILE [%s]. USE DEFAULT VALUE",
	"CTS005E": "CONFIG PARM IS NOT EXACT FORMAT. REASON [%s]",
	"CTS006W": "CONFIG PARM [%s / %s] USE DEFAULT VALUE [%v]",
	"CTS007E": "LICENSE VERIFICATION FAILED.",
	"CTS008I": "REMINDER: EVALUATION PERIOD WILL EXPIRE IN [%d] DAYS.",
	"CTS009I": "LICENSE VALID TILL %s/%s/%s.",
	"CTS010I": "JOB [%s] STARTED. INSTANCE [%d] ID [%s] PID [%d].",
	"CTS011I": "JOB [%s] ENDED. INSTANCE [%d] ID [%s] STATUS [%d] RC [%d].",
	"CTS012E": "NOT FOUND JOB FILE [%s].",
	"CTS013E": "ENABLE TO EXECUTE JOB [%s].",
	"CTS014I": "RECEIVED REQUEST MESSAGE.",
	"CTS015E": "RECEIVED MESSAGE IS NOT EXACT FORMAL.DETAIL[%s]",
	"CTS016I": "EXECUTION OF JOB, I STAND BY. INSTANCE [%d] ID [%s] NAME [%s].",
	"CTS017I": "SENDED RESPONSE MESSAGE.",
	"CTS018E": "ENABLE TO SEND RESPONSE MESSAGE.",
	"CTS019E": "EXCEPTION OCCURED - %s",
	"CTS020I": "SERVANT CHILD STARTED. PID[%d]",
	"CTS021I": "SERVANT CHILD ENDED. RC[%d]",
	"CTS022E": "UNABLE TO OUTPUT JOBLOG. MSG[%s]",
	"CTS023E": "COULD NOT INITIALIZE LOGGER. REASON[%s]",
	"2":       "",
	"CTU001I": "SHOW UTILITY STARTED. VERSION [%v]",
	"CTU002I": "SHOW UTILITY ENDED. RC [%d].",
	"CTU003E": "EXCEPTION OCCURED - %v",
	"CTU004E": "AN INTERNAL ERROR OCCURRED. - %v",
	"CTU005W": "FAILED TO JOB INFORMATION NID[%v]. - %v",
	"CTU006E": "NOT FOUND CONFIG FILE. - %v",
}

// 標準出力へメッセージコードcodeに対応したメッセージを表示する。
//
// param : code メッセージコードID。
//
// param : a... メッセージの書式制御文字に渡す内容。
//
// return : 出力文字数。
//
// return : エラー情報。
func Display(code string, a ...interface{}) (int, error) {
	msg := GetMessage(code, a...)
	log.Info(msg)

	for _, s := range stack_msg {
		if code == s {
			PrintStack()
		}
	}
	return fmt.Println(msg)
}

// 標準エラー出力へメッセージコードcodeに対応したメッセージを表示する。
//
// param : code メッセージコードID。
//
// param : a... メッセージの書式制御文字に渡す内容。
//
// return : 出力文字数。
//
// return : エラー情報。
func DisplayError(code string, a ...interface{}) (int, error) {
	msg := GetMessage(code, a...)
	log.Error(msg)

	for _, s := range stack_msg {
		if code == s {
			PrintStack()
		}
	}
	return fmt.Fprintln(os.Stderr, msg)
}

// 出力メッセージを文字列型で取得する。
//
// param : code メッセージコードID。
//
// param : a... メッセージの書式制御文字に渡す内容。
//
// return : 取得したメッセージ
func GetMessage(code string, a ...interface{}) string {
	return fmt.Sprintf("%s %s", code, fmt.Sprintf(msgs[code], a...))
}
