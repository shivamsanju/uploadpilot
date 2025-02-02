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
import { useState } from "react";
import { IconMenu2 } from "@tabler/icons-react";

const AppLayout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [opened, toggle] = useState(true);
  const { colorScheme } = useMantineColorScheme();
  const theme = useMantineTheme();

  const bg = colorScheme === "dark" ? "#141414" : theme.colors.gray[0];

  const appShellBorderColor =
    colorScheme === "dark" ? theme.colors.dark[8] : theme.colors.gray[1];

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
          <NavBar toggle={toggle} />
        </AppShell.Navbar>
        <AppShell.Main bg={bg} m={0}>
          <ScrollArea scrollbarSize={6} h="93vh">
            {children}
          </ScrollArea>
        </AppShell.Main>
      </AppShell>
    </AuthWrapper>
  );
};

export default AppLayout;
