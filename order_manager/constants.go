package order_manager

const (
	FieldID            = "id"
	FieldStatus        = "status"
	FieldCreatedAt     = "created_at"
	FieldUpdatedAt     = "updated_at"
	FieldCustomerID    = "customer_id"
	FieldCustomerName  = "customer_name"
	FieldCustomerEmail = "customer_email"
	FieldOrders        = "orders"
	FieldTotal         = "total"
)

// Condition field constants for the multi-filter system.
const (
	CondFieldStatus        = "status"
	CondFieldOrderID       = "order_id"
	CondFieldCustomerName  = "customer_name"
	CondFieldCustomerEmail = "customer_email"
	CondFieldDateFrom      = "date_from"
	CondFieldDateTo        = "date_to"
)

// condition is a single filter condition sent from the frontend.
type condition struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}
