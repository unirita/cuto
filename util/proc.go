package util

// IsProcessExists checks if process with pid running or not running.
func IsProcessExists(pid int) bool {
	return isProcessExists(pid)
}
