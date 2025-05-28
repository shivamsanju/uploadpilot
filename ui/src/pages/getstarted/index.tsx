import { Box, Group, Image, Paper, Stack, Text, Title } from '@mantine/core';
import { useViewportSize } from '@mantine/hooks';
import { IconCircles } from '@tabler/icons-react';
import { useEffect, useState } from 'react';
import GolangIcon from '../../assets/icons/golang.svg';
import JavaIcon from '../../assets/icons/java.svg';
import JavaScriptIcon from '../../assets/icons/javascript.svg';
import NodejsIcon from '../../assets/icons/nodejs.svg';
import PythonIcon from '../../assets/icons/python.svg';
import { useSetBreadcrumbs } from '../../hooks/breadcrumb';
import BrowserIntegrationPage from './browser';
import classes from './getstarted.module.css';
import GoIntegrationPage from './go';
import JavaIntegrationPage from './java';
import NodejsIntegrationPage from './nodejs';
import PythonIntegrationPage from './python';

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
    icon: JavaScriptIcon,
  },
  {
    name: 'Node.js',
    icon: NodejsIcon,
  },
  {
    name: 'Go',
    icon: GolangIcon,
  },
  {
    name: 'Python',
    icon: PythonIcon,
  },
  {
    name: 'Java',
    icon: JavaIcon,
  },
];

export const GetStartedPage = () => {
  const { width } = useViewportSize();
  const setBreadcrumbs = useSetBreadcrumbs();
  const [selectedFramework, setSelectedFramework] = useState('Browser');

  const s = style(width);

  useEffect(() => {
    setBreadcrumbs([
      { label: 'Workspaces', path: '/' },
      { label: 'Get started' },
    ]);
  }, [setBreadcrumbs]);

  return (
    <Box mr="xl">
      <Group mb="xl">
        <IconCircles size={24} />
        <Title order={3}>Get started</Title>
      </Group>
      <Stack justify="center" pt="sm" mb={50} w="100%">
        <Stack mb="md">
          <Group gap="xl" grow>
            {frameworks.map(framework => (
              <Paper
                key={framework.name}
                withBorder
                p="md"
                radius="md"
                className={`${classes.frameworkCard} ${
                  selectedFramework === framework.name ? classes.selected : ''
                }`}
                onClick={() => setSelectedFramework(framework.name)}
              >
                <Group align="center">
                  <Image
                    src={framework.icon}
                    alt={framework.name}
                    h={40}
                    w={40}
                  />
                  <Text c="dimmed" size="sm">
                    {framework.name}
                  </Text>
                </Group>
              </Paper>
            ))}
          </Group>
        </Stack>
        {selectedFramework === 'Browser' && (
          <BrowserIntegrationPage style={s} />
        )}
        {selectedFramework === 'Node.js' && <NodejsIntegrationPage style={s} />}
        {selectedFramework === 'Go' && <GoIntegrationPage style={s} />}
        {selectedFramework === 'Python' && <PythonIntegrationPage style={s} />}
        {selectedFramework === 'Java' && <JavaIntegrationPage style={s} />}
      </Stack>
    </Box>
  );
};
