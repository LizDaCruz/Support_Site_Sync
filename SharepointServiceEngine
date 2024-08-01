package sharepoint

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "bytes"
    "io/ioutil"
)

type SharePointService struct {
    baseURL   string
    accessToken string
}

func NewSharePointService(baseURL, accessToken string) *SharePointService {
    return &SharePointService{baseURL: baseURL, accessToken: accessToken}
}

// CreatePage creates a new page in SharePoint
func (s *SharePointService) CreatePage(ctx context.Context, page Page) (string, error) {
    url := fmt.Sprintf("%s/_api/web/lists/getbytitle('Site Pages')/items", s.baseURL)
    reqBody, _ := json.Marshal(map[string]interface{}{
        "__metadata": map[string]string{"type": "SP.Data.SitePagesItem"},
        "Title":      page.Title,
        "Content":    page.Content,
    })

    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
    req.Header.Set("Authorization", "Bearer " + s.accessToken)
    req.Header.Set("Content-Type", "application/json;odata=verbose")

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

    id, ok := result["Id"].(float64)
    if !ok {
        return "", fmt.Errorf("failed to parse page ID")
    }

    return fmt.Sprintf("%.0f", id), nil
}

// UpdatePage updates an existing page in SharePoint
func (s *SharePointService) UpdatePage(ctx context.Context, page Page) error {
    url := fmt.Sprintf("%s/_api/web/lists/getbytitle('Site Pages')/items(%s)", s.baseURL, page.ID)

    reqBody, _ := json.Marshal(map[string]interface{}{
        "__metadata": map[string]string{"type": "SP.Data.SitePagesItem"},
        "Title":      page.Title,
        "Content":    page.Content,
    })

    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
    req.Header.Set("Authorization", "Bearer " + s.accessToken)
    req.Header.Set("Content-Type", "application/json;odata=verbose")
    req.Header.Set("X-HTTP-Method", "MERGE")
    req.Header.Set("If-Match", "*")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusNoContent {
        body, _ := ioutil.ReadAll(resp.Body)
        return fmt.Errorf("failed to update page: %s - %s", resp.Status, body)
    }

    return nil
}

// DeletePage deletes a page in SharePoint
func (s *SharePointService) DeletePage(ctx context.Context, id string) error {
    url := fmt.Sprintf("%s/_api/web/lists/getbytitle('Site Pages')/items(%s)", s.baseURL, id)

    req, _ := http.NewRequest("DELETE", url, nil)
    req.Header.Set("Authorization", "Bearer " + s.accessToken)
    req.Header.Set("If-Match", "*")

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

// GetPage retrieves a page from SharePoint
func (s *SharePointService) GetPage(ctx context.Context, id string) (Page, error) {
    url := fmt.Sprintf("%s/_api/web/lists/getbytitle('Site Pages')/items(%s)", s.baseURL, id)

    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("Authorization", "Bearer " + s.accessToken)
    req.Header.Set("Accept", "application/json;odata=verbose")

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

    pageID := fmt.Sprintf("%.0f", result["Id"].(float64))
    title := result["Title"].(string)
    content := result["Content"].(string)

    return Page{
        ID:      pageID,
        Title:   title,
        Content: content,
    }, nil
}
