import { Stack, Text } from '@mantine/core';
import { TablerIcon } from '@tabler/icons-react';

type Props = {
  framework: string;
  Icon: TablerIcon;
  h: string;
  w: string;
};
export const FrameworkCard: React.FC<Props> = ({ framework, Icon, h, w }) => {
  return (
    <Stack align="center" justify="center" h={h} w={w}>
      <Icon size={50} opacity={0.7} />
      <Text fw={500} c="dimmed" fz="sm">
        {framework}
      </Text>
    </Stack>
  );
};
