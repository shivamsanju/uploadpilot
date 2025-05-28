import { AppShell, ScrollArea } from '@mantine/core';
import {
  IconAdjustments,
  IconCircles,
  IconGauge,
  IconServer2,
  IconTools,
} from '@tabler/icons-react';
import { useMemo } from 'react';
import { useParams } from 'react-router-dom';
import { useNavbar } from '../../context/NavbarContext';
import { Header } from '../Header';
import NavBar from '../Navigation/Navbar';
import TenancyManager from '../Tenancy/TenancyManager';

const WorkspaceLayout: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const { workspaceId } = useParams();
  const { opened } = useNavbar();

  const navItems = useMemo(
    () => [
      {
        label: 'Get Started',
        icon: IconCircles,
        link: `/workspace/${workspaceId}`,
      },
      {
        label: 'Uploads',
        icon: IconServer2,
        link: `/workspace/${workspaceId}/uploads`,
      },
      {
        label: 'Processors',
        icon: IconTools,
        link: `/workspace/${workspaceId}/processors`,
      },
      {
        label: 'Configuration',
        icon: IconAdjustments,
        link: `/workspace/${workspaceId}/configuration`,
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
        <AppShell.Navbar withBorder bg="#08060f">
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

export default WorkspaceLayout;
