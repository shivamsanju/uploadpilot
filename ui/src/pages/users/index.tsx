import {
  ActionIcon,
  Box,
  Breadcrumbs,
  Group,
  Text,
  Title,
} from '@mantine/core';
import { IconChevronLeft } from '@tabler/icons-react';
import { NavLink, useNavigate, useParams } from 'react-router-dom';
import { AppLoader } from '../../components/Loader/AppLoader';
import UserList from './List';

const WorkspaceUsersPage = () => {
  const { workspaceId } = useParams();
  const navigate = useNavigate();

  if (!workspaceId) {
    return <AppLoader h="70vh" />;
  }

  return (
    <Box mb={50}>
      <Breadcrumbs separator=">">
        <NavLink to="/" className="bredcrumb-link">
          <Text>Workspaces</Text>
        </NavLink>
        <Text>Users</Text>
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
        <Title order={3}>Users</Title>
      </Group>
      <UserList />
    </Box>
  );
};

export default WorkspaceUsersPage;
