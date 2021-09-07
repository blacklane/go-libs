package uhttp

//
// func ExampleHTTPDefaultMiddleware() {
// 	trackingID := "tracking_id_ExampleHTTPAddAll"
// 	livePath := "/live"
// 	// Set current time function so we can control the request duration
// 	logger.SetNowFunc(func() time.Time {
// 		return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
// 	})
// 	log := logger.New(prettyJSONWriter{}, "ExampleHTTPAddAll")
//
// 	respWriterBar := httptest.NewRecorder()
// 	requestBar := httptest.NewRequest(http.MethodGet, "http://example.com/foo?bar=foo", nil)
// 	requestBar.Header.Set(constants.HeaderTrackingID, trackingID) // This header is set to have predictable value in the log output
// 	requestBar.Header.Set(constants.HeaderForwardedFor, "localhost")
//
// 	respWriterLive := httptest.NewRecorder()
// 	requestLive := httptest.NewRequest(http.MethodGet, "http://example.com"+livePath, nil)
// 	requestLive.Header.Set(constants.HeaderTrackingID, trackingID) // This header is set to have predictable value in the log output
// 	requestLive.Header.Set(constants.HeaderForwardedFor, "localhost")
//
// 	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		logger.FromContext(r.Context()).Info().Msg("always logged")
// 		_, _ = fmt.Fprint(w, "ExampleHTTPAddAll")
// 	})
//
// 	allInOneMiddleware := HTTPDefaultMiddleware("/foo", opentracing.NoopTracer{}, log, livePath)
// 	h := allInOneMiddleware(handler)
//
// 	h.ServeHTTP(respWriterLive, requestLive)
// 	h.ServeHTTP(respWriterBar, requestBar)
//
// 	// Output:
// 	// {
// 	//   "application": "ExampleHTTPAddAll",
// 	//   "entry_point": false,
// 	//   "host": "example.com",
// 	//   "ip": "localhost",
// 	//   "level": "info",
// 	//   "message": "always logged",
// 	//   "params": "",
// 	//   "path": "/live",
// 	//   "request_depth": 0,
// 	//   "request_id": "tracking_id_ExampleHTTPAddAll",
// 	//   "route": "",
// 	//   "timestamp": "2009-11-10T23:00:00.000Z",
// 	//   "tracking_id": "tracking_id_ExampleHTTPAddAll",
// 	//   "tree_path": "",
// 	//   "user_agent": "",
// 	//   "verb": "GET"
// 	// }
// 	// {
// 	//   "application": "ExampleHTTPAddAll",
// 	//   "entry_point": false,
// 	//   "host": "example.com",
// 	//   "ip": "localhost",
// 	//   "level": "info",
// 	//   "message": "always logged",
// 	//   "params": "bar=foo",
// 	//   "path": "/foo",
// 	//   "request_depth": 0,
// 	//   "request_id": "tracking_id_ExampleHTTPAddAll",
// 	//   "route": "",
// 	//   "timestamp": "2009-11-10T23:00:00.000Z",
// 	//   "tracking_id": "tracking_id_ExampleHTTPAddAll",
// 	//   "tree_path": "",
// 	//   "user_agent": "",
// 	//   "verb": "GET"
// 	// }
// 	// {
// 	//   "application": "ExampleHTTPAddAll",
// 	//   "entry_point": false,
// 	//   "event": "request_finished",
// 	//   "host": "example.com",
// 	//   "ip": "localhost",
// 	//   "level": "info",
// 	//   "message": "GET /foo",
// 	//   "params": "bar=foo",
// 	//   "path": "/foo",
// 	//   "request_depth": 0,
// 	//   "request_duration": 0,
// 	//   "request_id": "tracking_id_ExampleHTTPAddAll",
// 	//   "route": "",
// 	//   "status": 200,
// 	//   "timestamp": "2009-11-10T23:00:00.000Z",
// 	//   "tracking_id": "tracking_id_ExampleHTTPAddAll",
// 	//   "tree_path": "",
// 	//   "user_agent": "",
// 	//   "verb": "GET"
// 	// }
// }
