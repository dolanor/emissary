package service

import (
	"io"
)

type Icons struct{}

func (service Icons) Get(name string) string {
	switch name {

	// App Actions and Behaviors
	case "add":
		return service.get("plus-lg")
	case "add-circle":
		return service.get("plus-circle")
	case "add-emoji":
		return service.get("emoji-smile")
	case "archive":
		return service.get("archive")
	case "archive-fill":
		return service.get("archive-fill")
	case "book":
		return service.get("book")
	case "book-fill":
		return service.get("book-fill")
	case "bookmark":
		return service.get("bookmark")
	case "bookmark-fill":
		return service.get("bookmark-fill")
	case "box":
		return service.get("box")
	case "box-fill":
		return service.get("box-fill")
	case "calendar":
		return service.get("calendar3")
	case "calendar-fill":
		return service.get("calendar3-week-fill")
	case "cancel":
		return service.get("x-lg")
	case "chat":
		return service.get("chat")
	case "chat-fill":
		return service.get("chat-fill")
	case "check-badge":
		return service.get("patch-check")
	case "check-badge-fill":
		return service.get("patch-check-fill")
	case "check-circle":
		return service.get("check-circle")
	case "check-circle-fill":
		return service.get("check-circle-fill")
	case "check-shield":
		return service.get("shield-check")
	case "check-shield-fill":
		return service.get("shield-check-fill")
	case "chevron-left":
		return service.get("chevron-left")
	case "chevron-right":
		return service.get("chevron-right")
	case "circle":
		return service.get("circle")
	case "circle-fill":
		return service.get("circle-fill")
	case "clipboard":
		return service.get("clipboard")
	case "clipboard-fill":
		return service.get("clipboard-fill")
	case "clock":
		return service.get("clock")
	case "clock-fill":
		return service.get("clock-fill")
	case "cloud":
		return service.get("cloud")
	case "cloud-fill":
		return service.get("cloud-fill")
	case "database":
		return service.get("database")
	case "database-fill":
		return service.get("database-fill")
	case "delete":
		return service.get("trash")
	case "delete-fill":
		return service.get("trash-fill")
	case "drag-handle":
		return service.get("grip-vertical")
	case "edit":
		return service.get("pencil-square")
	case "edit-fill":
		return service.get("pencil-square-fill")
	case "email":
		return service.get("envelope")
	case "email-fill":
		return service.get("envelope-fill")
	case "file":
		return service.get("file-earmark")
	case "file-fill":
		return service.get("file-earmark-fill")
	case "filter":
		return service.get("filter-circle")
	case "filter-fill":
		return service.get("filter-circle-fill")
	case "flag":
		return service.get("flag")
	case "flag-fill":
		return service.get("flag-fill")
	case "folder":
		return service.get("folder")
	case "folder-fill":
		return service.get("folder-fill")
	case "globe":
		return service.get("globe2")
	case "globe-fill":
		return service.get("globe2")
	case "grip-vertical":
		return service.get("grip-vertical")
	case "grip-horizontal":
		return service.get("grip-horizontal")
	case "hashtag":
		return service.get("hash")
	case "heart":
		return service.get("heart")
	case "heart-fill":
		return service.get("heart-fill")
	case "home":
		return service.get("house")
	case "home-fill":
		return service.get("house-fill")
	case "info":
		return service.get("info-circle")
	case "info-fill":
		return service.get("info-circle-fill")
	case "invisible":
		return service.get("eye-slash")
	case "invisible-fill":
		return service.get("eye-slash-fill")
	case "journal":
		return service.get("journal")
	case "link":
		return service.get("link-45deg")
	case "link-outbound":
		return service.get("box-arrow-up-right")
	case "location":
		return service.get("geo-alt")
	case "location-fill":
		return service.get("geo-alt-fill")
	case "lock":
		return service.get("lock")
	case "lock-fill":
		return service.get("lock-fill")
	case "loading":
		return service.get("arrow-clockwise")
	case "login":
		return service.get("box-arrow-in-right")
	case "megaphone":
		return service.get("megaphone")
	case "megaphone-fill":
		return service.get("megaphone-fill")
	case "mention":
		return service.get("at")
	case "more-horizontal":
		return service.get("three-dots")
	case "more-vertical":
		return service.get("three-dots-vertical")
	case "mute":
		return service.get("mic-mute")
	case "mute-fill":
		return service.get("mic-mute-fill")
	case "newspaper":
		return service.get("newspaper")
	case "person":
		return service.get("person")
	case "person-fill":
		return service.get("person-fill")
	case "people":
		return service.get("people")
	case "people-fill":
		return service.get("people-fill")
	case "reply":
		return service.get("reply")
	case "reply-fill":
		return service.get("reply-fill")
	case "repost":
		return service.get("repeat")
	case "rocket":
		return service.get("rocket-takeoff")
	case "rocket-fill":
		return service.get("rocket-takeoff-fill")
	case "rule":
		return service.get("funnel")
	case "rule-fill":
		return service.get("funnel-fill")
	case "save":
		return service.get("check-lg")
	case "search":
		return service.get("search")
	case "settings":
		return service.get("gear")
	case "settings-fill":
		return service.get("gear-fill")
	case "server":
		return service.get("hdd-stack")
	case "server-fill":
		return service.get("hdd-stack-fill")
	case "share":
		return service.get("arrow-up-right-square")
	case "share-fill":
		return service.get("arrow-up-right-square-fill")
	case "shield":
		return service.get("shield")
	case "shield-fill":
		return service.get("shield-fill")
	case "star":
		return service.get("star")
	case "star-fill":
		return service.get("star-fill")
	case "thumbs-down":
		return service.get("hand-thumbs-down")
	case "thumbs-down-fill":
		return service.get("hand-thumbs-down-fill")
	case "thumbs-up":
		return service.get("hand-thumbs-up")
	case "thumbs-up-fill":
		return service.get("hand-thumbs-up-fill")
	case "unlink":
		return service.get("link-45deg")
	case "upload":
		return service.get("upload")
	case "user":
		return service.get("person-circle")
	case "user-fill":
		return service.get("person-circle-fill")
	case "users":
		return service.get("people")
	case "users-fill":
		return service.get("people-fill")
	case "visible":
		return service.get("eye")
	case "visible-fill":
		return service.get("eye-fill")

		// Layouts
	case "layout-social":
		return service.get("list-ul")
	case "layout-social-fill":
		return service.get("list-ul")

	case "layout-chat":
		return service.get("chat-text")
	case "layout-chat-fill":
		return service.get("chat-text")

	case "layout-newspaper":
		return service.get("postcard")
	case "layout-newspaper-fill":
		return service.get("postcard")

	case "layout-magazine":
		return service.get("view-stacked")
	case "layout-magazine-fill":
		return service.get("view-stacked")

		// Services
	case "activitypub":
		return service.get("globe2")
	case "activitypub-fill":
		return service.get("globe2")
	case "facebook":
		return service.get("facebook")
	case "github":
		return service.get("github")
	case "google":
		return service.get("google")
	case "json":
		return service.get("braces")
	case "json-fill":
		return service.get("braces")
	case "instagram":
		return service.get("instagram")
	case "twitter":
		return service.get("twitter")
	case "rss":
		return service.get("rss")
	case "rss-fill":
		return service.get("rss-fill")
	case "rss-cloud":
		return service.get("cloud-arrow-down")
	case "rss-cloud-fill":
		return service.get("cloud-arrow-down-fill")
	case "stripe":
		return service.get("credit-card")
	case "stripe-fill":
		return service.get("credit-card-fill")
	case "websub":
		return service.get("cloud-arrow-down")
	case "websub-fill":
		return service.get("cloud-arrow-down-fill")

	// Content Types
	case "article":
		return service.get("file-text")
	case "article-fill":
		return service.get("file-text-fill")
	case "block":
		return service.get("slash-circle")
	case "block-fill":
		return service.get("slash-circle-fill")
	case "collection":
		return service.get("view-stacked")
	case "forward":
		return service.get("forward")
	case "forward-fill":
		return service.get("forward-fill")
	case "html":
		return service.get("code-slash")
	case "html-fill":
		return service.get("code-slash")
	case "inbox":
		return service.get("inbox")
	case "inbox-fill":
		return service.get("inbox-fill")
	case "markdown":
		return service.get("markdown")
	case "markdown-fill":
		return service.get("markdown-fill")
	case "message":
		return service.get("chat-left-text")
	case "message-fill":
		return service.get("chat-left-text-fill")
	case "outbox":
		return service.get("envelope")
	case "outbox-fill":
		return service.get("envelope-fill")
	case "picture":
		return service.get("image")
	case "picture-fill":
		return service.get("image-fill")
	case "pictures":
		return service.get("images")
	case "shopping-cart":
		return service.get("cart")
	case "shopping-cart-fill":
		return service.get("cart-fill")
	case "video":
		return service.get("camera-video")
	case "video-fill":
		return service.get("camera-video-fill")
	}

	return service.get(name)
}

func (service Icons) get(name string) string {
	return `<i class="bi bi-` + name + `"></i>`
}

func (service Icons) Write(name string, writer io.Writer) {
	// Okay to ignore write error
	// nolint:errcheck
	writer.Write([]byte(service.Get(name)))
}
