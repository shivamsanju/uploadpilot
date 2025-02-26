import {
  ActionIcon,
  AppShell,
  ScrollArea,
  useMantineColorScheme,
  useMantineTheme,
} from '@mantine/core';
import {
  IconActivity,
  IconBolt,
  IconMenu2,
  IconSettings,
} from '@tabler/icons-react';
import { useMemo, useState } from 'react';
import { useParams } from 'react-router-dom';
import AuthWrapper from '../AuthWrapper/AuthWrapper';
import { AdminHeader } from '../Header/Header';
import NavBar from '../Navigation/Navbar';

const ProcessorLayout: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const { workspaceId, processorId } = useParams();
  const [opened, toggle] = useState(true);
  const { colorScheme } = useMantineColorScheme();
  const theme = useMantineTheme();

  const bg = colorScheme === 'dark' ? '#141414' : theme.colors.gray[0];

  const appShellBorderColor =
    colorScheme === 'dark' ? theme.colors.dark[6] : theme.colors.gray[1];

  const headerNavBg = colorScheme === 'dark' ? '#161616' : '#fff';

  const navItems = useMemo(
    () => [
      {
        label: 'Workflow',
        icon: IconActivity,
        link: `/workspace/${workspaceId}/processors/${processorId}/workflow`,
      },
      {
        label: 'Runs',
        icon: IconBolt,
        link: `/workspace/${workspaceId}/processors/${processorId}/runs`,
      },
      {
        label: 'Settings',
        icon: IconSettings,
        link: `/workspace/${workspaceId}/processors/${processorId}/settings`,
      },
    ],
    [workspaceId, processorId],
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
          <NavBar toggle={toggle} items={navItems} />
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

export default ProcessorLayout;
