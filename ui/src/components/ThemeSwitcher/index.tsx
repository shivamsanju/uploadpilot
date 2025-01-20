import { useMantineTheme, rem, useMantineColorScheme, Menu, Text } from '@mantine/core';
import { IconSun, IconMoonStars, IconCheck } from '@tabler/icons-react';
import { useCallback } from 'react';


const ThemeSwitcher = () => {
    const theme = useMantineTheme();
    const { colorScheme, setColorScheme } = useMantineColorScheme();

    const changeTheme = useCallback((scheme: "dark" | "light" | "auto") => {
        setColorScheme(scheme);
    }, [setColorScheme]);

    const sunIcon = (
        <IconSun
            style={{ width: rem(16), height: rem(16) }}
            stroke={2.5}
            color={theme.colors.yellow[4]}
        />
    );

    const moonIcon = (
        <IconMoonStars
            style={{ width: rem(16), height: rem(16) }}
            stroke={2.5}
            color={theme.colors.grape[6]}
        />
    );

    const checkIcon = (
        <IconCheck
            style={{ width: rem(16), height: rem(16) }}
            stroke={2.5}
            color={theme.colors.green[4]}
        />
    )

    return (
        <Menu
            width={200}
            position="left"
            trigger='click'
            closeOnItemClick={false}
        >
            <Menu.Target>
                <Text size='sm'>Theme</Text>
            </Menu.Target>

            <Menu.Dropdown>
                <Menu.Item
                    leftSection={sunIcon}
                    onClick={() => changeTheme("light")}
                    rightSection={colorScheme === "light" && checkIcon}
                >
                    Light
                </Menu.Item>
                <Menu.Item
                    leftSection={moonIcon}
                    onClick={() => changeTheme("dark")}
                    rightSection={colorScheme === "dark" && checkIcon}
                >
                    Dark
                </Menu.Item>
            </Menu.Dropdown>

        </Menu >
    );
}

export default ThemeSwitcher