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

var postMocked = entities.Post{
	Title:      "Test Post",
	Content:    "new post",
	AuthorID:   "1",
	AuthorNick: "1",
	Likes:      0,
}

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
			name:                    "Success on UpdatePost",
			input:                   `{"title": "Wow", "content": "updated post"}`,
			urlId:                   "64a399cdb6a0487490ed730c",
			validToken:              ValidToken,
			userId:                  "1",
			expectedStatusCode:      204,
			expectedPostWithIdError: nil,
			expectedUpdatedResult:   nil,
		},
		{
			name:                    "Error on UpdatePost, unexistent url ID",
			input:                   `{"title": "Wow", "content": "updated post"}`,
			urlId:                   "2222",
			validToken:              ValidToken,
			userId:                  "1",
			expectedStatusCode:      500,
			expectedPostWithIdError: assert.AnError,
			expectedUpdatedResult:   assert.AnError,
		},
		{
			name:                    "Error on UpdatePost, tokenId != PostAuthorNick",
			input:                   `{"title": "Wow", "content": "updated post"}`,
			urlId:                   "64a399cdb6a0487490ed730c",
			validToken:              DiffToken,
			userId:                  "1",
			expectedStatusCode:      403,
			expectedPostWithIdError: nil,
			expectedUpdatedResult:   assert.AnError,
		},
		{
			name:                    "Error on UpdatePost, empty bodyReq",
			input:                   "",
			urlId:                   "64a399cdb6a0487490ed730c",
			validToken:              ValidToken,
			userId:                  "1",
			expectedStatusCode:      400,
			expectedPostWithIdError: nil,
			expectedUpdatedResult:   assert.AnError,
		},
		{
			name:                    "Error on call UpdatePost",
			input:                   `{"title": "Wow", "content": "updated post"}`,
			urlId:                   "1",
			validToken:              ValidToken,
			userId:                  "1",
			expectedStatusCode:      500,
			expectedPostWithIdError: nil,
			expectedUpdatedResult:   assert.AnError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repositoryMock := mocks.NewPostsRepositoryMock()

			repositoryMock.On("GetPostWithId", test.urlId).Return(postMocked, test.expectedPostWithIdError)
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

func TestDeletePost(t *testing.T) {

	tests := []struct {
		name                     string
		urlId                    string
		validToken               string
		userId                   string
		expectedPostWithIdResult entities.Post
		expectedPostWithIdError  error
		expectedStatusCode       int
		expectedDeleteResult     error
	}{
		{
			name:                    "Success on DeletePost",
			urlId:                   "64a399cdb6a0487490ed730c",
			validToken:              ValidToken,
			userId:                  "1",
			expectedStatusCode:      204,
			expectedPostWithIdError: nil,
			expectedDeleteResult:    nil,
		},
		{
			name:                    "Error on DeletePost",
			urlId:                   "64a399cdb6a0487490ed730c",
			validToken:              ValidToken,
			userId:                  "1",
			expectedStatusCode:      500,
			expectedPostWithIdError: nil,
			expectedDeleteResult:    assert.AnError,
		},
		{
			name:                    "Error on Delete, incorrect postID",
			urlId:                   "2222",
			validToken:              ValidToken,
			userId:                  "1",
			expectedStatusCode:      500,
			expectedPostWithIdError: assert.AnError,
			expectedDeleteResult:    nil,
		},
		{
			name:                    "Error on Delete, invalid authToken",
			urlId:                   "64a399cdb6a0487490ed730c",
			validToken:              DiffToken,
			userId:                  "1",
			expectedStatusCode:      403,
			expectedPostWithIdError: nil,
			expectedDeleteResult:    nil,
		},
		{
			name:                    "Error on Delete, empty postID",
			urlId:                   "",
			validToken:              ValidToken,
			userId:                  "1",
			expectedStatusCode:      500,
			expectedPostWithIdError: assert.AnError,
			expectedDeleteResult:    nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repositoryMock := mocks.NewPostsRepositoryMock()

			repositoryMock.On("GetPostWithId", test.urlId).Return(postMocked, test.expectedPostWithIdError)
			repositoryMock.On("DeletePost", test.urlId).Return(test.expectedDeleteResult)

			postsController := NewPostsController(repositoryMock)

			req, _ := http.NewRequest("DELETE", "/posts/", nil)
			req.Header.Add("Authorization", "Bearer "+test.validToken)
			params := map[string]string{
				"postID": test.urlId,
			}
			req = mux.SetURLVars(req, params)

			rr := httptest.NewRecorder()

			controller := http.HandlerFunc(postsController.DeletePost)
			controller.ServeHTTP(rr, req)

			assert.Equal(t, test.expectedStatusCode, rr.Code)
		})
	}

}
