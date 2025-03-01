import { Alert, Badge, Box, Group, Title } from '@mantine/core';
import { useLocalStorage } from '@mantine/hooks';
import { IconCloudUpload, IconInfoCircle } from '@tabler/icons-react';
import { useEffect, useState } from 'react';
import { useSetBreadcrumbs } from '../../hooks/breadcrumb';
import UploadList from './List';

const UploadsPage = () => {
  const [totalRecords, setTotalRecords] = useState(0);
  const setBreadcrumbs = useSetBreadcrumbs();
  const [alert, setAlert] = useLocalStorage({
    key: 'alert',
    defaultValue: 1,
    getInitialValueInEffect: false,
  });

  useEffect(() => {
    setBreadcrumbs([{ label: 'Workspaces', path: '/' }, { label: 'Uploads' }]);
  }, [setBreadcrumbs]);

  return (
    <Box>
      <Group mb="xl">
        <IconCloudUpload size={24} />
        <Title order={3}>Uploads</Title>
        <Badge variant="outline" radius="xl">
          {totalRecords}
        </Badge>
      </Group>

      {alert === 1 && (
        <Alert
          opacity={0.7}
          icon={<IconInfoCircle size={16} />}
          withCloseButton
          onClose={() => setAlert(0)}
        >
          Uploaded files may take a few seconds to appear here after upload.
          Please click the refresh button if you don't see any new files.
        </Alert>
      )}
      <UploadList setTotalRecords={setTotalRecords} />
    </Box>
  );
};

export default UploadsPage;
