package job

import (
	"fmt"

	pbEndpoint "github.com/bayu-aditya/ideagate/backend/model/gen-go/core/endpoint"
)

type IJob interface {
	Start() (StartOutput, error) // TODO rename to Process
}

func New(jobType pbEndpoint.StepType, input StartInput) (IJob, error) {
	switch jobType {
	case pbEndpoint.StepType_STEP_TYPE_START:
		return &start{Input: input}, nil

	case pbEndpoint.StepType_STEP_TYPE_END:
		return &end{input: input}, nil

	case pbEndpoint.StepType_STEP_TYPE_SLEEP:
		return &sleep{input: input}, nil

	case pbEndpoint.StepType_STEP_TYPE_REST:
		return &rest{Input: input}, nil

	case pbEndpoint.StepType_STEP_TYPE_MYSQL:
		return &mysql{Input: input}, nil
	}

	return nil, fmt.Errorf("unknown job type '%s'", jobType.String())
}
