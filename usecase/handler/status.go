package handler

type StepStatusType int

const (
	StepStatusUnknown StepStatusType = iota
	StepStatusQueue
	StepStatusInit
	StepStatusSkip
	StepStatusProgress
	StepStatusError
	StepStatusFinish
)
