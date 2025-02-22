import { Group, Paper, Stack, Title } from "@mantine/core";
import { useViewportSize } from "@mantine/hooks";
import {
  IconBrandGolang,
  IconBrandJavascript,
  IconBrandPython,
  IconBrandReact,
  IconBrandTypescript,
} from "@tabler/icons-react";
import { FrameworkCard } from "./FrameworkCard";
import { useState } from "react";
import ReactUploaderPreviewPage from "./react";
import GoIntegrationPage from "./go";
import classes from "./getstarted.module.css";

const style = (width: number) => {
  if (width > 768) {
    return {};
  }

  let scale = 1;
  if (width < 768 && width > 700) {
    scale = width / 768;
  } else if (width < 700 && width > 500) {
    scale = (width / 768) * 1.1;
  } else {
    scale = (width / 768) * 1.35;
  }

  return {
    transform: `scale(${scale})`,
    transformOrigin: "top left",
  };
};

const frameworks = [
  {
    name: "React",
    icon: IconBrandReact,
  },
  {
    name: "Go",
    icon: IconBrandGolang,
  },
  {
    name: "Python",
    icon: IconBrandPython,
  },
  {
    name: "JavaScript",
    icon: IconBrandJavascript,
  },
  {
    name: "TypeScript",
    icon: IconBrandTypescript,
  },
];

export const GetStartedPage = () => {
  const [selectedFramework, setSelectedFramework] = useState("React");
  const { width } = useViewportSize();
  const s = style(width);

  return (
    <Stack justify="center" align="center" pt="sm" mb={50}>
      <Title order={4} c="dimmed">
        Choose your framework
      </Title>
      <Group justify="center" gap="xl" mb="md">
        {frameworks.map((framework) => (
          <Paper
            className={`${classes.frameworkCard} ${
              selectedFramework === framework.name ? classes.selected : ""
            }`}
            key={framework.name}
            p="md"
            withBorder
            onClick={() => setSelectedFramework(framework.name)}
          >
            <FrameworkCard
              framework={framework.name}
              Icon={framework.icon}
              h="100%"
              w="100px"
            />
          </Paper>
        ))}
      </Group>
      {selectedFramework === "React" && <ReactUploaderPreviewPage style={s} />}
      {selectedFramework === "Go" && <GoIntegrationPage style={s} />}
      {selectedFramework === "Python" && <GoIntegrationPage style={s} />}
      {selectedFramework === "JavaScript" && <GoIntegrationPage style={s} />}
      {selectedFramework === "TypeScript" && <GoIntegrationPage style={s} />}
    </Stack>
  );
};
