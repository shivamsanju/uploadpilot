import { CodeHighlight } from "@mantine/code-highlight";
import { Modal } from "@mantine/core";
import React from "react";
import { UseGetUploadLogs } from "../../apis/upload";
import { RefreshButton } from "../../components/Buttons/RefreshButton/RefreshButton";
import { ErrorLoadingWrapper } from "../../components/ErrorLoadingWrapper";

const formatLogs = (logs: any[]) => {
  let str = "";
  for (const log of logs) {
    str += `${log?.timestamp} | [${log?.level?.toUpperCase()}] | ${
      log?.processorId
        ? log?.taskId &&
          `(processor: '${log?.processorId}', task: '${log?.taskId}')`
        : ""
    } ${log?.message}\n`;
  }
  return str;
};
type LogsModalProps = {
  open: boolean;
  onClose: () => void;
  uploadId: string;
  workspaceId: string;
};
export const LogsModal: React.FC<LogsModalProps> = ({
  open,
  onClose,
  uploadId,
  workspaceId,
}) => {
  const { isPending, error, data, invalidate } = UseGetUploadLogs(
    workspaceId,
    uploadId
  );

  return (
    <Modal opened={open} fullScreen title="Logs" onClose={onClose} size="100%">
      <ErrorLoadingWrapper isPending={isPending} error={error}>
        <RefreshButton onClick={invalidate} my="sm" />
        <CodeHighlight mih={300} code={formatLogs(data || [])} />
      </ErrorLoadingWrapper>
    </Modal>
  );
};
