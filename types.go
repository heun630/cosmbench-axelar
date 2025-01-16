package main

type LogEntry struct {
	TxIdx     int    `json:"txIdx"`
	Timestamp int64  `json:"timestamp"`
	TxHash    string `json:"txHash"`
	Height    int    `json:"height,omitempty"`
}
