package helpjuice

import "context"

// HelpjuiceService defines the methods for interacting with Helpjuice
type HelpjuiceService interface {
    CreatePage(ctx context.Context, page Page) (string, error)
    UpdatePage(ctx context.Context, page Page) error
    DeletePage(ctx context.Context, id string) error
    GetPage(ctx context.Context, id string) (Page, error)
}
