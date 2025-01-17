import {
    IconAdjustments,
    IconCircles,
    IconDatabase,
    IconGauge,
    IconWebhook,
} from '@tabler/icons-react';
import { ScrollArea } from '@mantine/core';
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
        label: 'Storage Connectors',
        icon: IconDatabase,
        link: '/storageConnectors',
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