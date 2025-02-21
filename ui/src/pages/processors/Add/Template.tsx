import { Paper, Text, useMantineTheme } from "@mantine/core";
import { IconActivity } from "@tabler/icons-react";
import classes from "./Utils.module.css";

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
      h="180"
      key={template.key}
      p="sm"
      withBorder
      className={`${classes.card} ${selected ? classes.selected : ""}`}
    >
      <IconActivity size={30} stroke={2} color={theme.colors.appcolor[6]} />
      <Text fz="lg" fw={500} className={classes.cardTitle} mt="md">
        {template.label}
      </Text>
      <Text fz="sm" c="dimmed" mt="sm">
        {template.description}
      </Text>
    </Paper>
  );
};
