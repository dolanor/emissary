package model

import (
	"encoding/json"

	"github.com/benpate/convert"
	"github.com/benpate/datatype"
	"github.com/benpate/derp"
	"github.com/whisperverse/whisperverse/model/step"
)

// Action holds the data for actions that can be performed on any Stream from a particular Template.
type Action struct {
	Roles      []string            `path:"roles"      json:"roles"      bson:"roles"`      // List of roles required to execute this Action.  If empty, then none are required.
	States     []string            `path:"states"     json:"states"     bson:"states"`     // List of states required to execute this Action.  If empty, then one are required.
	StateRoles map[string][]string `path:"stateRoles" json:"stateRoles" bson:"stateRoles"` // Map of states -> list of roles that can perform this Action (when a stream is in the given state)
	Steps      []step.Step         `path:"steps"      json:"steps"      bson:"steps"`      // List of steps to execute when GET-ing or POST-ing this Action.
}

// NewAction returns a fully initialized Action
func NewAction() Action {
	return Action{
		Roles:      make([]string, 0),
		States:     make([]string, 0),
		StateRoles: make(map[string][]string),
		Steps:      make([]step.Step, 0),
	}
}

// UserCan returns TRUE if this action is permitted on a stream (using the provided authorization)
func (action Action) UserCan(stream *Stream, authorization *Authorization) bool {

	// Get the list of request roles that the user has
	userRoles := stream.Roles(authorization)

	// If one or more of these matches the allowed roles, then the request is granted.
	return matchAny(userRoles, action.AllowedRoles(stream))
}

// AllowedRoles returns a string of all page request roles that are allowed to
// perform this action.  This includes system roles slike "anonymous", "authenticated", "author", and "owner".
func (action *Action) AllowedRoles(stream *Stream) []string {

	// If present, "States" limits the states where this action can take place at all.
	if len(action.States) > 0 {
		// If the current state is not present in the list of allowed states, then this action cannot be
		// taken until the stream is moved into a new state.
		if !matchOne(action.States, stream.StateID) {
			return make([]string, 0)
		}
	}

	// If no rules are present, then this action is available to everybody
	if (len(action.Roles) == 0) && (len(action.StateRoles) == 0) {
		return []string{"anonymous"}
	}

	// result will collect all of the allowable roles.
	result := []string{"owner"}

	// Add roles to the result
	result = append(result, action.Roles...)

	// If there's a corresponding entry in stateRoles, add that to the result, too.
	if stateRoles, ok := action.StateRoles[stream.StateID]; ok {
		result = append(result, stateRoles...)
	}

	return result
}

func (action *Action) UnmarshalJSON(data []byte) error {
	var asMap map[string]any

	if err := json.Unmarshal(data, &asMap); err != nil {
		return derp.Wrap(err, "model.Action.UnmarshalJSON", "Invalid JSON")
	}

	return action.UnmarshalMap(asMap)
}

func (action *Action) UnmarshalMap(data map[string]any) error {

	// Import easy values
	action.Roles = convert.SliceOfString(data["roles"])
	action.States = convert.SliceOfString(data["states"])

	// Import stateRoles
	action.StateRoles = make(map[string][]string, 0)
	stateRoles := convert.MapOfInterface(data["stateRoles"])
	for key, value := range stateRoles {
		action.StateRoles[key] = convert.SliceOfString(value)
	}

	// Import steps
	stepsInfo := convert.SliceOfMap(data["steps"])
	if pipeline, err := step.NewPipeline(stepsInfo); err == nil {
		action.Steps = pipeline
	} else {
		return derp.Wrap(err, "model.action.UnmarshalMap", "Error reading steps", stepsInfo)
	}

	// If no steps configued, then try the "step" alias
	if len(action.Steps) == 0 {
		if name := convert.String(data["step"]); name != "" {
			action.Steps, _ = step.NewPipeline([]datatype.Map{{"step": name}})
		}
	}

	return nil
}

// matchOne returns TRUE if the value matches one (or more) of the values in the slice
func matchOne(slice []string, value string) bool {
	for index := range slice {
		if slice[index] == value {
			return true
		}
	}

	return false
}

// matchAny returns TRUE if any of the values in slice1 are equal to any of the values in slice2
func matchAny(slice1 []string, slice2 []string) bool {

	for index1 := range slice1 {
		for index2 := range slice2 {
			if slice1[index1] == slice2[index2] {
				return true
			}
		}
	}

	return false
}
