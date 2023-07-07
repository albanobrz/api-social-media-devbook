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

			req := httptest.NewRequest("POST", "/users", test.input)
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

			req, _ := http.NewRequest("GET", "/users/", nil)

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

			req, _ := http.NewRequest("GET", "/users/", nil)
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
		expectedUpdatedResult error
	}{
		{
			name:                  "Success on UpdateUser",
			input:                 `{"name":"updated", "nick":"testupdated", "email":"user1@email.com"}`,
			urlId:                 "1",
			validToken:            ValidToken,
			userId:                "1",
			expectedStatusCode:    204,
			expectedUpdatedResult: nil,
		},
		{
			name:                  "Error on UpdateUser, unexistent url ID",
			input:                 `{"username":"updated", "nick":"testupdated", "email":"user1@email.com"}`,
			urlId:                 "",
			validToken:            ValidToken,
			userId:                "1",
			expectedStatusCode:    403,
			expectedUpdatedResult: assert.AnError,
		},
		{
			name:                  "Error on UpdateUser, ExtractUserID",
			input:                 `{"username":"updated", "nick":"testupdated", "email":"user1@email.com"}`,
			urlId:                 "1",
			validToken:            ValidToken + "invalidate token",
			userId:                "1",
			expectedStatusCode:    401,
			expectedUpdatedResult: assert.AnError,
		},
		{
			name:                  "Error on UpdateUser, tokenId != requestId",
			input:                 `{"username":"updated", "nick":"testupdated", "email":"user1@email.com"}`,
			urlId:                 "1",
			validToken:            DiffToken,
			userId:                "1",
			expectedStatusCode:    403,
			expectedUpdatedResult: assert.AnError,
		},
		{
			name:                  "Error on UpdateUser, empty bodyReq",
			input:                 "",
			urlId:                 "1",
			validToken:            ValidToken,
			userId:                "1",
			expectedStatusCode:    400,
			expectedUpdatedResult: assert.AnError,
		},
		{
			name:                  "Error on UpdateUser, broken bodyReq",
			input:                 `{"usernameupdated", "nick":"testupdated", "email":"user1@email.com"}`,
			urlId:                 "1",
			validToken:            ValidToken,
			userId:                "1",
			expectedStatusCode:    400,
			expectedUpdatedResult: assert.AnError,
		},
		{
			name:                  "Error on UpdateUser, incorrect field on bodyReq",
			input:                 `{"username":"updated", "nick":"testupdated", "email":"user1@email.com"}`,
			urlId:                 "1",
			validToken:            ValidToken,
			userId:                "1",
			expectedStatusCode:    400,
			expectedUpdatedResult: assert.AnError,
		},
		{
			name:                  "Error on call UpdateUser",
			input:                 `{"invalidField":"updated", "nick":"testupdated", "email":"user1@email.com"}`,
			urlId:                 "1",
			validToken:            ValidToken,
			userId:                "1",
			expectedStatusCode:    400,
			expectedUpdatedResult: assert.AnError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repositoryMock := mocks.NewUsersRepositoryMock()
			repositoryMock.On("UpdateUser", test.userId, mock.AnythingOfType("entities.User")).Return(test.expectedUpdatedResult)

			usersController := NewUsersController(repositoryMock)

			req, _ := http.NewRequest("PUT", "/users/", strings.NewReader(test.input))
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

func TestDeleteUser(t *testing.T) {

	tests := []struct {
		name                     string
		expectedStatusCode       int
		userId                   string
		validToken               string
		expectedDeleteUserResult error
	}{
		{
			name:                     "Success on Delete",
			expectedStatusCode:       204,
			userId:                   "1",
			validToken:               ValidToken,
			expectedDeleteUserResult: nil,
		},
		{
			name:                     "Error on Delete",
			expectedStatusCode:       500,
			userId:                   "1",
			validToken:               ValidToken,
			expectedDeleteUserResult: assert.AnError,
		},
		{
			name:                     "Error on Delete, incorrect userId",
			expectedStatusCode:       403,
			userId:                   "122",
			validToken:               ValidToken,
			expectedDeleteUserResult: assert.AnError,
		},
		{
			name:                     "Error on Delete, invalid authToken",
			expectedStatusCode:       401,
			userId:                   "1",
			validToken:               ValidToken + "invalidate",
			expectedDeleteUserResult: assert.AnError,
		},
		{
			name:                     "Error on Delete, empty userId",
			expectedStatusCode:       403,
			userId:                   "",
			validToken:               ValidToken,
			expectedDeleteUserResult: assert.AnError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repositoryMock := mocks.NewUsersRepositoryMock()
			repositoryMock.On("DeleteUser", mock.AnythingOfType("string")).Return(test.expectedDeleteUserResult)

			usersController := NewUsersController(repositoryMock)

			req, _ := http.NewRequest("DELETE", "/users/", nil)
			parameters := map[string]string{
				"userID": test.userId,
			}
			req = mux.SetURLVars(req, parameters)
			req.Header.Add("Authorization", "Bearer "+test.validToken)

			rr := httptest.NewRecorder()

			controller := http.HandlerFunc(usersController.DeleteUser)
			controller.ServeHTTP(rr, req)

			assert.Equal(t, test.expectedStatusCode, rr.Code)
		})
	}

}

func TestFollowUser(t *testing.T) {

	tests := []struct {
		name                 string
		expectedStatusCode   int
		followedId           string
		validToken           string
		expectedFollowResult error
	}{
		{
			name:                 "Succcess on FollowUser",
			expectedStatusCode:   204,
			followedId:           "test",
			validToken:           ValidToken,
			expectedFollowResult: nil,
		},
		{
			name:                 "Error on FollowUser",
			expectedStatusCode:   500,
			followedId:           "test",
			validToken:           ValidToken,
			expectedFollowResult: assert.AnError,
		},
		{
			name:                 "Error on FollowUser, empty followedId",
			expectedStatusCode:   500,
			followedId:           "",
			validToken:           ValidToken,
			expectedFollowResult: assert.AnError,
		},
		{
			name:                 "Error on FollowUser, invalid authToken",
			expectedStatusCode:   401,
			followedId:           "test",
			validToken:           ValidToken + "invalidate",
			expectedFollowResult: nil,
		},
		{
			name:                 "Error on FollowUser, follow the own user",
			expectedStatusCode:   403,
			followedId:           "1",
			validToken:           ValidToken,
			expectedFollowResult: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repositoryMock := mocks.NewUsersRepositoryMock()
			repositoryMock.On("Follow", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(test.expectedFollowResult)

			usersController := NewUsersController(repositoryMock)

			req, _ := http.NewRequest("POST", "/users/test/follow", nil)
			req.Header.Add("Authorization", "Bearer "+test.validToken)
			parameters := map[string]string{
				"userID": test.followedId,
			}
			req = mux.SetURLVars(req, parameters)

			rr := httptest.NewRecorder()

			controller := http.HandlerFunc(usersController.FollowUser)
			controller.ServeHTTP(rr, req)

			assert.Equal(t, test.expectedStatusCode, rr.Code)

		})
	}
}

func TestUnfollowUser(t *testing.T) {
	tests := []struct {
		name                   string
		expectedStatusCode     int
		followedId             string
		validToken             string
		expectedUnfollowResult error
	}{
		{
			name:                   "Succcess on UnfollowUser",
			expectedStatusCode:     204,
			followedId:             "test",
			validToken:             ValidToken,
			expectedUnfollowResult: nil,
		},
		{
			name:                   "Error on UnfollowUser",
			expectedStatusCode:     500,
			followedId:             "test",
			validToken:             ValidToken,
			expectedUnfollowResult: assert.AnError,
		},
		{
			name:                   "Error on UnfollowUser, invalid token",
			expectedStatusCode:     401,
			followedId:             "test",
			validToken:             ValidToken + "Invalidate",
			expectedUnfollowResult: assert.AnError,
		},
		{
			name:                   "Error on UnfollowUser, empty userId",
			expectedStatusCode:     500,
			followedId:             "",
			validToken:             ValidToken,
			expectedUnfollowResult: assert.AnError,
		},
		{
			name:                   "Error on UnfollowUser, wrong userId",
			expectedStatusCode:     403,
			followedId:             "1",
			validToken:             ValidToken,
			expectedUnfollowResult: assert.AnError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repositoryMock := mocks.NewUsersRepositoryMock()
			repositoryMock.On("Unfollow", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(test.expectedUnfollowResult)

			usersController := NewUsersController(repositoryMock)

			req, _ := http.NewRequest("POST", "/users/test/unfollow", nil)
			req.Header.Add("Authorization", "Bearer "+test.validToken)
			parameters := map[string]string{
				"userID": test.followedId,
			}
			req = mux.SetURLVars(req, parameters)

			rr := httptest.NewRecorder()

			controller := http.HandlerFunc(usersController.UnfollowUser)
			controller.ServeHTTP(rr, req)

			assert.Equal(t, test.expectedStatusCode, rr.Code)
		})
	}
}

func TestGetFollowers(t *testing.T) {
	tests := []struct {
		name                       string
		expectedStatusCode         int
		userId                     string
		expectedGetFollowersError  error
		expectedGetFollowersResult []string
	}{
		{
			name:                       "Success on GetFollowers",
			expectedStatusCode:         200,
			userId:                     "1",
			expectedGetFollowersError:  nil,
			expectedGetFollowersResult: []string{},
		},
		{
			name:                       "Error on GetFollowers",
			expectedStatusCode:         500,
			userId:                     "1",
			expectedGetFollowersError:  assert.AnError,
			expectedGetFollowersResult: []string{},
		},
		{
			name:                       "Error on GetFollowers, empty userId",
			expectedStatusCode:         400,
			userId:                     "",
			expectedGetFollowersError:  nil,
			expectedGetFollowersResult: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repositoryMock := mocks.NewUsersRepositoryMock()
			repositoryMock.On("GetFollowers", mock.AnythingOfType("string")).Return(test.expectedGetFollowersResult, test.expectedGetFollowersError)

			usersController := NewUsersController(repositoryMock)

			req, _ := http.NewRequest("GET", "/users/test/followers", nil)
			parameters := map[string]string{
				"userID": test.userId,
			}
			req = mux.SetURLVars(req, parameters)

			rr := httptest.NewRecorder()

			controller := http.HandlerFunc(usersController.GetFollowers)
			controller.ServeHTTP(rr, req)

			assert.Equal(t, test.expectedStatusCode, rr.Code)
		})
	}
}

