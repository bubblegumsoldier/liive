package tests

import (
    "net/http"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestUserRegistration(t *testing.T) {
    client := NewTestClient()

    t.Run("successful registration", func(t *testing.T) {
        email, username := generateUniqueCredentials()
        body := map[string]string{
            "email":    email,
            "username": username,
            "password": "testpass123",
        }

        var response map[string]interface{}
        err := client.DoJSON(http.MethodPost, userBaseURL+"/register", body, &response)
        require.NoError(t, err)
        assert.NotEmpty(t, response["id"])
        assert.Equal(t, email, response["email"])
        assert.Equal(t, username, response["username"])
    })

    t.Run("duplicate email", func(t *testing.T) {
        email, _ := generateUniqueCredentials()
        
        // First registration
        body := map[string]string{
            "email":    email,
            "username": "user1",
            "password": "testpass123",
        }
        err := client.DoJSON(http.MethodPost, userBaseURL+"/register", body, nil)
        require.NoError(t, err)

        // Second registration with same email
        body["username"] = "user2"
        resp, err := client.DoRequest(http.MethodPost, userBaseURL+"/register", body)
        require.NoError(t, err)
        assert.Equal(t, http.StatusConflict, resp.StatusCode)
    })

    t.Run("duplicate username", func(t *testing.T) {
        _, username := generateUniqueCredentials()
        
        // First registration
        body := map[string]string{
            "email":    "test1@example.com",
            "username": username,
            "password": "testpass123",
        }
        err := client.DoJSON(http.MethodPost, userBaseURL+"/register", body, nil)
        require.NoError(t, err)

        // Second registration with same username
        body["email"] = "test2@example.com"
        resp, err := client.DoRequest(http.MethodPost, userBaseURL+"/register", body)
        require.NoError(t, err)
        assert.Equal(t, http.StatusConflict, resp.StatusCode)
    })

    t.Run("invalid data", func(t *testing.T) {
        testCases := []struct {
            name string
            body map[string]string
        }{
            {
                name: "missing email",
                body: map[string]string{
                    "username": "testuser",
                    "password": "testpass123",
                },
            },
            {
                name: "invalid email format",
                body: map[string]string{
                    "email":    "invalid-email",
                    "username": "testuser",
                    "password": "testpass123",
                },
            },
            {
                name: "password too short",
                body: map[string]string{
                    "email":    "test@example.com",
                    "username": "testuser",
                    "password": "short",
                },
            },
        }

        for _, tc := range testCases {
            t.Run(tc.name, func(t *testing.T) {
                resp, err := client.DoRequest(http.MethodPost, userBaseURL+"/register", tc.body)
                require.NoError(t, err)
                assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
            })
        }
    })
}

func TestUserProfile(t *testing.T) {
    client := NewTestClient()

    // Register and login a test user
    token, _ := registerTestUser(t, client)
    client.SetToken(token)

    t.Run("get profile", func(t *testing.T) {
        var response map[string]interface{}
        err := client.DoJSON(http.MethodGet, userBaseURL+"/profile", nil, &response)
        require.NoError(t, err)
        assert.NotEmpty(t, response["email"])
        assert.NotEmpty(t, response["username"])
    })

    t.Run("update profile", func(t *testing.T) {
        updateBody := map[string]interface{}{
            "username": "updated_username",
            "bio":     "Test bio",
        }

        var response map[string]interface{}
        err := client.DoJSON(http.MethodPut, userBaseURL+"/profile", updateBody, &response)
        require.NoError(t, err)
        assert.Equal(t, updateBody["username"], response["username"])
        assert.Equal(t, updateBody["bio"], response["bio"])

        // Verify the changes are persisted
        err = client.DoJSON(http.MethodGet, userBaseURL+"/profile", nil, &response)
        require.NoError(t, err)
        assert.Equal(t, updateBody["username"], response["username"])
        assert.Equal(t, updateBody["bio"], response["bio"])
    })

    t.Run("update profile with invalid data", func(t *testing.T) {
        testCases := []struct {
            name string
            body map[string]interface{}
        }{
            {
                name: "username too short",
                body: map[string]interface{}{
                    "username": "a",
                },
            },
            {
                name: "bio too long",
                body: map[string]interface{}{
                    "bio": string(make([]byte, 1001)), // 1001 characters
                },
            },
        }

        for _, tc := range testCases {
            t.Run(tc.name, func(t *testing.T) {
                resp, err := client.DoRequest(http.MethodPut, userBaseURL+"/profile", tc.body)
                require.NoError(t, err)
                assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
            })
        }
    })
}

func TestChangePassword(t *testing.T) {
    client := NewTestClient()
    oldPassword := "testpass123"
    newPassword := "newtestpass123"

    // Register and login a test user
    token, email := registerTestUser(t, client)
    client.SetToken(token)

    t.Run("successful password change", func(t *testing.T) {
        changeBody := map[string]string{
            "old_password": oldPassword,
            "new_password": newPassword,
        }

        // Change password
        resp, err := client.DoRequest(http.MethodPost, userBaseURL+"/change-password", changeBody)
        require.NoError(t, err)
        assert.Equal(t, http.StatusOK, resp.StatusCode)

        // Try logging in with new password
        loginBody := map[string]string{
            "email":    email,
            "password": newPassword,
        }
        var loginResponse map[string]string
        err = client.DoJSON(http.MethodPost, authBaseURL+"/login", loginBody, &loginResponse)
        require.NoError(t, err)
        assert.NotEmpty(t, loginResponse["token"])
    })

    t.Run("invalid old password", func(t *testing.T) {
        changeBody := map[string]string{
            "old_password": "wrongpass",
            "new_password": "newpass123",
        }

        resp, err := client.DoRequest(http.MethodPost, userBaseURL+"/change-password", changeBody)
        require.NoError(t, err)
        assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
    })

    t.Run("invalid new password", func(t *testing.T) {
        changeBody := map[string]string{
            "old_password": newPassword,
            "new_password": "short",
        }

        resp, err := client.DoRequest(http.MethodPost, userBaseURL+"/change-password", changeBody)
        require.NoError(t, err)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })
} 