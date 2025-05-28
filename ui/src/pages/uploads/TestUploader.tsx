import { Box, Button, Group, Modal, Paper } from '@mantine/core';
import { FilePond } from 'react-filepond';
import '../../style/filepond.css';

// Import FilePond styles
import 'filepond/dist/filepond.min.css';
import { useParams } from 'react-router-dom';
import { getTenantId } from '../../utils/config';
import { Uploader } from '../../utils/uploader';

type Props = {
  open: boolean;
  onClose: () => void;
};
export const TestUploaderModal: React.FC<Props> = ({ open, onClose }) => {
  const { workspaceId } = useParams();

  return (
    <Modal
      component={Paper}
      opened={open}
      onClose={onClose}
      size="xl"
      title="Test uploader"
      closeOnClickOutside={false}
      withCloseButton={false}
      closeOnEscape={false}
    >
      <Box>
        <FilePond
          allowMultiple={true}
          maxFiles={10}
          server={{
            process: (
              fieldName,
              file,
              metadata,
              load,
              error,
              progress,
              abort,
              transfer,
              options,
            ) => {
              const uploader = new Uploader(
                file as File,
                getTenantId()!,
                workspaceId!,
                '',
                metadata,
              );

              uploader.onComplete((serverFileId: string) => {
                load(serverFileId);
              });

              uploader.onError((message: string) => {
                error(message);
              });

              uploader.onProgress((computable, loaded, total) => {
                progress(computable, loaded, total);
              });

              uploader.start();

              return {
                abort: () => {
                  uploader.cancel();
                  abort();
                },
              };
            },
          }}
          allowRevert={false}
          instantUpload={false}
          credits={false}
          labelFileProcessingError={(err: any) => err.body}
        />
        <Group justify="flex-end" mt="50">
          <Button onClick={onClose}>Close</Button>
        </Group>
      </Box>
    </Modal>
  );
};
