package job

type end struct {
	input StartInput
}

func (e *end) Start() (output StartOutput, err error) {
	var (
		dataCtx   = e.input.DataCtx
		step      = e.input.Step
		actionEnd = e.input.Step.Action.End
	)

	if actionEnd == nil {
		err = &ErrActionConfigEmpty{jobType: step.Type.String(), stepId: step.Id}
		return
	}

	// construct data
	var data map[string]any // key is step id
	for _, stepId := range actionEnd.ReturnDataFromStepIds {
		stepDataCtx := dataCtx.GetStep(stepId)
		data[stepId] = stepDataCtx.Data.Body
	}

	output.Data = data

	return
}
