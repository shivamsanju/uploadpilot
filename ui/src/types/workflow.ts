export type Workflow = {
  variables: Record<string, string>;
  root: Statement;
};

export type Statement = {
  activity?: ActivityInvocation;
  sequence?: Sequence;
  parallel?: Parallel;
};

export type Sequence = {
  elements: Statement[];
};

export type Parallel = {
  branches: Statement[];
};

export type ActivityInvocation = {
  id: string;
  label: string;
  key: string;
  arguments: string[];
  result: string;
  retries: number;
  timeout: number;
  continueOnError: boolean;
  hasErrors?: boolean;
  error?: string;
};

export type EditableProperties = {
  label: string;
  retries: number;
  timeoutMs: number;
  continueOnError: boolean;
};
