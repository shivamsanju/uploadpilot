import { Box, Group, Title } from '@mantine/core';
import ApiKeyList from './List';

const ApiKeyPage = () => {
  return (
    <Box>
      <Group align="center" gap="xs" h="10%" mb="md">
        <Title order={3}>API Keys</Title>
      </Group>
      <ApiKeyList />
    </Box>
  );
};

export default ApiKeyPage;
