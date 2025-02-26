import { Group, Loader } from '@mantine/core';

export const AppLoader = ({ h }: { h?: string }) => (
  <Group p="xl" align="center" justify="center" h={h}>
    <Loader />
  </Group>
);
