package websocket

const (
	AuthComplete = "AuthComplete"
)

type ReadMessage struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

type WriteMessage struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

// 位置情報メッセージ
type GeoMessage struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
