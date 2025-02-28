package cmd

// Event はデータセット内の各記念日イベントを表現する構造体です。
type Event struct {
	ID        int    `json:"id" yaml:"id"`
	Date      string `json:"date" yaml:"date"`           // "MM-DD"形式で記録（例: "02-22"）
	Frequency string `json:"frequency" yaml:"frequency"` // 毎年繰り返すイベントの場合、"yearly" と記録
	Title     string `json:"title" yaml:"title"`
}
