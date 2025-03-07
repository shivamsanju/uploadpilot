import { Box, Group, Modal, Text, Timeline } from '@mantine/core';
import { useGetProcessorRunLogs } from '../../apis/processors';
import { RefreshButton } from '../../components/Buttons/RefreshButton/RefreshButton';
import { ErrorLoadingWrapper } from '../../components/ErrorLoadingWrapper';
import classes from './Runs.module.css';

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
        <Box px="md">
          <RefreshButton onClick={invalidate} mb="md" variant="outline" />
          <Timeline bulletSize={14}>
            {logs?.length > 0 &&
              logs.map((item: any, index: number) => (
                <Timeline.Item
                  key={index}
                  title={
                    <Group align="center" gap="md">
                      <Text c="dimmed">{item?.timestamp}</Text>
                      {item?.eventType === 'ActivityTaskScheduled' ? (
                        <Text fw="bold" c="green">
                          {`Scheduled -> ${item?.details?.split(',')[0]}`}
                        </Text>
                      ) : (
                        <Text c="blue">{item?.eventType}</Text>
                      )}
                    </Group>
                  }
                >
                  <code className={classes.codetext}>{item?.details}</code>
                </Timeline.Item>
              ))}
          </Timeline>
        </Box>
      </ErrorLoadingWrapper>
    </Modal>
  );
};
