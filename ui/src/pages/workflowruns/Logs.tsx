import { Box, Modal, Stack } from '@mantine/core';
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
  finishedRun: boolean;
};
export const LogsModal: React.FC<LogsModalProps> = ({
  open,
  onClose,
  workspaceId,
  processorId,
  workflowId,
  runId,
  finishedRun,
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
          {!finishedRun && (
            <RefreshButton onClick={invalidate} mb="md" variant="filled" />
          )}
          <Stack>
            {logs?.length > 0 &&
              logs.map((item: any, index: number) => (
                <Box
                  key={index}
                  className={
                    item?.eventType === 'ActivityTaskScheduled'
                      ? classes.logScheduled
                      : ''
                  }
                >
                  <span className={classes.logTimestamp}>
                    {item?.timestamp}
                  </span>
                  {item?.eventType === 'ActivityTaskScheduled' ? (
                    <span className={classes.logEventScheduled}>
                      {`${item?.details?.split('"')[1]}`.replace(' ', '\n')}
                    </span>
                  ) : (
                    <span className={classes.logEvent}>
                      {item?.eventType.replace(' ', '\n')}
                    </span>
                  )}
                  <span className={classes.logDetails}>{item?.details}</span>
                </Box>
              ))}
          </Stack>
          {!finishedRun && (
            <RefreshButton onClick={invalidate} mt="md" variant="filled" />
          )}
        </Box>
      </ErrorLoadingWrapper>
    </Modal>
  );
};
