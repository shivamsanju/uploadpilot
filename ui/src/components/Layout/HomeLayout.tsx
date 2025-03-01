import { AppShell, ScrollArea } from '@mantine/core';
import {
  IconCreditCardFilled,
  IconDeviceLaptop,
  IconKey,
} from '@tabler/icons-react';
import { useNavbar } from '../../context/NavbarContext';
import { Header } from '../Header';
import NavBar from '../Navigation/Navbar';
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
  const { opened } = useNavbar();

  return (
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

export default HomeLayout;
