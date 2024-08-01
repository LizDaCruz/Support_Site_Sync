package freshdesk

import "context"

// FreshdeskService defines the methods for interacting with Freshdesk
type FreshdeskService interface {
    CreatePage(ctx context.Context, page Page) (string, error)
    UpdatePage(ctx context.Context, page Page) error
    DeletePage(ctx context.Context, id string) error
    GetPage(ctx context.Context, id string) (Page, error)
}
