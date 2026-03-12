package notification

// Channel represents the delivery channel.
type Channel string

const (
	ChannelEmail     Channel = "email"
	ChannelSMS       Channel = "sms"
	ChannelWebSocket Channel = "websocket"
)

// EmailPayload holds email notification data.
type EmailPayload struct {
	To       string   // recipient email
	Subject  string
	Body     string   // plain text
	HTML     string   // optional HTML body
	Cc       []string // optional
	Bcc      []string // optional
	UserID   string   // optional: for audit/correlation (tenant-less or tenant flow)
	TenantID string   // optional: for tenant-scoped flows
}

// SMSPayload holds SMS notification data.
type SMSPayload struct {
	To       string // phone number (E.164)
	Body     string
	UserID   string // optional: for audit/correlation
	TenantID string // optional: for tenant-scoped flows
}

// WebSocketPayload holds real-time in-app notification data.
// Delivery is typically via Pub/Sub (e.g. Redis) to connected WebSocket servers.
// Tenant-less flow: UserIDs only. Tenant flow: UserIDs + TenantID for context.
type WebSocketPayload struct {
	UserIDs  []string               // target user IDs (required for delivery)
	TenantID string                 // optional: tenant context when multi-tenant
	Event    string                 // event name, e.g. "notification.new"
	Data     map[string]interface{} // JSON-serializable payload
}
