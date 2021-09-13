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
	itemForCreate = &models.Item{
		ID:          0,
		ListID:      listForCreate.ID,
		Title:       "item title",
		Description: "item description",
		Done:        false,
	}
)

//expected success item create
func createItem(t *testing.T, listID int64, r http.Handler, input string, headers map[string]string) int64 {
	code, data := helpers.MakeRequest(
		r,
		t,
		http.MethodPost,
		fmt.Sprintf("/api/lists/%d/items", listID),
		bytes.NewBuffer([]byte(input)),
		headers,
	)
	require.Equal(t, http.StatusOK, code)

	crResp := &handler.ItemCreateResponse{}
	err := json.Unmarshal(data, crResp)
	require.NoError(t, err)

	return crResp.ItemID
}

func (suite *TestingSuite) TestGetItemByID() {
	listInput := fmt.Sprintf(
		`{"title": "%s", "description": "%s"}`,
		listForCreate.Title, listForCreate.Description,
	)

	itemInput := fmt.Sprintf(
		`{"title": "%s", "description": "%s"}`,
		itemForCreate.Title, itemForCreate.Description,
	)

	signInInput := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		getItemUser.Email, defaultPassword,
	)

	authResp := makeSignIn(suite.T(), suite.router, signInInput)

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", authResp.AccessToken),
	}

	listID := createList(suite.T(), suite.router, listInput, headers)
	itemID := createItem(suite.T(), listID, suite.router, itemInput, headers)

	tests := []struct {
		name      string
		code      int
		listID    int64
		itemID    int64
		expErrMsg string
		expItem   *models.Item
	}{
		{
			name:      "List not found",
			code:      http.StatusNotFound,
			listID:    100000,
			itemID:    itemID,
			expErrMsg: models.ErrNoList.Error(),
			expItem:   nil,
		},
		{
			name:      "Item not found",
			code:      http.StatusNotFound,
			listID:    listID,
			itemID:    itemID + 77777,
			expErrMsg: models.ErrNoItem.Error(),
			expItem:   nil,
		},
		{
			name:      "Success get",
			code:      http.StatusOK,
			listID:    listID,
			itemID:    itemID,
			expErrMsg: "",
			expItem:   itemForCreate,
		},
	}

	for _, tc := range tests {
		suite.T().Run(tc.name, func(t *testing.T) {
			code, responseData := helpers.MakeRequest(
				suite.router,
				t,
				http.MethodGet,
				fmt.Sprintf("/api/lists/%d/items/%d", tc.listID, tc.itemID),
				bytes.NewBuffer([]byte{}),
				headers,
			)

			require.Equal(t, tc.code, code)

			if tc.expErrMsg != "" {
				errResp := &handler.ErrorResponse{}
				err := json.Unmarshal(responseData, errResp)
				require.NoError(t, err)
				require.Equal(t, tc.expErrMsg, errResp.Message)
			} else {
				item := &models.Item{}
				err := json.Unmarshal(responseData, item)
				require.NoError(t, err)
				require.Equal(t, itemID, item.ID)
				require.Equal(t, tc.expItem.Title, item.Title)
				require.Equal(t, tc.expItem.Description, item.Description)
				require.Equal(t, tc.expItem.Done, item.Done)
			}
		})
	}

	makeLogout(suite.T(), suite.router, authResp)
}

