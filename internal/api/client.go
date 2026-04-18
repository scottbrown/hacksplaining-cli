package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const BaseURL = "https://hacksplaining.com/api/v1"

type Client struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type User struct {
	Email               string            `json:"email"`
	Progress            int               `json:"progress"`
	Complete            bool              `json:"complete"`
	CompletionDate      *string           `json:"completion_date"`
	CreatedAt           string            `json:"created_at"`
	LastCommunication   *string           `json:"last_communication"`
	InvitationAccepted  *string           `json:"invitation_accepted_at"`
	Group               *Group            `json:"group,omitempty"`
	Exercises           map[string]string `json:"exercises"`
}

type AddUserRequest struct {
	Subject string `json:"subject,omitempty"`
	Message string `json:"message,omitempty"`
	GroupID int    `json:"groupId,omitempty"`
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    BaseURL,
	}
}

func (c *Client) authHeader() string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(c.apiKey+":"))
}

func (c *Client) doRequest(method, path string, body io.Reader) (*http.Response, error) {
	url := c.baseURL + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", c.authHeader())
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.httpClient.Do(req)
}

func (c *Client) ListUsers() ([]User, error) {
	resp, err := c.doRequest(http.MethodGet, "/users", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, parseError(resp)
	}

	var users []User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return users, nil
}

func (c *Client) GetUser(email string) (*User, error) {
	resp, err := c.doRequest(http.MethodGet, "/users/"+url.PathEscape(email), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, parseError(resp)
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &user, nil
}

func (c *Client) AddUser(email string, req *AddUserRequest) (int, error) {
	var body io.Reader
	if req != nil {
		data, err := json.Marshal(req)
		if err != nil {
			return 0, fmt.Errorf("encoding request body: %w", err)
		}
		body = strings.NewReader(string(data))
	}

	resp, err := c.doRequest(http.MethodPut, "/users/"+url.PathEscape(email), body)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		return resp.StatusCode, nil
	default:
		return resp.StatusCode, parseError(resp)
	}
}

func (c *Client) RemoveUser(email string) error {
	resp, err := c.doRequest(http.MethodDelete, "/users/"+url.PathEscape(email), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return parseError(resp)
	}
	return nil
}

func (c *Client) RemindUser(email string, subject, message string) (int, error) {
	var body io.Reader
	if subject != "" || message != "" {
		payload := map[string]string{}
		if subject != "" {
			payload["subject"] = subject
		}
		if message != "" {
			payload["message"] = message
		}
		data, err := json.Marshal(payload)
		if err != nil {
			return 0, fmt.Errorf("encoding request body: %w", err)
		}
		body = strings.NewReader(string(data))
	}

	resp, err := c.doRequest(http.MethodPut, "/users/"+url.PathEscape(email)+"/reminder", body)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return resp.StatusCode, nil
	case 302:
		return resp.StatusCode, fmt.Errorf("user has already completed training")
	default:
		return resp.StatusCode, parseError(resp)
	}
}

func parseError(resp *http.Response) error {
	return fmt.Errorf("HTTP %d", resp.StatusCode)
}
