import {
    AppShell,
    Burger,
    useMantineColorScheme,
    useMantineTheme,
} from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { AdminHeader } from "../Header/Header";
import NavBar from "../Navigation/Navbar";
import AuthWrapper from "../AuthWrapper/AuthWrapper";

const AppLayout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [opened, { toggle }] = useDisclosure(true);
    const { colorScheme } = useMantineColorScheme();
    const theme = useMantineTheme();

    const bg =
        colorScheme === "dark" ? theme.colors.dark[7] : theme.colors.gray[0];

    const navbarHeaderBg =
        colorScheme === "dark" ? theme.colors.dark[6] : "";

    return (
        <AuthWrapper>
            <AppShell
                header={{ height: "5vh" }}
                navbar={{
                    width: 220,
                    breakpoint: 'sm',
                    collapsed: { mobile: !opened, desktop: !opened },
                }}
                padding="md"
                transitionDuration={500}
                transitionTimingFunction="ease"
            >
                <AppShell.Header bg={navbarHeaderBg}>
                    <AdminHeader
                        burger={
                            <Burger
                                opened={false}
                                onClick={toggle}
                                size="xs"
                            />
                        }
                    />
                </AppShell.Header>
                <AppShell.Navbar bg={navbarHeaderBg} withBorder={false}>
                    <NavBar />
                </AppShell.Navbar>
                <AppShell.Main bg={bg}>
                    {children}
                </AppShell.Main>
            </AppShell>
        </AuthWrapper>
    );
}

export default AppLayout