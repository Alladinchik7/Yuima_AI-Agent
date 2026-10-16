package todo

type checkAddTodoArgs struct {
	Task string `json:"task"`
	Confirmed bool `json:"confirmed"`
}

type checkTaskID struct {
	TaskID int `json:"task_id"`
	Confirmed bool `json:"confirmed"`
}