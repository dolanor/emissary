package render

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/EmissarySocial/emissary/model"
	"github.com/EmissarySocial/emissary/service"
	"github.com/benpate/data"
	"github.com/benpate/data/option"
	"github.com/benpate/derp"
	"github.com/benpate/exp"
	builder "github.com/benpate/exp-builder"
	"github.com/benpate/rosetta/schema"
	"github.com/benpate/steranko"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Profile struct {
	user *model.User
	Common
}

func NewProfile(factory Factory, ctx *steranko.Context, user *model.User, actionID string) (Profile, error) {

	// Load the Template
	templateService := factory.Template()
	template, err := templateService.Load("user-profile")

	if err != nil {
		return Profile{}, derp.Wrap(err, "render.NewProfile", "Error loading template")
	}

	// Verify the requested action is valid for this template
	action, ok := template.Action(actionID)

	if !ok {
		return Profile{}, derp.NewBadRequestError("render.NewProfile", "Invalid action", actionID)
	}

	return Profile{
		user:   user,
		Common: NewCommon(factory, ctx, template, action, actionID),
	}, nil
}

/******************************************
 * RENDERER INTERFACE
 ******************************************/

// Render generates the string value for this Profile
func (w Profile) Render() (template.HTML, error) {

	var buffer bytes.Buffer

	// Execute step (write HTML to buffer, update context)
	if err := Pipeline(w.action.Steps).Get(w._factory, &w, &buffer); err != nil {
		return "", derp.Report(derp.Wrap(err, "render.Profile.Render", "Error generating HTML", w._context.Request().URL.String()))

	}

	// Success!
	return template.HTML(buffer.String()), nil
}

// View executes a separate view for this Profile
func (w Profile) View(actionID string) (template.HTML, error) {

	renderer, err := NewProfile(w._factory, w._context, w.user, actionID)

	if err != nil {
		return template.HTML(""), derp.Wrap(err, "render.Profile.View", "Error creating Profile renderer")
	}

	return renderer.Render()
}

// NavigationID returns the ID to use for highlighing navigation menus
func (w Profile) NavigationID() string {

	// TODO: This is returning incorrect values when we CREATE a new outbox item.
	// Is there a better way to handle this that doesn't just HARDCODE stuff in here?

	// If the user is viewing their own profile, then the top-level ID is the user's own ID
	if w.UserID() == w.Common.AuthenticatedID().Hex() {

		switch w.ActionID() {
		case "inbox", "inbox-folder":
			return "inbox"
		default:
			return "profile"
		}
	}

	return ""
}

func (w Profile) PageTitle() string {
	return w.user.DisplayName
}

func (w Profile) Permalink() string {
	return w.Host() + "/@" + w.user.UserID.Hex()
}

func (w Profile) Token() string {
	return "users"
}

func (w Profile) object() data.Object {
	return w.user
}

func (w Profile) objectID() primitive.ObjectID {
	return w.user.UserID
}

func (w Profile) objectType() string {
	return "User"
}

func (w Profile) schema() schema.Schema {
	return schema.New(model.UserSchema())
}

func (w Profile) service() service.ModelService {
	return w._factory.User()
}

func (w Profile) templateRole() string {
	return "outbox"
}

// UserCan returns TRUE if this Request is authorized to access the requested view
func (w Profile) UserCan(actionID string) bool {

	action, ok := w._template.Action(actionID)

	if !ok {
		return false
	}

	authorization := w.authorization()

	return action.UserCan(w.user, &authorization)
}

/******************************************
 * Data Accessors
 ******************************************/

func (w Profile) UserID() string {
	return w.user.UserID.Hex()
}

// Myself returns TRUE if the current user is viewing their own profile
func (w Profile) Myself() bool {
	authorization := getAuthorization(w._context)

	if err := authorization.Valid(); err == nil {
		return authorization.UserID == w.user.UserID
	}

	return false
}

func (w Profile) Username() string {
	return w.user.Username
}

func (w Profile) FollowerCount() int {
	return w.user.FollowerCount
}

func (w Profile) FollowingCount() int {
	return w.user.FollowingCount
}

