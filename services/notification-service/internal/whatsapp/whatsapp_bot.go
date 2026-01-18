// Package whatsapp implements the WhatsApp Business API integration for OmniRoute
// Layer 4: Accessibility - WhatsApp Commerce Bot
package whatsapp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// =============================================================================
// DOMAIN MODELS
// =============================================================================

// Message represents an incoming WhatsApp message
type Message struct {
	From        string               `json:"from"`
	ID          string               `json:"id"`
	Timestamp   time.Time            `json:"timestamp"`
	Type        string               `json:"type"`
	Text        *TextMessage         `json:"text,omitempty"`
	Interactive *InteractiveResponse `json:"interactive,omitempty"`
	Location    *LocationMessage     `json:"location,omitempty"`
	Image       *MediaMessage        `json:"image,omitempty"`
}

type TextMessage struct {
	Body string `json:"body"`
}

type InteractiveResponse struct {
	Type        string       `json:"type"`
	ButtonReply *ButtonReply `json:"button_reply,omitempty"`
	ListReply   *ListReply   `json:"list_reply,omitempty"`
}

type ButtonReply struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type ListReply struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type LocationMessage struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name,omitempty"`
	Address   string  `json:"address,omitempty"`
}

type MediaMessage struct {
	ID       string `json:"id"`
	MimeType string `json:"mime_type"`
	SHA256   string `json:"sha256"`
}

// Session maintains conversation state
type Session struct {
	PhoneNumber string
	State       ConversationState
	Data        map[string]interface{}
	LastActive  time.Time
	CartItems   []CartItem
}

type CartItem struct {
	ProductID   string
	ProductName string
	Quantity    int
	UnitPrice   float64
}

type ConversationState string

const (
	StateWelcome        ConversationState = "welcome"
	StateMainMenu       ConversationState = "main_menu"
	StateOrderFlow      ConversationState = "order_flow"
	StateCategorySelect ConversationState = "category_select"
	StateProductSelect  ConversationState = "product_select"
	StateQuantityInput  ConversationState = "quantity_input"
	StateCartView       ConversationState = "cart_view"
	StateCheckout       ConversationState = "checkout"
	StatePayment        ConversationState = "payment"
	StateOrderStatus    ConversationState = "order_status"
	StateSupport        ConversationState = "support"
)

// =============================================================================
// WHATSAPP BOT SERVICE
// =============================================================================

type WhatsAppBotService struct {
	logger        *zap.Logger
	sessions      map[string]*Session
	apiToken      string
	phoneNumberID string
	webhookToken  string
}

func NewWhatsAppBotService(logger *zap.Logger, apiToken, phoneNumberID, webhookToken string) *WhatsAppBotService {
	return &WhatsAppBotService{
		logger:        logger,
		sessions:      make(map[string]*Session),
		apiToken:      apiToken,
		phoneNumberID: phoneNumberID,
		webhookToken:  webhookToken,
	}
}

// =============================================================================
// MESSAGE HANDLERS
// =============================================================================

func (s *WhatsAppBotService) HandleIncomingMessage(ctx context.Context, msg Message) error {
	session := s.getOrCreateSession(msg.From)

	// Determine message content
	var userInput string
	switch {
	case msg.Text != nil:
		userInput = strings.ToLower(strings.TrimSpace(msg.Text.Body))
	case msg.Interactive != nil:
		if msg.Interactive.ButtonReply != nil {
			userInput = msg.Interactive.ButtonReply.ID
		} else if msg.Interactive.ListReply != nil {
			userInput = msg.Interactive.ListReply.ID
		}
	}

	// Route to appropriate handler
	return s.routeMessage(ctx, session, userInput)
}

func (s *WhatsAppBotService) routeMessage(ctx context.Context, session *Session, input string) error {
	// Global commands
	switch input {
	case "hi", "hello", "start", "menu":
		return s.sendMainMenu(ctx, session.PhoneNumber)
	case "help", "support":
		return s.sendSupportMenu(ctx, session.PhoneNumber)
	case "cart":
		return s.sendCartView(ctx, session)
	case "0", "back":
		return s.goBack(ctx, session)
	}

	// State-based routing
	switch session.State {
	case StateWelcome, StateMainMenu:
		return s.handleMainMenuSelection(ctx, session, input)
	case StateCategorySelect:
		return s.handleCategorySelection(ctx, session, input)
	case StateProductSelect:
		return s.handleProductSelection(ctx, session, input)
	case StateQuantityInput:
		return s.handleQuantityInput(ctx, session, input)
	case StateCheckout:
		return s.handleCheckout(ctx, session, input)
	case StateOrderStatus:
		return s.handleOrderStatusQuery(ctx, session, input)
	default:
		return s.sendMainMenu(ctx, session.PhoneNumber)
	}
}

