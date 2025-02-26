import { Anchor, Breadcrumbs } from '@mantine/core';
import { Link, useLocation } from 'react-router-dom';

const Breadcrumb = () => {
  const location = useLocation();

  const pathParts = location.pathname.split('/').filter(Boolean);

  const breadcrumbItems = pathParts.map((part, index) => {
    const href = `/${pathParts.slice(0, index + 1).join('/')}`;

    const title = part.charAt(0).toUpperCase() + part.slice(1);

    return (
      <Anchor key={index} component={Link} to={href}>
        {title}
      </Anchor>
    );
  });

  const items = [
    <Anchor key="home" component={Link} to="/">
      Home
    </Anchor>,
    ...breadcrumbItems,
  ];

  return <Breadcrumbs>{items}</Breadcrumbs>;
};

export default Breadcrumb;
