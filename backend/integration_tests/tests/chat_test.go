package tests

import (
    "fmt"
    "net/http"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestChatManagement(t *testing.T) {
    // Create two test users
    client1 := NewTestClient()
    token1, _ := registerTestUser(t, client1)
    client1.SetToken(token1)

    client2 := NewTestClient()
    token2, _ := registerTestUser(t, client2)
    client2.SetToken(token2)

    var chatID string

    t.Run("create direct chat", func(t *testing.T) {
        createBody := map[string]interface{}{
            "title":    "Test Direct Chat",
            "members": []string{"user2"},
        }

        var response map[string]interface{}
        err := client1.DoJSON(http.MethodPost, restBaseURL+"/chats", createBody, &response)
        require.NoError(t, err)
        assert.NotEmpty(t, response["id"])
        assert.Equal(t, createBody["title"], response["title"])
        assert.False(t, response["is_group"].(bool))

        chatID = response["id"].(string)
    })

    t.Run("get chat list", func(t *testing.T) {
        var response []map[string]interface{}
        err := client1.DoJSON(http.MethodGet, restBaseURL+"/chats", nil, &response)
        require.NoError(t, err)
        assert.NotEmpty(t, response)

        found := false
        for _, chat := range response {
            if chat["id"].(string) == chatID {
                found = true
                break
            }
        }
        assert.True(t, found, "Created chat should be in the list")
    })

    t.Run("get single chat", func(t *testing.T) {
        var response map[string]interface{}
        err := client1.DoJSON(http.MethodGet, fmt.Sprintf("%s/chats/%s", restBaseURL, chatID), nil, &response)
        require.NoError(t, err)
        assert.Equal(t, chatID, response["id"])
    })

    t.Run("update chat title", func(t *testing.T) {
        updateBody := map[string]string{
            "title": "Updated Chat Title",
        }

        var response map[string]interface{}
        err := client1.DoJSON(http.MethodPut, fmt.Sprintf("%s/chats/%s", restBaseURL, chatID), updateBody, &response)
        require.NoError(t, err)
        assert.Equal(t, updateBody["title"], response["title"])

        // Verify the change is persisted
        err = client1.DoJSON(http.MethodGet, fmt.Sprintf("%s/chats/%s", restBaseURL, chatID), nil, &response)
        require.NoError(t, err)
        assert.Equal(t, updateBody["title"], response["title"])
    })

    t.Run("create group chat", func(t *testing.T) {
        // Create a third user for the group
        client3 := NewTestClient()
        token3, _ := registerTestUser(t, client3)
        client3.SetToken(token3)

        createBody := map[string]interface{}{
            "title":    "Test Group Chat",
            "members": []string{"user2", "user3"},
        }

        var response map[string]interface{}
        err := client1.DoJSON(http.MethodPost, restBaseURL+"/chats", createBody, &response)
        require.NoError(t, err)
        assert.NotEmpty(t, response["id"])
        assert.Equal(t, createBody["title"], response["title"])
        assert.True(t, response["is_group"].(bool))

        groupChatID := response["id"].(string)

        // Test adding a new member
        addMemberBody := map[string]interface{}{
            "username": "user1", // Add the first user again (should fail)
        }
        resp, err := client1.DoRequest(http.MethodPost, fmt.Sprintf("%s/chats/%s/members", restBaseURL, groupChatID), addMemberBody)
        require.NoError(t, err)
        assert.Equal(t, http.StatusConflict, resp.StatusCode)

        // Test removing a member
        resp, err = client1.DoRequest(http.MethodDelete, fmt.Sprintf("%s/chats/%s/members/%s", restBaseURL, groupChatID, "user2"), nil)
        require.NoError(t, err)
        assert.Equal(t, http.StatusOK, resp.StatusCode)

        // Verify member was removed
        var chatResponse map[string]interface{}
        err = client1.DoJSON(http.MethodGet, fmt.Sprintf("%s/chats/%s", restBaseURL, groupChatID), nil, &chatResponse)
        require.NoError(t, err)
        members := chatResponse["members"].([]interface{})
        found := false
        for _, member := range members {
            if member.(map[string]interface{})["username"] == "user2" {
                found = true
                break
            }
        }
        assert.False(t, found, "Removed member should not be in the list")
    })

    t.Run("leave chat", func(t *testing.T) {
        // Create a new chat for leaving
        createBody := map[string]interface{}{
            "title":    "Chat to Leave",
            "members": []string{"user2"},
        }

        var response map[string]interface{}
        err := client1.DoJSON(http.MethodPost, restBaseURL+"/chats", createBody, &response)
        require.NoError(t, err)
        leaveChatID := response["id"].(string)

        // Leave the chat
        resp, err := client1.DoRequest(http.MethodPost, fmt.Sprintf("%s/chats/%s/leave", restBaseURL, leaveChatID), nil)
        require.NoError(t, err)
        assert.Equal(t, http.StatusOK, resp.StatusCode)

        // Verify chat is not accessible anymore
        resp, err = client1.DoRequest(http.MethodGet, fmt.Sprintf("%s/chats/%s", restBaseURL, leaveChatID), nil)
        require.NoError(t, err)
        assert.Equal(t, http.StatusNotFound, resp.StatusCode)
    })

    t.Run("invalid operations", func(t *testing.T) {
        testCases := []struct {
            name           string
            method        string
            endpoint      string
            body          interface{}
            expectedCode int
        }{
            {
                name:      "create chat with non-existent user",
                method:   http.MethodPost,
                endpoint: restBaseURL + "/chats",
                body: map[string]interface{}{
                    "title":    "Invalid Chat",
                    "members": []string{"nonexistent"},
                },
                expectedCode: http.StatusNotFound,
            },
            {
                name:      "get non-existent chat",
                method:   http.MethodGet,
                endpoint: restBaseURL + "/chats/nonexistent",
                expectedCode: http.StatusNotFound,
            },
            {
                name:      "update non-existent chat",
                method:   http.MethodPut,
                endpoint: restBaseURL + "/chats/nonexistent",
                body: map[string]string{
                    "title": "Updated Title",
                },
                expectedCode: http.StatusNotFound,
            },
        }

        for _, tc := range testCases {
            t.Run(tc.name, func(t *testing.T) {
                resp, err := client1.DoRequest(tc.method, tc.endpoint, tc.body)
                require.NoError(t, err)
                assert.Equal(t, tc.expectedCode, resp.StatusCode)
            })
        }
    })
} 