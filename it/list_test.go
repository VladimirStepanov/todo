package it

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

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

func (suite *TestingSuite) TestCreateAndGetList() {
	listInput := fmt.Sprintf(
		`{"title": "%s", "description": "%s"}`,
		listForCreate.Title, listForCreate.Description,
	)

	siginInput := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		createListUser.Email, defaultPassword,
	)
	code, signInData := helpers.MakeRequest(
		suite.router,
		suite.T(),
		http.MethodPost,
		"/auth/sign-in",
		bytes.NewBuffer([]byte(siginInput)),
		nil,
	)
	require.Equal(suite.T(), http.StatusOK, code)

	authResp := &handler.TokensResponse{}
	err := json.Unmarshal(signInData, authResp)
	require.NoError(suite.T(), err)

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", authResp.AccessToken),
	}

	code, listCreateData := helpers.MakeRequest(
		suite.router,
		suite.T(),
		http.MethodPost,
		"/api/lists",
		bytes.NewBuffer([]byte(listInput)),
		headers,
	)
	require.Equal(suite.T(), http.StatusOK, code)

	crResp := &handler.ListCreateResponse{}
	err = json.Unmarshal(listCreateData, crResp)
	require.NoError(suite.T(), err)

	code, listGetData := helpers.MakeRequest(
		suite.router,
		suite.T(),
		http.MethodGet,
		fmt.Sprintf("/api/lists/%d", crResp.ListID),
		bytes.NewBuffer([]byte{}),
		headers,
	)
	require.Equal(suite.T(), http.StatusOK, code)

	userList := &models.List{}
	err = json.Unmarshal(listGetData, userList)
	require.NoError(suite.T(), err)

	require.Equal(suite.T(), listForCreate.Title, userList.Title)
	require.Equal(suite.T(), listForCreate.Description, userList.Description)

	makeLogout(suite.T(), suite.router, signInData)
}
