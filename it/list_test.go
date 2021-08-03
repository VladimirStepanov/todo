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

func (suite *TestingSuite) TestCreateAndGetList() {
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

	code, listGetData := helpers.MakeRequest(
		suite.router,
		suite.T(),
		http.MethodGet,
		fmt.Sprintf("/api/lists/%d", listID),
		bytes.NewBuffer([]byte{}),
		headers,
	)
	require.Equal(suite.T(), http.StatusOK, code)

	userList := &models.List{}
	err := json.Unmarshal(listGetData, userList)
	require.NoError(suite.T(), err)

	require.Equal(suite.T(), listForCreate.Title, userList.Title)
	require.Equal(suite.T(), listForCreate.Description, userList.Description)

	makeLogout(suite.T(), suite.router, authResp)
}
