package tests

import (
    "net/http"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestAuthServiceHealth(t *testing.T) {
    client := NewTestClient()
    resp, err := client.DoRequest(http.MethodGet, authBaseURL+"/health", nil)
    require.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUserLogin(t *testing.T) {
    client := NewTestClient()
    email, username := generateUniqueCredentials()
    password := "testpass123"

    // Register a user first
    registerBody := map[string]string{
        "email":    email,
        "username": username,
        "password": password,
    }
    err := client.DoJSON(http.MethodPost, userBaseURL+"/register", registerBody, nil)
    require.NoError(t, err)

    t.Run("successful login", func(t *testing.T) {
        loginBody := map[string]string{
            "email":    email,
            "password": password,
        }

        var response map[string]string
        err := client.DoJSON(http.MethodPost, authBaseURL+"/login", loginBody, &response)
        require.NoError(t, err)
        assert.NotEmpty(t, response["token"])
    })

    t.Run("invalid credentials", func(t *testing.T) {
        testCases := []struct {
            name     string
            email    string
            password string
        }{
            {
                name:     "wrong password",
                email:    email,
                password: "wrongpass",
            },
            {
                name:     "non-existent email",
                email:    "nonexistent@example.com",
                password: password,
            },
        }

        for _, tc := range testCases {
            t.Run(tc.name, func(t *testing.T) {
                loginBody := map[string]string{
                    "email":    tc.email,
                    "password": tc.password,
                }

                resp, err := client.DoRequest(http.MethodPost, authBaseURL+"/login", loginBody)
                require.NoError(t, err)
                assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
            })
        }
    })
} 