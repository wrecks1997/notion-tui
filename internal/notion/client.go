package notion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseURL        = "https://api.notion.com/v1"
	notionVersion  = "2022-06-28"
	defaultTimeout = 30 * time.Second
)

type Client struct {
	token      string
	httpClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		token: token,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

func (c *Client) do(method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, baseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Notion-Version", notionVersion)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("notion API error (%d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (c *Client) Search(query string) (*SearchResponse, error) {
	req := SearchRequest{
		Query: query,
		Sort: &SearchSort{
			Direction: "descending",
			Timestamp: "last_edited_time",
		},
		PageSize: 50,
	}

	data, err := c.do("POST", "/search", req)
	if err != nil {
		return nil, err
	}

	var result SearchResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode search response: %w", err)
	}

	return &result, nil
}

func (c *Client) GetPage(pageID string) (*Page, error) {
	data, err := c.do("GET", "/pages/"+pageID, nil)
	if err != nil {
		return nil, err
	}

	var page Page
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("decode page: %w", err)
	}

	return &page, nil
}

func (c *Client) GetDatabase(dbID string) (*Database, error) {
	data, err := c.do("GET", "/databases/"+dbID, nil)
	if err != nil {
		return nil, err
	}

	var db Database
	if err := json.Unmarshal(data, &db); err != nil {
		return nil, fmt.Errorf("decode database: %w", err)
	}

	return &db, nil
}

func (c *Client) GetBlockChildren(blockID string) ([]Block, error) {
	var allBlocks []Block
	cursor := ""

	for {
		path := "/blocks/" + blockID + "/children?page_size=100"
		if cursor != "" {
			path += "&start_cursor=" + cursor
		}

		data, err := c.do("GET", path, nil)
		if err != nil {
			return nil, err
		}

		var resp BlocksResponse
		if err := json.Unmarshal(data, &resp); err != nil {
			return nil, fmt.Errorf("decode blocks: %w", err)
		}

		allBlocks = append(allBlocks, resp.Results...)

		if !resp.HasMore {
			break
		}
		cursor = resp.NextCursor
	}

	return allBlocks, nil
}
