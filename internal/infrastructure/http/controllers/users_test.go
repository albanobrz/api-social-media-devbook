package controllers

import (
	"api/internal/domain/entities"
	"api/internal/infrastructure/database/repositories/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

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
		expectedCreateUserResult uint64
		expectedStatusCode       int
		expectedErrorMessage     string
		responseIsAnError        bool
		expectedError            error
	}{
		{
			name:                     "Success on CreateUser",
			input:                    bytes.NewBuffer(userSerialized),
			expectedCreateUserResult: 1,
			expectedStatusCode:       http.StatusCreated,
			responseIsAnError:        false,
			expectedErrorMessage:     "",
			expectedError:            nil,
		},
		{
			name:                     "Error on CreateUser",
			input:                    bytes.NewBuffer(userSerialized),
			expectedCreateUserResult: 0,
			expectedStatusCode:       http.StatusBadRequest,
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
			repositoryMock.On("CreateUser", mock.AnythingOfType("entities.User")).Return(test.expectedCreateUserResult, test.expectedError)

			usersController := NewUsersController(repositoryMock)

			req := httptest.NewRequest("POST", "/users", test.input)
			rr := httptest.NewRecorder()

			controller := http.HandlerFunc(usersController.CreateUser)
			controller.ServeHTTP(rr, req)

			responseBody, _ := ioutil.ReadAll(rr.Result().Body)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Error on status code got %d; expected %d", rr.Result().StatusCode, test.expectedStatusCode)
			} else {
				if !test.responseIsAnError {
					assert.Equal(t, fmt.Sprint(test.expectedCreateUserResult), strings.TrimSpace(string(responseBody)))
				} else {
					assert.Equal(t, test.expectedErrorMessage, strings.TrimSpace(string(responseBody)))
				}
			}
		})
	}
}
