package models

type KeyGetListResponse []struct {
	Key Key `json:"key"`
}
type Key struct {
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
	Type        string `json:"type"`
	Size        int    `json:"size"`
	Data        string `json:"data"`
}
