package controllers

import (
	"api/internal/domain/entities"
	"api/internal/domain/repositories/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	userSerialized, err := os.ReadFile("../../../../test/resources/user.json")
	invalidUserSerialized, err := os.ReadFile("../../../../test/resources/invalid_user.json")

	if err != nil {
		t.Errorf("json")
	}

	var user entities.User
	json.Unmarshal(userSerialized, &user)

	tests := []struct {
		name                     string
		input                    *bytes.Buffer
		expectedCreateUserResult entities.User
		expectedStatusCode       int
		expectedErrorMessage     string
		responseIsAnError        bool
		expectedError            error
	}{
		{
			name:                     "Success on CreateUser",
			input:                    bytes.NewBuffer(userSerialized),
			expectedCreateUserResult: entities.User{},
			expectedStatusCode:       http.StatusCreated,
			responseIsAnError:        false,
			expectedErrorMessage:     "",
			expectedError:            nil,
		},
		{
			name:                     "Error on CreateUser",
			input:                    bytes.NewBuffer(userSerialized),
			expectedCreateUserResult: entities.User{},
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
		{
			name:                 "Error on CreateUser, invalid user data",
			input:                bytes.NewBuffer(invalidUserSerialized),
			expectedStatusCode:   http.StatusBadRequest,
			responseIsAnError:    true,
			expectedErrorMessage: "{\"error\":\"nick is empty\"}",
			expectedError:        assert.AnError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repositoryMock := mocks.NewUsersRepositoryMock()
			repositoryMock.On("Create", mock.AnythingOfType("entities.User")).Return(test.expectedCreateUserResult, test.expectedError)

			usersController := NewUsersController(repositoryMock)

			req := httptest.NewRequest("POST", "/mongo/users", test.input)
			rr := httptest.NewRecorder()

			controller := http.HandlerFunc(usersController.CreateUser)
			controller.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Error on status code got %d; expected %d", rr.Result().StatusCode, test.expectedStatusCode)
			}
		})
	}
}

func TestGetAllUsers(t *testing.T) {

	tests := []struct {
		name                      string
		input                     string
		expectedGetAllUsersReturn []entities.User
		expectedGetAllUsersError  error
		expectedStatusCode        int
	}{
		{
			name:                      "Success on GetUsers",
			input:                     "",
			expectedGetAllUsersReturn: []entities.User{},
			expectedGetAllUsersError:  nil,
			expectedStatusCode:        200,
		},
		{
			name:                      "Error on GetUsers",
			input:                     "",
			expectedGetAllUsersReturn: []entities.User{},
			expectedGetAllUsersError:  assert.AnError,
			expectedStatusCode:        500,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repositoryMock := mocks.NewUsersRepositoryMock()
			repositoryMock.On("GetAllUsers").Return(test.expectedGetAllUsersReturn, test.expectedGetAllUsersError)

			usersController := NewUsersController(repositoryMock)

			req, _ := http.NewRequest("GET", "/mongo/users/", nil)

			rr := httptest.NewRecorder()

			controller := http.HandlerFunc(usersController.GetAllUsers)
			controller.ServeHTTP(rr, req)

			assert.Equal(t, test.expectedStatusCode, rr.Code)
		})
	}
}

func TestGetUser(t *testing.T) {

	var returnedUser entities.User
	userSerialized, _ := os.ReadFile("../../../../../test/resources/user.json")
	json.Unmarshal(userSerialized, &returnedUser)

	tests := []struct {
		name                  string
		requestID             string
		expectedStatusCode    int
		input                 string
		expectedGetUserReturn entities.User
		expectedGetUserError  error
	}{
		{
			name:                  "Success on GetUser",
			requestID:             "1",
			expectedStatusCode:    200,
			input:                 "1",
			expectedGetUserReturn: returnedUser,
			expectedGetUserError:  nil,
		},
		{
			name:                  "Error on GetUser",
			requestID:             "1",
			expectedStatusCode:    500,
			input:                 "1",
			expectedGetUserReturn: entities.User{},
			expectedGetUserError:  assert.AnError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repositoryMock := mocks.NewUsersRepositoryMock()
			repositoryMock.On("GetUserByNick", test.input).Return(test.expectedGetUserReturn, test.expectedGetUserError)

			usersController := NewUsersController(repositoryMock)

			req, _ := http.NewRequest("GET", "/mongo/users/", nil)
			params := map[string]string{
				"userID": test.requestID,
			}
			req = mux.SetURLVars(req, params)
			rr := httptest.NewRecorder()

			controller := http.HandlerFunc(usersController.GetUser)

			controller.ServeHTTP(rr, req)

			assert.Equal(t, test.expectedStatusCode, rr.Code)
		})
	}
}

