package domain

import notifv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/notification"

type Notification struct {
	UserID      int64  `json:"user_id"`
	Message     string `json:"message"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Severity    string `json:"severity"`
	Source      string `json:"source"`
	DisplayType string `json:"display_type"`
	CreatedAt   string `json:"created_at"`
	ExpiryAt    string `json:"expiry_at"`
}

func NotificationProtoToDomain(n *notifv1.NotificationResponse) Notification {
	return Notification{
		UserID:      n.GetUserId(),
		Message:     n.GetMessage(),
		Description: n.GetDescription(),
		Status:      n.GetStatus(),
		Severity:    n.GetSeverity(),
		Source:      n.GetSource(),
		DisplayType: n.GetDisplayType(),
		CreatedAt:   n.GetCreatedAt(),
		ExpiryAt:    n.GetExpiryAt(),
	}
}
