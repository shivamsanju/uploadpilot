import { Alert, Box, Button, Group, Title } from '@mantine/core';
import { useLocalStorage } from '@mantine/hooks';
import {
  IconCancel,
  IconCircleCheck,
  IconExclamationCircle,
  IconInfoCircle,
  IconNumber,
  IconServer2,
  IconTestPipe,
} from '@tabler/icons-react';
import { useEffect, useMemo, useState } from 'react';
import { SimpleKPICard } from '../../components/SimpleKpi';
import { useSetBreadcrumbs } from '../../hooks/breadcrumb';
import UploadList from './List';
import { TestUploaderModal } from './TestUploader';

const UploadsPage = () => {
  const [totalRecords, setTotalRecords] = useState(0);
  const setBreadcrumbs = useSetBreadcrumbs();
  const [openTestUploaderModal, setOpenTestUploaderModal] = useState(false);

  const [alert, setAlert] = useLocalStorage({
    key: 'alert',
    defaultValue: 1,
    getInitialValueInEffect: false,
  });

  useEffect(() => {
    setBreadcrumbs([{ label: 'Workspaces', path: '/' }, { label: 'Uploads' }]);
  }, [setBreadcrumbs]);

  const kpis = useMemo(() => {
    return [
      {
        label: 'Total',
        value: totalRecords,
        Icon: IconNumber,
        iconColor: 'gray',
      },
      {
        label: 'Successful',
        value: totalRecords,
        Icon: IconCircleCheck,
        iconColor: 'green',
      },
      {
        label: 'Failed',
        value: 0,
        Icon: IconExclamationCircle,
        iconColor: 'red',
      },
      {
        label: 'Cancelled',
        value: 0,
        Icon: IconCancel,
        iconColor: 'gray',
      },
    ];
  }, [totalRecords]);

  return (
    <Box mr="md">
      <Group mb="xl" justify="space-between">
        <Group align="center">
          <IconServer2 size={24} />
          <Title order={3}>Uploads</Title>
        </Group>
        <Button
          leftSection={<IconTestPipe size={16} />}
          onClick={() => setOpenTestUploaderModal(true)}
        >
          Test uploader
        </Button>
      </Group>
      <Group w="100%" grow>
        {kpis.map(kpi => (
          <SimpleKPICard
            key={kpi.label}
            title={kpi.label}
            value={kpi.value.toString()}
            Icon={kpi.Icon}
            iconColor={kpi.iconColor}
          />
        ))}
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
      <TestUploaderModal
        open={openTestUploaderModal}
        onClose={() => setOpenTestUploaderModal(false)}
      />
    </Box>
  );
};

export default UploadsPage;
