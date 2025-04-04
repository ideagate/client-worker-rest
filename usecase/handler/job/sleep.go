package job

import "time"

type sleep struct {
	input StartInput
}

func (s *sleep) Start() (output StartOutput, err error) {
	step := s.input.Step
	action := step.Action.Sleep

	if action == nil {
		err = &ErrActionConfigEmpty{jobType: step.Type.String(), stepId: step.Id}
		return
	}

	sleepMs := action.TimeoutMs
	time.Sleep(time.Duration(sleepMs) * time.Millisecond)
	return
}
