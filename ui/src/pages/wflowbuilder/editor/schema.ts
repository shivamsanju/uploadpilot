import Ajv, { ErrorObject } from 'ajv';
import { parse } from 'yaml';

const schema = {
  $schema: 'http://json-schema.org/draft-07/schema#',
  type: 'object',
  title: 'Workflow',
  properties: {
    variables: {
      type: 'object',
      description: 'Map of variables that can be used as input to Activities.',
    },
    root: {
      $ref: '#/definitions/Statement',
    },
  },
  required: ['variables', 'root'],
  definitions: {
    Statement: {
      type: 'object',
      description:
        'A building block of the workflow, which can be an activity, sequence, parallel execution, condition, or loop.',
      properties: {
        activity: { $ref: '#/definitions/ActivityInvocation' },
        sequence: { $ref: '#/definitions/Sequence' },
        parallel: { $ref: '#/definitions/Parallel' },
        condition: { $ref: '#/definitions/Condition' },
        loop: { $ref: '#/definitions/Loop' },
      },
      oneOf: [
        { required: ['activity'] },
        { required: ['sequence'] },
        { required: ['parallel'] },
        { required: ['condition'] },
        { required: ['loop'] },
      ],
    },
    Sequence: {
      type: 'object',
      description: 'A collection of statements that run sequentially.',
      properties: {
        elements: {
          type: 'array',
          items: { $ref: '#/definitions/Statement' },
        },
      },
      required: ['elements'],
    },
    Parallel: {
      type: 'object',
      description: 'A collection of statements that run in parallel.',
      properties: {
        branches: {
          type: 'array',
          items: { $ref: '#/definitions/Statement' },
        },
      },
      required: ['branches'],
    },
    Condition: {
      type: 'object',
      description: 'Defines a conditional execution block.',
      properties: {
        variable: {
          type: 'string',
          description: 'The variable to check.',
        },
        value: {
          type: 'string',
          description: 'The expected value to satisfy the condition.',
        },
        then: { $ref: '#/definitions/Statement' },
        else: { $ref: '#/definitions/Statement' },
      },
      required: ['variable', 'value', 'then'],
    },
    Loop: {
      type: 'object',
      description: 'Defines a loop execution block.',
      properties: {
        iterations: {
          type: 'integer',
          minimum: 1,
          description: 'Number of iterations to run.',
        },
        body: { $ref: '#/definitions/Statement' },
        breakVariable: {
          type: 'string',
          description: 'The variable to check for breaking the loop.',
        },
        breakValue: {
          type: 'string',
          description: 'The expected value to break the loop.',
        },
      },
      required: ['iterations', 'body'],
    },
    ActivityInvocation: {
      type: 'object',
      description:
        'Defines an activity invocation with arguments and execution properties.',
      properties: {
        key: {
          type: 'string',
          description: 'A unique key identifying the activity.',
        },
        uses: {
          type: 'string',
          description: 'The name of the activity to invoke.',
        },
        with: {
          type: 'object',
          description: 'Key-value pairs for activity parameters.',
        },
        input: {
          type: 'string',
          description: 'Optional input string for the activity.',
        },
        scheduleToCloseTimeoutSeconds: { type: 'integer', minimum: 1 },
        scheduleToStartTimeoutSeconds: { type: 'integer', minimum: 1 },
        startToCloseTimeoutSeconds: { type: 'integer', minimum: 1 },
        maxRetries: { type: 'integer', minimum: 0 },
        retryBackoffCoefficient: { type: 'number', minimum: 0 },
        retryMaxIntervalSeconds: { type: 'integer', minimum: 1 },
        retryInitialIntervalSeconds: { type: 'integer', minimum: 1 },
      },
      required: ['key', 'uses'],
    },
  },
};

export const validateWorkflowContent = (content: string): string | null => {
  try {
    const formattedContent = content.replace(/\t/g, '  ');
    const parsedData = parse(formattedContent) as unknown;
    validateUniqueKeys((parsedData as any).root);

    const ajv = new Ajv();
    const validate = ajv.compile(schema);
    const valid = validate(parsedData);

    if (!valid) {
      return (
        validate.errors?.map((err: ErrorObject) => err.message).join(', ') ||
        'Invalid YAML'
      );
    } else {
      return null;
    }
  } catch (e) {
    return (e as Error).message;
  }
};

function validateUniqueKeys(statement: any, seenKeys = new Set()) {
  if (!statement) return true;

  // Check activity key uniqueness
  if (statement.activity) {
    if (seenKeys.has(statement.activity.key)) {
      throw new Error(
        `Duplicate activity key found: ${statement.activity.key}`,
      );
    }
    seenKeys.add(statement.activity.key);
  }

  // Recursively check nested structures
  if (statement.sequence) {
    for (const element of statement.sequence.elements) {
      validateUniqueKeys(element, seenKeys);
    }
  }
  if (statement.parallel) {
    for (const branch of statement.parallel.branches) {
      validateUniqueKeys(branch, seenKeys);
    }
  }
  if (statement.condition) {
    validateUniqueKeys(statement.condition.then, seenKeys);
    if (statement.condition.else) {
      validateUniqueKeys(statement.condition.else, seenKeys);
    }
  }
  if (statement.loop) {
    validateUniqueKeys(statement.loop.body, seenKeys);
  }

  return true;
}
