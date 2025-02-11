package tasks

type Task[Data any] struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Data        Data   `json:"message"`
}
