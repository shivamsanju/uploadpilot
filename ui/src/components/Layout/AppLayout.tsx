import {
  ActionIcon,
  AppShell,
  ScrollArea,
  useMantineColorScheme,
  useMantineTheme,
} from "@mantine/core";
import { AdminHeader } from "../Header/Header";
import NavBar from "../Navigation/Navbar";
import AuthWrapper from "../AuthWrapper/AuthWrapper";
import { useMemo, useState } from "react";
import {
  IconMenu2,
  IconAdjustments,
  IconCircles,
  IconDatabase,
  IconRoute,
  IconGauge,
  IconUsers,
  IconShoppingCartBolt,
} from "@tabler/icons-react";
import { useParams } from "react-router-dom";
import WorkspaceSwitcher from "../WorkspaceSwitcher";

const AppLayout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [opened, toggle] = useState(true);
  const { colorScheme } = useMantineColorScheme();
  const theme = useMantineTheme();
  const { workspaceId } = useParams();

  const bg = colorScheme === "dark" ? "#141414" : theme.colors.gray[0];

  const appShellBorderColor =
    colorScheme === "dark" ? theme.colors.dark[8] : theme.colors.gray[1];

  const navItems = useMemo(
    () => [
      {
        label: "Get Started",
        icon: IconCircles,
        link: `/workspaces/${workspaceId}`,
      },
      {
        label: "Uploads",
        icon: IconDatabase,
        link: `/workspaces/${workspaceId}/uploads`,
      },
      {
        label: "Processors",
        icon: IconRoute,
        link: `/workspaces/${workspaceId}/processors`,
      },
      {
        label: "Configuration",
        icon: IconAdjustments,
        link: `/workspaces/${workspaceId}/configuration`,
      },
      {
        label: "Users",
        icon: IconUsers,
        link: `/workspaces/${workspaceId}/users`,
      },
      {
        label: "Marketplace",
        icon: IconShoppingCartBolt,
        link: `/workspaces/${workspaceId}/tools`,
      },
      {
        label: "Analytics",
        icon: IconGauge,
        link: `/workspaces/${workspaceId}/analytics`,
      },
    ],
    [workspaceId]
  );

  return (
    <AuthWrapper>
      <AppShell
        header={{ height: "7vh" }}
        navbar={{
          width: 250,
          breakpoint: "sm",
          collapsed: { mobile: opened, desktop: !opened },
        }}
        padding="md"
        transitionDuration={500}
        transitionTimingFunction="ease"
      >
        <AppShell.Header style={{ borderColor: appShellBorderColor }}>
          <AdminHeader
            burger={
              <ActionIcon
                variant="default"
                size="lg"
                onClick={() => toggle((o) => !o)}
              >
                <IconMenu2 color="gray" />
              </ActionIcon>
            }
          />
        </AppShell.Header>
        <AppShell.Navbar style={{ borderColor: appShellBorderColor }}>
          <NavBar
            toggle={toggle}
            items={navItems}
            footer={<WorkspaceSwitcher />}
          />
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

export default AppLayout;
