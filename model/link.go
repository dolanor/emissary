package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// LinkRelationAuthor refers to the context's author.
// https://html.spec.whatwg.org/multipage/links.html#link-type-author
const LinkRelationAuthor = "author"

// LinkRelationBookmark gives a permanent link to use for bookmarking purposes
// https://html.spec.whatwg.org/multipage/links.html#link-type-bookmark
const LinkRelationBookmark = "bookmark"

// LinkRelationInReplyTo indicates that this document is a reply to another document
// https://www.rfc-editor.org/rfc/rfc4685.html
const LinkRelationInReplyTo = "in-reply-to"

// LinkRelationOriginal points to an Original Resource.
// https://www.rfc-editor.org/rfc/rfc7089.html#section-2.2.1
const LinkRelationOriginal = "original"

// LinkRelationProfile refers to the context's author.
// https://www.rfc-editor.org/rfc/rfc6906.html
const LinkRelationProfile = "profile"

// LinkSourceActivityPub indicates that this linked document was generated by an ActivityPub source
const LinkSourceActivityPub = "ACTIVITYPUB"

// LinkSourceInternal indicates that this linked document was generated by Emissary
const LinkSourceInternal = "INTERNAL"

// LinkSourceRSS indicates that that this linked document was generated by an RSS Feed
const LinkSourceRSS = "RSS"

// LinkSourceTwitter indicates that this linked document was generated by Twitter
const LinkSourceTwitter = "TWITTER"

// Link represents a link to another document on the Internet.
type Link struct {
	Relation   string             `path:"rel"        json:"rel"        bson:"rel"`                  // The relationship of the linked document, per https://www.iana.org/assignments/link-relations/link-relations.xhtml
	Source     string             `path:"source"     json:"source"     bson:"source"`               // The source of the link.  This could be "ACTIVITYPUB", "RSS", "TWITTER", or "EMAIL"
	InternalID primitive.ObjectID `path:"internalId" json:"internalId" bson:"internalId,omitempty"` // Unique ID of a document in this database
	Label      string             `path:"label"      json:"label"      bson:"label,omitempty"`      // Label of the link
	URL        string             `path:"url"        json:"url"        bson:"url,omitempty"`        // Public URL of the document
	UpdateDate int64              `path:"updateDate" json:"updateDate" bson:"updateDate"`           // Unix timestamp of the date/time when this link was last updated.
}

// ID implements the Set.ID interface
func (link Link) ID() string {
	return link.Relation
}

// IsEmpty returns TRUE if this link is empty
func (link Link) IsEmpty() bool {
	return (link.URL == "" && link.InternalID.IsZero())
}

// IsPresent returns TRUE if this link has a valid value
func (link Link) IsPresent() bool {
	return !link.IsEmpty()
}

func (link Link) HTML() string {
	if link.URL == "" {
		return ""
	}

	return "<link rel=\"" + link.Relation + "\" href=\"" + link.URL + "\">"
}
