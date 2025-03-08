import { AppShell, ScrollArea } from '@mantine/core';
import { IconBolt, IconFileStack, IconSettings } from '@tabler/icons-react';
import { useMemo } from 'react';
import { useParams } from 'react-router-dom';
import { useNavbar } from '../../context/NavbarContext';
import { Header } from '../Header';
import NavBar from '../Navigation/Navbar';
import TenancyManager from '../Tenancy/TenancyManager';

const ProcessorLayout: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const { workspaceId, processorId } = useParams();
  const { opened } = useNavbar();

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
    <TenancyManager>
      <AppShell
        navbar={{
          width: opened ? 250 : 75,
          breakpoint: 'sm',
          collapsed: { mobile: !opened },
        }}
        padding="md"
        transitionDuration={500}
        transitionTimingFunction="ease"
      >
        <AppShell.Navbar withBorder>
          <NavBar items={navItems} />
        </AppShell.Navbar>
        <AppShell.Main m={0} pr={0}>
          <ScrollArea scrollbarSize={6} h="100vh" px="md">
            <Header />
            {children}
          </ScrollArea>
        </AppShell.Main>
      </AppShell>
    </TenancyManager>
  );
};

export default ProcessorLayout;
