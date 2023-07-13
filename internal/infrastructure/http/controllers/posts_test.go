package controllers

import (
	"api/internal/domain/entities"
	"api/internal/domain/repositories/mocks"
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreatePost(t *testing.T) {
	postSerialized, err := os.ReadFile("../../../../test/resources/post.json")

	if err != nil {
		t.Errorf("json")
	}

	tests := []struct {
		name                     string
		input                    *bytes.Buffer
		expectedCreatePostResult entities.Post
		expectedStatusCode       int
		expectedErrorMessage     string
		responseIsAnError        bool
		expectedError            error
	}{
		{
			name:                     "Success on CreateUser",
			input:                    bytes.NewBuffer(postSerialized),
			expectedCreatePostResult: entities.Post{},
			expectedStatusCode:       http.StatusCreated,
			responseIsAnError:        false,
			expectedErrorMessage:     "",
			expectedError:            nil,
		},
		{
			name:                     "Error on CreateUser",
			input:                    bytes.NewBuffer(postSerialized),
			expectedCreatePostResult: entities.Post{},
			expectedStatusCode:       http.StatusInternalServerError,
			responseIsAnError:        true,
			expectedErrorMessage:     "{\"error\":\"error ocurred\"}",
			expectedError:            errors.New("error ocurred"),
		},
		{
			name:                 "Error on CreateUser, empty input",
			input:                bytes.NewBuffer([]byte{}),
			expectedStatusCode:   http.StatusBadRequest,
			responseIsAnError:    true,
			expectedErrorMessage: "{\"error\":\"unexpected end of JSON input\"}",
			expectedError:        assert.AnError,
		},
		// {
		// 	name:                 "Error on CreateUser, invalid user data",
		// 	input:                bytes.NewBuffer(invalidUserSerialized),
		// 	expectedStatusCode:   http.StatusBadRequest,
		// 	responseIsAnError:    true,
		// 	expectedErrorMessage: "{\"error\":\"nick is empty\"}",
		// 	expectedError:        assert.AnError,
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repositoryMock := mocks.NewPostsRepositoryMock()
			repositoryMock.On("CreatePost", mock.AnythingOfType("entities.Post")).Return(test.expectedCreatePostResult, test.expectedError)

			postsController := NewPostsController(repositoryMock)

			req := httptest.NewRequest("POST", "/posts", test.input)
			rr := httptest.NewRecorder()

			controller := http.HandlerFunc(postsController.CreatePost)
			controller.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Error on status code got %d; expected %d", rr.Result().StatusCode, test.expectedStatusCode)
			}
		})
	}
}

// func TestGetAllUsers(t *testing.T) {

// 	tests := []struct {
// 		name                      string
// 		input                     string
// 		expectedGetAllUsersReturn []entities.User
// 		expectedGetAllUsersError  error
// 		expectedStatusCode        int
// 	}{
// 		{
// 			name:                      "Success on GetUsers",
// 			input:                     "",
// 			expectedGetAllUsersReturn: []entities.User{},
// 			expectedGetAllUsersError:  nil,
// 			expectedStatusCode:        200,
// 		},
// 		{
// 			name:                      "Error on GetUsers",
// 			input:                     "",
// 			expectedGetAllUsersReturn: []entities.User{},
// 			expectedGetAllUsersError:  assert.AnError,
// 			expectedStatusCode:        500,
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			repositoryMock := mocks.NewUsersRepositoryMock()
// 			repositoryMock.On("GetAllUsers").Return(test.expectedGetAllUsersReturn, test.expectedGetAllUsersError)

// 			usersController := NewUsersController(repositoryMock)

// 			req, _ := http.NewRequest("GET", "/users/", nil)

// 			rr := httptest.NewRecorder()

// 			controller := http.HandlerFunc(usersController.GetAllUsers)
// 			controller.ServeHTTP(rr, req)

// 			assert.Equal(t, test.expectedStatusCode, rr.Code)
// 		})
// 	}
// }

// func TestGetUser(t *testing.T) {

// 	var returnedUser entities.User
// 	userSerialized, _ := os.ReadFile("../../../../../test/resources/user.json")
// 	json.Unmarshal(userSerialized, &returnedUser)

// 	tests := []struct {
// 		name                  string
// 		requestID             string
// 		expectedStatusCode    int
// 		input                 string
// 		expectedGetUserReturn entities.User
// 		expectedGetUserError  error
// 	}{
// 		{
// 			name:                  "Success on GetUser",
// 			requestID:             "1",
// 			expectedStatusCode:    200,
// 			input:                 "1",
// 			expectedGetUserReturn: returnedUser,
// 			expectedGetUserError:  nil,
// 		},
// 		{
// 			name:                  "Error on GetUser",
// 			requestID:             "1",
// 			expectedStatusCode:    500,
// 			input:                 "1",
// 			expectedGetUserReturn: entities.User{},
// 			expectedGetUserError:  assert.AnError,
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			repositoryMock := mocks.NewUsersRepositoryMock()
// 			repositoryMock.On("GetUserByNick", test.input).Return(test.expectedGetUserReturn, test.expectedGetUserError)

// 			usersController := NewUsersController(repositoryMock)

// 			req, _ := http.NewRequest("GET", "/users/", nil)
// 			params := map[string]string{
// 				"userID": test.requestID,
// 			}
// 			req = mux.SetURLVars(req, params)
// 			rr := httptest.NewRecorder()

// 			controller := http.HandlerFunc(usersController.GetUser)

// 			controller.ServeHTTP(rr, req)

// 			assert.Equal(t, test.expectedStatusCode, rr.Code)
// 		})
// 	}
// }
