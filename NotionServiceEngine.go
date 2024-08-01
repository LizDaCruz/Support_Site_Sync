package notion

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "bytes"
    "io/ioutil"
)

// NotionServiceImpl is the implementation of the NotionService interface
type NotionServiceImpl struct {
    baseURL string
    apiKey  string
}

// NewNotionService creates a new instance of NotionService
func NewNotionService(baseURL, apiKey string) *NotionServiceImpl {
    return &NotionServiceImpl{baseURL: baseURL, apiKey: apiKey}
}

// CreatePage creates a new page in Notion
func (s *NotionServiceImpl) CreatePage(ctx context.Context, page Page) (string, error) {
    url := fmt.Sprintf("%s/pages", s.baseURL)
    reqBody, _ := json.Marshal(map[string]interface{}{
        "parent": map[string]interface{}{
            "database_id": "YOUR_DATABASE_ID", // Replace with your actual database ID
        },
        "properties": map[string]interface{}{
            "title": map[string]interface{}{
                "title": []map[string]interface{}{
                    {
                        "text": map[string]interface{}{
                            "content": page.Title,
                        },
                    },
                },
            },
        },
        "children": []map[string]interface{}{
            {
                "object": "block",
                "type":   "paragraph",
                "paragraph": map[string]interface{}{
                    "text": []map[string]interface{}{
                        {
                            "type": "text",
                            "text": map[string]interface{}{
                                "content": page.Content,
                            },
                        },
                    },
                },
            },
        },
    })

    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := ioutil.ReadAll(resp.Body)
        return "", fmt.Errorf("failed to create page: %s - %s", resp.Status, body)
    }

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)

    pageID, ok := result["id"].(string)
    if !ok {
        return "", fmt.Errorf("failed to parse page ID")
    }

    return pageID, nil
}

// UpdatePage updates an existing page in Notion
func (s *NotionServiceImpl) UpdatePage(ctx context.Context, page Page) error {
    url := fmt.Sprintf("%s/pages/%s", s.baseURL, page.ID)

    reqBody, _ := json.Marshal(map[string]interface{}{
        "properties": map[string]interface{}{
            "title": map[string]interface{}{
                "title": []map[string]interface{}{
                    {
                        "text": map[string]interface{}{
                            "content": page.Title,
                        },
                    },
                },
            },
        },
        "children": []map[string]interface{}{
            {
                "object": "block",
                "type":   "paragraph",
                "paragraph": map[string]interface{}{
                    "text": []map[string]interface{}{
                        {
                            "type": "text",
                            "text": map[string]interface{}{
                                "content": page.Content,
                            },
                        },
                    },
                },
            },
        },
    })

    req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(reqBody))
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
    req.Header.Set("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := ioutil.ReadAll(resp.Body)
        return fmt.Errorf("failed to update page: %s - %s", resp.Status, body)
    }

    return nil
}

// DeletePage deletes a page in Notion
func (s *NotionServiceImpl) DeletePage(ctx context.Context, id string) error {
    url := fmt.Sprintf("%s/pages/%s", s.baseURL, id)

    req, _ := http.NewRequest("DELETE", url, nil)
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusNoContent {
        body, _ := ioutil.ReadAll(resp.Body)
        return fmt.Errorf("failed to delete page: %s - %s", resp.Status, body)
    }

    return nil
}

// GetPage retrieves a page from Notion
func (s *NotionServiceImpl) GetPage(ctx context.Context, id string) (Page, error) {
    url := fmt.Sprintf("%s/pages/%s", s.baseURL, id)

    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return Page{}, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := ioutil.ReadAll(resp.Body)
        return Page{}, fmt.Errorf("failed to get page: %s - %s", resp.Status, body)
    }

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)

    properties := result["properties"].(map[string]interface{})
    title := properties["title"].(map[string]interface{})["title"].([]interface{})[0].(map[string]interface{})["text"].(map[string]interface{})["content"].(string)
    content := result["children"].([]interface{})[0].(map[string]interface{})["paragraph"].(map[string]interface{})["text"].([]interface{})[0].(map[string]interface{})["text"].(map[string]interface{})["content"].(string)

    return Page{
        ID:      id,
        Title:   title,
        Content: content,
    }, nil
}