// =============================================================================
// MENU HANDLERS
// =============================================================================

func (s *WhatsAppBotService) sendMainMenu(ctx context.Context, to string) error {
	menu := InteractiveListMessage{
		To:   to,
		Type: "interactive",
		Interactive: Interactive{
			Type:   "list",
			Header: &Header{Type: "text", Text: "🛒 OmniRoute Commerce"},
			Body:   Body{Text: "Welcome! How can we help you today?"},
			Footer: &Footer{Text: "Reply with number or tap option"},
			Action: Action{
				Button: "View Options",
				Sections: []Section{
					{
						Title: "Shopping",
						Rows: []Row{
							{ID: "order", Title: "📦 Place Order", Description: "Browse and order products"},
							{ID: "reorder", Title: "🔄 Quick Reorder", Description: "Repeat your last order"},
							{ID: "cart", Title: "🛒 View Cart", Description: "See items in your cart"},
						},
					},
					{
						Title: "Account",
						Rows: []Row{
							{ID: "status", Title: "📍 Track Order", Description: "Check order status"},
							{ID: "balance", Title: "💰 Balance", Description: "View wallet & credit"},
							{ID: "support", Title: "💬 Support", Description: "Get help"},
						},
					},
				},
			},
		},
	}

	return s.sendInteractiveMessage(ctx, menu)
}

func (s *WhatsAppBotService) handleMainMenuSelection(ctx context.Context, session *Session, input string) error {
	switch input {
	case "order", "1":
		session.State = StateCategorySelect
		return s.sendCategories(ctx, session.PhoneNumber)
	case "reorder", "2":
		return s.sendQuickReorder(ctx, session)
	case "cart", "3":
		return s.sendCartView(ctx, session)
	case "status", "4":
		session.State = StateOrderStatus
		return s.sendOrderStatusPrompt(ctx, session.PhoneNumber)
	case "balance", "5":
		return s.sendBalance(ctx, session.PhoneNumber)
	case "support", "6":
		return s.sendSupportMenu(ctx, session.PhoneNumber)
	default:
		return s.sendTextMessage(ctx, session.PhoneNumber, "❌ Invalid option. Please try again.")
	}
}

func (s *WhatsAppBotService) sendCategories(ctx context.Context, to string) error {
	menu := InteractiveListMessage{
		To:   to,
		Type: "interactive",
		Interactive: Interactive{
			Type:   "list",
			Header: &Header{Type: "text", Text: "📂 Categories"},
			Body:   Body{Text: "Select a product category:"},
			Action: Action{
				Button: "Browse Categories",
				Sections: []Section{
					{
						Title: "Product Categories",
						Rows: []Row{
							{ID: "cat_beverages", Title: "🥤 Beverages", Description: "Drinks, water, juices"},
							{ID: "cat_food", Title: "🍚 Food Items", Description: "Rice, noodles, canned goods"},
							{ID: "cat_personal", Title: "🧴 Personal Care", Description: "Soap, toiletries"},
							{ID: "cat_household", Title: "🏠 Household", Description: "Cleaning, kitchen items"},
							{ID: "cat_electronics", Title: "📱 Electronics", Description: "Phones, accessories"},
						},
					},
				},
			},
		},
	}

	return s.sendInteractiveMessage(ctx, menu)
}

func (s *WhatsAppBotService) handleCategorySelection(ctx context.Context, session *Session, input string) error {
	// Store selected category
	session.Data["category"] = input
	session.State = StateProductSelect

	// Send products for category (mock data)
	products := s.getProductsByCategory(input)
	return s.sendProductList(ctx, session.PhoneNumber, products)
}

func (s *WhatsAppBotService) sendProductList(ctx context.Context, to string, products []Product) error {
	var rows []Row
	for _, p := range products {
		rows = append(rows, Row{
			ID:          p.ID,
			Title:       p.Name,
			Description: fmt.Sprintf("₦%.0f", p.Price),
		})
	}

	menu := InteractiveListMessage{
		To:   to,
		Type: "interactive",
		Interactive: Interactive{
			Type:   "list",
			Header: &Header{Type: "text", Text: "🛍️ Products"},
			Body:   Body{Text: "Select a product to add to cart:"},
			Action: Action{
				Button:   "View Products",
				Sections: []Section{{Title: "Available Products", Rows: rows}},
			},
		},
	}

	return s.sendInteractiveMessage(ctx, menu)
}

