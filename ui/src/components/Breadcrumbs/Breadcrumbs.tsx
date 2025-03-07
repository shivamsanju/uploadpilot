import { Breadcrumbs, Group, Text } from '@mantine/core';
import { IconArrowLeft } from '@tabler/icons-react';
import { useEffect } from 'react';
import { NavLink, useNavigate } from 'react-router-dom';
import { useBreadcrumbs } from '../../context/BreadcrumbContext';

export const BreadcrumbsComponent = () => {
  const navigate = useNavigate();
  const { breadcrumbs } = useBreadcrumbs();

  useEffect(() => {
    document.title =
      'UploadPilot | ' + (breadcrumbs[breadcrumbs.length - 1]?.label || '');
  }, [breadcrumbs]);

  if (breadcrumbs.length === 0) {
    return <></>;
  }

  return (
    <Group mb="lg">
      <IconArrowLeft
        size={18}
        onClick={() => navigate(`${breadcrumbs[breadcrumbs.length - 2].path}`)}
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
