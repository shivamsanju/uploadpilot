import { AppShell, ScrollArea } from '@mantine/core';

const HeaderAuthNoTenancyLayout: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  return (
    <AppShell
      header={{ height: '7vh' }}
      padding="md"
      transitionDuration={500}
      transitionTimingFunction="ease"
    >
      <AppShell.Main m={0} pr={0}>
        <ScrollArea scrollbarSize={6} h="100vh" pr="md">
          {children}
        </ScrollArea>
      </AppShell.Main>
    </AppShell>
  );
};

export default HeaderAuthNoTenancyLayout;