func (s *WhatsAppBotService) handleProductSelection(ctx context.Context, session *Session, input string) error {
	session.Data["selected_product"] = input
	session.State = StateQuantityInput

	return s.sendTextMessage(ctx, session.PhoneNumber,
		"📝 Enter quantity (number only):\n\nExample: 5")
}

func (s *WhatsAppBotService) handleQuantityInput(ctx context.Context, session *Session, input string) error {
	// Parse quantity
	var qty int
	_, err := fmt.Sscanf(input, "%d", &qty)
	if err != nil || qty <= 0 {
		return s.sendTextMessage(ctx, session.PhoneNumber, "❌ Please enter a valid number (e.g., 5)")
	}

	// Add to cart
	productID := session.Data["selected_product"].(string)
	product := s.getProductByID(productID)

	session.CartItems = append(session.CartItems, CartItem{
		ProductID:   productID,
		ProductName: product.Name,
		Quantity:    qty,
		UnitPrice:   product.Price,
	})

	session.State = StateMainMenu

	// Send confirmation with buttons
	return s.sendButtonMessage(ctx, session.PhoneNumber,
		fmt.Sprintf("✅ Added to cart:\n%d x %s\n\nTotal: ₦%.0f", qty, product.Name, float64(qty)*product.Price),
		[]Button{
			{Type: "reply", Reply: Reply{ID: "cart", Title: "🛒 View Cart"}},
			{Type: "reply", Reply: Reply{ID: "order", Title: "➕ Add More"}},
			{Type: "reply", Reply: Reply{ID: "checkout", Title: "💳 Checkout"}},
		},
	)
}

func (s *WhatsAppBotService) sendCartView(ctx context.Context, session *Session) error {
	if len(session.CartItems) == 0 {
		return s.sendTextMessage(ctx, session.PhoneNumber, "🛒 Your cart is empty.\n\nReply with *order* to start shopping!")
	}

	var total float64
	var cartText strings.Builder
	cartText.WriteString("🛒 *Your Cart*\n\n")

	for i, item := range session.CartItems {
		itemTotal := float64(item.Quantity) * item.UnitPrice
		total += itemTotal
		cartText.WriteString(fmt.Sprintf("%d. %s\n   %d x ₦%.0f = ₦%.0f\n\n",
			i+1, item.ProductName, item.Quantity, item.UnitPrice, itemTotal))
	}

	cartText.WriteString(fmt.Sprintf("─────────────\n*Total: ₦%.0f*", total))

	return s.sendButtonMessage(ctx, session.PhoneNumber, cartText.String(),
		[]Button{
			{Type: "reply", Reply: Reply{ID: "checkout", Title: "💳 Checkout"}},
			{Type: "reply", Reply: Reply{ID: "order", Title: "➕ Add More"}},
			{Type: "reply", Reply: Reply{ID: "clear_cart", Title: "🗑️ Clear Cart"}},
		},
	)
}

func (s *WhatsAppBotService) sendBalance(ctx context.Context, to string) error {
	// Mock balance data
	balance := `💰 *Account Balance*

Wallet: ₦45,230.00
Available Credit: ₦150,000.00
Credit Used: ₦35,000.00

Next Payment Due: Jan 25, 2026
Amount Due: ₦35,000.00`

	return s.sendButtonMessage(ctx, to, balance,
		[]Button{
			{Type: "reply", Reply: Reply{ID: "topup", Title: "➕ Top Up"}},
			{Type: "reply", Reply: Reply{ID: "pay_credit", Title: "💳 Pay Credit"}},
			{Type: "reply", Reply: Reply{ID: "menu", Title: "🏠 Main Menu"}},
		},
	)
}

// =============================================================================
// API MESSAGE SENDING
// =============================================================================

type InteractiveListMessage struct {
	To          string      `json:"to"`
	Type        string      `json:"type"`
	Interactive Interactive `json:"interactive"`
}

type Interactive struct {
	Type   string  `json:"type"`
	Header *Header `json:"header,omitempty"`
	Body   Body    `json:"body"`
	Footer *Footer `json:"footer,omitempty"`
	Action Action  `json:"action"`
}

type Header struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Body struct {
	Text string `json:"text"`
}

type Footer struct {
	Text string `json:"text"`
}

type Action struct {
	Button   string    `json:"button,omitempty"`
	Buttons  []Button  `json:"buttons,omitempty"`
	Sections []Section `json:"sections,omitempty"`
}

