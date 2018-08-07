package backend

import (
	"encoding/json"
	"errors"
	"thirdopinion/internal/pkg/config"
	"thirdopinion/internal/pkg/psql"
)

// ListArguments lists arguments from the database
func ListArguments(argument *config.Argument) ([]byte, error) {
	var filter string
	if argument.ID != 0 {
		filter = "specificPost"
	}
	m, err := psql.View(filter, argument.ID)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// NewArgument writes an argument to the database
func NewArgument(argument *config.Argument) (*config.WSResponse, error) {
	cr := verifyArgument(argument)
	if cr.Error != "" {
		return nil, errors.New(cr.Error)
	}
	resp, err := psql.Create(argument)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func verifyArgument(arg *config.Argument) config.CreationResult {
	if len(arg.Title) < 5 {
		return config.CreationResult{
			Error:    "Title too short",
			Position: "title",
		}
	}
	if len(arg.Opinions) <= 1 {
		return config.CreationResult{
			Error:    "You need at least two opinions for an argument",
			Position: "argument",
		}
	}
	return config.CreationResult{
		Error: "",
	}

}
