import {
    IconAdjustments,
    IconCircles,
    IconDatabase,
    IconGauge,
    IconUsers,
    IconWebhook,
} from '@tabler/icons-react';
import { ScrollArea } from '@mantine/core';
import { LinksGroup } from './NavLinksGroup';
import classes from './Navbar.module.css';
import { useLocation, useParams } from 'react-router-dom';
import WorkspaceSwitcher from '../WorkspaceSwitcher';
import { useMemo } from 'react';



const NavBar = () => {
    const { pathname } = useLocation();
    const { workspaceId } = useParams();

    const links = useMemo(() => {
        return [
            {
                label: 'Get Started',
                icon: IconCircles,
                link: `/workspaces/${workspaceId}`,
            },
            {
                label: 'Imports',
                icon: IconDatabase,
                link: `/workspaces/${workspaceId}/imports`,
            },
            {
                label: 'Configuration',
                icon: IconAdjustments,
                link: `/workspaces/${workspaceId}/configuration`,
            },
            {
                label: 'Hooks',
                icon: IconWebhook,
                link: `/workspaces/${workspaceId}/hooks`,
            },
            {
                label: 'Webhooks',
                icon: IconWebhook,
                link: `/workspaces/${workspaceId}/webhooks`,
            },
            {
                label: 'Users',
                icon: IconUsers,
                link: `/workspaces/${workspaceId}/users`,
            },
            {
                label: 'Analytics',
                icon: IconGauge,
                link: `/workspaces/${workspaceId}/analytics`,
            },
        ].map((item) => <LinksGroup {...item} key={item.label} active={pathname === item.link} />)
    }, [workspaceId, pathname])


    return (
        <nav className={classes.navbar}>
            <ScrollArea className={classes.links}>
                <div className={classes.linksInner}>{links}</div>
            </ScrollArea>
            <div className={classes.footer}>
                <WorkspaceSwitcher />
            </div>
        </nav>
    );
}

export default NavBar