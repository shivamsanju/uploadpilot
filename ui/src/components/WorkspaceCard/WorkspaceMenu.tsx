import { ActionIcon, Menu, Text } from '@mantine/core';
import {
  IconAdjustments,
  IconArrowRightDashed,
  IconDotsVertical,
  IconTools,
  IconTrash,
} from '@tabler/icons-react';
import { useNavigate } from 'react-router-dom';

const WorkspaceMenu = ({ id }: { id: string }) => {
  const navigate = useNavigate();
  return (
    <Menu
      trigger="click"
      transitionProps={{ transition: 'pop' }}
      position="left-start"
      trapFocus={false}
    >
      <Menu.Target>
        <ActionIcon variant="subtle">
          <IconDotsVertical size={18} stroke={2} />
        </ActionIcon>
      </Menu.Target>

      <Menu.Dropdown>
        <Menu.Item
          leftSection={<IconArrowRightDashed size={16} />}
          onClick={() => navigate(`/workspace/${id}`)}
        >
          <Text size="xs">Open</Text>
        </Menu.Item>
        <Menu.Item
          leftSection={<IconAdjustments size={16} />}
          onClick={() => navigate(`/workspace/${id}/configuration`)}
        >
          <Text size="xs">Configuration</Text>
        </Menu.Item>
        <Menu.Item
          leftSection={<IconTools size={16} />}
          onClick={() => navigate(`/workspace/${id}/processors`)}
        >
          <Text size="xs">Processors</Text>
        </Menu.Item>
        <Menu.Divider />
        <Menu.Item c="red" leftSection={<IconTrash size={16} />}>
          <Text size="xs">Delete</Text>
        </Menu.Item>
      </Menu.Dropdown>
    </Menu>
  );
};

export default WorkspaceMenu;
