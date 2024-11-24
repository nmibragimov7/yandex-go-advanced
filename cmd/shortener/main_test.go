package main

//func TestAbs(t *testing.T) {
//	tests := []struct {
//		name   string
//		values float64
//		want   float64
//	}{
//		{
//			name:   "-3 to 3",
//			values: -3,
//			want:   3,
//		},
//		{
//			name:   "-2.000001 to 2.000001",
//			values: -2.000001,
//			want:   2.000001,
//		},
//		{
//			name:   "-0.000000003 to 0.000000003",
//			values: -0.000000003,
//			want:   0.000000003,
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			assert.Equal(t, test.want, Abs(test.values))
//			//if res := Abs(test.values); res != test.want {
//			//	t.Errorf("Abs() = %f, want %f", res, test.want)
//			//}
//		})
//	}
//}
//
//func TestUser_FullName(t *testing.T) {
//	tests := []struct {
//		name   string
//		values Client
//		want   string
//	}{
//		{
//			name: "1 test",
//			values: Client{
//				FirstName: "Nurlan",
//				LastName:  "Ibragimov",
//			},
//			want: "Nurlan Ibragimov",
//		},
//		{
//			name: "2 test",
//			values: Client{
//				FirstName: "Petr",
//				LastName:  "Petrov",
//			},
//			want: "Petr Petrov",
//		},
//	}
//
//	for _, test := range tests { // цикл по всем тестам
//		t.Run(test.name, func(t *testing.T) {
//			assert.Equal(t, test.want, test.values.FullName())
//			//if res := test.values.FullName(); res != test.want {
//			//	t.Errorf("user.FullName() = %s, want %s", res, test.want)
//			//}
//		})
//	}
//}
//
//func TestFamily_AddNew(t *testing.T) {
//	f := Family{}
//
//	tests := []struct {
//		name         string
//		relationship Relationship
//		values       Person
//		want         bool
//	}{
//		{
//			name:         "add father",
//			relationship: Father,
//			values: Person{
//				FirstName: "Nurlan",
//				LastName:  "Ibragimov",
//				Age:       31,
//			},
//			want: false,
//		},
//		{
//			name:         "double add father",
//			relationship: Father,
//			values: Person{
//				FirstName: "Petr",
//				LastName:  "Petrov",
//				Age:       31,
//			},
//			want: true,
//		},
//	}
//
//	for _, test := range tests { // цикл по всем тестам
//		t.Run(test.name, func(t *testing.T) {
//			assert.Equal(t, test.want, f.AddNew(test.relationship, test.values) != nil)
//			//if res := f.AddNew(test.relationship, test.values); (res != nil) != test.want {
//			//	t.Errorf("f.AddNew() = %s, want %t", res, test.want)
//			//}
//		})
//	}
//}
//
//func TestUserViewHandler(t *testing.T) {
//	type want struct {
//		code        int
//		response    string
//		contentType string
//	}
//	tests := []struct {
//		name string
//		url  string
//		want want
//	}{
//		{
//			name: "negative test #1",
//			url:  "/users",
//			want: want{
//				code:        400,
//				response:    `{"message": "user_id is empty"}`,
//				contentType: "application/json",
//			},
//		},
//		{
//			name: "negative test #2",
//			url:  "/users?user_id=u3",
//			want: want{
//				code:        404,
//				response:    `{"message": "user not found"}`,
//				contentType: "application/json",
//			},
//		},
//		{
//			name: "positive test #3",
//			url:  "/users?user_id=u1",
//			want: want{
//				code:        200,
//				response:    `{"ID":"u1","FirstName":"Misha","LastName":"Popov"}`,
//				contentType: "application/json",
//			},
//		},
//	}
//
//	users := make(map[string]User)
//	u1 := User{
//		ID:        "u1",
//		FirstName: "Misha",
//		LastName:  "Popov",
//	}
//	u2 := User{
//		ID:        "u2",
//		FirstName: "Sasha",
//		LastName:  "Popov",
//	}
//	users["u1"] = u1
//	users["u2"] = u2
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			request := httptest.NewRequest(http.MethodPost, test.url, nil)
//			// создаём новый Recorder
//			w := httptest.NewRecorder()
//			handler := UserViewHandler(users)
//			handler.ServeHTTP(w, request)
//
//			res := w.Result()
//			// проверяем код ответа
//			assert.Equal(t, test.want.code, res.StatusCode)
//
//			// получаем и проверяем тело запроса
//			defer res.Body.Close()
//			resBody, err := io.ReadAll(res.Body)
//
//			require.NoError(t, err)
//			assert.JSONEq(t, test.want.response, string(resBody))
//			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
//		})
//	}
//}
