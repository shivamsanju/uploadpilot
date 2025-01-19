import { Avatar, Group, Menu, UnstyledButton, Text } from '@mantine/core';
import { IconUser, IconLogout, IconChevronDown, IconSun } from '@tabler/icons-react';
import { useNavigate } from "react-router-dom";
import ThemeSwitcher from '../ThemeSwitcher';
import { useGetSession } from '../../apis/user';
import { getApiDomain } from '../../utils/config';


const UserButton = () => {
    const navigate = useNavigate();
    const { isPending, error, session } = useGetSession();


    const handleSignOut = async () => {
        localStorage.removeItem("uploadpilottoken");
        window.location.href = getApiDomain() + `/auth/logout`
    }


    const handleProfileClick = () => {
        navigate("/profile");
    }

    if (error) {
        navigate("/auth");
    }

    if (isPending) {
        return <></>;
    }

    return (session.email || session.name) ? (
        <Menu
            width={200}
            position="bottom"
            trapFocus={false}
        >
            <Menu.Target>
                <UnstyledButton>
                    <Group gap={7}>
                        <Avatar src={session.avatarUrl} alt={session.name[0]} radius="xl" size={25} />
                        <Text fw={500} size="sm" lh={1} mr={3}>
                            {session.name || session.firstName + " " + session.lastName}
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