type Section struct {
	Title string `json:"title"`
	Rows  []Row  `json:"rows"`
}

type Row struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

type Button struct {
	Type  string `json:"type"`
	Reply Reply  `json:"reply"`
}

type Reply struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type Product struct {
	ID    string
	Name  string
	Price float64
}

func (s *WhatsAppBotService) sendTextMessage(ctx context.Context, to, text string) error {
	s.logger.Info("Sending text message", zap.String("to", to))
	// Implementation would call WhatsApp Cloud API
	return nil
}

func (s *WhatsAppBotService) sendButtonMessage(ctx context.Context, to, text string, buttons []Button) error {
	s.logger.Info("Sending button message", zap.String("to", to))
	// Implementation would call WhatsApp Cloud API
	return nil
}

func (s *WhatsAppBotService) sendInteractiveMessage(ctx context.Context, msg InteractiveListMessage) error {
	s.logger.Info("Sending interactive message", zap.String("to", msg.To))
	// Implementation would call WhatsApp Cloud API
	return nil
}

func (s *WhatsAppBotService) getOrCreateSession(phone string) *Session {
	if session, ok := s.sessions[phone]; ok {
		session.LastActive = time.Now()
		return session
	}

	session := &Session{
		PhoneNumber: phone,
		State:       StateWelcome,
		Data:        make(map[string]interface{}),
		LastActive:  time.Now(),
		CartItems:   []CartItem{},
	}
	s.sessions[phone] = session
	return session
}

func (s *WhatsAppBotService) getProductsByCategory(category string) []Product {
	// Mock product data
	return []Product{
		{ID: "prod-001", Name: "Peak Milk 400g", Price: 2500},
		{ID: "prod-002", Name: "Indomie Chicken (Carton)", Price: 5500},
		{ID: "prod-003", Name: "Golden Penny Rice 50kg", Price: 65000},
	}
}

func (s *WhatsAppBotService) getProductByID(id string) Product {
	products := s.getProductsByCategory("")
	for _, p := range products {
		if p.ID == id {
			return p
		}
	}
	return Product{}
}

func (s *WhatsAppBotService) goBack(ctx context.Context, session *Session) error {
	session.State = StateMainMenu
	return s.sendMainMenu(ctx, session.PhoneNumber)
}

func (s *WhatsAppBotService) sendQuickReorder(ctx context.Context, session *Session) error {
	return s.sendTextMessage(ctx, session.PhoneNumber, "🔄 Quick Reorder feature coming soon!")
}

func (s *WhatsAppBotService) sendOrderStatusPrompt(ctx context.Context, to string) error {
	return s.sendTextMessage(ctx, to, "📍 Enter your Order ID to track:\n\nExample: ORD-2026-00123")
}

func (s *WhatsAppBotService) handleOrderStatusQuery(ctx context.Context, session *Session, input string) error {
	session.State = StateMainMenu
	return s.sendTextMessage(ctx, session.PhoneNumber, fmt.Sprintf("📦 Order %s\n\nStatus: In Transit\nETA: Today, 3:00 PM\nDriver: John (0803-XXX-1234)", input))
}

func (s *WhatsAppBotService) sendSupportMenu(ctx context.Context, to string) error {
	return s.sendTextMessage(ctx, to, "💬 *Support*\n\n📞 Call: 0800-OMNIROUTE\n📧 Email: support@omniroute.io\n💬 WhatsApp: +234 800 123 4567")
}

func (s *WhatsAppBotService) handleCheckout(ctx context.Context, session *Session, input string) error {
	return s.sendTextMessage(ctx, session.PhoneNumber, "💳 Checkout feature coming soon!")
}

// =============================================================================
// HTTP HANDLERS
// =============================================================================

func (s *WhatsAppBotService) SetupRoutes(r *gin.Engine) {
	webhook := r.Group("/webhook/whatsapp")
	{
		webhook.GET("", s.VerifyWebhook)
		webhook.POST("", s.HandleWebhook)
	}
}

func (s *WhatsAppBotService) VerifyWebhook(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if mode == "subscribe" && token == s.webhookToken {
		c.String(http.StatusOK, challenge)
		return
	}
	c.String(http.StatusForbidden, "Forbidden")
}

func (s *WhatsAppBotService) HandleWebhook(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Process webhook payload asynchronously
	go s.processWebhookPayload(payload)

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

func (s *WhatsAppBotService) processWebhookPayload(payload map[string]interface{}) {
	// Parse and handle webhook payload
	data, _ := json.Marshal(payload)
	s.logger.Info("Received webhook", zap.String("payload", string(data)))
}