func (suite *TestingSuite) TestDeleteItem() {
	listInput := fmt.Sprintf(
		`{"title": "%s", "description": "%s"}`,
		listForCreate.Title, listForCreate.Description,
	)

	itemInput := fmt.Sprintf(
		`{"title": "%s", "description": "%s"}`,
		itemForCreate.Title, itemForCreate.Description,
	)

	signInInput := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		deleteItemUser.Email, defaultPassword,
	)

	authResp := makeSignIn(suite.T(), suite.router, signInInput)

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", authResp.AccessToken),
	}

	listID := createList(suite.T(), suite.router, listInput, headers)
	itemID := createItem(suite.T(), listID, suite.router, itemInput, headers)

	tests := []struct {
		name      string
		code      int
		listID    int64
		itemID    int64
		expErrMsg string
	}{
		{
			name:      "Item not found",
			code:      http.StatusNotFound,
			listID:    listID,
			itemID:    itemID + 77777,
			expErrMsg: models.ErrNoItem.Error(),
		},
		{
			name:      "Success delete",
			code:      http.StatusOK,
			listID:    listID,
			itemID:    itemID,
			expErrMsg: "",
		},
	}

	for _, tc := range tests {
		suite.T().Run(tc.name, func(t *testing.T) {
			code, responseData := helpers.MakeRequest(
				suite.router,
				t,
				http.MethodDelete,
				fmt.Sprintf("/api/lists/%d/items/%d", tc.listID, tc.itemID),
				bytes.NewBuffer([]byte{}),
				headers,
			)

			require.Equal(t, tc.code, code)

			if tc.expErrMsg != "" {
				errResp := &handler.ErrorResponse{}
				err := json.Unmarshal(responseData, errResp)
				require.NoError(t, err)
				require.Equal(t, tc.expErrMsg, errResp.Message)
			} else {
				t.Run("Check delete item result", func(t *testing.T) {
					code, _ := helpers.MakeRequest(
						suite.router,
						t,
						http.MethodGet,
						fmt.Sprintf("/api/lists/%d/items/%d", tc.listID, tc.itemID),
						bytes.NewBuffer([]byte{}),
						headers,
					)
					require.Equal(suite.T(), http.StatusNotFound, code)
				})
			}
		})
	}
}

func (suite *TestingSuite) TestUpdateItem() {
	uitem := models.Item{
		Title:       "Title for update",
		Description: "Description",
	}

	updateInput := fmt.Sprintf(
		`{"title": "%s", "description": "%s"}`,
		uitem.Title, uitem.Description,
	)

	listInput := fmt.Sprintf(
		`{"title": "%s", "description": "%s"}`,
		listForCreate.Title, listForCreate.Description,
	)

	itemInput := fmt.Sprintf(
		`{"title": "%s", "description": "%s"}`,
		itemForCreate.Title, itemForCreate.Description,
	)

	signInInput := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		updateItemUser.Email, defaultPassword,
	)

	authResp := makeSignIn(suite.T(), suite.router, signInInput)

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", authResp.AccessToken),
	}

	listID := createList(suite.T(), suite.router, listInput, headers)
	itemID := createItem(suite.T(), listID, suite.router, itemInput, headers)

	tests := []struct {
		name      string
		code      int
		listID    int64
		itemID    int64
		expErrMsg string
	}{
		{
			name:      "Item not found",
			code:      http.StatusNotFound,
			listID:    listID,
			itemID:    itemID + 77777,
			expErrMsg: models.ErrNoItem.Error(),
		},
		{
			name:      "Success update",
			code:      http.StatusOK,
			listID:    listID,
			itemID:    itemID,
			expErrMsg: "",
		},
	}

	for _, tc := range tests {
		suite.T().Run(tc.name, func(t *testing.T) {
			code, responseData := helpers.MakeRequest(
				suite.router,
				t,
				http.MethodPatch,
				fmt.Sprintf("/api/lists/%d/items/%d", tc.listID, tc.itemID),
				bytes.NewBuffer([]byte(updateInput)),
				headers,
			)

			require.Equal(t, tc.code, code)

			if tc.expErrMsg != "" {
				errResp := &handler.ErrorResponse{}
				err := json.Unmarshal(responseData, errResp)
				require.NoError(t, err)
				require.Equal(t, tc.expErrMsg, errResp.Message)
			} else {
				t.Run("Check update result", func(t *testing.T) {
					code, getData := helpers.MakeRequest(
						suite.router,
						t,
						http.MethodGet,
						fmt.Sprintf("/api/lists/%d/items/%d", tc.listID, tc.itemID),
						bytes.NewBuffer([]byte{}),
						headers,
					)
					require.Equal(t, http.StatusOK, code)
					item := &models.Item{}
					err := json.Unmarshal(getData, item)
					require.NoError(t, err)
					require.Equal(t, uitem.Title, item.Title)
					require.Equal(t, uitem.Description, item.Description)
					require.Equal(t, itemID, item.ID)
				})
			}
		})
	}
}
