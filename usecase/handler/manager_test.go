package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/bayu-aditya/ideagate/backend/core/model/constant"
	pbEndpoint "github.com/bayu-aditya/ideagate/backend/model/gen-go/core/endpoint"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manager - Process", func() {
	Context("Linear", func() {
		It("start - sleep - sleep - end", func() {
			t := GinkgoT()

			gin.SetMode(gin.TestMode)
			resp := httptest.NewRecorder()
			mockCtxGin, _ := gin.CreateTestContext(resp)
			mockCtxGin.Request = &http.Request{
				URL: &url.URL{},
			}

			mockEndpoint := &pbEndpoint.Endpoint{
				Id: "mock_endpoint_id",
			}

			mockWorkflow := &pbEndpoint.Workflow{
				Steps: []*pbEndpoint.Step{
					{
						Id:   constant.StepIdStart,
						Type: pbEndpoint.StepType_STEP_TYPE_START,
					},
					{
						Id:   "sleep_1",
						Type: pbEndpoint.StepType_STEP_TYPE_SLEEP,
						Action: &pbEndpoint.Action{
							Sleep: &pbEndpoint.ActionSleep{
								TimeoutMs: 1000,
							},
						},
					},
					{
						Id:   "sleep_2",
						Type: pbEndpoint.StepType_STEP_TYPE_SLEEP,
						Action: &pbEndpoint.Action{
							Sleep: &pbEndpoint.ActionSleep{
								TimeoutMs: 500,
							},
						},
					},
					{
						Id:   "end",
						Type: pbEndpoint.StepType_STEP_TYPE_END,
						Action: &pbEndpoint.Action{
							End: &pbEndpoint.ActionEnd{},
						},
					},
				},
				Edges: []*pbEndpoint.Edge{
					{Id: "edge_1", Source: constant.StepIdStart, Dest: "sleep_1"},
					{Id: "edge_2", Source: "sleep_1", Dest: "sleep_2"},
					{Id: "edge_3", Source: "sleep_2", Dest: constant.StepIdEnd},
				},
			}

			manager, err := newManager(mockCtxGin, mockEndpoint, mockWorkflow)
			if err != nil {
				t.Error("new manager failed", err)
			}

			manager.RunHandler()
			Expect(resp.Code).To(Equal(http.StatusOK))
		})
	})
	Context("Parallel", func() {
		// start -- sleep1(3000) -- sleep2(3000) -- sleep5(500) -- sleep6(1500) -- end
		//       |- sleep3(1000) -- sleep4(1000) -|
		Context("start - sleep - end", func() {
			var (
				mockCtxGin       *gin.Context
				t                = GinkgoT()
				httpRecorder     = httptest.NewRecorder()
				mockDataEndpoint = &pbEndpoint.Endpoint{
					Id: "mock_endpoint_id",
					Settings: &pbEndpoint.Endpoint_SettingRest{
						SettingRest: &pbEndpoint.SettingRest{
							NumWorkers: 1,
							TimeoutMs:  8100,
						},
					},
				}
				mockWorkflow = &pbEndpoint.Workflow{
					Steps: []*pbEndpoint.Step{
						{
							Id:   constant.StepIdStart,
							Type: pbEndpoint.StepType_STEP_TYPE_START,
						},
						{
							Id:   "sleep_1",
							Type: pbEndpoint.StepType_STEP_TYPE_SLEEP,
							Action: &pbEndpoint.Action{
								Sleep: &pbEndpoint.ActionSleep{
									TimeoutMs: 3000,
								},
							},
						},
						{
							Id:   "sleep_2",
							Type: pbEndpoint.StepType_STEP_TYPE_SLEEP,
							Action: &pbEndpoint.Action{
								Sleep: &pbEndpoint.ActionSleep{
									TimeoutMs: 3000,
								},
							},
						},
						{
							Id:   "sleep_3",
							Type: pbEndpoint.StepType_STEP_TYPE_SLEEP,
							Action: &pbEndpoint.Action{
								Sleep: &pbEndpoint.ActionSleep{
									TimeoutMs: 1000,
								},
							},
						},
						{
							Id:   "sleep_4",
							Type: pbEndpoint.StepType_STEP_TYPE_SLEEP,
							Action: &pbEndpoint.Action{
								Sleep: &pbEndpoint.ActionSleep{
									TimeoutMs: 1000,
								},
							},
						},
						{
							Id:   "sleep_5",
							Type: pbEndpoint.StepType_STEP_TYPE_SLEEP,
							Action: &pbEndpoint.Action{
								Sleep: &pbEndpoint.ActionSleep{
									TimeoutMs: 500,
								},
							},
						},
						{
							Id:   "sleep_6",
							Type: pbEndpoint.StepType_STEP_TYPE_SLEEP,
							Action: &pbEndpoint.Action{
								Sleep: &pbEndpoint.ActionSleep{
									TimeoutMs: 1500,
								},
							},
						},
						{
							Id:   "end",
							Type: pbEndpoint.StepType_STEP_TYPE_END,
							Action: &pbEndpoint.Action{
								End: &pbEndpoint.ActionEnd{},
							},
						},
					},
					Edges: []*pbEndpoint.Edge{
						{Id: "edge_1", Source: constant.StepIdStart, Dest: "sleep_1"},
						{Id: "edge_2", Source: "sleep_1", Dest: "sleep_2"},
						{Id: "edge_3", Source: "sleep_2", Dest: "sleep_5"},
						{Id: "edge_4", Source: constant.StepIdStart, Dest: "sleep_3"},
						{Id: "edge_5", Source: "sleep_3", Dest: "sleep_4"},
						{Id: "edge_6", Source: "sleep_4", Dest: "sleep_5"},
						{Id: "edge_7", Source: "sleep_5", Dest: "sleep_6"},
						{Id: "edge_4", Source: "sleep_6", Dest: constant.StepIdEnd},
					},
				}
			)

			BeforeEach(func() {
				gin.SetMode(gin.TestMode)
				mockCtxGin, _ = gin.CreateTestContext(httpRecorder)
				mockCtxGin.Request = &http.Request{
					URL: &url.URL{},
				}
			})

			It("single worker: 10 secs", func() {
				mockDataEndpoint.Settings = &pbEndpoint.Endpoint_SettingRest{
					SettingRest: &pbEndpoint.SettingRest{
						NumWorkers: 1,
						TimeoutMs:  10100, // must be finished around 10 secs
					},
				}

				manager, err := newManager(mockCtxGin, mockDataEndpoint, mockWorkflow)
				if err != nil {
					t.Error("new manager failed", err)
				}

				manager.RunHandler()
				Expect(httpRecorder.Code).To(Equal(http.StatusOK))
			})

			It("multiple worker: 8 secs", func() {
				mockDataEndpoint.Settings = &pbEndpoint.Endpoint_SettingRest{
					SettingRest: &pbEndpoint.SettingRest{
						NumWorkers: 3,
						TimeoutMs:  8100, // must be finished around 8 secs
					},
				}

				manager, err := newManager(mockCtxGin, mockDataEndpoint, mockWorkflow)
				if err != nil {
					t.Error("new manager failed", err)
				}

				manager.RunHandler()
				Expect(httpRecorder.Code).To(Equal(http.StatusOK))
			})
		})
	})
})
