package docsify

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "bytes"
    "io/ioutil"
)

// DocsifyServiceImpl is the implementation of the DocsifyService interface
type DocsifyServiceImpl struct {
    repoOwner    string
    repoName     string
    apiKey        string
}

// NewDocsifyService creates a new instance of DocsifyService
func NewDocsifyService(repoOwner, repoName, apiKey string) *DocsifyServiceImpl {
    return &DocsifyServiceImpl{repoOwner: repoOwner, repoName: repoName, apiKey: apiKey}
}

// CreatePage creates a new page in the Docsify repository
func (s *DocsifyServiceImpl) CreatePage(ctx context.Context, page Page) (string, error) {
    // GitHub API URL to create a new file
    url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", s.repoOwner, s.repoName, page.ID)

    reqBody, _ := json.Marshal(map[string]interface{}{
        "message": "Create page " + page.Title,
        "content": encodeBase64(page.Content),
    })

    req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(reqBody))
    req.Header.Set("Authorization", fmt.Sprintf("token %s", s.apiKey))
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

    return page.ID, nil
}

// UpdatePage updates an existing page in the Docsify repository
func (s *DocsifyServiceImpl) UpdatePage(ctx context.Context, page Page) error {
    // GitHub API URL to update a file
    url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", s.repoOwner, s.repoName, page.ID)

    // Fetch the current file to get the SHA
    sha, err := s.getFileSHA(ctx, page.ID)
    if err != nil {
        return err
    }

    reqBody, _ := json.Marshal(map[string]interface{}{
        "message": "Update page " + page.Title,
        "content": encodeBase64(page.Content),
        "sha":     sha,
    })

    req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(reqBody))
    req.Header.Set("Authorization", fmt.Sprintf("token %s", s.apiKey))
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

// DeletePage deletes a page from the Docsify repository
func (s *DocsifyServiceImpl) DeletePage(ctx context.Context, id string) error {
    // GitHub API URL to delete a file
    url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", s.repoOwner, s.repoName, id)

    // Fetch the current file to get the SHA
    sha, err := s.getFileSHA(ctx, id)
    if err != nil {
        return err
    }

    reqBody, _ := json.Marshal(map[string]interface{}{
        "message": "Delete page",
        "sha":     sha,
    })

    req, _ := http.NewRequest("DELETE", url, bytes.NewBuffer(reqBody))
    req.Header.Set("Authorization", fmt.Sprintf("token %s", s.apiKey))
    req.Header.Set("Content-Type", "application/json")

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

// GetPage retrieves a page from the Docsify repository
func (s *DocsifyServiceImpl) GetPage(ctx context.Context, id string) (Page, error) {
    url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", s.repoOwner, s.repoName, id)

    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("Authorization", fmt.Sprintf("token %s", s.apiKey))

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

    content, ok := result["content"].(string)
    if !ok {
        return Page{}, fmt.Errorf("failed to parse content")
    }

    return Page{
        ID:      id,
        Title:   result["name"].(string),
        Content: decodeBase64(content),
    }, nil
}

// getFileSHA fetches the SHA of the file
func (s *DocsifyServiceImpl) getFileSHA(ctx context.Context, path string) (string, error) {
    url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", s.repoOwner, s.repoName, path)

    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("Authorization", fmt.Sprintf("token %s", s.apiKey))

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := ioutil.ReadAll(resp.Body)
        return "", fmt.Errorf("failed to get file SHA: %s - %s", resp.Status, body)
    }

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)

    sha, ok := result["sha"].(string)
    if !ok {
        return "", fmt.Errorf("failed to parse file SHA")
    }

    return sha, nil
}

// encodeBase64 encodes a string to base64
func encodeBase64(content string) string {
    return base64.StdEncoding.EncodeToString([]byte(content))
}

// decodeBase64 decodes a base64 string
func decodeBase64(content string) string {
    decoded, _ := base64.StdEncoding.DecodeString(content)
    return string(decoded)
}
