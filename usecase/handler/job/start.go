package job

import (
	"github.com/bayu-aditya/ideagate/backend/core/model/endpoint"
	pbEndpoint "github.com/bayu-aditya/ideagate/backend/model/gen-go/core/endpoint"
	"github.com/gin-gonic/gin"
)

type start struct {
	Input StartInput
}

func (j *start) Start() (output StartOutput, err error) {
	// Parse request query and json
	queries, jsons, err := j.parseReqQueryJson(j.Input.GinCtx, j.Input.Endpoint.GetSettingRest())
	if err != nil {
		return
	}
	j.Input.DataCtx.SetRequestQuery(queries)
	j.Input.DataCtx.SetRequestJson(jsons)

	// Construct variables
	variables := make(map[string]interface{})
	for name, variable := range j.Input.Step.Variables {
		variableParsed := endpoint.Variable(*variable)
		variables[name], err = variableParsed.GetValue(j.Input.Step.Id, j.Input.DataCtx)
		if err != nil {
			return
		}
	}
	j.Input.DataCtx.SetStepVariable(j.Input.Step.Id, variables)

	// Construct output
	outputs := make(map[string]interface{})
	for name, out := range j.Input.Step.Outputs {
		outVariable := endpoint.Variable(*out)
		outputs[name], err = outVariable.GetValue(j.Input.Step.Id, j.Input.DataCtx)
		if err != nil {
			return
		}
	}
	j.Input.DataCtx.SetStepOutput(j.Input.Step.Id, outputs)

	// Construct result
	if len(j.Input.Step.Returns) == 0 {
		return
	}

	// Determine next step
	for _, returnItem := range j.Input.Step.Returns {
		output.NextStepIds = append(output.NextStepIds, returnItem.NextStepId)
	}

	return
}

func (j *start) parseReqQueryJson(c *gin.Context, setting *pbEndpoint.SettingRest) (dataQuery, dataJson map[string]any, err error) {
	// construct dataQuery parameters
	query := map[string]string{}
	if err = c.BindQuery(&query); err != nil {
		return
	}

	for fieldName, variable := range setting.GetQuery() {
		variableParsed := endpoint.Variable(*variable)
		dataQuery[fieldName], _ = variableParsed.GetValueString(j.Input.Step.Id, j.Input.DataCtx)
	}

	// construct body json
	dataJson = map[string]any{}
	if c.Request.Body != nil {
		if err = c.ShouldBindJSON(&dataJson); err != nil {
			return
		}
	}

	for fieldName, variable := range setting.GetJson() {
		variableParsed := endpoint.Variable(*variable)
		dataJson[fieldName], _ = variableParsed.GetValue(j.Input.Step.Id, j.Input.DataCtx)
	}

	return
}
