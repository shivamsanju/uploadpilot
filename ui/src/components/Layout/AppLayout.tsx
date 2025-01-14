import {
    AppShell,
    Burger,
    useMantineColorScheme,
    useMantineTheme,
} from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { AdminHeader } from "../Header/Header";
import NavBar from "../Navigation/Navbar";
import { useEffect } from "react";

const AppLayout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [opened, { toggle }] = useDisclosure(true);
    const { colorScheme } = useMantineColorScheme();
    const theme = useMantineTheme();

    const bg =
        colorScheme === "dark" ? theme.colors.dark[7] : theme.colors.gray[0];

    const navbarHeaderBg =
        colorScheme === "dark" ? theme.colors.dark[6] : "";

    useEffect(() => {
        const token = localStorage.getItem('token')
        if (!token) {
            window.location.href = '/auth'
        }
    }, [])


    return (
        <AppShell
            header={{ height: "7vh" }}
            navbar={{
                width: 250,
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
                            size="sm"
                            mr="sm"
                        />
                    }
                />
            </AppShell.Header>
            <AppShell.Navbar bg={navbarHeaderBg}>
                <NavBar />
            </AppShell.Navbar>
            <AppShell.Main bg={bg}>
                {children}
            </AppShell.Main>
        </AppShell>
    );
}

export default AppLayout