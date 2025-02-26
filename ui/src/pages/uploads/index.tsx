import {
  ActionIcon,
  Alert,
  Badge,
  Box,
  Breadcrumbs,
  Group,
  Text,
  Title,
} from '@mantine/core';
import { useLocalStorage } from '@mantine/hooks';
import { IconChevronLeft, IconInfoCircle } from '@tabler/icons-react';
import { useState } from 'react';
import { NavLink, useNavigate } from 'react-router-dom';
import UploadList from './List';

const UploadsPage = () => {
  const [totalRecords, setTotalRecords] = useState(0);
  const navigate = useNavigate();
  const [alert, setAlert] = useLocalStorage({
    key: 'alert',
    defaultValue: 1,
    getInitialValueInEffect: false,
  });

  return (
    <Box>
      <Breadcrumbs separator=">">
        <NavLink to="/" className="bredcrumb-link">
          <Text>Workspaces</Text>
        </NavLink>
        <Text>Uploads</Text>
      </Breadcrumbs>
      <Group mt="xs" mb="xl">
        <ActionIcon
          variant="default"
          radius="xl"
          size="sm"
          onClick={() => navigate(`/`)}
        >
          <IconChevronLeft size={16} />
        </ActionIcon>
        <Group align="center" gap="xs" h="10%">
          <Title order={3}>Uploads</Title>
          <Badge variant="outline" radius="xl">
            {totalRecords}
          </Badge>
        </Group>
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
