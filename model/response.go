package model

import (
	"github.com/benpate/data/journal"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/first"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/toot/object"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Reponse defines a single Actor's response to an Object.  The actor may be a local or remote user, and the
// Object may be a local stream or an inbox message.
type Response struct {
	ResponseID primitive.ObjectID `json:"responseId" bson:"_id"`               // Unique identifier for this Response
	UserID     primitive.ObjectID `json:"userId"     bson:"userId"`            // ID of the INTERNAL user who made this response (Zero for external users)
	URL        string             `json:"url"        bson:"url"`               // URL of this Response document
	ActorID    string             `json:"actorId"    bson:"actorId"`           // URL of the Actor who made the response
	ObjectID   string             `json:"objectId"   bson:"objectId"`          // URL of the Object that the actor responded to
	Type       string             `json:"type"       bson:"type"`              // Type of Response (e.g. "like", "dislike", "favorite", "bookmark", "share", "reply", "repost", "follow", "subscribe", "tag", "flag", "comment", "mention", "react", "rsvpYes", "rsvpNo", "rsvpMaybe", "review")
	Summary    string             `json:"summary"    bson:"summary,omitempty"` // Summary of the response (e.g. "I liked this post because...")
	Content    string             `json:"content"    bson:"content,omitempty"` // Custom value assigned to the response (emoji, vote, etc.)

	journal.Journal `json:"-" bson:",inline"`
}

// NewReponse returns a fully initialized Response object
func NewResponse() Response {
	return Response{
		ResponseID: primitive.NewObjectID(),
	}
}

/******************************************
 * data.Object Interface
 ******************************************/

// ID returns the unique identifier for this Response (in string format)
func (response Response) ID() string {
	return response.ResponseID.Hex()
}

func (response Response) Fields() []string {
	return []string{"responseId", "url", "objectId", "type", "content", "createDate"}
}

/******************************************
 * Other Data Methods
 ******************************************/

// GetJSONLD returns the JSON-LD representation of this Response
func (response Response) GetJSONLD() mapof.Any {
	return mapof.Any{
		"@context": vocab.ContextTypeActivityStreams,
		"id":       response.URL,
		"type":     response.ActivityPubType(),
		"actor":    response.ActorID,
		"object":   response.ObjectID,
		"summary":  response.Summary,
		"content":  response.Content,
	}
}

// ActivityPubType converts a ResponseType into an ActivityStreams vocabulary type
func (response Response) ActivityPubType() string {

	switch response.Type {

	case ResponseTypeLike:
		return vocab.ActivityTypeLike

	case ResponseTypeDislike:
		return vocab.ActivityTypeDislike

	default:
		return vocab.ActivityTypeAnnounce
	}
}

func (response Response) EnglishType() string {

	switch response.Type {

	case ResponseTypeAnnounce:
		return "Shared"

	case ResponseTypeLike:
		return "Liked"

	case ResponseTypeDislike:
		return "Disliked"

	case ResponseTypeMention:
		return "Mentioned"

	default:
		return "Responded to"
	}
}

// IsEqual returns TRUE if two responses match urls, actors, objects, types, and values
func (response Response) IsEqual(other Response) bool {
	return (response.URL == other.URL) &&
		(response.ActorID == other.ActorID) &&
		(response.ObjectID == other.ObjectID) &&
		(response.Type == other.Type) &&
		(response.Content == other.Content)
}

// CalcContent sets the content of the response to a default value, if it is not already set.
func (response *Response) CalcContent() {

	// RULE: If the type is empty, then this is a "DELETE", so make the content is empty too.
	if response.Type == "" {
		response.Content = ""
		return
	}

	// Otherwise, set default content based on the response type.
	switch response.Type {

	case ResponseTypeMention:
		response.Content = "@"

	case ResponseTypeDislike:
		response.Content = first.String(response.Content, "👎")

	case ResponseTypeLike:
		response.Content = first.String(response.Content, "👍")

	default:
		response.Content = first.String(response.Content, "👍")
	}

	// Nothin to return.
}

// CreateDateSeconds returns the CreateDate in Unix Epoch seconds (instead of milliseconds)
func (response Response) CreateDateSeconds() int64 {
	return response.CreateDate / 1000
}

// IsEmpty returns TRUE if this Response has no data in it.
func (response Response) IsEmpty() bool {
	return response.Type == ""
}

// NotEmpty returns TRUE if this Response has data in it.
func (response Response) NotEmpty() bool {
	return !response.IsEmpty()
}

/******************************************
 * RoleStateEnumerator Methods
 ******************************************/

// State returns the current state of this Stream.  It is
// part of the implementation of the RoleStateEmulator interface
func (response Response) State() string {
	return ""
}

// Roles returns a list of all roles that match the provided authorization
func (response Response) Roles(authorization *Authorization) []string {

	// Everyone has "anonymous" access
	result := []string{MagicRoleAnonymous}

	if authorization.IsAuthenticated() {

		// Owners are hard-coded to do everything, so no other roles need to be returned.
		if authorization.DomainOwner {
			return []string{MagicRoleOwner}
		}

		result = append(result, MagicRoleAuthenticated)

		// Authors sometimes have special permissions, too.
		if response.UserID == authorization.UserID {
			result = append(result, MagicRoleAuthor)
		}
	}

	return result
}

/******************************************
 * Mastodon API
 ******************************************/

func (response Response) Toot() object.Status {

	return object.Status{
		ID:  response.URL,
		URI: response.URL,
		Account: object.Account{
			ID: response.ActorID,
		},
	}
}
