import { Avatar, Group, Menu, UnstyledButton, Text } from '@mantine/core';
import { IconUser, IconLogout, IconChevronDown, IconSun } from '@tabler/icons-react';
import { signOut } from "supertokens-auth-react/recipe/session";
import { useNavigate } from "react-router-dom";
import axios from 'axios';
import { getApiDomain } from '../../config';
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
        await signOut();
        navigate("/auth");
    }
    const getUserInfo = async () => {
        let response = await axios.get(getApiDomain() + "/auth/userinfo");
        return response.data;
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
                sessionStorage.setItem("usermetadata", JSON.stringify(u))
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