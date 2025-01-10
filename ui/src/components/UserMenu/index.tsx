import { Avatar, Group, Menu, UnstyledButton, Text } from '@mantine/core';
import { IconUser, IconLogout, IconChevronDown, IconSun } from '@tabler/icons-react';
import { useNavigate } from "react-router-dom";
import { useEffect, useState } from 'react';
import ThemeSwitcher from '../ThemeSwitcher';
import axiosInstance from '../../utils/axios';
import { notifications } from '@mantine/notifications';

type User = {
    email: string;
    firstName: string;
    lastName: string;
    image: string;
}

const UserButton = () => {
    const navigate = useNavigate();
    const [user, setUser] = useState<User>({} as User);

    const handleSignOut = async () => {
        localStorage.removeItem("token");
        localStorage.removeItem("refreshToken");
        navigate("/auth");
    }

    const getUserInfo = async () => {
        try {
            const response = await axiosInstance.get("/users/me")
            if (response.status === 200) {
                return response.data;
            }
            notifications.show({
                title: 'Error',
                message: response.data.message,
                color: 'red',
            })
        } catch (error: any) {
            console.log(error);
            notifications.show({
                title: 'Error',
                message: error.toString(),
                color: 'red',
            })
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
                    firstName: data.firstName,
                    lastName: data.lastName,
                    image: "https://avatar.iran.liara.run/public/33"
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
                            {user.firstName + " " + user.lastName}
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