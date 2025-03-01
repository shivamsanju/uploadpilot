import { Breadcrumbs, Group, Text } from '@mantine/core';
import { IconArrowLeft } from '@tabler/icons-react';
import { NavLink, useNavigate } from 'react-router-dom';
import { useBreadcrumbs } from '../../context/BreadcrumbContext';

export const BreadcrumbsComponent = () => {
  const navigate = useNavigate();
  const { breadcrumbs } = useBreadcrumbs();

  if (breadcrumbs.length === 0) {
    return <></>;
  }

  return (
    <Group mb="lg">
      <IconArrowLeft
        size={18}
        onClick={() => navigate(-1)}
        style={{ cursor: 'pointer' }}
      />
      <Breadcrumbs separator="/">
        {breadcrumbs.map((breadcrumb, index) =>
          breadcrumb.path ? (
            <NavLink
              key={index}
              to={breadcrumb.path}
              className="breadcrumb-link"
            >
              <Text>{breadcrumb.label}</Text>
            </NavLink>
          ) : (
            <Text key={index}>{breadcrumb.label}</Text>
          ),
        )}
      </Breadcrumbs>
    </Group>
  );
};
