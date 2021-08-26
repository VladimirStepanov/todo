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

func (suite *TestingSuite) TestEditRole() {
	listInput := fmt.Sprintf(
		`{"title": "%s", "description": "%s"}`,
		listForCreate.Title, listForCreate.Description,
	)

	siginInputUser1 := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		editRoleUser1.Email, defaultPassword,
	)
	siginInputUser2 := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		editRoleUser2.Email, defaultPassword,
	)

	authRespUser1 := makeSignIn(suite.T(), suite.router, siginInputUser1)
	authRespUser2 := makeSignIn(suite.T(), suite.router, siginInputUser2)

	headersUser1 := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", authRespUser1.AccessToken),
	}
	headersUser2 := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", authRespUser2.AccessToken),
	}

	ListID := createList(suite.T(), suite.router, listInput, headersUser1)

	code, _ := helpers.MakeRequest(
		suite.router,
		suite.T(),
		http.MethodGet,
		fmt.Sprintf("/api/lists/%d", ListID),
		bytes.NewBuffer([]byte{}),
		headersUser2,
	)

	require.Equal(suite.T(), http.StatusNotFound, code)

	tests := []struct {
		name      string
		code      int
		input     string
		expErrMsg string
	}{
		{
			name:      "User not found",
			code:      http.StatusNotFound,
			input:     `{"user_id": 77777, "is_admin":true}`,
			expErrMsg: models.ErrUserNotFound.Error(),
		},
		{
			name:      "EditRole success",
			code:      http.StatusOK,
			input:     fmt.Sprintf(`{"user_id": %d, "is_admin":true}`, GetUserID(editRoleUser2.Email)),
			expErrMsg: "",
		},
	}

	for _, tc := range tests {
		suite.T().Run(tc.name, func(t *testing.T) {
			code, editRoleData := helpers.MakeRequest(
				suite.router,
				t,
				http.MethodPatch,
				fmt.Sprintf("/api/lists/%d/edit-role", ListID),
				bytes.NewBuffer([]byte(tc.input)),
				headersUser1,
			)
			require.Equal(t, tc.code, code)

			if tc.expErrMsg != "" {
				errResp := &handler.ErrorResponse{}
				err := json.Unmarshal(editRoleData, errResp)
				require.NoError(t, err)
				require.Equal(t, tc.expErrMsg, errResp.Message)
			} else {
				actResp := map[string]interface{}{}
				err := json.Unmarshal(editRoleData, &actResp)
				require.NoError(t, err)
				require.Equal(t, "success", actResp["status"])
			}
		})
	}

	code, _ = helpers.MakeRequest(
		suite.router,
		suite.T(),
		http.MethodGet,
		fmt.Sprintf("/api/lists/%d", ListID),
		bytes.NewBuffer([]byte{}),
		headersUser2,
	)

	require.Equal(suite.T(), http.StatusOK, code)
}

func (suite *TestingSuite) TestDeleteList() {
	listInput := fmt.Sprintf(
		`{"title": "%s", "description": "%s"}`,
		listForCreate.Title, listForCreate.Description,
	)
	siginInputUser := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		deleteUser.Email, defaultPassword,
	)

	authRespUser := makeSignIn(suite.T(), suite.router, siginInputUser)

	headersUser := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", authRespUser.AccessToken),
	}

	ListID := createList(suite.T(), suite.router, listInput, headersUser)

	tests := []struct {
		name        string
		code        int
		paramListId int
		expErrMsg   string
	}{
		{
			name:        "List not found",
			code:        http.StatusNotFound,
			paramListId: 10000,
			expErrMsg:   models.ErrNoList.Error(),
		},
		{
			name:        "Success delete",
			code:        http.StatusOK,
			paramListId: int(ListID),
			expErrMsg:   "",
		},
	}

	for _, tc := range tests {
		suite.T().Run(tc.name, func(t *testing.T) {
			code, editRoleData := helpers.MakeRequest(
				suite.router,
				t,
				http.MethodDelete,
				fmt.Sprintf("/api/lists/%d", tc.paramListId),
				bytes.NewBuffer([]byte{}),
				headersUser,
			)
			require.Equal(t, tc.code, code)

			if tc.expErrMsg != "" {
				errResp := &handler.ErrorResponse{}
				err := json.Unmarshal(editRoleData, errResp)
				require.NoError(t, err)
				require.Equal(t, tc.expErrMsg, errResp.Message)
			} else {
				actResp := map[string]interface{}{}
				err := json.Unmarshal(editRoleData, &actResp)
				require.NoError(t, err)
				require.Equal(t, "success", actResp["status"])
			}
		})
	}

	suite.T().Run("Check delete result", func(t *testing.T) {
		code, _ := helpers.MakeRequest(
			suite.router,
			t,
			http.MethodGet,
			fmt.Sprintf("/api/lists/%d", ListID),
			bytes.NewBuffer([]byte{}),
			headersUser,
		)
		require.Equal(suite.T(), http.StatusNotFound, code)
	})
}

