package gofed

import (
	"context"

	"github.com/EmissarySocial/emissary/model"
	"github.com/benpate/derp"
	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
)

func ToGoFed(item *model.ActivityStream) (vocab.Type, error) {
	return streams.ToType(context.TODO(), item.Content)
}

func ToModel(item vocab.Type, container model.ActivityStreamContainer) (model.ActivityStream, error) {
	result := model.NewActivityStream(container)
	data, err := streams.Serialize(item)

	if err != nil {
		return result, derp.Wrap(err, "gofed.ToModel", "Unable to serialize item", item)
	}

	result.Content = data

	// TODO: CRITICAL: Map from ActivityStreamID and UserID from the original data.

	return result, nil
}
