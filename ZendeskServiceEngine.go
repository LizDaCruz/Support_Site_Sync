package zendesk

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "bytes"
    "io/ioutil"
)

// ZendeskServiceImpl is the implementation of the ZendeskService interface
type ZendeskServiceImpl struct {
    baseURL    string
    email      string
    apiToken    string
}

// NewZendeskService creates a new instance of ZendeskService
func NewZendeskService(baseURL, email, apiToken string) *ZendeskServiceImpl {
    return &ZendeskServiceImpl{baseURL: baseURL, email: email, apiToken: apiToken}
}

// CreatePage creates a new page in Zendesk
func (s *ZendeskServiceImpl) CreatePage(ctx context.Context, page Page) (string, error) {
    url := fmt.Sprintf("%s/api/v2/help_center/articles.json", s.baseURL)
    reqBody, _ := json.Marshal(map[string]interface{}{
        "article": map[string]interface{}{
            "title":   page.Title,
            "body":    page.Content,
            "locale":  "en-us",
        },
    })

    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
    req.SetBasicAuth(s.email+"/token", s.apiToken)
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

    articleID, ok := result["article"].(map[string]interface{})["id"].(float64)
    if !ok {
        return "", fmt.Errorf("failed to parse article ID")
    }

    return fmt.Sprintf("%.0f", articleID), nil
}

// UpdatePage updates an existing page in Zendesk
func (s *ZendeskServiceImpl) UpdatePage(ctx context.Context, page Page) error {
    url := fmt.Sprintf("%s/api/v2/help_center/articles/%s.json", s.baseURL, page.ID)

    // Fetch the current version to increment
    req, _ := http.NewRequest("GET", url, nil)
    req.SetBasicAuth(s.email+"/token", s.apiToken)
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

    currentVersion := result["article"].(map[string]interface{})["version"].(float64)
    newVersion := int(currentVersion) + 1

    reqBody, _ := json.Marshal(map[string]interface{}{
        "article": map[string]interface{}{
            "title":   page.Title,
            "body":    page.Content,
            "locale":  "en-us",
            "version": newVersion,
        },
    })

    req, _ = http.NewRequest("PUT", url, bytes.NewBuffer(reqBody))
    req.SetBasicAuth(s.email+"/token", s.apiToken)
    req.Header.Set("Content-Type", "application/json")

    resp, err = http.DefaultClient.Do(req)
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

// DeletePage deletes a page in Zendesk
func (s *ZendeskServiceImpl) DeletePage(ctx context.Context, id string) error {
    url := fmt.Sprintf("%s/api/v2/help_center/articles/%s.json", s.baseURL, id)

    req, _ := http.NewRequest("DELETE", url, nil)
    req.SetBasicAuth(s.email+"/token", s.apiToken)

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

// GetPage retrieves a page from Zendesk
func (s *ZendeskServiceImpl) GetPage(ctx context.Context, id string) (Page, error) {
    url := fmt.Sprintf("%s/api/v2/help_center/articles/%s.json", s.baseURL, id)

    req, _ := http.NewRequest("GET", url, nil)
    req.SetBasicAuth(s.email+"/token", s.apiToken)

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
    content := article["body"].(string)

    return Page{
        ID:      fmt.Sprintf("%.0f", pageID),
        Title:   title,
        Content: content,
    }, nil
}
