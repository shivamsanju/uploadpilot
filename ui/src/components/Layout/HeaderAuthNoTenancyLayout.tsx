import {
  AppShell,
  ScrollArea,
  useMantineColorScheme,
  useMantineTheme,
} from '@mantine/core';
import SessionManager from '../SessionManager/SessionManager';

const HeaderAuthNoTenancyLayout: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const { colorScheme } = useMantineColorScheme();
  const theme = useMantineTheme();

  const bg = colorScheme === 'dark' ? '#141414' : theme.colors.gray[0];

  return (
    <SessionManager>
      <AppShell
        header={{ height: '7vh' }}
        padding="md"
        transitionDuration={500}
        transitionTimingFunction="ease"
      >
        <AppShell.Main bg={bg} m={0} pr={0}>
          <ScrollArea scrollbarSize={6} h="100vh" pr="md">
            {children}
          </ScrollArea>
        </AppShell.Main>
      </AppShell>
    </SessionManager>
  );
};

export default HeaderAuthNoTenancyLayout;
