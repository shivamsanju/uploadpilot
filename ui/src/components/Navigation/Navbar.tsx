import {
    IconAdjustments,
    IconCircles,
    IconDatabase,
    IconGauge,
    IconPackageImport,
    IconWebhook,
} from '@tabler/icons-react';
import { Code, ScrollArea } from '@mantine/core';
import { LinksGroup } from './NavLinksGroup';
import classes from './Navbar.module.css';

const mockdata = [
    { label: 'Dashboard', icon: IconGauge, link: "/" },
    {
        label: 'Workflows',
        icon: IconCircles,
        link: "/workflows",
    },
    {
        label: 'Import',
        icon: IconPackageImport,
        links: [
            { label: 'Policies', link: '/importPolicies' },
            { label: 'Validations', link: '/validations' },
        ],
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

            <div className={classes.footer}>
                <Code fw={700} p="xs">Version : v0.0.1</Code>
            </div>
        </nav>
    );
}

export default NavBar