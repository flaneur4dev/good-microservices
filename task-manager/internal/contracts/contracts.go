package contracts

type Task struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Timestamp   int64  `json:"timestamp"`
}
