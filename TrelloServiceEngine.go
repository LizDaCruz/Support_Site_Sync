package trello

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "bytes"
    "io/ioutil"
)

// TrelloServiceImpl is the implementation of the TrelloService interface
type TrelloServiceImpl struct {
    apiKey    string
    apiToken  string
    baseURL   string
}

// NewTrelloService creates a new instance of TrelloService
func NewTrelloService(apiKey, apiToken string) *TrelloServiceImpl {
    return &TrelloServiceImpl{
        apiKey:   apiKey,
        apiToken: apiToken,
        baseURL:  "https://api.trello.com/1",
    }
}

// CreatePage creates a new card in Trello
func (s *TrelloServiceImpl) CreatePage(ctx context.Context, page Page) (string, error) {
    url := fmt.Sprintf("%s/cards?key=%s&token=%s", s.baseURL, s.apiKey, s.apiToken)
    reqBody, _ := json.Marshal(map[string]interface{}{
        "name":        page.Title,
        "desc":        page.Content,
        "idList":      "YOUR_LIST_ID", // Replace with your actual list ID
        "keepFromSource": "all",
    })

    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
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

// UpdatePage updates an existing card in Trello
func (s *TrelloServiceImpl) UpdatePage(ctx context.Context, page Page) error {
    url := fmt.Sprintf("%s/cards/%s?key=%s&token=%s", s.baseURL, page.ID, s.apiKey, s.apiToken)

    reqBody, _ := json.Marshal(map[string]interface{}{
        "name":        page.Title,
        "desc":        page.Content,
    })

    req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(reqBody))
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

// DeletePage deletes a card from Trello
func (s *TrelloServiceImpl) DeletePage(ctx context.Context, id string) error {
    url := fmt.Sprintf("%s/cards/%s?key=%s&token=%s", s.baseURL, id, s.apiKey, s.apiToken)

    req, _ := http.NewRequest("DELETE", url, nil)
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

// GetPage retrieves a card from Trello
func (s *TrelloServiceImpl) GetPage(ctx context.Context, id string) (Page, error) {
    url := fmt.Sprintf("%s/cards/%s?key=%s&token=%s", s.baseURL, id, s.apiKey, s.apiToken)

    req, _ := http.NewRequest("GET", url, nil)
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

    return Page{
        ID:      result["id"].(string),
        Title:   result["name"].(string),
        Content: result["desc"].(string),
    }, nil
}
