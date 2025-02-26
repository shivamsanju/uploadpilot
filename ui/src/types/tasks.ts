export type Task = {
  id: string;
  processorId: string;
  key: string;
  position: number;
  label: string;
  name: string;
  data: Record<string, any>;
  retries: number;
  timeoutMs: number;
  continueOnError: boolean;
  enabled: boolean;
  hasErrors?: boolean;
  error?: string;
};

export type EditableProperties = {
  name: string;
  retries: number;
  timeoutMs: number;
  continueOnError: boolean;
};
