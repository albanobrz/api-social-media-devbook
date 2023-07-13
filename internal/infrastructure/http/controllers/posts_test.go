package controllers

import (
	"api/internal/domain/entities"
	"api/internal/domain/repositories/mocks"
	"bytes"
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

func TestCreatePost(t *testing.T) {
	postSerialized, err := os.ReadFile("../../../../test/resources/post.json")

	if err != nil {
		t.Errorf("json")
	}

	tests := []struct {
		name                     string
		input                    *bytes.Buffer
		urlId                    string
		validToken               string
		expectedCreatePostResult entities.Post
		expectedStatusCode       int
		expectedErrorMessage     string
		responseIsAnError        bool
		expectedError            error
	}{
		{
			name:                     "Success on CreateUser",
			input:                    bytes.NewBuffer(postSerialized),
			urlId:                    "1",
			validToken:               ValidToken,
			expectedCreatePostResult: entities.Post{},
			expectedStatusCode:       http.StatusCreated,
			responseIsAnError:        false,
			expectedErrorMessage:     "",
			expectedError:            nil,
		},
		{
			name:                     "Error on CreatePost",
			input:                    bytes.NewBuffer(postSerialized),
			urlId:                    "1",
			validToken:               ValidToken,
			expectedCreatePostResult: entities.Post{},
			expectedStatusCode:       http.StatusInternalServerError,
			responseIsAnError:        true,
			expectedErrorMessage:     "{\"error\":\"error ocurred\"}",
			expectedError:            errors.New("error ocurred"),
		},
		{
			name:                 "Error on CreatePost, empty input",
			input:                bytes.NewBuffer([]byte{}),
			urlId:                "1",
			validToken:           ValidToken,
			expectedStatusCode:   http.StatusBadRequest,
			responseIsAnError:    true,
			expectedErrorMessage: "{\"error\":\"unexpected end of JSON input\"}",
			expectedError:        assert.AnError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repositoryMock := mocks.NewPostsRepositoryMock()
			repositoryMock.On("CreatePost", mock.AnythingOfType("entities.Post")).Return(test.expectedCreatePostResult, test.expectedError)

			postsController := NewPostsController(repositoryMock)

			req := httptest.NewRequest("POST", "/posts", test.input)
			req.Header.Add("Authorization", "Bearer "+test.validToken)
			params := map[string]string{
				"userID": test.urlId,
			}
			req = mux.SetURLVars(req, params)
			rr := httptest.NewRecorder()

			controller := http.HandlerFunc(postsController.CreatePost)
			controller.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Error on status code got %d; expected %d", rr.Result().StatusCode, test.expectedStatusCode)
			}
		})
	}
}

func TestGetPosts(t *testing.T) {

	tests := []struct {
		name                      string
		input                     string
		expectedGetAllPostsResult []entities.Post
		expectedError             error
		expectedStatusCode        int
	}{
		{
			name:                      "Success on GetAllPosts",
			input:                     "",
			expectedGetAllPostsResult: []entities.Post{},
			expectedError:             nil,
			expectedStatusCode:        200,
		},
		{
			name:                      "Error on GetAllPosts",
			input:                     "",
			expectedGetAllPostsResult: []entities.Post{},
			expectedError:             assert.AnError,
			expectedStatusCode:        500,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repositoryMock := mocks.NewPostsRepositoryMock()
			repositoryMock.On("GetAllPosts").Return(test.expectedGetAllPostsResult, test.expectedError)

			postsController := NewPostsController(repositoryMock)

			req, _ := http.NewRequest("GET", "/posts/", nil)

			rr := httptest.NewRecorder()

			controller := http.HandlerFunc(postsController.GetAllPosts)
			controller.ServeHTTP(rr, req)

			assert.Equal(t, test.expectedStatusCode, rr.Code)
		})
	}
}

