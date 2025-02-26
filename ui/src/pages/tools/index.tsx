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
import { ToolsGrid } from './Grid';

const ToolsPage = () => {
  const { workspaceId } = useParams();
  const navigate = useNavigate();

  if (!workspaceId) {
    return <AppLoader h="70vh" />;
  }

  return (
    <Box mb={50} mr="sm">
      <Breadcrumbs separator=">">
        <NavLink to="/" className="bredcrumb-link">
          <Text>Workspaces</Text>
        </NavLink>
        <Text>Marketplace</Text>
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
        <Title order={3}>Marketplace</Title>
      </Group>
      <ToolsGrid />
    </Box>
  );
};

export default ToolsPage;
