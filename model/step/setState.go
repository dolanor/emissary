package step

import "github.com/benpate/rosetta/maps"

// SetState represents an action-step that can change a Stream's state
type SetState struct {
	StateID string
}

func NewSetState(stepInfo maps.Map) (SetState, error) {

	return SetState{
		StateID: stepInfo.GetString("state"),
	}, nil
}

// AmStep is here only to verify that this struct is a render pipeline step
func (step SetState) AmStep() {}
