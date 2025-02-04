package entity

type Transaction struct {
    ID            int64     `json:"id"`
    UserID        int64     `json:"user_id"`
    Amount        float64   `json:"amount"`
    OperationType string    `json:"operation_type"`
    Description   string    `json:"description"`
    CreatedAt     string    `json:"created_at"`
    RelatedUserID int64     `json:"related_user_id"`
}