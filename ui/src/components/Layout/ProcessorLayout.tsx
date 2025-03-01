import { AppShell, ScrollArea } from '@mantine/core';
import { IconBolt, IconFileStack, IconSettings } from '@tabler/icons-react';
import { useMemo, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Header } from '../Header';
import NavBar from '../Navigation/Navbar';
import SessionManager from '../SessionManager/SessionManager';
import TenancyManager from '../Tenancy/TenancyManager';

const ProcessorLayout: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [opened, toggle] = useState(true);
  const { workspaceId, processorId } = useParams();

  const navItems = useMemo(
    () => [
      {
        label: 'Workflow',
        icon: IconFileStack,
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
          <AppShell.Navbar withBorder>
            <NavBar collapsed={!opened} toggle={toggle} items={navItems} />
          </AppShell.Navbar>
          <AppShell.Main m={0} pr={0}>
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

export default ProcessorLayout;
