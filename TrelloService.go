package trello

import "context"

// TrelloService defines the methods for interacting with Trello
type TrelloService interface {
    CreatePage(ctx context.Context, page Page) (string, error)
    UpdatePage(ctx context.Context, page Page) error
    DeletePage(ctx context.Context, id string) error
    GetPage(ctx context.Context, id string) (Page, error)
}
