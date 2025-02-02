import {
  Avatar,
  Group,
  Menu,
  UnstyledButton,
  Text,
  Button,
} from "@mantine/core";
import {
  IconLogout,
  IconChevronDown,
  IconSun,
  IconStopwatch,
} from "@tabler/icons-react";
import { useNavigate } from "react-router-dom";
import ThemeSwitcher from "../ThemeSwitcher";
import { useGetSession } from "../../apis/user";
import { getApiDomain } from "../../utils/config";

const UserButton = () => {
  const navigate = useNavigate();
  const { isPending, error, session } = useGetSession();

  const handleSignOut = async () => {
    localStorage.removeItem("uploadpilottoken");
    window.location.href = getApiDomain() + `/auth/logout`;
  };

  if (error) {
    navigate("/auth");
  }

  if (isPending) {
    return <></>;
  }

  return session.email || session.name ? (
    <Group gap="md" align="center">
      <Button
        variant="light"
        radius="sm"
        mr="40"
        leftSection={<IconStopwatch />}
        visibleFrom="md"
      >
        Trial expires in {Math.round((session?.trialExpiresIn || 0) / 24)} days
      </Button>
      <Menu
        trigger="click"
        transitionProps={{ transition: "pop" }}
        width={200}
        position="bottom"
        trapFocus={false}
      >
        <Menu.Target>
          <UnstyledButton>
            <Group gap={7}>
              <Avatar
                src={session.avatarUrl}
                alt={session.name[0]}
                radius="xl"
                size={30}
              />
              <Text fw={500} size="sm" lh={1} mr={3} visibleFrom="md">
                {session.name || session.firstName + " " + session.lastName}
              </Text>
              <IconChevronDown size={12} stroke={1.5} />
            </Group>
          </UnstyledButton>
        </Menu.Target>

        <Menu.Dropdown>
          <Menu.Item
            leftSection={<IconSun size={16} />}
            closeMenuOnClick={false}
          >
            {" "}
            <ThemeSwitcher />
          </Menu.Item>
          <Menu.Divider />
          <Menu.Item
            c="red"
            leftSection={<IconLogout size={16} />}
            onClick={handleSignOut}
          >
            <Text size="sm">Logout</Text>
          </Menu.Item>
        </Menu.Dropdown>
      </Menu>
    </Group>
  ) : (
    <></>
  );
};

export default UserButton;
