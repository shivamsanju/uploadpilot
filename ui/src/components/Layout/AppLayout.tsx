import {
    AppShell,
    Burger,
    ScrollArea,
    useMantineColorScheme,
    useMantineTheme,
} from "@mantine/core";
import { AdminHeader } from "../Header/Header";
import NavBar from "../Navigation/Navbar";
import AuthWrapper from "../AuthWrapper/AuthWrapper";
import { useState } from "react";

const AppLayout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [opened, toggle] = useState(true);
    const { colorScheme } = useMantineColorScheme();
    const theme = useMantineTheme();

    const bg =
        colorScheme === "dark" ? "#0A0A0A" : theme.colors.gray[0];

    const navbarHeaderBg =
        colorScheme === "dark" ? theme.colors.dark[9] : "";

    const appShellBorderColor =
        colorScheme === "dark" ? theme.colors.dark[8] : theme.colors.gray[1];

    return (
        <AuthWrapper>
            <AppShell
                header={{ height: "5vh" }}
                navbar={{
                    width: 220,
                    breakpoint: 'sm',
                    collapsed: { mobile: opened, desktop: !opened },
                }}
                padding="md"
                transitionDuration={500}
                transitionTimingFunction="ease"
            >
                <AppShell.Header bg={navbarHeaderBg} style={{ borderColor: appShellBorderColor }}>
                    <AdminHeader
                        burger={
                            <Burger
                                opened={false}
                                onClick={() => toggle((o) => !o)}
                            />
                        }
                    />
                </AppShell.Header>
                <AppShell.Navbar bg={navbarHeaderBg} style={{ borderColor: appShellBorderColor }}>
                    <NavBar toggle={toggle} />
                </AppShell.Navbar>
                <AppShell.Main bg={bg} m={0}>
                    <ScrollArea scrollbarSize={6} h="95vh">
                        {children}
                    </ScrollArea>
                </AppShell.Main>
            </AppShell>
        </AuthWrapper>
    );
}

export default AppLayout