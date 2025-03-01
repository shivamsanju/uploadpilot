import { Group } from '@mantine/core';
import { BreadcrumbsComponent } from '../Breadcrumbs/Breadcrumbs';

export const Header = () => {
  return (
    <Group p={0} m={0} align="center" justify="space-between">
      <BreadcrumbsComponent />
    </Group>
  );
};
