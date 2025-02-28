import {
  AppShell,
  ScrollArea,
  useMantineColorScheme,
  useMantineTheme,
} from '@mantine/core';
import { AdminHeader } from '../Header/Header';
import SessionManager from '../SessionManager/SessionManager';

const HeaderAuthNoTenancyLayout: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const { colorScheme } = useMantineColorScheme();
  const theme = useMantineTheme();

  const bg = colorScheme === 'dark' ? '#141414' : theme.colors.gray[0];

  const appShellBorderColor =
    colorScheme === 'dark' ? theme.colors.dark[6] : theme.colors.gray[2];

  const headerNavBg = colorScheme === 'dark' ? '#161616' : '#fff';

  return (
    <SessionManager>
      <AppShell
        header={{ height: '7vh' }}
        padding="md"
        transitionDuration={500}
        transitionTimingFunction="ease"
      >
        <AppShell.Header
          style={{ borderColor: appShellBorderColor }}
          bg={headerNavBg}
        >
          <AdminHeader />
        </AppShell.Header>
        <AppShell.Main bg={bg} m={0} pr={0}>
          <ScrollArea scrollbarSize={6} h="93vh" pr="md">
            {children}
          </ScrollArea>
        </AppShell.Main>
      </AppShell>
    </SessionManager>
  );
};

export default HeaderAuthNoTenancyLayout;
