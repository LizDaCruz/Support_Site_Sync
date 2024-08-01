package servicenow

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "bytes"
    "io/ioutil"
)

// ServiceNowServiceImpl is the implementation of the ServiceNowService interface
type ServiceNowServiceImpl struct {
    baseURL   string
    username   string
    password   string
}

// NewServiceNowService creates a new instance of ServiceNowService
func NewServiceNowService(baseURL, username, password string) *ServiceNowServiceImpl {
    return &ServiceNowServiceImpl{baseURL: baseURL, username: username, password: password}
}

// CreatePage creates a new page in ServiceNow
func (s *ServiceNowServiceImpl) CreatePage(ctx context.Context, page Page) (string, error) {
    url := fmt.Sprintf("%s/api/now/table/kb_knowledge", s.baseURL)
    reqBody, _ := json.Marshal(map[string]interface{}{
        "short_description": page.Title,
        "text":              page.Content,
        "workflow_state":    "published",
    })

    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
    req.SetBasicAuth(s.username, s.password)
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

    sysID, ok := result["result"].(map[string]interface{})["sys_id"].(string)
    if !ok {
        return "", fmt.Errorf("failed to parse sys_id")
    }

    return sysID, nil
}

// UpdatePage updates an existing page in ServiceNow
func (s *ServiceNowServiceImpl) UpdatePage(ctx context.Context, page Page) error {
    url := fmt.Sprintf("%s/api/now/table/kb_knowledge/%s", s.baseURL, page.ID)

    reqBody, _ := json.Marshal(map[string]interface{}{
        "short_description": page.Title,
        "text":              page.Content,
        "workflow_state":    "published",
    })

    req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(reqBody))
    req.SetBasicAuth(s.username, s.password)
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

// DeletePage deletes a page in ServiceNow
func (s *ServiceNowServiceImpl) DeletePage(ctx context.Context, id string) error {
    url := fmt.Sprintf("%s/api/now/table/kb_knowledge/%s", s.baseURL, id)

    req, _ := http.NewRequest("DELETE", url, nil)
    req.SetBasicAuth(s.username, s.password)

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

// GetPage retrieves a page from ServiceNow
func (s *ServiceNowServiceImpl) GetPage(ctx context.Context, id string) (Page, error) {
    url := fmt.Sprintf("%s/api/now/table/kb_knowledge/%s", s.baseURL, id)

    req, _ := http.NewRequest("GET", url, nil)
    req.SetBasicAuth(s.username, s.password)

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

    kbArticle := result["result"].(map[string]interface{})
    pageID := kbArticle["sys_id"].(string)
    title := kbArticle["short_description"].(string)
    content := kbArticle["text"].(string)

    return Page{
        ID:      pageID,
        Title:   title,
        Content: content,
    }, nil
}
