package scheduler

type TaskType string

const (
	ProcessAmericanTournamentsTaskType          = "process_american_tournaments_task"
	ProcessNotUrgentEuropeanTournamentsTaskType = "process_not_urgent_european_tournaments_task"
	ProcessUrgentEuropeanTournamentsTaskType    = "process_urgent_european_tournaments_task"
)
