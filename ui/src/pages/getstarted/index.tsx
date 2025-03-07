import { ActionIcon, Box, Group, Stack, Text } from '@mantine/core';
import { useViewportSize } from '@mantine/hooks';
import {
  IconBrandGolang,
  IconBrandNodejs,
  IconWorldWww,
} from '@tabler/icons-react';
import { useEffect, useState } from 'react';
import { useSetBreadcrumbs } from '../../hooks/breadcrumb';
import BrowserIntegrationPage from './browser';
import classes from './getstarted.module.css';
import GoIntegrationPage from './go';
import NodejsIntegrationPage from './nodejs';

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
    name: 'Browser',
    icon: IconWorldWww,
  },
  {
    name: 'Node.js',
    icon: IconBrandNodejs,
  },
  {
    name: 'Go',
    icon: IconBrandGolang,
  },
];

export const GetStartedPage = () => {
  const { width } = useViewportSize();
  const setBreadcrumbs = useSetBreadcrumbs();
  const [selectedFramework, setSelectedFramework] = useState('Browser');

  const s = style(width);

  useEffect(() => {
    setBreadcrumbs([]);
  }, [setBreadcrumbs]);

  return (
    <Stack justify="center" pt="sm" mb={50}>
      <Stack justify="center" align="center" mb="md">
        <Box w={{ sm: '100vw', md: '70vw', lg: '60vw' }}>
          <Text size="sm" mb="sm">
            Select your platform to get started...
          </Text>
          <Group gap="xl">
            {frameworks.map(framework => (
              <ActionIcon
                key={framework.name}
                title={framework.name}
                size="80"
                p="lg"
                variant={
                  selectedFramework === framework.name ? 'light' : 'outline'
                }
                onClick={() => setSelectedFramework(framework.name)}
                className={
                  selectedFramework === framework.name ? classes.selected : ''
                }
              >
                <Stack align="center" gap="5">
                  {<framework.icon size={30} />}
                  <Text fw={500}>{framework.name}</Text>
                </Stack>
              </ActionIcon>
            ))}
          </Group>
        </Box>
      </Stack>
      {selectedFramework === 'Browser' && <BrowserIntegrationPage style={s} />}
      {selectedFramework === 'Go' && <GoIntegrationPage style={s} />}
      {selectedFramework === 'Node.js' && <NodejsIntegrationPage style={s} />}
    </Stack>
  );
};
