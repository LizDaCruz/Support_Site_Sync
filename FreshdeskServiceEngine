package freshdesk

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "bytes"
    "io/ioutil"
)

// FreshdeskServiceImpl is the implementation of the FreshdeskService interface
type FreshdeskServiceImpl struct {
    baseURL    string
    apiKey      string
}

// NewFreshdeskService creates a new instance of FreshdeskService
func NewFreshdeskService(baseURL, apiKey string) *FreshdeskServiceImpl {
    return &FreshdeskServiceImpl{baseURL: baseURL, apiKey: apiKey}
}

// CreatePage creates a new page in Freshdesk
func (s *FreshdeskServiceImpl) CreatePage(ctx context.Context, page Page) (string, error) {
    url := fmt.Sprintf("%s/api/v2/solutions/articles", s.baseURL)
    reqBody, _ := json.Marshal(map[string]interface{}{
        "title":   page.Title,
        "description": page.Content,
        "status":  2, // Published
    })

    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
    req.SetBasicAuth(s.apiKey, "X")
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusCreated {
        body, _ := ioutil.ReadAll(resp.Body)
        return "", fmt.Errorf("failed to create page: %s - %s", resp.Status, body)
    }

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)

    articleID, ok := result["id"].(float64)
    if !ok {
        return "", fmt.Errorf("failed to parse article ID")
    }

    return fmt.Sprintf("%.0f", articleID), nil
}

// UpdatePage updates an existing page in Freshdesk
func (s *FreshdeskServiceImpl) UpdatePage(ctx context.Context, page Page) error {
    url := fmt.Sprintf("%s/api/v2/solutions/articles/%s", s.baseURL, page.ID)

    reqBody, _ := json.Marshal(map[string]interface{}{
        "title":       page.Title,
        "description": page.Content,
        "status":      2, // Published
    })

    req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(reqBody))
    req.SetBasicAuth(s.apiKey, "X")
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

// DeletePage deletes a page in Freshdesk
func (s *FreshdeskServiceImpl) DeletePage(ctx context.Context, id string) error {
    url := fmt.Sprintf("%s/api/v2/solutions/articles/%s", s.baseURL, id)

    req, _ := http.NewRequest("DELETE", url, nil)
    req.SetBasicAuth(s.apiKey, "X")

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

// GetPage retrieves a page from Freshdesk
func (s *FreshdeskServiceImpl) GetPage(ctx context.Context, id string) (Page, error) {
    url := fmt.Sprintf("%s/api/v2/solutions/articles/%s", s.baseURL, id)

    req, _ := http.NewRequest("GET", url, nil)
    req.SetBasicAuth(s.apiKey, "X")

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

    article := result["article"].(map[string]interface{})
    pageID := article["id"].(float64)
    title := article["title"].(string)
    content := article["description"].(string)

    return Page{
        ID:      fmt.Sprintf("%.0f", pageID),
        Title:   title,
        Content: content,
    }, nil
}
