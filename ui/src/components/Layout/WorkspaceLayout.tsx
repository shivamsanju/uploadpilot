import {
  ActionIcon,
  AppShell,
  ScrollArea,
  useMantineColorScheme,
  useMantineTheme,
} from '@mantine/core';
import {
  IconCategory,
  IconCreditCardFilled,
  IconKey,
  IconMenu2,
  IconUser,
} from '@tabler/icons-react';
import { useMemo, useState } from 'react';
import AuthWrapper from '../AuthWrapper/AuthWrapper';
import { AdminHeader } from '../Header/Header';
import NavBar from '../Navigation/Navbar';

const WorkspacesLayout: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
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
        label: 'Workspaces',
        icon: IconCategory,
        link: `/`,
      },
      {
        label: 'Billing',
        icon: IconCreditCardFilled,
        link: `/billing`,
      },
      {
        label: 'API Keys',
        icon: IconKey,
        link: `/api-keys`,
      },
      {
        label: 'Profile',
        icon: IconUser,
        link: `/profile`,
      },
    ],
    [],
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

export default WorkspacesLayout;
