package register

//
//func registerUserTest(t *testing.T) {
//	s := server.New()
//
//	testCases := []struct {
//		name         string
//		payload      interface{}
//		expectedCode int
//	}{
//		{
//			name: "valid",
//			payload: map[string]string{
//				"email":    "user@gmail.com",
//				"password": "password",
//			},
//			expectedCode: http.StatusCreated,
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			rec := httptest.NewRecorder()
//			b := &bytes.Buffer{}
//			json.NewEncoder(b).Encode(tc.payload)
//			req, _ := http.NewRequest(http.MethodPost, "/register", b)
//		})
//	}
//}
