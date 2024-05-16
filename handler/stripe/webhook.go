package stripe

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/EmissarySocial/emissary/domain"
	"github.com/EmissarySocial/emissary/model"
	"github.com/EmissarySocial/emissary/server"
	"github.com/EmissarySocial/emissary/service"
	"github.com/benpate/derp"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/client"
	"github.com/stripe/stripe-go/v78/webhook"
)

func PostWebhook(serverFactory *server.Factory) echo.HandlerFunc {

	const location = "handler.stripe.PostWebhook"

	return func(ctx echo.Context) error {

		//////////////////////////////////////
		// 1. PREPARE AND VALIDATE THE REQUEST

		// Get the factory for this Domain
		factory, err := serverFactory.ByContext(ctx)

		if err != nil {
			return derp.ReportAndReturn(derp.Wrap(err, location, "Error getting domain factory"))
		}

		// Load the Domain record
		domainService := factory.Domain()
		domain, err := domainService.LoadDomain()

		if err != nil {
			return derp.ReportAndReturn(derp.Wrap(err, location, "Error loading domain record"))
		}

		// RULE: Require that a registration form has been defined
		if !domain.HasRegistrationForm() {
			return derp.ReportAndReturn(derp.NewNotFoundError(location, "Stripe Webhook not defined (no registration form)"))
		}

		// Collect Registration Settings
		secret := domain.RegistrationData.GetString("stripe_webhook_secret")
		if secret == "" {
			return derp.ReportAndReturn(derp.NewInternalError(location, "Stripe Webhook Secret not defined"))
		}

		restrictedKey := domain.RegistrationData.GetString("stripe_restricted_key")
		if restrictedKey == "" {
			return derp.ReportAndReturn(derp.NewInternalError(location, "Stripe Restricted Key not defined"))
		}

		////////////////////////////////
		// 2. READ DATA FROM THE WEBHOOK

		// Read the request body
		payload, err := io.ReadAll(ctx.Request().Body)

		if err != nil {
			return derp.ReportAndReturn(derp.Wrap(err, location, "Error reading request body"))
		}

		// Verify the WebHook signature
		signatureHeader := ctx.Request().Header.Get("Stripe-Signature")
		event, err := webhook.ConstructEvent(payload, signatureHeader, secret)

		if err != nil {
			return derp.ReportAndReturn(derp.Wrap(err, location, "Error verifying webhook signature"))
		}

		// Require that the event is a "subscription" event
		eventType := string(event.Type)

		if !strings.HasPrefix(eventType, "customer.subscription.") {
			log.Trace().Str("event", eventType).Msg("Ignoring Stripe Webhook")
			return nil
		}

		log.Trace().Str("event", eventType).Msg("Processing Stripe Webhook")

		////////////////////////////////
		// 3. DRAW THE REST OF THE OWL HERE
		// Moved to an async function so that our Webhook will respond to the server quickly.
		// Whatever else happens, it's on us from here on out.
		derp.Report(finishWebhook(factory, restrictedKey, event))

		// Success?
		return ctx.NoContent(http.StatusOK)
	}
}

func finishWebhook(factory *domain.Factory, restrictedKey string, event stripe.Event) error {

	const location = "handler.stripe.finishWebhook"

	// Get the subscription from the event details
	subscription := stripe.Subscription{}

	if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
		return derp.Wrap(err, "handler.getSubscription", "Error unmarshalling event data")
	}

	// This is the price that was paid, but it doesn't include the metadata we need.
	// So, use the API to look up the productID first.

	price := getSubscriptionPrice(&subscription)

	if price == nil {
		return derp.NewBadRequestError(location, "No price found in subscription", subscription)
	}

	if err := loadStripeProduct(restrictedKey, price.Product); err != nil {
		return derp.ReportAndReturn(derp.Wrap(err, location, "Error getting product details"))
	}

	// Get ready to create/update a user
	userService := factory.User()
	user := model.NewUser()

	switch subscription.Status {

	// If the subscription is ACTIVE, then add the user and their group memberships
	case stripe.SubscriptionStatusActive,
		stripe.SubscriptionStatusTrialing:

		// Try to load/create the user
		if err := loadOrCreateUser(restrictedKey, userService, subscription.Customer, &user); err != nil {
			return derp.Wrap(err, location, "Error creating customer", subscription.Customer)
		}

		// Add the user to the designated groups
		addGroups(factory, &user, price.Product, "add_groups")

		// Remove the user from the designated groups
		removeGroups(factory, &user, price.Product, "remove_groups")

	// Otherwise, CANCEL the user's subscription
	default:

		// If the user doesn't exists, then we don't have to cancel their access here.
		if err := loadUser(userService, subscription.Customer, &user); err != nil {
			return nil
		}

		// Since this subscription is no longer active, remove the user from the designated groups
		removeGroups(factory, &user, price.Product, "add_groups")
	}

	// Save the new User to the database.  Yay!
	if err := userService.Save(&user, "Created by Stripe Webhook"); err != nil {
		return derp.Wrap(err, location, "Error saving user record")
	}

	// Success!
	return nil
}

