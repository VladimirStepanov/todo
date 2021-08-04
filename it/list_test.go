package it

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/VladimirStepanov/todo-app/internal/handler"
	"github.com/VladimirStepanov/todo-app/internal/helpers"
	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/stretchr/testify/require"
)

var (
	listForCreate = &models.List{
		Title:       "title",
		Description: "description",
	}
)

//expected success list create
func createList(t *testing.T, r http.Handler, input string, headers map[string]string) int64 {
	code, listCreateData := helpers.MakeRequest(
		r,
		t,
		http.MethodPost,
		"/api/lists",
		bytes.NewBuffer([]byte(input)),
		headers,
	)
	require.Equal(t, http.StatusOK, code)

	crResp := &handler.ListCreateResponse{}
	err := json.Unmarshal(listCreateData, crResp)
	require.NoError(t, err)

	return crResp.ListID
}

func (suite *TestingSuite) TestGetListByID() {
	listInput := fmt.Sprintf(
		`{"title": "%s", "description": "%s"}`,
		listForCreate.Title, listForCreate.Description,
	)

	siginInput := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		createListUser.Email, defaultPassword,
	)
	authResp := makeSignIn(suite.T(), suite.router, siginInput)

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", authResp.AccessToken),
	}

	listID := createList(suite.T(), suite.router, listInput, headers)

	tests := []struct {
		name        string
		code        int
		inputListID int64
		expErrMsg   string
		expList     *models.List
	}{
		{
			name:        "List not found",
			code:        http.StatusNotFound,
			inputListID: 100000,
			expErrMsg:   models.ErrNoList.Error(),
			expList:     nil,
		},
		{
			name:        "Success get",
			code:        http.StatusOK,
			inputListID: listID,
			expErrMsg:   "",
			expList:     listForCreate,
		},
	}

	for _, tc := range tests {
		suite.T().Run(tc.name, func(t *testing.T) {
			code, listGetData := helpers.MakeRequest(
				suite.router,
				t,
				http.MethodGet,
				fmt.Sprintf("/api/lists/%d", tc.inputListID),
				bytes.NewBuffer([]byte{}),
				headers,
			)

			require.Equal(t, tc.code, code)

			if tc.expErrMsg != "" {
				errResp := &handler.ErrorResponse{}
				err := json.Unmarshal(listGetData, errResp)
				require.NoError(t, err)
				require.Equal(t, tc.expErrMsg, errResp.Message)
			} else {
				userList := &models.List{}
				err := json.Unmarshal(listGetData, userList)
				require.NoError(t, err)
				require.Equal(t, tc.expList.Title, userList.Title)
				require.Equal(t, tc.expList.Description, userList.Description)
			}
		})
	}

	makeLogout(suite.T(), suite.router, authResp)
}
