import { ActionIcon, Group, Stack, Text } from '@mantine/core';
import { useViewportSize } from '@mantine/hooks';
import {
  IconBrandGolang,
  IconBrandJavascript,
  IconBrandPython,
  IconBrandReact,
  IconBrandTypescript,
} from '@tabler/icons-react';
import { useEffect, useState } from 'react';
import { useSetBreadcrumbs } from '../../hooks/breadcrumb';
import classes from './getstarted.module.css';
import GoIntegrationPage from './go';
import ReactUploaderPreviewPage from './react';

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
    transformOrigin: 'top left',
  };
};

const frameworks = [
  {
    name: 'React',
    icon: IconBrandReact,
  },
  {
    name: 'Go',
    icon: IconBrandGolang,
  },
  {
    name: 'Python',
    icon: IconBrandPython,
  },
  {
    name: 'JavaScript',
    icon: IconBrandJavascript,
  },
  {
    name: 'TypeScript',
    icon: IconBrandTypescript,
  },
  {
    name: 'React',
    icon: IconBrandReact,
  },
  {
    name: 'Go',
    icon: IconBrandGolang,
  },
  {
    name: 'Python',
    icon: IconBrandPython,
  },
  {
    name: 'JavaScript',
    icon: IconBrandJavascript,
  },
  {
    name: 'TypeScript',
    icon: IconBrandTypescript,
  },
  {
    name: 'React',
    icon: IconBrandReact,
  },
  {
    name: 'Go',
    icon: IconBrandGolang,
  },
  {
    name: 'Python',
    icon: IconBrandPython,
  },
  {
    name: 'JavaScript',
    icon: IconBrandJavascript,
  },
  {
    name: 'TypeScript',
    icon: IconBrandTypescript,
  },
];

export const GetStartedPage = () => {
  const { width } = useViewportSize();
  const setBreadcrumbs = useSetBreadcrumbs();
  const [selectedFramework, setSelectedFramework] = useState('React');

  const s = style(width);

  useEffect(() => {
    setBreadcrumbs([]);
  }, [setBreadcrumbs]);

  return (
    <Stack justify="center" align="center" pt="sm" mb={50}>
      <Text ta="center" fw={700} fz="25px" mb="sm">
        Choose your framework
      </Text>
      <Group justify="center" gap="xl" mb="md" w="50%">
        {frameworks.map(framework => (
          <ActionIcon
            key={framework.name}
            size="60"
            p="sm"
            radius="xl"
            variant={
              selectedFramework === framework.name ? 'filled' : 'outline'
            }
            onClick={() => setSelectedFramework(framework.name)}
            className={
              selectedFramework === framework.name ? classes.selected : ''
            }
          >
            {<framework.icon size={30} />}
          </ActionIcon>
        ))}
      </Group>
      {selectedFramework === 'React' && <ReactUploaderPreviewPage style={s} />}
      {selectedFramework === 'Go' && <GoIntegrationPage style={s} />}
      {selectedFramework === 'Python' && <GoIntegrationPage style={s} />}
      {selectedFramework === 'JavaScript' && <GoIntegrationPage style={s} />}
      {selectedFramework === 'TypeScript' && <GoIntegrationPage style={s} />}
    </Stack>
  );
};
