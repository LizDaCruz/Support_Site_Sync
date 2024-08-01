package main

import (
    "context"
    "log"
    "sync"
    "time"
)

// Page represents a document or page managed by the services
type Page struct {
    ID        string
    Title     string
    Content   string
    Timestamp time.Time // Added to keep track of the last updated time
}

// ServiceInterface defines the methods that all services must implement
type ServiceInterface interface {
    CreatePage(ctx context.Context, page Page) (string, error)
    UpdatePage(ctx context.Context, page Page) error
    DeletePage(ctx context.Context, id string) error
    GetPage(ctx context.Context, id string) (Page, error)
}

func syncPages(ctx context.Context, services []ServiceInterface, page Page) {
    var wg sync.WaitGroup
    errors := make(map[string]error)
    results := make(map[string]Page)

    // Create Page
    for _, svc := range services {
        wg.Add(1)
        go func(svc ServiceInterface) {
            defer wg.Done()
            _, err := svc.CreatePage(ctx, page)
            if err != nil {
                log.Printf("Error creating page in service: %v", err)
                errors[svcName(svc)] = err
            }
        }(svc)
    }
    wg.Wait()

    // Update Page
    for _, svc := range services {
        wg.Add(1)
        go func(svc ServiceInterface) {
            defer wg.Done()
            err := svc.UpdatePage(ctx, page)
            if err != nil {
                log.Printf("Error updating page in service: %v", err)
                errors[svcName(svc)] = err
            }
        }(svc)
    }
    wg.Wait()

    // Delete Page
    for _, svc := range services {
        wg.Add(1)
        go func(svc ServiceInterface) {
            defer wg.Done()
            err := svc.DeletePage(ctx, page.ID)
            if err != nil {
                log.Printf("Error deleting page in service: %v", err)
                errors[svcName(svc)] = err
            }
        }(svc)
    }
    wg.Wait()

    // Get Page and Compare Versions
    for _, svc := range services {
        wg.Add(1)
        go func(svc ServiceInterface) {
            defer wg.Done()
            retrievedPage, err := svc.GetPage(ctx, page.ID)
            if err != nil {
                log.Printf("Error getting page from service: %v", err)
                errors[svcName(svc)] = err
                return
            }
            results[svcName(svc)] = retrievedPage
        }(svc)
    }
    wg.Wait()

    // Compare Results to Identify the Most Current Version
    var latestPage Page
    for svc, p := range results {
        if p.Timestamp.After(latestPage.Timestamp) {
            latestPage = p
        }
        log.Printf("Service %s returned page: %+v", svc, p)
    }

    if latestPage.ID != "" {
        log.Printf("Most current version of the page is from service with ID: %s", latestPage.ID)
        log.Printf("Page Details: %+v", latestPage)
    } else {
        log.Println("No page data retrieved from services.")
    }

    if len(errors) > 0 {
        for svc, err := range errors {
            log.Printf("Service %s encountered an error: %v", svc, err)
        }
    }
}

func svcName(svc ServiceInterface) string {
    switch svc.(type) {
    case *confluence.ConfluenceService:
        return "Confluence"
    case *sharepoint.SharePointService:
        return "SharePoint"
    case *zendesk.ZendeskService:
        return "Zendesk"
    case *freshdesk.FreshdeskService:
        return "Freshdesk"
    case *servicenow.ServiceNowService:
        return "ServiceNow"
    case *helpjuice.HelpjuiceService:
        return "Helpjuice"
    case *notion.NotionService:
        return "Notion"
    case *docsify.DocsifyService:
        return "Docsify"
    case *guru.GuruService:
        return "Guru"
    case *trello.TrelloService:
        return "Trello"
    default:
        return "Unknown Service"
    }
}

func main() {
    // Initialize services
    confluenceService := NewConfluenceService("your-confluence-url", "your-confluence-api-key")
    sharepointService := NewSharePointService("your-sharepoint-url", "your-sharepoint-api-key")
    zendeskService := NewZendeskService("your-zendesk-url", "your-zendesk-api-key")
    freshdeskService := NewFreshdeskService("your-freshdesk-url", "your-freshdesk-api-key")
    servicenowService := NewServiceNowService("your-servicenow-url", "your-servicenow-api-key")
    helpjuiceService := NewHelpjuiceService("your-helpjuice-url", "your-helpjuice-api-key")
    notionService := NewNotionService("your-notion-url", "your-notion-api-key")
    docsifyService := NewDocsifyService("your-docsify-repo-owner", "your-docsify-repo-name", "your-docsify-api-key")
    guruService := NewGuruService("your-guru-base-url", "your-guru-api-key")
    trelloService := NewTrelloService("your-trello-api-key", "your-trello-api-token")

    services := []ServiceInterface{
        confluenceService,
        sharepointService,
        zendeskService,
        freshdeskService,
        servicenowService,
        helpjuiceService,
        notionService,
        docsifyService,
        guruService,
        trelloService,
    }

    // Create a page and sync
    page := Page{
        ID:        "example-page-id",
        Title:     "Example Page",
        Content:   "This is an example page content.",
        Timestamp: time.Now(), // Set the current time as the initial timestamp
    }

    ctx := context.Background()
    syncPages(ctx, services, page)
}
