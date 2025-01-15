import {
    IconAdjustments,
    IconCircles,
    IconDatabase,
    IconGauge,
    IconWebhook,
} from '@tabler/icons-react';
import { Code, ScrollArea } from '@mantine/core';
import { LinksGroup } from './NavLinksGroup';
import classes from './Navbar.module.css';

const mockdata = [
    { label: 'Dashboard', icon: IconGauge, link: "/" },
    {
        label: 'Uploaders',
        icon: IconCircles,
        link: "/uploaders",
    },
    {
        label: 'Storage',
        icon: IconDatabase,
        links: [
            { label: 'Datastores', link: '/storage/datastores' },
            { label: 'Connectors', link: '/storage/connectors' },
        ],
    },
    { label: 'Hooks', icon: IconWebhook, link: "/hooks" },
    { label: 'Settings', icon: IconAdjustments, link: "/settings" },
];

const NavBar = () => {
    const links = mockdata.map((item) => <LinksGroup {...item} key={item.label} />);

    return (
        <nav className={classes.navbar}>
            <ScrollArea className={classes.links}>
                <div className={classes.linksInner}>{links}</div>
            </ScrollArea>
        </nav>
    );
}

export default NavBar