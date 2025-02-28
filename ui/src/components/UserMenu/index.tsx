import { Avatar, Group, Menu, Text, UnstyledButton } from '@mantine/core';
import {
  IconChevronDown,
  IconLogout,
  IconSun,
  IconSwitch,
} from '@tabler/icons-react';
import { useNavigate } from 'react-router-dom';
import { handleSignout } from '../../apis/auth';
import { useGetUserDetails } from '../../apis/user';
import ThemeSwitcher from '../ThemeSwitcher';

const UserButton = () => {
  const navigate = useNavigate();
  const { isPending, error, user } = useGetUserDetails();

  if (error) {
    navigate('/auth');
  }

  if (isPending) {
    return <></>;
  }

  return user.email || user.name ? (
    <Group gap="md" align="center">
      <Menu
        trigger="click"
        transitionProps={{ transition: 'pop' }}
        width={200}
        position="bottom"
        trapFocus={false}
      >
        <Menu.Target>
          <UnstyledButton>
            <Group gap={7}>
              <Avatar
                src={user.avatar}
                alt={user.name ? user.name[0] : user.email[0]}
                radius="xl"
                size={30}
              />
              <Text fw={500} size="sm" lh={1} mr={3} visibleFrom="md">
                {user.name}
              </Text>
              <IconChevronDown size={12} stroke={1.5} />
            </Group>
          </UnstyledButton>
        </Menu.Target>

        <Menu.Dropdown>
          <Menu.Item
            leftSection={<IconSun size={16} />}
            closeMenuOnClick={false}
          >
            {' '}
            <ThemeSwitcher />
          </Menu.Item>
          <Menu.Item
            leftSection={<IconSwitch size={16} />}
            onClick={() => navigate('/tenants')}
          >
            <Text size="sm">Switch Tenant</Text>
          </Menu.Item>
          <Menu.Divider />
          <Menu.Item
            c="red"
            leftSection={<IconLogout size={16} />}
            onClick={handleSignout}
          >
            <Text size="sm">Logout</Text>
          </Menu.Item>
        </Menu.Dropdown>
      </Menu>
    </Group>
  ) : (
    <></>
  );
};

export default UserButton;
