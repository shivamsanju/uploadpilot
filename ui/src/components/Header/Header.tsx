import { Badge, Box, Group, Text, Title } from "@mantine/core";
import UserMenu from "../UserMenu";
import { useNavigate } from "react-router-dom";
import { useViewportSize } from "@mantine/hooks";

interface Props {
  burger?: React.ReactNode;
}

export function AdminHeader({ burger }: Props) {
  const { width } = useViewportSize();
  const navigate = useNavigate();
  return (
    <Group
      justify="space-between"
      align="center"
      px="md"
      h="100%"
      gap={width > 768 ? "xl" : "md"}
      wrap="nowrap"
    >
      {burger}
      <Box onClick={() => navigate("/")} style={{ cursor: "pointer" }}>
        <Group gap="md" align="center">
          <Title order={4} opacity={0.7}>
            Upload Pilot{" "}
          </Title>
          <Badge
            variant="gradient"
            gradient={{ from: "appcolor", to: "orange" }}
          >
            beta
          </Badge>
        </Group>
        {width > 768 && <Text c="dimmed">Set your uploads on autopilot</Text>}
      </Box>
      <Box style={{ flex: 1 }} />
      <UserMenu />
    </Group>
  );
}
