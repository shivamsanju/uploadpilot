package dsl

const DSLSchema = `{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"type": "object",
	"title": "Workflow",
	"properties": {
	  "variables": {
		"type": "object",
		"additionalProperties": {
		  "type": "string"
		},
		"description": "Map of variables that can be used as input to Activities."
	  },
	  "root": {
		"$ref": "#/definitions/Statement"
	  }
	},
	"required": ["variables", "root"],
	"definitions": {
	  "Statement": {
		"type": "object",
		"description": "A building block of the workflow, which can be an activity, sequence, or parallel execution.",
		"properties": {
		  "activity": { "$ref": "#/definitions/ActivityInvocation" },
		  "sequence": { "$ref": "#/definitions/Sequence" },
		  "parallel": { "$ref": "#/definitions/Parallel" }
		},
		"oneOf": [
		  { "required": ["activity"] },
		  { "required": ["sequence"] },
		  { "required": ["parallel"] }
		]
	  },
	  "Sequence": {
		"type": "object",
		"description": "A collection of statements that run sequentially.",
		"properties": {
		  "elements": {
			"type": "array",
			"items": { "$ref": "#/definitions/Statement" }
		  }
		},
		"required": ["elements"]
	  },
	  "Parallel": {
		"type": "object",
		"description": "A collection of statements that run in parallel.",
		"properties": {
		  "branches": {
			"type": "array",
			"items": { "$ref": "#/definitions/Statement" }
		  }
		},
		"required": ["branches"]
	  },
	  "ActivityInvocation": {
		"type": "object",
		"description": "Defines an activity invocation with arguments and execution properties.",
		"properties": {
		  "name": {
			"type": "string",
			"description": "The name of the activity to invoke."
		  },
		  "label": {
			"type": "string",
			"description": "A label for the activity."
		  },
		  "arguments": {
			"type": "array",
			"items": { "type": "string" },
			"description": "List of arguments passed to the activity."
		  },
		  "result": {
			"type": "string",
			"description": "The name of the variable where the result will be stored."
		  },
		  "retries": {
			"type": "integer",
			"minimum": 0,
			"description": "Number of retry attempts allowed."
		  },
		  "timeoutMs": {
			"type": "integer",
			"minimum": 0,
			"description": "Timeout for the activity in milliseconds."
		  },
		  "continueOnError": {
			"type": "boolean",
			"description": "Specifies if execution should continue even after an error."
		  }
		},
		"required": ["name", "arguments"]
	  }
	}
  }`
