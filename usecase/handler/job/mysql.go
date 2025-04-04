package job

import (
	"github.com/bayu-aditya/ideagate/backend/core/model/endpoint"
	"github.com/bayu-aditya/ideagate/backend/core/utils/errors"
)

type mysql struct {
	Input StartInput
}

func (j *mysql) Start() (output StartOutput, err error) {
	var (
		ctx        = j.Input.Ctx
		ctxData    = j.Input.DataCtx
		dataSource = j.Input.DataSource
		step       = j.Input.Step
		action     = j.Input.Step.Action.Mysql
	)

	if dataSource.MysqlConn == nil {
		err = errors.New("empty mysql connection")
		return
	}

	if action == nil {
		err = &ErrActionConfigEmpty{jobType: step.Type.String(), stepId: step.Id}
		return
	}

	session := dataSource.MysqlConn.WithContext(ctx)

	if len(action.Queries) > 1 {
		session.Begin()
		defer session.Commit()
	}

	for _, queryItem := range action.Queries {
		// run template for query template
		queryVariable := endpoint.Variable(*queryItem.Query)
		query, errQuery := queryVariable.GetValueString(step.Id, ctxData)
		if errQuery != nil {
			err = errQuery
			return
		}

		// run template for parameters
		var paramsParsed []interface{}
		for _, param := range queryItem.Parameters {
			paramVariable := endpoint.Variable(*param)
			paramParsed, errParsed := paramVariable.GetValue(step.Id, ctxData)
			if errParsed != nil {
				err = errParsed
				return
			}

			paramsParsed = append(paramsParsed, paramParsed)
		}

		// execute query
		var rows []map[string]interface{}
		if err = session.Raw(query, paramsParsed...).Scan(&rows).Error; err != nil {
			return
		}

		// write into context data
		ctxData.SetStepDataBody(step.Id, rows)
	}

	return
}
