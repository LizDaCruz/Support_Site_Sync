package servicenow

import "context"

// ServiceNowService defines the methods for interacting with ServiceNow
type ServiceNowService interface {
    CreatePage(ctx context.Context, page Page) (string, error)
    UpdatePage(ctx context.Context, page Page) error
    DeletePage(ctx context.Context, id string) error
    GetPage(ctx context.Context, id string) (Page, error)
}