// addGroups adds groups to the User's list, as specified by the Product metadata
func addGroups(factory *domain.Factory, user *model.User, product *stripe.Product, token string) {

	if user == nil {
		return
	}

	if product == nil {
		return
	}

	groupService := factory.Group()
	groupIDs := strings.Split(product.Metadata[token], ",")

	for _, groupToken := range groupIDs {
		group := model.NewGroup()
		if err := groupService.LoadByToken(groupToken, &group); err == nil {
			user.AddGroup(group.GroupID)
		}
	}
}

// removeGroups removes groups from the User's list, as specified by the Product metadata
func removeGroups(factory *domain.Factory, user *model.User, product *stripe.Product, token string) {

	if user == nil {
		return
	}

	if product == nil {
		return
	}

	groupService := factory.Group()
	groupIDs := strings.Split(product.Metadata[token], ",")

	for _, groupToken := range groupIDs {
		group := model.NewGroup()
		if err := groupService.LoadByToken(groupToken, &group); err == nil {
			user.RemoveGroup(group.GroupID)
		}
	}
}

func getSubscriptionPrice(subscription *stripe.Subscription) *stripe.Price {

	if items := subscription.Items; items != nil {
		for _, item := range items.Data {
			if item.Price != nil {
				return item.Price
			}
		}
	}

	return nil
}

func loadUser(userService *service.User, customer *stripe.Customer, user *model.User) error {

	if customer == nil {
		return derp.NewBadRequestError("handler.stripe.loadUser", "Customer must not be nil")
	}

	// Try to load the user by their email address
	if err := userService.LoadByMapID(model.UserMapIDStripe, customer.ID, user); err != nil {
		return derp.Wrap(err, "handler.stripe.loadUser", "Error loading user record")
	}

	return nil
}

func loadOrCreateUser(apiKey string, userService *service.User, customer *stripe.Customer, user *model.User) error {

	err := loadUser(userService, customer, user)

	if err == nil {
		return nil
	}

	if derp.NotFound(err) {

		if err := loadStripeCustomer(apiKey, customer); err != nil {
			return derp.Wrap(err, "handler.stripe.loadOrCreateUser", "Error loading customer from Stripe API")
		}

		if customer.Name != "" {
			user.DisplayName = customer.Name
		} else if customer.Description != "" {
			user.DisplayName = customer.Description
		}

		user.EmailAddress = customer.Email
		user.MapIDs[model.UserMapIDStripe] = customer.ID

		return nil
	}

	return derp.Wrap(err, "handler.stripe.loadOrCreateUser", "Error loading user record")
}

func loadStripeCustomer(apiKey string, customer *stripe.Customer) error {

	const location = "handler.stripe.loadStripeCustomer"

	if customer == nil {
		return derp.NewBadRequestError(location, "Customer must not be nil")
	}

	if customer.ID == "" {
		return derp.NewBadRequestError(location, "Customer.ID must not be empty")
	}

	// Create an API client
	stripeClient := client.API{}
	stripeClient.Init(apiKey, nil)

	// Load the Customer
	params := stripe.CustomerParams{}
	value, err := stripeClient.Customers.Get(customer.ID, &params)

	if err != nil {
		return derp.Wrap(err, location, "Error loading customer from Stripe API")
	}

	// Copy the value from the API call into the original customer
	*customer = *value

	// Success
	return nil
}

func loadStripeProduct(apiKey string, product *stripe.Product) error {

	const location = "handler.stripe.loadStripeProduct"

	if product == nil {
		return derp.NewBadRequestError(location, "Product must not be nil")
	}

	if product.ID == "" {
		return derp.NewBadRequestError(location, "Product.ID must not be empty")
	}

	// Create an API client
	stripeClient := client.API{}
	stripeClient.Init(apiKey, nil)

	// Load the Product
	params := stripe.ProductParams{}
	value, err := stripeClient.Products.Get(product.ID, &params)

	if err != nil {
		return derp.Wrap(err, location, "Error loading product from Stripe API")
	}

	// Copy the value from the API call into the original product
	*product = *value

	// Success
	return nil
}