func TestGetFollowing(t *testing.T) {
	tests := []struct {
		name                       string
		expectedStatusCode         int
		userId                     string
		expectedGetFollowingError  error
		expectedGetFollowingResult []string
	}{
		{
			name:                       "Success on GetFollowing",
			expectedStatusCode:         200,
			userId:                     "test",
			expectedGetFollowingError:  nil,
			expectedGetFollowingResult: []string{},
		},
		{
			name:                       "Error on GetFollowing",
			expectedStatusCode:         500,
			userId:                     "1",
			expectedGetFollowingError:  assert.AnError,
			expectedGetFollowingResult: nil,
		},
		{
			name:                       "Error on GetFollowing, empty userId",
			expectedStatusCode:         400,
			userId:                     "",
			expectedGetFollowingError:  assert.AnError,
			expectedGetFollowingResult: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repositoryMock := mocks.NewUsersRepositoryMock()
			repositoryMock.On("GetFollowing", mock.AnythingOfType("string")).Return(test.expectedGetFollowingResult, test.expectedGetFollowingError)

			usersController := NewUsersController(repositoryMock)

			req, _ := http.NewRequest("GET", "/users/test/following", nil)
			parameters := map[string]string{
				"userID": test.userId,
			}
			req = mux.SetURLVars(req, parameters)

			rr := httptest.NewRecorder()

			controller := http.HandlerFunc(usersController.GetFollowing)
			controller.ServeHTTP(rr, req)

			assert.Equal(t, test.expectedStatusCode, rr.Code)
		})
	}
}

