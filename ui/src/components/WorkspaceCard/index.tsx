import { ActionIcon, Group, Paper, Stack, Text, Title } from '@mantine/core';
import { IconChevronRightPipe, IconDeviceLaptop } from '@tabler/icons-react';
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
    <Paper withBorder w={{ base: '100%', md: '400px' }} h={200}>
      <Stack h="100%" justify="space-between" gap={0}>
        <Stack
          gap="xs"
          p="sm"
          pb={3}
          justify="space-between"
          style={{ flex: 1 }}
        >
          <Group gap="xs">
            <IconDeviceLaptop size={35} />
            <Title order={3}>{name}</Title>
          </Group>
          <Text size="xs" c="dimmed" lineClamp={3}>
            {description}
          </Text>
          <Group gap="xs">
            {tags.map(tag => (
              <Text key={tag} size="xs" c="dimmed">
                #{tag}
              </Text>
            ))}
          </Group>
        </Stack>

        <Group
          justify="space-between"
          align="center"
          style={{
            borderTop:
              '1px solid light-dark(var(--mantine-color-gray-4), var(--mantine-color-dark-8))',
          }}
          p="sm"
        >
          <Group gap="xs">
            <Text size="sm" c="dimmed">
              Uploads: <b>{uploads}</b>
            </Text>
            <Text size="sm" c="dimmed">
              Storage: <b>{storage} GB</b>
            </Text>
          </Group>
          <ActionIcon
            variant="outline"
            onClick={() => navigate(`/workspace/${id}`)}
          >
            <IconChevronRightPipe size={20} />
          </ActionIcon>
        </Group>
      </Stack>
    </Paper>
  );
};
