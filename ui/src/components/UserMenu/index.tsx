import {
  Avatar,
  Box,
  Group,
  Menu,
  Text,
  Transition,
  UnstyledButton,
} from '@mantine/core';
import { IconChevronRight, IconLogout, IconSwitch } from '@tabler/icons-react';
import { useNavigate } from 'react-router-dom';
import { handleSignout } from '../../apis/auth';
import { useGetUserDetails } from '../../apis/user';

const UserButton = ({ collapsed }: { collapsed: boolean }) => {
  const navigate = useNavigate();
  const { isPending, error, user } = useGetUserDetails();

  if (error) {
    navigate('/auth');
  }

  if (isPending) {
    return <></>;
  }

  return user.email || user.name ? (
    <Menu
      trigger="click"
      transitionProps={{ transition: 'pop' }}
      position="right-end"
      trapFocus={false}
    >
      <Menu.Target>
        <UnstyledButton w="100%">
          <Group justify="space-between" wrap="nowrap">
            <Group gap={7} wrap="nowrap">
              <Avatar
                src={
                  user.avatar ||
                  'https://raw.githubusercontent.com/mantinedev/mantine/master/.demo/avatars/avatar-1.png'
                }
                alt={user.name ? user.name[0] : user.email[0]}
                radius="xl"
                size={30}
              />
              <Transition
                mounted={!collapsed}
                transition="fade"
                duration={200}
                timingFunction="ease"
              >
                {styles =>
                  collapsed ? (
                    <></>
                  ) : (
                    <Box style={{ ...styles }}>
                      <Text size="sm" fw={500} lineClamp={1}>
                        {user.name || 'User'}
                      </Text>

                      <Text c="dimmed" size="xs" lineClamp={1}>
                        {user.email}
                      </Text>
                    </Box>
                  )
                }
              </Transition>
            </Group>

            <IconChevronRight size={12} stroke={1.5} />
          </Group>
        </UnstyledButton>
      </Menu.Target>

      <Menu.Dropdown>
        {/* <Menu.Item leftSection={<IconSun size={16} />} closeMenuOnClick={false}>
          {' '}
          <ThemeSwitcher />
        </Menu.Item> */}
        <Menu.Item
          leftSection={<IconSwitch size={16} />}
          onClick={() => navigate('/tenants')}
        >
          <Text size="xs">Switch Tenant</Text>
        </Menu.Item>
        <Menu.Divider />
        <Menu.Item
          c="red"
          leftSection={<IconLogout size={16} />}
          onClick={handleSignout}
        >
          <Text size="xs">Logout</Text>
        </Menu.Item>
      </Menu.Dropdown>
    </Menu>
  ) : (
    <></>
  );
};

export default UserButton;
