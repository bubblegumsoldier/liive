package tests

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "testing"
    "time"

    "github.com/stretchr/testify/require"
)

const (
    authBaseURL      = "http://localhost:8082"
    userBaseURL      = "http://localhost:8083"
    restBaseURL      = "http://localhost:8081"
    defaultTimeout   = 10 * time.Second
)

type TestClient struct {
    httpClient *http.Client
    token     string
}

func NewTestClient() *TestClient {
    return &TestClient{
        httpClient: &http.Client{
            Timeout: defaultTimeout,
        },
    }
}

func (c *TestClient) SetToken(token string) {
    c.token = token
}

func (c *TestClient) DoRequest(method, url string, body interface{}) (*http.Response, error) {
    var reqBody []byte
    var err error

    if body != nil {
        reqBody, err = json.Marshal(body)
        if err != nil {
            return nil, fmt.Errorf("failed to marshal request body: %v", err)
        }
    }

    req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %v", err)
    }

    if body != nil {
        req.Header.Set("Content-Type", "application/json")
    }

    if c.token != "" {
        req.Header.Set("Authorization", "Bearer "+c.token)
    }

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to do request: %v", err)
    }

    return resp, nil
}

func (c *TestClient) DoJSON(method, url string, body, response interface{}) error {
    resp, err := c.DoRequest(method, url, body)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 400 {
        return fmt.Errorf("request failed with status %d", resp.StatusCode)
    }

    if response != nil {
        if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
            return fmt.Errorf("failed to decode response: %v", err)
        }
    }

    return nil
}

// Helper function to generate unique email and username for tests
func generateUniqueCredentials() (string, string) {
    timestamp := time.Now().UnixNano()
    return fmt.Sprintf("test%d@example.com", timestamp),
           fmt.Sprintf("user%d", timestamp)
}

// Helper function to register a test user and return the token and email
func registerTestUser(t *testing.T, client *TestClient) (string, string) {
    email, username := generateUniqueCredentials()
    password := "testpass123"

    // Register using user-manager service
    registerURL := userBaseURL + "/register"
    body := map[string]string{
        "email": email,
        "username": username,
        "password": password,
    }

    var response map[string]interface{}
    err := client.DoJSON(http.MethodPost, registerURL, body, &response)
    require.NoError(t, err, "Failed to register test user")

    // Login using auth service
    loginURL := authBaseURL + "/login"
    loginBody := map[string]string{
        "email": email,
        "password": password,
    }

    var loginResponse map[string]string
    err = client.DoJSON(http.MethodPost, loginURL, loginBody, &loginResponse)
    require.NoError(t, err, "Failed to login test user")

    token := loginResponse["token"]
    require.NotEmpty(t, token, "Token should not be empty")

    return token, email
} 