package render

import (
	"github.com/EmissarySocial/emissary/config"
	"github.com/EmissarySocial/emissary/model"
	"github.com/EmissarySocial/emissary/queue"
	"github.com/EmissarySocial/emissary/service"
	"github.com/EmissarySocial/emissary/tools/set"
	"github.com/benpate/form"
	"github.com/benpate/icon"
	"github.com/benpate/mediaserver"
	"github.com/stripe/stripe-go/v72/client"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Factory is used to locate all necessary services
type Factory interface {
	// Model Services
	Attachment() *service.Attachment
	Folder() *service.Folder
	Following() *service.Following
	Follower() *service.Follower
	Group() *service.Group
	Inbox() *service.Inbox
	Layout() *service.Layout
	Mention() *service.Mention
	Stream() *service.Stream
	StreamDraft() *service.StreamDraft
	Template() *service.Template
	User() *service.User

	// Other data services
	Config() config.Domain
	Content() *service.Content
	Domain() *service.Domain
	Host() string
	Hostname() string
	Icons() icon.Provider
	MediaServer() mediaserver.MediaServer
	LookupProvider() form.LookupProvider
	Providers() set.Slice[config.Provider]
	Publisher() service.Publisher
	Queue() *queue.Queue
	StreamUpdateChannel() chan model.Stream
	StripeClient() (client.API, error)

	WebSubInbox(primitive.ObjectID) service.WebSubInbox
	WebSubOutbox(primitive.ObjectID) service.WebSubOutbox
}