func TestUpdateUser(t *testing.T) {

	tests := []struct {
		name                  string
		input                 string
		urlId                 string
		validToken            string
		userId                string
		expectedStatusCode    int
		expectedUpdatedReturn error
		expectedUpdatedError  error
	}{
		{
			name:                  "Success on UpdateUser",
			input:                 `{"name":"updated", "nick":"testupdated", "email":"user1@email.com"}`,
			urlId:                 "1",
			validToken:            ValidToken,
			userId:                "1",
			expectedStatusCode:    204,
			expectedUpdatedReturn: nil,
			expectedUpdatedError:  nil,
		},
		{
			name:                  "Error on UpdateUser, unexistent url ID",
			input:                 `{"username":"updated", "nick":"testupdated", "email":"user1@email.com"}`,
			urlId:                 "",
			validToken:            ValidToken,
			userId:                "1",
			expectedStatusCode:    403,
			expectedUpdatedReturn: assert.AnError,
			expectedUpdatedError:  assert.AnError,
		},
		{
			name:                  "Error on UpdateUser, ExtractUserID",
			input:                 `{"username":"updated", "nick":"testupdated", "email":"user1@email.com"}`,
			urlId:                 "1",
			validToken:            ValidToken + "invalidate token",
			userId:                "1",
			expectedStatusCode:    401,
			expectedUpdatedReturn: assert.AnError,
			expectedUpdatedError:  assert.AnError,
		},
		{
			name:                  "Error on UpdateUser, tokenId != requestId",
			input:                 `{"username":"updated", "nick":"testupdated", "email":"user1@email.com"}`,
			urlId:                 "1",
			validToken:            DiffToken,
			userId:                "1",
			expectedStatusCode:    403,
			expectedUpdatedReturn: assert.AnError,
			expectedUpdatedError:  assert.AnError,
		},
		{
			name:                  "Error on UpdateUser, empty bodyReq",
			input:                 "",
			urlId:                 "1",
			validToken:            ValidToken,
			userId:                "1",
			expectedStatusCode:    400,
			expectedUpdatedReturn: assert.AnError,
			expectedUpdatedError:  assert.AnError,
		},
		{
			name:                  "Error on UpdateUser, broken bodyReq",
			input:                 `{"usernameupdated", "nick":"testupdated", "email":"user1@email.com"}`,
			urlId:                 "1",
			validToken:            ValidToken,
			userId:                "1",
			expectedStatusCode:    400,
			expectedUpdatedReturn: assert.AnError,
			expectedUpdatedError:  assert.AnError,
		},
		{
			name:                  "Error on UpdateUser, incorrect field on bodyReq",
			input:                 `{"username":"updated", "nick":"testupdated", "email":"user1@email.com"}`,
			urlId:                 "1",
			validToken:            ValidToken,
			userId:                "1",
			expectedStatusCode:    400,
			expectedUpdatedReturn: assert.AnError,
			expectedUpdatedError:  assert.AnError,
		},
		{
			name:                  "Error on call UpdateUser",
			input:                 `{"invalidField":"updated", "nick":"testupdated", "email":"user1@email.com"}`,
			urlId:                 "1",
			validToken:            ValidToken,
			userId:                "1",
			expectedStatusCode:    400,
			expectedUpdatedReturn: assert.AnError,
			expectedUpdatedError:  assert.AnError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repositoryMock := mocks.NewUsersRepositoryMock()
			repositoryMock.On("UpdateUser", test.userId, mock.AnythingOfType("entities.User")).Return(test.expectedUpdatedReturn, test.expectedUpdatedError)

			usersController := NewUsersController(repositoryMock)

			req, _ := http.NewRequest("PUT", "/mongo/users/", strings.NewReader(test.input))
			req.Header.Add("Authorization", "Bearer "+test.validToken)
			params := map[string]string{
				"userID": test.urlId,
			}
			req = mux.SetURLVars(req, params)

			rr := httptest.NewRecorder()

			controller := http.HandlerFunc(usersController.UpdateUser)
			controller.ServeHTTP(rr, req)

			assert.Equal(t, test.expectedStatusCode, rr.Code)

		})
	}
}
