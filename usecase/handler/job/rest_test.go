package job

import (
	"context"
	"reflect"
	"testing"

	entityContext "github.com/bayu-aditya/ideagate/backend/core/model/entity/context"
	entityDataSource "github.com/bayu-aditya/ideagate/backend/core/model/entity/datasource"
	pbEndpoint "github.com/bayu-aditya/ideagate/backend/model/gen-go/core/endpoint"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func Test_rest_Start(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockStepId := "mockStepId"
	mockNextStepId := "mockNextStepId"

	mockResponseString := `{
		"data": {
			"id": "mock-user-id",
			"name": "mock-user-name"
		}
	}`
	mockResponseMap := map[string]any{
		"data": map[string]interface{}{
			"id":   "mock-user-id",
			"name": "mock-user-name",
		},
	}

	type fields struct {
		Input StartInput
	}
	var tests = []struct {
		name        string
		fields      fields
		wantOutput  StartOutput
		wantErr     bool
		wantCtxData *entityContext.ContextData
		mockFunc    func()
	}{
		{
			name: "negative - nil context",
			fields: fields{
				Input: StartInput{
					Ctx:     nil,
					DataCtx: &entityContext.ContextData{},
					DataSource: &entityDataSource.DataSource{
						Config: entityDataSource.Config{Host: "https://mockhost.com/api"},
					},
					Step: &pbEndpoint.Step{
						Id: mockStepId,
						Action: &pbEndpoint.Action{
							Rest: &pbEndpoint.ActionRest{
								Method: "GET",
								Path:   &pbEndpoint.Variable{Value: "/user/detail"},
							},
						},
					},
				},
			},
			mockFunc:    func() {},
			wantErr:     true,
			wantCtxData: &entityContext.ContextData{},
		},
		{
			name: "negative - invalid method",
			fields: fields{
				Input: StartInput{
					Ctx:     context.TODO(),
					DataCtx: &entityContext.ContextData{},
					DataSource: &entityDataSource.DataSource{
						Config: entityDataSource.Config{Host: "https://mockhost.com/api"},
					},
					Step: &pbEndpoint.Step{
						Id: mockStepId,
						Action: &pbEndpoint.Action{
							Rest: &pbEndpoint.ActionRest{
								Method: "unknown",
								Path:   &pbEndpoint.Variable{Value: "/user/detail"},
							},
						},
					},
				},
			},
			mockFunc:    func() {},
			wantErr:     true,
			wantCtxData: &entityContext.ContextData{},
		},
		{
			name: "success",
			fields: fields{
				Input: StartInput{
					Ctx: context.TODO(),
					DataCtx: &entityContext.ContextData{
						Step: map[string]entityContext.ContextStepData{
							mockStepId: {
								Var: map[string]any{
									"user_id": 123,
								},
							},
						},
					},
					DataSource: &entityDataSource.DataSource{
						Config: entityDataSource.Config{Host: "https://mockhost.com/api"},
					},
					Step: &pbEndpoint.Step{
						Id: mockStepId,
						Action: &pbEndpoint.Action{
							Rest: &pbEndpoint.ActionRest{
								Method: "GET",
								Path:   &pbEndpoint.Variable{Value: "/user/detail?user_id={{.Var.user_id}}"},
							},
						},
						Returns: []*pbEndpoint.Return{
							{NextStepId: mockNextStepId},
						},
					},
				},
			},
			mockFunc: func() {
				httpmock.RegisterResponder("GET", "https://mockhost.com/api/user/detail?user_id=123", httpmock.NewStringResponder(200, mockResponseString))
			},
			wantCtxData: &entityContext.ContextData{
				Step: map[string]entityContext.ContextStepData{
					mockStepId: {
						Var: map[string]any{
							"user_id": 123,
						},
						Data: entityContext.ContextStepDataBody{
							StatusCode: 200,
							Body:       mockResponseMap,
						},
					},
				},
			},
		},
		{
			name: "success - 500 from http server",
			fields: fields{
				Input: StartInput{
					Ctx: context.TODO(),
					DataCtx: &entityContext.ContextData{
						Step: map[string]entityContext.ContextStepData{
							mockStepId: {
								Var: map[string]any{
									"user_id": 123,
								},
							},
						},
					},
					DataSource: &entityDataSource.DataSource{
						Config: entityDataSource.Config{Host: "https://mockhost.com/api"},
					},
					Step: &pbEndpoint.Step{
						Id: mockStepId,
						Action: &pbEndpoint.Action{
							Rest: &pbEndpoint.ActionRest{
								Method: "GET",
								Path:   &pbEndpoint.Variable{Value: "/user/detail?user_id={{.Var.user_id}}"},
							},
						},
						Returns: []*pbEndpoint.Return{
							{NextStepId: mockNextStepId},
						},
					},
				},
			},
			mockFunc: func() {
				httpmock.RegisterResponder("GET", "https://mockhost.com/api/user/detail?user_id=123", httpmock.NewBytesResponder(500, nil))
			},
			wantCtxData: &entityContext.ContextData{
				Step: map[string]entityContext.ContextStepData{
					mockStepId: {
						Var: map[string]any{
							"user_id": 123,
						},
						Data: entityContext.ContextStepDataBody{
							StatusCode: 500,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &rest{
				Input: tt.fields.Input,
			}

			tt.mockFunc()

			gotOutput, err := j.Start()
			if (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("Start() gotOutput = %v, want %v", gotOutput, tt.wantOutput)
			}

			assert.EqualExportedValues(t, tt.wantCtxData, tt.fields.Input.DataCtx)
		})
	}
}
