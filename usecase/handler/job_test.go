package handler

//func Test_jobMysql_Start(t *testing.T) {
//	db, dbmock, err := sqlmock.New()
//	if err != nil {
//		assert.NoError(t, err)
//	}
//	defer db.Close()
//
//	mysqldb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
//	assert.NoError(t, err)
//
//	dataSource := datasource.DataSource{
//		MysqlConn: mysqldb,
//	}
//
//	type args struct {
//		ctx        context.Context
//		ctxData    *entityContext.ContextData
//		dataSource datasource.DataSource
//		step       endpoint.Step
//		action     endpoint.Action
//	}
//	tests := []struct {
//		name        string
//		fields      jobMysql
//		args        args
//		wantErr     bool
//		wantCtxData *entityContext.ContextData
//		mockFunc    func(args)
//	}{
//		{
//			name:   "positive - single query",
//			fields: jobMysql{},
//			args: args{
//				ctx: context.TODO(),
//				ctxData: &entityContext.ContextData{
//					Step: map[string]entityContext.ContextStepData{
//						"session": {
//							Output: map[string]interface{}{
//								"user_id": 123,
//							},
//						},
//					},
//				},
//				dataSource: dataSource,
//				step:       endpoint.Step{Id: "fetch:user"},
//				action: endpoint.Action{
//					Id: "query1",
//					Queries: []endpoint.Query{
//						{
//							Query: "SELECT id, name FROM user where id = ? and status = ?",
//							Parameters: []string{
//								"{{.step.session.output.user_id}}",
//								"accepted",
//							},
//						},
//					},
//				},
//			},
//			wantCtxData: &entityContext.ContextData{
//				Step: map[string]entityContext.ContextStepData{
//					"session": {
//						Output: map[string]interface{}{
//							"user_id": 123,
//						},
//					},
//					"fetch:user": {
//						Data: map[string]interface{}{
//							"query1": map[string]interface{}{
//								"Body": []map[string]interface{}{
//									{"id": 123, "name": "hupla"},
//								},
//							},
//						},
//					},
//				},
//			},
//			mockFunc: func(args args) {
//				dbmock.ExpectQuery(regexp.QuoteMeta("SELECT id, name FROM user where id = ? and status = ?")).
//					WithArgs(123, "accepted").
//					WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(123, "hupla"))
//			},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			j := &tt.fields
//			tt.mockFunc(tt.args)
//			if err = j.Start(tt.args.ctx, tt.args.ctxData, tt.args.dataSource, tt.args.step, tt.args.action); (err != nil) != tt.wantErr {
//				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
//			}
//			assert.EqualExportedValues(t, tt.wantCtxData, tt.args.ctxData)
//		})
//	}
//}
