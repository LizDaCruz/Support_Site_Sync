package guru

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "bytes"
    "io/ioutil"
)

// GuruServiceImpl is the implementation of the GuruService interface
type GuruServiceImpl struct {
    baseURL   string
    apiKey    string
}

// NewGuruService creates a new instance of GuruService
func NewGuruService(baseURL, apiKey string) *GuruServiceImpl {
    return &GuruServiceImpl{baseURL: baseURL, apiKey: apiKey}
}

// CreatePage creates a new card in Guru
func (s *GuruServiceImpl) CreatePage(ctx context.Context, page Page) (string, error) {
    url := fmt.Sprintf("%s/v1/cards", s.baseURL)
    reqBody, _ := json.Marshal(map[string]interface{}{
        "title":       page.Title,
        "content":     page.Content,
        "category_id": "YOUR_CATEGORY_ID", // Replace with your actual category ID
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
        return "", fmt.Errorf("failed to create card: %s - %s", resp.Status, body)
    }

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)

    cardID, ok := result["id"].(string)
    if !ok {
        return "", fmt.Errorf("failed to parse card ID")
    }

    return cardID, nil
}

// UpdatePage updates an existing card in Guru
func (s *GuruServiceImpl) UpdatePage(ctx context.Context, page Page) error {
    url := fmt.Sprintf("%s/v1/cards/%s", s.baseURL, page.ID)

    reqBody, _ := json.Marshal(map[string]interface{}{
        "title":   page.Title,
        "content": page.Content,
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
        return fmt.Errorf("failed to update card: %s - %s", resp.Status, body)
    }

    return nil
}

// DeletePage deletes a card from Guru
func (s *GuruServiceImpl) DeletePage(ctx context.Context, id string) error {
    url := fmt.Sprintf("%s/v1/cards/%s", s.baseURL, id)

    req, _ := http.NewRequest("DELETE", url, nil)
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusNoContent {
        body, _ := ioutil.ReadAll(resp.Body)
        return fmt.Errorf("failed to delete card: %s - %s", resp.Status, body)
    }

    return nil
}

// GetPage retrieves a card from Guru
func (s *GuruServiceImpl) GetPage(ctx context.Context, id string) (Page, error) {
    url := fmt.Sprintf("%s/v1/cards/%s", s.baseURL, id)

    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return Page{}, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := ioutil.ReadAll(resp.Body)
        return Page{}, fmt.Errorf("failed to get card: %s - %s", resp.Status, body)
    }

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)

    cardID, _ := result["id"].(string)
    title, _ := result["title"].(string)
    content, _ := result["content"].(string)

    return Page{
        ID:      cardID,
        Title:   title,
        Content: content,
    }, nil
}