func (w Profile) BlockCount() int {
	return w.user.BlockCount
}

func (w Profile) DisplayName() string {
	return w.user.DisplayName
}

func (w Profile) StatusMessage() string {
	return w.user.StatusMessage
}

func (w Profile) ProfileURL() string {
	return w.user.ProfileURL
}

func (w Profile) ImageURL() string {
	return w.user.ActivityPubAvatarURL()
}

func (w Profile) Location() string {
	return w.user.Location
}

func (w Profile) Links() []model.PersonLink {
	return w.user.Links
}

func (w Profile) ActivityPubURL() string {
	return w.user.ActivityPubURL()
}

func (w Profile) ActivityPubAvatarURL() string {
	return w.user.ActivityPubAvatarURL()
}

func (w Profile) ActivityPubInboxURL() string {
	return w.user.ActivityPubInboxURL()
}

func (w Profile) ActivityPubOutboxURL() string {
	return w.user.ActivityPubOutboxURL()
}

func (w Profile) ActivityPubFollowersURL() string {
	return w.user.ActivityPubFollowersURL()
}

func (w Profile) ActivityPubFollowingURL() string {
	return w.user.ActivityPubFollowingURL()
}

func (w Profile) ActivityPubLikedURL() string {
	return w.user.ActivityPubLikedURL()
}

func (w Profile) ActivityPubPublicKeyURL() string {
	return w.user.ActivityPubPublicKeyURL()
}

/******************************************
 * Profile / Outbox Methods
 ******************************************/

func (w Profile) Outbox() *QueryBuilder[model.StreamSummary] {

	queryBuilder := builder.NewBuilder().
		Int("publishDate")

	criteria := exp.And(
		queryBuilder.Evaluate(w._context.Request().URL.Query()),
		exp.Equal("parentId", w.AuthenticatedID()),
	)

	result := NewQueryBuilder[model.StreamSummary](w._factory.Stream(), criteria)

	return &result
}

func (w Profile) Followers() *QueryBuilder[model.FollowerSummary] {

	queryBuilder := builder.NewBuilder().
		String("displayName")

	criteria := exp.And(
		queryBuilder.Evaluate(w._context.Request().URL.Query()),
		exp.Equal("parentId", w.AuthenticatedID()),
	)

	result := NewQueryBuilder[model.FollowerSummary](w._factory.Follower(), criteria)

	return &result
}

func (w Profile) Following() ([]model.FollowingSummary, error) {

	userID := w.AuthenticatedID()

	if userID.IsZero() {
		return nil, derp.NewUnauthorizedError("render.Profile.Following", "Must be signed in to view following")
	}

	followingService := w._factory.Following()

	return followingService.QueryByUser(userID)
}

func (w Profile) FollowingByFolder(token string) ([]model.FollowingSummary, error) {

	// Get the UserID from the authentication scope
	userID := w.AuthenticatedID()

	if userID.IsZero() {
		return nil, derp.NewUnauthorizedError("render.Profile.FollowingByFolder", "Must be signed in to view following")
	}

	// Get the followingID from the token
	followingID, err := primitive.ObjectIDFromHex(token)

	if err != nil {
		return nil, derp.Wrap(err, "render.Profile.FollowingByFolder", "Invalid following ID", token)
	}

	// Try to load the matching records
	followingService := w._factory.Following()
	return followingService.QueryByFolder(userID, followingID)

}

/******************************************
 * Inbox Methods
 ******************************************/

// Inbox returns a slice of messages in the current User's inbox
func (w Profile) Inbox() ([]model.Message, error) {

	// Must be authenticated to view any Inbox messages
	if !w.IsAuthenticated() {
		return []model.Message{}, derp.NewForbiddenError("render.Profile.Inbox", "Not authenticated")
	}

	expBuilder := builder.NewBuilder().
		ObjectID("origin.internalId").
		ObjectID("folderId").
		Int("rank")

	criteria := expBuilder.Evaluate(w._context.Request().URL.Query())

	return w._factory.Inbox().QueryByUserID(w.AuthenticatedID(), criteria, option.MaxRows(12), option.SortAsc("publishDate"))
}