func TestUpdatePassword(t *testing.T) {
	tests := []struct {
		name                        string
		expectedStatusCode          int
		userId                      string
		expectedGetPasswordError    error
		expectedGetPasswordResult   string
		expectedUpdatePasswordError error
	}{
		{
			name:                        "Success on UpdatePassword",
			expectedStatusCode:          200,
			userId:                      "1",
			expectedGetPasswordError:    nil,
			expectedGetPasswordResult:   "password",
			expectedUpdatePasswordError: nil,
		},
		// {
		// 	name:                       "Error on GetWhoAnUserFollow",
		// 	expectedStatusCode:         500,
		// 	userId:                     "1",
		// 	expectedGetPasswordError:  assert.AnError,
		// 	expectedGetPasswordResult: "password",
		// 	expectedUpdatePasswordError: nil,
		// },
		// {
		// 	name:                       "Error on GetWhoAnUserFollow, empty userId",
		// 	expectedStatusCode:         400,
		// 	userId:                     "",
		// 	expectedGetPasswordError:  nil,
		// 	expectedGetPasswordResult: "password",
		// 	expectedUpdatePasswordError: nil,
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repositoryMock := mocks.NewUsersRepositoryMock()
			repositoryMock.On("GetPassword", mock.AnythingOfType("string")).Return(test.expectedGetPasswordResult, test.expectedGetPasswordError)
			repositoryMock.On("UpdatePassword", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(test.expectedUpdatePasswordError)

			usersController := NewUsersController(repositoryMock)

			req, _ := http.NewRequest("GET", "/users/test/update-password", nil)
			parameters := map[string]string{
				"userID": test.userId,
			}
			req = mux.SetURLVars(req, parameters)

			rr := httptest.NewRecorder()

			controller := http.HandlerFunc(usersController.UpdatePassword)
			controller.ServeHTTP(rr, req)

			assert.Equal(t, test.expectedStatusCode, rr.Code)
		})
	}
}