func (suite *TestingSuite) TestUpdateList() {
	updateList := &models.List{
		ID:          0,
		Title:       "New title",
		Description: "New description",
	}

	listInput := fmt.Sprintf(
		`{"title": "%s", "description": "%s"}`,
		listForCreate.Title, listForCreate.Description,
	)
	siginInputUser := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		updateUser.Email, defaultPassword,
	)

	authRespUser := makeSignIn(suite.T(), suite.router, siginInputUser)

	headersUser := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", authRespUser.AccessToken),
	}

	ListID := createList(suite.T(), suite.router, listInput, headersUser)

	tests := []struct {
		name        string
		code        int
		input       string
		paramListId int64
		expErrMsg   string
	}{
		{
			name:        "Empty arguments",
			code:        http.StatusBadRequest,
			input:       "{}",
			paramListId: ListID,
			expErrMsg:   models.ErrUpdateEmptyArgs.Error(),
		},
		{
			name:        "Title too short",
			code:        http.StatusBadRequest,
			input:       `{"title": "1"}`,
			paramListId: ListID,
			expErrMsg:   models.ErrTitleTooShort.Error(),
		},
		{
			name:        "Update list not found",
			code:        http.StatusNotFound,
			input:       `{"title": "1"}`,
			paramListId: 77777,
			expErrMsg:   models.ErrNoList.Error(),
		},
		{
			name: "Success update",
			code: http.StatusOK,
			input: fmt.Sprintf(
				`{"title": "%s", "description": "%s"}`,
				updateList.Title, updateList.Description,
			),
			paramListId: ListID,
			expErrMsg:   "",
		},
	}

	for _, tc := range tests {
		suite.T().Run(tc.name, func(t *testing.T) {
			code, editRoleData := helpers.MakeRequest(
				suite.router,
				t,
				http.MethodPatch,
				fmt.Sprintf("/api/lists/%d", tc.paramListId),
				bytes.NewBuffer([]byte(tc.input)),
				headersUser,
			)
			require.Equal(t, tc.code, code)

			if tc.expErrMsg != "" {
				errResp := &handler.ErrorResponse{}
				err := json.Unmarshal(editRoleData, errResp)
				require.NoError(t, err)
				require.Equal(t, tc.expErrMsg, errResp.Message)
			} else {
				actResp := map[string]interface{}{}
				err := json.Unmarshal(editRoleData, &actResp)
				require.NoError(t, err)
				require.Equal(t, "success", actResp["status"])
			}
		})
	}

	suite.T().Run("Check update result", func(t *testing.T) {
		code, getRequestData := helpers.MakeRequest(
			suite.router,
			t,
			http.MethodGet,
			fmt.Sprintf("/api/lists/%d", ListID),
			bytes.NewBuffer([]byte{}),
			headersUser,
		)
		require.Equal(suite.T(), http.StatusOK, code)
		l := &models.List{}
		require.NoError(t, json.Unmarshal(getRequestData, l))

		require.Equal(t, updateList.Title, l.Title)
		require.Equal(t, updateList.Description, l.Description)
	})
}

func (suite *TestingSuite) TestGetUserLists() {
	siginInputUser := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		updateUser.Email, defaultPassword,
	)

	authRespUser := makeSignIn(suite.T(), suite.router, siginInputUser)

	headersUser := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", authRespUser.AccessToken),
	}

	for _, l := range helpers.ExpLists {
		listInput := fmt.Sprintf(
			`{"title": "%s", "description": "%s"}`,
			l.Title, l.Description,
		)
		createList(suite.T(), suite.router, listInput, headersUser)
	}

	suite.T().Run("Check get users lists", func(t *testing.T) {
		code, data := helpers.MakeRequest(
			suite.router,
			t,
			http.MethodGet,
			"/api/lists",
			bytes.NewBuffer([]byte{}),
			headersUser,
		)
		require.Equal(t, http.StatusOK, code)
		resp := handler.UserListsResponse{}
		require.NoError(t, json.Unmarshal(data, &resp))
		require.Equal(t, "success", resp.Status)
		for i, l := range resp.Result {
			require.Equal(t, helpers.ExpLists[i].Title, l.Title)
			require.Equal(t, helpers.ExpLists[i].Description, l.Description)
		}
	})
}
