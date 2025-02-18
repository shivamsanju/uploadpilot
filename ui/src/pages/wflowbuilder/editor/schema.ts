import { parse } from "yaml";
import Ajv, { ErrorObject } from "ajv";

const schema = {
  $schema: "http://json-schema.org/draft-07/schema#",
  type: "object",
  title: "Workflow",
  properties: {
    variables: {
      type: "object",
      additionalProperties: {
        type: "string",
      },
      description: "Map of variables that can be used as input to Activities.",
    },
    root: {
      $ref: "#/definitions/Statement",
    },
  },
  required: ["variables", "root"],
  definitions: {
    Statement: {
      type: "object",
      description:
        "A building block of the workflow, which can be an activity, sequence, parallel execution, condition, or loop.",
      properties: {
        activity: { $ref: "#/definitions/ActivityInvocation" },
        sequence: { $ref: "#/definitions/Sequence" },
        parallel: { $ref: "#/definitions/Parallel" },
        condition: { $ref: "#/definitions/Condition" },
        loop: { $ref: "#/definitions/Loop" },
      },
      oneOf: [
        { required: ["activity"] },
        { required: ["sequence"] },
        { required: ["parallel"] },
        { required: ["condition"] },
        { required: ["loop"] },
      ],
    },
    Sequence: {
      type: "object",
      description: "A collection of statements that run sequentially.",
      properties: {
        elements: {
          type: "array",
          items: { $ref: "#/definitions/Statement" },
        },
      },
      required: ["elements"],
    },
    Parallel: {
      type: "object",
      description: "A collection of statements that run in parallel.",
      properties: {
        branches: {
          type: "array",
          items: { $ref: "#/definitions/Statement" },
        },
      },
      required: ["branches"],
    },
    Condition: {
      type: "object",
      description: "Defines a conditional execution block.",
      properties: {
        variable: {
          type: "string",
          description: "The variable to check.",
        },
        value: {
          type: "string",
          description: "The expected value to satisfy the condition.",
        },
        then: { $ref: "#/definitions/Statement" },
        else: { $ref: "#/definitions/Statement" },
      },
      required: ["variable", "value", "then", "else"],
    },
    Loop: {
      type: "object",
      description: "Defines a loop execution block.",
      properties: {
        breakvariable: {
          type: "string",
          description: "The variable to check for breaking the loop.",
        },
        breakvalue: {
          type: "string",
          description: "The expected value to break the loop.",
        },
        iterations: {
          type: "integer",
          minimum: 1,
          description: "Number of iterations to run.",
        },
        body: { $ref: "#/definitions/Statement" },
      },
      required: ["iterations", "body"],
    },
    ActivityInvocation: {
      type: "object",
      description:
        "Defines an activity invocation with arguments and execution properties.",
      properties: {
        name: {
          type: "string",
          description: "The name of the activity to invoke.",
        },
        arguments: {
          type: "array",
          items: { type: "string" },
          description: "List of arguments passed to the activity.",
        },
        result: {
          type: "string",
          description:
            "The name of the variable where the result will be stored.",
        },
      },
      required: ["name", "arguments"],
    },
  },
};

export const validateYaml = (content: string): string | null => {
  try {
    const parsedData = parse(content) as unknown;
    const ajv = new Ajv();
    const validate = ajv.compile(schema);
    const valid = validate(parsedData);

    if (!valid) {
      return (
        validate.errors?.map((err: ErrorObject) => err.message).join(", ") ||
        "Invalid YAML"
      );
    } else {
      return null;
    }
  } catch (e) {
    return (e as Error).message;
  }
};
