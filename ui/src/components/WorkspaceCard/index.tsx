import { ActionIcon, Badge, Group, Paper, Stack, Title } from '@mantine/core';
import { IconArrowRightDashed } from '@tabler/icons-react';
import { useNavigate } from 'react-router-dom';

type Props = {
  id: string;
  name: string;
  description: string;
  uploads: number;
  storage: number;
  tags: string[];
};

export const WorkspaceCard: React.FC<Props> = ({
  id,
  name,
  description,
  uploads,
  storage,
  tags,
}) => {
  const navigate = useNavigate();

  return (
    <Paper withBorder w="100%" p="xs" h="80">
      <Stack gap="xs" justify="space-between">
        <Group justify="space-between">
          <Title order={4}>{name}</Title>
          <ActionIcon
            variant="subtle"
            onClick={() => navigate(`/workspace/${id}`)}
          >
            <IconArrowRightDashed size={18} stroke={2} />
          </ActionIcon>
        </Group>
        <Group gap="xs">
          <Badge size="xs" c="dimmed">
            {uploads} uploads
          </Badge>
          <Badge size="xs" c="dimmed">
            {storage} GB
          </Badge>
        </Group>
      </Stack>
    </Paper>
  );
};
