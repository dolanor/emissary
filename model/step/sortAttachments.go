package step

import (
	"github.com/benpate/rosetta/first"
	"github.com/benpate/rosetta/maps"
)

// SortAttachments represents an action-step that can update multiple records at once
type SortAttachments struct {
	Keys    string
	Values  string
	Message string
}

func NewSortAttachments(stepInfo maps.Map) (SortAttachments, error) {

	return SortAttachments{
		Keys:    first.String(stepInfo.GetString("keys"), "_id"),
		Values:  first.String(stepInfo.GetString("values"), "rank"),
		Message: stepInfo.GetString("message"),
	}, nil
}

// AmStep is here only to verify that this struct is a render pipeline step
func (step SortAttachments) AmStep() {}
