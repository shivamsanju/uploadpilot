import { CodeHighlight } from '@mantine/code-highlight';
import { Modal } from '@mantine/core';
import React from 'react';
import { useGetProcessorRunLogs } from '../../apis/processors';
import { RefreshButton } from '../../components/Buttons/RefreshButton/RefreshButton';
import { ErrorLoadingWrapper } from '../../components/ErrorLoadingWrapper';

const formatLogs = (logs: any[]) => {
  let str = '';
  for (const log of logs) {
    str += `${log?.timestamp} | ${log?.eventType?.toUpperCase()} | ${log?.details}\n`;
  }
  return str;
};
type LogsModalProps = {
  open: boolean;
  onClose: () => void;
  workspaceId: string;
  processorId: string;
  workflowId: string;
  runId: string;
};
export const LogsModal: React.FC<LogsModalProps> = ({
  open,
  onClose,
  workspaceId,
  processorId,
  workflowId,
  runId,
}) => {
  const { isPending, error, logs, invalidate } = useGetProcessorRunLogs(
    workspaceId,
    processorId,
    workflowId,
    runId,
  );

  return (
    <Modal opened={open} fullScreen title="Logs" onClose={onClose} size="100%">
      <ErrorLoadingWrapper isPending={isPending} error={error}>
        <RefreshButton onClick={invalidate} my="sm" />
        <CodeHighlight mih={300} code={formatLogs(logs || [])} />
      </ErrorLoadingWrapper>
    </Modal>
  );
};
