import { useLocation, Link } from 'react-router-dom';
import { Breadcrumbs, Anchor } from '@mantine/core';

function Breadcrumb() {
    const location = useLocation();

    const pathParts = location.pathname.split('/').filter(Boolean);

    const breadcrumbItems = pathParts.map((part, index) => {
        const href = `/${pathParts.slice(0, index + 1).join('/')}`;

        const title = part.charAt(0).toUpperCase() + part.slice(1);

        return (
            <Anchor size="xs" key={index} component={Link} to={href}>
                {title}
            </Anchor>
        );
    });

    const items = [
        <Anchor size="xs" key="home" component={Link} to="/">
            Home
        </Anchor>,
        ...breadcrumbItems,
    ];

    return <Breadcrumbs>{items}</Breadcrumbs>;
}

export default Breadcrumb;
