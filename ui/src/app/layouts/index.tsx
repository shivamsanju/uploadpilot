import {
    AppShell,
    Burger,
    useMantineColorScheme,
    useMantineTheme,
} from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { AdminHeader } from "../../components/Header/Header";
import classes from "./Layout.module.css";

const HomeLayout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [opened, { toggle }] = useDisclosure();
    const { colorScheme } = useMantineColorScheme();
    const theme = useMantineTheme();

    const bg =
        colorScheme === "dark" ? theme.colors.dark[7] : theme.colors.gray[0];

    return (
        <AppShell
            header={{ height: "7vh" }}
            padding="md"
            transitionDuration={500}
            transitionTimingFunction="ease"
        >
            <AppShell.Header>
                <AdminHeader
                    burger={
                        <Burger
                            opened={opened}
                            onClick={toggle}
                            hiddenFrom="sm"
                            size="sm"
                            mr="xl"
                        />
                    }
                />
            </AppShell.Header>
            <AppShell.Main bg={bg} className={classes.main}>
                {children}
            </AppShell.Main>
        </AppShell>
    );
}

export default HomeLayout