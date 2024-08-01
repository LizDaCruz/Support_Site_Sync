package confluence

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "bytes"
    "io/ioutil"
)

// ConfluenceServiceImpl is the implementation of the ConfluenceService interface
type ConfluenceServiceImpl struct {
    baseURL   string
    username   string
    apiToken   string
}

// NewConfluenceService creates a new instance of ConfluenceService
func NewConfluenceService(baseURL, username, apiToken string) *ConfluenceServiceImpl {
    return &ConfluenceServiceImpl{baseURL: baseURL, username: username, apiToken: apiToken}
}

// CreatePage creates a new page in Confluence
func (s *ConfluenceServiceImpl) CreatePage(ctx context.Context, page Page) (string, error) {
    url := fmt.Sprintf("%s/wiki/rest/api/content", s.baseURL)
    reqBody, _ := json.Marshal(map[string]interface{}{
        "type":    "page",
        "title":   page.Title,
        "body":    map[string]interface{}{"storage": map[string]string{"value": page.Content, "representation": "storage"}},
        "version": map[string]int{"number": 1},
    })

    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
    req.SetBasicAuth(s.username, s.apiToken)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("failed to create page: %s", resp.Status)
    }

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)

    id, ok := result["id"].(string)
    if !ok {
        return "", fmt.Errorf("failed to parse page ID")
    }

    return id, nil
}

// UpdatePage updates an existing page in Confluence
func (s *ConfluenceServiceImpl) UpdatePage(ctx context.Context, page Page) error {
    pageID := page.ID
    url := fmt.Sprintf("%s/wiki/rest/api/content/%s", s.baseURL, pageID)

    // Fetch current version
    req, _ := http.NewRequest("GET", url, nil)
    req.SetBasicAuth(s.username, s.apiToken)
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to get page: %s", resp.Status)
    }

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)

    currentVersion := result["version"].(map[string]interface{})["number"].(float64)
    newVersion := int(currentVersion) + 1

    reqBody, _ := json.Marshal(map[string]interface{}{
        "version": map[string]int{"number": newVersion},
        "title":   page.Title,
        "body":    map[string]interface{}{"storage": map[string]string{"value": page.Content, "representation": "storage"}},
    })

    req, _ = http.NewRequest("PUT", url, bytes.NewBuffer(reqBody))
    req.SetBasicAuth(s.username, s.apiToken)
    req.Header.Set("Content-Type", "application/json")

    resp, err = http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to update page: %s", resp.Status)
    }

    return nil
}

// DeletePage deletes a page in Confluence
func (s *ConfluenceServiceImpl) DeletePage(ctx context.Context, id string) error {
    url := fmt.Sprintf("%s/wiki/rest/api/content/%s", s.baseURL, id)

    req, _ := http.NewRequest("DELETE", url, nil)
    req.SetBasicAuth(s.username, s.apiToken)

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusNoContent {
        return fmt.Errorf("failed to delete page: %s", resp.Status)
    }

    return nil
}

// GetPage retrieves a page from Confluence
func (s *ConfluenceServiceImpl) GetPage(ctx context.Context, id string) (Page, error) {
    url := fmt.Sprintf("%s/wiki/rest/api/content/%s", s.baseURL, id)

    req, _ := http.NewRequest("GET", url, nil)
    req.SetBasicAuth(s.username, s.apiToken)

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return Page{}, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return Page{}, fmt.Errorf("failed to get page: %s", resp.Status)
    }

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)

    pageID := result["id"].(string)
    title := result["title"].(string)
    content := result["body"].(map[string]interface{})["storage"].(map[string]interface{})["value"].(string)

    return Page{
        ID:      pageID,
        Title:   title,
        Content: content,
    }, nil
}
