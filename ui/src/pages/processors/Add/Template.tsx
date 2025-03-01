import {
  Box,
  Group,
  Paper,
  Text,
  ThemeIcon,
  useMantineTheme,
} from '@mantine/core';
import { IconFileStack } from '@tabler/icons-react';
import classes from './Utils.module.css';

type Props = {
  template: {
    label: string;
    key: string;
    description: string;
  };
  selected?: boolean;
};
export const Template: React.FC<Props> = ({ template, selected }) => {
  const theme = useMantineTheme();

  return (
    <Paper
      h="150"
      key={template.key}
      p="sm"
      withBorder
      className={`${classes.card} ${selected ? classes.selected : ''}`}
    >
      <Box className={classes.cardTitle}>
        <Group p={0} m={0} align="center" wrap="nowrap">
          <ThemeIcon size={30} radius="xl" variant="light">
            <IconFileStack
              size={20}
              stroke={2}
              color={theme.colors.appcolor[6]}
            />
          </ThemeIcon>
          <Text fz="lg" fw={500}>
            {template.label}
          </Text>
        </Group>
      </Box>
      <Text fz="sm" c="dimmed" mt="sm">
        {template.description}
      </Text>
    </Paper>
  );
};
