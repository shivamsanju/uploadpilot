import { Group, Paper, Text, ThemeIcon } from '@mantine/core';
import type { FC } from 'react';
type Props = {
  title: string;
  value: string;
  Icon: React.FC<any>;
  iconColor?: string;
};

export const SimpleKPICard: FC<Props> = ({ title, value, Icon, iconColor }) => {
  return (
    <Paper withBorder radius="md" p="xs" key={title}>
      <Group>
        <ThemeIcon size="xl" variant="light" color={iconColor}>
          <Icon size={30} stroke={1.5} />
        </ThemeIcon>
        <div>
          <Text c="dimmed" size="xs" tt="uppercase" fw={700}>
            {title}
          </Text>
          <Text fw={700} size="lg">
            {value}
          </Text>
        </div>
      </Group>
    </Paper>
  );
};
