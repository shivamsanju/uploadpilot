import { Breadcrumbs, Group, Text } from '@mantine/core';
import { IconArrowLeft, IconMenu4 } from '@tabler/icons-react';
import { useEffect } from 'react';
import { NavLink, useNavigate } from 'react-router-dom';
import { useBreadcrumbs } from '../../context/BreadcrumbContext';
import { useNavbar } from '../../context/NavbarContext';

export const BreadcrumbsComponent = () => {
  const navigate = useNavigate();
  const { breadcrumbs } = useBreadcrumbs();
  const { toggle, opened } = useNavbar();

  useEffect(() => {
    if (breadcrumbs.length === 0) {
      document.title = 'UploadPilot';
      return;
    }

    if (breadcrumbs.length === 1) {
      document.title = 'UploadPilot | ' + breadcrumbs[0].label;
      return;
    }

    document.title =
      'UploadPilot | ' + (breadcrumbs[breadcrumbs.length - 1]?.label || '');
  }, [breadcrumbs]);

  if (breadcrumbs.length === 0) {
    return (
      <>
        <IconMenu4
          size={18}
          onClick={toggle}
          style={{ cursor: 'pointer' }}
          display={opened || window.innerWidth > 768 ? 'none' : 'block'}
        />
      </>
    );
  }

  return (
    <Group mb="lg">
      <IconMenu4
        size={18}
        onClick={toggle}
        style={{ cursor: 'pointer' }}
        display={opened || window.innerWidth > 768 ? 'none' : 'block'}
      />
      <IconArrowLeft
        display={breadcrumbs.length > 1 ? 'block' : 'none'}
        size={18}
        onClick={() =>
          breadcrumbs.length > 1 &&
          navigate(`${breadcrumbs[breadcrumbs.length - 2].path}`)
        }
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