// IsInboxEmpty returns TRUE if the inbox has no results and there are no filters applied
// This corresponds to there being NOTHING in the inbox, instead of just being filtered out.
func (w Profile) IsInboxEmpty(inbox []model.Message) bool {

	if len(inbox) > 0 {
		return false
	}

	if w._context.Request().URL.Query().Get("rank") != "" {
		return false
	}

	return true
}

// FIlteredByFollowing returns the Following record that is being used to filter the Inbox
func (w Profile) FilteredByFollowing() model.Following {

	result := model.NewFollowing()

	if !w.IsAuthenticated() {
		return result
	}

	token := w._context.QueryParam("origin.internalId")

	if followingID, err := primitive.ObjectIDFromHex(token); err == nil {
		followingService := w._factory.Following()

		if err := followingService.LoadByID(w.AuthenticatedID(), followingID, &result); err == nil {
			return result
		}
	}

	return result
}

// Folders returns a slice of all folders owned by the current User
func (w Profile) Folders() (model.FolderList, error) {

	result := model.NewFolderList()

	// User must be authenticated to view any folders
	if !w.IsAuthenticated() {
		return result, derp.NewForbiddenError("render.Profile.Folders", "Not authenticated")
	}

	folderService := w._factory.Folder()
	folders, err := folderService.QueryByUserID(w.AuthenticatedID())

	if err != nil {
		return result, derp.Wrap(err, "render.Profile.Folders", "Error loading folders")
	}

	result.Folders = folders
	return result, nil
}

func (w Profile) FoldersWithSelection() (model.FolderList, error) {

	// Get Folder List
	result, err := w.Folders()

	if err != nil {
		return result, derp.Wrap(err, "render.Profile.FoldersWithSelection", "Error loading folders")
	}

	// Get Selected FolderID
	token := w._context.QueryParam("folderId")

	if folderID, err := primitive.ObjectIDFromHex(token); err == nil {
		result.SelectedID = folderID
		return result, nil
	}

	if len(result.Folders) > 0 {
		result.SelectedID = result.Folders[0].FolderID
		return result, nil
	}

	return result, derp.NewInternalError("render.Profile.FoldersWithSelection", "No folders found", nil)
}

// Message uses the `messageId` URL parameter to load an individual message from the Inbox
func (w Profile) Message() (model.Message, error) {

	const location = "render.Profile.Message"

	result := model.NewMessage()

	// Guarantee that the user is signed in
	if !w.IsAuthenticated() {
		return result, derp.NewForbiddenError(location, "Not authenticated")
	}

	// Get Inbox Service
	inboxService := w._factory.Inbox()

	// Try to parse the messageID from the URL
	if messageID, err := primitive.ObjectIDFromHex(w._context.QueryParam("messageId")); err == nil {

		// Try to load an Activity record from the Inbox
		if err := inboxService.LoadByID(w.AuthenticatedID(), messageID, &result); err != nil {
			return result, derp.Wrap(err, location, "Error loading inbox item")
		}

		return result, nil
	}

	// Otherwise, look for folder/rank search parameters
	if folderToken := w._context.QueryParam("folderId"); folderToken != "" {
		if folderID, err := primitive.ObjectIDFromHex(folderToken); err == nil {

			var sort option.Option

			if strings.HasPrefix(w._context.QueryParam("rank"), "GT:") {
				sort = option.SortAsc("rank")
			} else {
				sort = option.SortDesc("rank")
			}

			expBuilder := builder.NewBuilder().
				ObjectID("origin.internalId").
				Int("rank")

			rank := expBuilder.Evaluate(w._context.Request().URL.Query())

			if err := inboxService.LoadByRank(w.AuthenticatedID(), folderID, rank, &result, sort); err != nil {
				return result, derp.Wrap(err, location, "Error loading inbox item")
			}

			return result, nil
		}
	}

	// Fall through means no valid parameters were found
	return result, derp.NewBadRequestError(location, "Invalid message ID", w._context.QueryParam("messageId"))
}
