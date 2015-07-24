package network

import "encoding/json"

type Result struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	PID     int           `json:"pid"`
	Network networkResult `json:"network"`
}

type networkResult struct {
	Instance int    `json:"instance"`
	Name     string `json:"name"`
}

const encodeErrorResult = `{
    "status":2,
    "message":"Result message encode error.",
    "pid":0,
	"network":{"instance":0,"name":""}
}`

const (
	status_SUCCESS = iota
	status_MASTER_ERROR
	status_REALTIME_ERROR
)

// SuccessResult generates success result message as JSON.
func SuccessResult(pid int, instanceID int, networkName string) string {
	result := new(Result)
	result.Status = status_SUCCESS
	result.Message = "Success."
	result.PID = pid
	result.Network.Instance = instanceID
	result.Network.Name = networkName
	return result.Encode()
}

// MasterErrorResult generates master error result message as JSON.
func MasterErrorResult(msg string, pid int) string {
	result := new(Result)
	result.Status = status_MASTER_ERROR
	result.Message = msg
	result.PID = pid
	return result.Encode()
}

// RealtimeErrorResult generates realtime error result message as JSON.
func RealtimeErrorResult(msg string) string {
	result := new(Result)
	result.Status = status_REALTIME_ERROR
	result.Message = msg
	return result.Encode()
}

// Encode generates JSON string from object data.
func (r *Result) Encode() string {
	b, err := json.Marshal(r)
	if err != nil {
		return encodeErrorResult
	}

	return string(b)
}
