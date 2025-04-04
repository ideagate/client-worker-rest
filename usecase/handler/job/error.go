package job

import (
	"fmt"
)

type ErrActionConfigEmpty struct {
	jobType string
	stepId  string
}

func (e *ErrActionConfigEmpty) Error() string {
	return fmt.Sprintf("action config empty, job type: %s, step id: %s", e.jobType, e.stepId)
}