func TestUpdatePost(t *testing.T) {

	// postSerialized, err := os.ReadFile("../../../../test/resources/post.json")

	// if err != nil {
	// 	t.Errorf("json")
	// }

	tests := []struct {
		name                     string
		input                    string
		urlId                    string
		validToken               string
		userId                   string
		expectedPostWithIdResult entities.Post
		expectedPostWithIdError  error
		expectedStatusCode       int
		expectedUpdatedResult    error
	}{
		{
			name:                  "Success on UpdateUser",
			input:                 `{"title": "Wow", "content": "updated post"}`,
			urlId:                 "64a399cdb6a0487490ed730c",
			validToken:            ValidToken,
			userId:                "1",
			expectedStatusCode:    204,
			expectedUpdatedResult: nil,
		},
		// {
		// 	name:                  "Error on UpdateUser, unexistent url ID",
		// 	input:                 `{"username":"updated", "nick":"testupdated", "email":"user1@email.com"}`,
		// 	urlId:                 "",
		// 	validToken:            ValidToken,
		// 	userId:                "1",
		// 	expectedStatusCode:    403,
		// 	expectedUpdatedResult: assert.AnError,
		// },
		// {
		// 	name:                  "Error on UpdateUser, ExtractUserID",
		// 	input:                 `{"username":"updated", "nick":"testupdated", "email":"user1@email.com"}`,
		// 	urlId:                 "1",
		// 	validToken:            ValidToken + "invalidate token",
		// 	userId:                "1",
		// 	expectedStatusCode:    401,
		// 	expectedUpdatedResult: assert.AnError,
		// },
		// {
		// 	name:                  "Error on UpdateUser, tokenId != requestId",
		// 	input:                 `{"username":"updated", "nick":"testupdated", "email":"user1@email.com"}`,
		// 	urlId:                 "1",
		// 	validToken:            DiffToken,
		// 	userId:                "1",
		// 	expectedStatusCode:    403,
		// 	expectedUpdatedResult: assert.AnError,
		// },
		// {
		// 	name:                  "Error on UpdateUser, empty bodyReq",
		// 	input:                 "",
		// 	urlId:                 "1",
		// 	validToken:            ValidToken,
		// 	userId:                "1",
		// 	expectedStatusCode:    400,
		// 	expectedUpdatedResult: assert.AnError,
		// },
		// {
		// 	name:                  "Error on UpdateUser, broken bodyReq",
		// 	input:                 `{"usernameupdated", "nick":"testupdated", "email":"user1@email.com"}`,
		// 	urlId:                 "1",
		// 	validToken:            ValidToken,
		// 	userId:                "1",
		// 	expectedStatusCode:    400,
		// 	expectedUpdatedResult: assert.AnError,
		// },
		// {
		// 	name:                  "Error on UpdateUser, incorrect field on bodyReq",
		// 	input:                 `{"username":"updated", "nick":"testupdated", "email":"user1@email.com"}`,
		// 	urlId:                 "1",
		// 	validToken:            ValidToken,
		// 	userId:                "1",
		// 	expectedStatusCode:    400,
		// 	expectedUpdatedResult: assert.AnError,
		// },
		// {
		// 	name:                  "Error on call UpdateUser",
		// 	input:                 `{"invalidField":"updated", "nick":"testupdated", "email":"user1@email.com"}`,
		// 	urlId:                 "1",
		// 	validToken:            ValidToken,
		// 	userId:                "1",
		// 	expectedStatusCode:    400,
		// 	expectedUpdatedResult: assert.AnError,
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repositoryMock := mocks.NewPostsRepositoryMock()
			repositoryMock.On("GetPostWithId", test.urlId).Return(test.expectedPostWithIdResult, test.expectedPostWithIdError)
			repositoryMock.On("UpdatePost", test.urlId, mock.AnythingOfType("entities.Post")).Return(test.expectedUpdatedResult)

			postsController := NewPostsController(repositoryMock)

			req, _ := http.NewRequest("PUT", "/posts/", strings.NewReader(test.input))
			req.Header.Add("Authorization", "Bearer "+test.validToken)
			params := map[string]string{
				"postID": test.urlId,
			}
			req = mux.SetURLVars(req, params)

			rr := httptest.NewRecorder()

			controller := http.HandlerFunc(postsController.UpdatePost)
			controller.ServeHTTP(rr, req)

			assert.Equal(t, test.expectedStatusCode, rr.Code)

		})
	}
}
