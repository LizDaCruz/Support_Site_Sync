package helpjuice

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "bytes"
    "io/ioutil"
)

// HelpjuiceServiceImpl is the implementation of the HelpjuiceService interface
type HelpjuiceServiceImpl struct {
    baseURL   string
    apiKey    string
}

// NewHelpjuiceService creates a new instance of HelpjuiceService
func NewHelpjuiceService(baseURL, apiKey string) *HelpjuiceServiceImpl {
    return &HelpjuiceServiceImpl{baseURL: baseURL, apiKey: apiKey}
}

// CreatePage creates a new page in Helpjuice
func (s *HelpjuiceServiceImpl) CreatePage(ctx context.Context, page Page) (string, error) {
    url := fmt.Sprintf("%s/api/v1/articles", s.baseURL)
    reqBody, _ := json.Marshal(map[string]interface{}{
        "title":   page.Title,
        "content": page.Content,
        "status":  "published",
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

    if resp.StatusCode != http.StatusCreated {
        body, _ := ioutil.ReadAll(resp.Body)
        return "", fmt.Errorf("failed to create page: %s - %s", resp.Status, body)
    }

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)

    articleID, ok := result["id"].(string)
    if !ok {
        return "", fmt.Errorf("failed to parse article ID")
    }

    return articleID, nil
}

// UpdatePage updates an existing page in Helpjuice
func (s *HelpjuiceServiceImpl) UpdatePage(ctx context.Context, page Page) error {
    url := fmt.Sprintf("%s/api/v1/articles/%s", s.baseURL, page.ID)

    reqBody, _ := json.Marshal(map[string]interface{}{
        "title":   page.Title,
        "content": page.Content,
        "status":  "published",
    })

    req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(reqBody))
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

// DeletePage deletes a page in Helpjuice
func (s *HelpjuiceServiceImpl) DeletePage(ctx context.Context, id string) error {
    url := fmt.Sprintf("%s/api/v1/articles/%s", s.baseURL, id)

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

// GetPage retrieves a page from Helpjuice
func (s *HelpjuiceServiceImpl) GetPage(ctx context.Context, id string) (Page, error) {
    url := fmt.Sprintf("%s/api/v1/articles/%s", s.baseURL, id)

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

    article := result["article"].(map[string]interface{})
    pageID := article["id"].(string)
    title := article["title"].(string)
    content := article["content"].(string)

    return Page{
        ID:      pageID,
        Title:   title,
        Content: content,
    }, nil
}
