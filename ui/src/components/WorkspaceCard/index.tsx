import {
  Avatar,
  Button,
  Card,
  Group,
  Stack,
  Text,
  ThemeIcon,
  Title,
} from '@mantine/core';
import {
  IconArrowRightDashed,
  IconPinFilled,
  IconServer2,
  IconStopwatch,
  IconTools,
} from '@tabler/icons-react';
import { useNavigate } from 'react-router-dom';
import WorkspaceMenu from './WorkspaceMenu';

type Props = {
  id: string;
  name: string;
  description: string;
  uploads: number;
  totalUploads: number;
  storage: number;
  storageUsed: number;
  processors: number;
  tags: string[];
  pinned?: boolean;
};

export const WorkspaceCard: React.FC<Props> = ({
  id,
  name,
  description,
  uploads,
  totalUploads,
  storage,
  storageUsed,
  processors,
  tags,
  pinned = false,
}) => {
  const navigate = useNavigate();

  return (
    <Card w="100%" p="md" h="auto" radius="md" shadow="sm">
      <Stack gap="md">
        <Group justify="space-between">
          <Group align="center">
            <Avatar radius="xs" size="md">
              {name[0].toUpperCase()}
            </Avatar>
            <Stack gap={3}>
              <Title order={4}>{name}</Title>
              <Group>
                {tags.map(tag => (
                  <Text
                    key={tag}
                    size={'calc(var(--mantine-font-size-xs) * 0.85)'}
                    c="dimmed"
                  >
                    #{tag}
                  </Text>
                ))}
              </Group>
            </Stack>
          </Group>
          <Group align="flex-end">
            {pinned && (
              <ThemeIcon variant="subtle" c="dimmed">
                <IconPinFilled size={18} stroke={2} />
              </ThemeIcon>
            )}

            <WorkspaceMenu id={id} />
          </Group>
        </Group>

        <Text size="xs" c="dimmed" lineClamp={1}>
          {description}
        </Text>

        <Group justify="space-between" align="flex-end">
          <Group wrap="nowrap" gap="md">
            <KPI
              value={uploads.toString()}
              label="Recent Uploads"
              icon={IconStopwatch}
            />
            <KPI
              value={totalUploads.toString()}
              label="Total Uploads"
              icon={IconServer2}
            />
            <KPI
              value={processors.toString()}
              label="Processors"
              icon={IconTools}
            />
          </Group>

          <Button
            variant="filled"
            size="xs"
            onClick={() => navigate(`/workspace/${id}/uploads`)}
            rightSection={<IconArrowRightDashed size={15} />}
          >
            Open
          </Button>
        </Group>
      </Stack>
    </Card>
  );
};

export const KPI = ({
  value,
  label,
  icon: Icon,
}: {
  value: string;
  label: string;
  icon: React.FC<any>;
}) => {
  return (
    <Group gap="0">
      <ThemeIcon c="dimmed" size="xl" variant="subtle">
        <Icon size={25} stroke={1} />
      </ThemeIcon>
      <Stack gap={0}>
        <Text size="xs" fw="700">
          {value}
        </Text>
        <Text
          c="dimmed"
          style={{ fontSize: 'calc(var(--mantine-font-size-xs) * 0.8)' }}
        >
          {label}
        </Text>
      </Stack>
    </Group>
  );
};
