import {
  AppShell,
  ScrollArea,
  useMantineColorScheme,
  useMantineTheme,
} from '@mantine/core';
import {
  IconCreditCardFilled,
  IconDeviceLaptop,
  IconKey,
} from '@tabler/icons-react';
import { useState } from 'react';
import { Header } from '../Header';
import NavBar from '../Navigation/Navbar';
import SessionManager from '../SessionManager/SessionManager';
import TenancyManager from '../Tenancy/TenancyManager';

const navItems = [
  {
    label: 'Workspaces',
    icon: IconDeviceLaptop,
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
];

const HomeLayout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [opened, toggle] = useState(true);
  const { colorScheme } = useMantineColorScheme();
  const theme = useMantineTheme();

  const bg = colorScheme === 'dark' ? '#121212' : theme.colors.gray[0];
  const appShellBorderColor =
    colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[2];

  return (
    <SessionManager>
      <TenancyManager>
        <AppShell
          navbar={{
            width: opened ? 250 : 75,
            breakpoint: 'sm',
            collapsed: { mobile: opened },
          }}
          padding="md"
          transitionDuration={500}
          transitionTimingFunction="ease"
        >
          <AppShell.Navbar
            withBorder
            style={{
              borderColor: appShellBorderColor,
              transition: 'width 0.2s ease',
            }}
            bg={bg}
          >
            <NavBar collapsed={!opened} toggle={toggle} items={navItems} />
          </AppShell.Navbar>
          <AppShell.Main bg={bg} m={0} pr={0}>
            <ScrollArea scrollbarSize={6} h="100vh" px="md">
              <Header />
              {children}
            </ScrollArea>
          </AppShell.Main>
        </AppShell>
      </TenancyManager>
    </SessionManager>
  );
};

export default HomeLayout;
