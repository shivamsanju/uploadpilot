package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type (
	// Workflow is the type used to express the workflow definition. Variables are a map of valuables. Variables can be
	// used as input to Activity.
	Workflow struct {
		Variables map[string]string `json:"variables"`
		Root      Statement         `json:"root"`
	}

	// Statement is the building block of dsl workflow. A Statement can be a simple ActivityInvocation or it
	// could be a Sequence or Parallel.
	Statement struct {
		Activity *ActivityInvocation `json:"activity"`
		Sequence *Sequence           `json:"sequence"`
		Parallel *Parallel           `json:"parallel"`
	}

	// Sequence consist of a collection of Statements that runs in sequential.
	Sequence struct {
		Elements []*Statement `json:"elements"`
	}

	// Parallel can be a collection of Statements that runs in parallel.
	Parallel struct {
		Branches []*Statement `json:"branches"`
	}

	// ActivityInvocation is used to express invoking an Activity. The Arguments defined expected arguments as input to
	// the Activity, the result specify the name of variable that it will store the result as which can then be used as
	// arguments to subsequent ActivityInvocation.
	ActivityInvocation struct {
		Name            string   `json:"name"`
		Label           string   `json:"label"`
		Arguments       []string `json:"arguments"`
		Result          string   `json:"result"`
		Retries         uint64   `json:"retries"`
		TimeoutMs       uint64   `json:"timeoutMs"`
		ContinueOnError bool     `json:"continueOnError"`
	}
)

func (s Statement) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan implements the sql.Scanner interface.
func (s *Statement) Scan(value interface{}) error {
	if value == nil {
		*s = Statement{}
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan struct: invalid type")
	}

	return json.Unmarshal(b, s)
}
