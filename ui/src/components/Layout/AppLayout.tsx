import {
  ActionIcon,
  AppShell,
  ScrollArea,
  useMantineColorScheme,
  useMantineTheme,
} from '@mantine/core';
import {
  IconAdjustments,
  IconCircles,
  IconDatabase,
  IconGauge,
  IconKey,
  IconMenu2,
  IconRoute,
  IconShoppingCartBolt,
  IconUsers,
} from '@tabler/icons-react';
import { useMemo, useState } from 'react';
import { useParams } from 'react-router-dom';
import AuthWrapper from '../AuthWrapper/AuthWrapper';
import { AdminHeader } from '../Header/Header';
import NavBar from '../Navigation/Navbar';
import WorkspaceSwitcher from '../WorkspaceSwitcher';

const AppLayout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [opened, toggle] = useState(true);
  const { colorScheme } = useMantineColorScheme();
  const theme = useMantineTheme();
  const { workspaceId } = useParams();

  const bg = colorScheme === 'dark' ? '#141414' : theme.colors.gray[0];

  const appShellBorderColor =
    colorScheme === 'dark' ? theme.colors.dark[6] : theme.colors.gray[1];

  const headerNavBg = colorScheme === 'dark' ? '#161616' : '#fff';

  const navItems = useMemo(
    () => [
      {
        label: 'Get Started',
        icon: IconCircles,
        link: `/workspace/${workspaceId}`,
      },
      {
        label: 'Uploads',
        icon: IconDatabase,
        link: `/workspace/${workspaceId}/uploads`,
      },
      {
        label: 'Processors',
        icon: IconRoute,
        link: `/workspace/${workspaceId}/processors`,
      },
      {
        label: 'Configuration',
        icon: IconAdjustments,
        link: `/workspace/${workspaceId}/configuration`,
      },
      {
        label: 'Users',
        icon: IconUsers,
        link: `/workspace/${workspaceId}/users`,
      },
      {
        label: 'API Keys',
        icon: IconKey,
        link: `/workspace/${workspaceId}/apikeys`,
      },
      {
        label: 'Marketplace',
        icon: IconShoppingCartBolt,
        link: `/workspace/${workspaceId}/tools`,
      },
      {
        label: 'Analytics',
        icon: IconGauge,
        link: `/workspace/${workspaceId}/analytics`,
      },
    ],
    [workspaceId],
  );

  return (
    <AuthWrapper>
      <AppShell
        header={{ height: '7vh' }}
        navbar={{
          width: 250,
          breakpoint: 'sm',
          collapsed: { mobile: opened, desktop: !opened },
        }}
        padding="md"
        transitionDuration={500}
        transitionTimingFunction="ease"
      >
        <AppShell.Header
          style={{ borderColor: appShellBorderColor }}
          bg={headerNavBg}
        >
          <AdminHeader
            burger={
              <ActionIcon
                variant="default"
                size="lg"
                onClick={() => toggle(o => !o)}
              >
                <IconMenu2 color="gray" />
              </ActionIcon>
            }
          />
        </AppShell.Header>
        <AppShell.Navbar
          style={{ borderColor: appShellBorderColor }}
          bg={headerNavBg}
        >
          <NavBar
            toggle={toggle}
            items={navItems}
            footer={<WorkspaceSwitcher />}
          />
        </AppShell.Navbar>
        <AppShell.Main bg={bg} m={0} pr={0}>
          <ScrollArea scrollbarSize={6} h="93vh" pr="md">
            {children}
          </ScrollArea>
        </AppShell.Main>
      </AppShell>
    </AuthWrapper>
  );
};

export default AppLayout;
