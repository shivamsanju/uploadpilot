import { Avatar, Group, Menu, UnstyledButton, Text } from '@mantine/core';
import { IconUser, IconLogout, IconChevronDown, IconSun } from '@tabler/icons-react';
import { useNavigate } from "react-router-dom";
import { useEffect, useState } from 'react';
import ThemeSwitcher from '../ThemeSwitcher';

type User = {
    email: string;
    image: string;
}

const UserButton = () => {
    const navigate = useNavigate();
    const [user, setUser] = useState<User>({} as User);

    const handleSignOut = async () => {
        sessionStorage.removeItem("usermetadata");
    }
    const getUserInfo = async () => {
        return {
            email: "john.doe@example.com",
            image: "https://avatar.iran.liara.run/public/33"
        }
    }

    const handleProfileClick = () => {
        navigate("/profile");
    }

    useEffect(() => {
        getUserInfo()
            .then((data: any) => {
                const u = {
                    email: data.email,
                    image: data.image
                }
                setUser(u)
            })
            .catch(e => console.log(e));
    }, [])

    return user.email ? (
        <Menu
            width={200}
            position="bottom"
            trapFocus={false}
        >
            <Menu.Target>
                <UnstyledButton>
                    <Group gap={7}>
                        <Avatar src={user.image} alt={user.email} radius="xl" size={25} />
                        <Text fw={500} size="sm" lh={1} mr={3}>
                            {user.email.split("@")[0]}
                        </Text>
                        <IconChevronDown size={12} stroke={1.5} />
                    </Group>
                </UnstyledButton>
            </Menu.Target>

            <Menu.Dropdown>
                <Menu.Item leftSection={<IconUser size={16} />} onClick={handleProfileClick}>Profile</Menu.Item>
                <Menu.Item leftSection={<IconSun size={16} />}> <ThemeSwitcher /></Menu.Item>
                <Menu.Divider />
                <Menu.Item leftSection={<IconLogout size={16} />} onClick={handleSignOut}>Sign Out</Menu.Item>
            </Menu.Dropdown>

        </Menu >
    ) : <></>;
};

export default UserButton;