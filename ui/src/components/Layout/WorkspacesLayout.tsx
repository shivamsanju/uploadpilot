import {
    AppShell,
    ScrollArea,
    useMantineColorScheme,
    useMantineTheme,
} from "@mantine/core";
import { AdminHeader } from "../Header/Header";
import AuthWrapper from "../AuthWrapper/AuthWrapper";

const WorkspacesLayout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
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
                padding="md"
                transitionDuration={500}
                transitionTimingFunction="ease"
            >
                <AppShell.Header bg={navbarHeaderBg} style={{ borderColor: appShellBorderColor }}>
                    <AdminHeader />
                </AppShell.Header>
                <AppShell.Main bg={bg} m={0}>
                    <ScrollArea scrollbarSize={6} h="95vh">
                        {children}
                    </ScrollArea>
                </AppShell.Main>
            </AppShell>
        </AuthWrapper>
    );
}

export default WorkspacesLayout