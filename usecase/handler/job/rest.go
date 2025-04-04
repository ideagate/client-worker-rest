package job

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/bayu-aditya/ideagate/backend/core/model/endpoint"
)

type rest struct {
	Input StartInput
}

func (j *rest) Start() (output StartOutput, err error) {
	var (
		ctx        = j.Input.Ctx
		dataCtx    = j.Input.DataCtx
		dataSource = j.Input.DataSource
		step       = j.Input.Step
		actionRest = j.Input.Step.Action.Rest
	)

	if actionRest == nil {
		err = &ErrActionConfigEmpty{jobType: step.Type.String(), stepId: step.Id}
		return
	}

	//construct request url, path is getting by template
	pathVariable := endpoint.Variable(*actionRest.Path)
	path, err := pathVariable.GetValueString(step.Id, dataCtx)
	if err != nil {
		return
	}
	reqUrl := dataSource.Config.Host + path

	// doing request
	req, err := http.NewRequestWithContext(ctx, actionRest.Method, reqUrl, nil)
	if err != nil {
		return
	}

	// construct request header
	for headerKey, headerVar := range actionRest.Headers {
		var headerValue string

		headerVariable := endpoint.Variable(*headerVar)
		headerValue, err = headerVariable.GetValueString(step.Id, dataCtx)
		if err != nil {
			return
		}

		req.Header.Set(headerKey, headerValue)
	}

	// doing http call
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	// set status code into context data
	dataCtx.SetStepStatusCode(step.Id, resp.StatusCode)

	// parse response json into context data
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var respData any
	if len(respBytes) > 0 {
		if err = json.Unmarshal(respBytes, &respData); err != nil {
			return
		}
	}
	dataCtx.SetStepDataBody(step.Id, respData)

	return
}
