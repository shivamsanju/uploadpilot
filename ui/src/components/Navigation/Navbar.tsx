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



const NavBar = ({ toggle }: { toggle: React.Dispatch<React.SetStateAction<boolean>> }) => {
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
                label: 'Uploads',
                icon: IconDatabase,
                link: `/workspaces/${workspaceId}/uploads`,
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
        ].map((item) => <LinksGroup {...item} key={item.label} active={pathname === item.link} toggle={toggle} />)
    }, [workspaceId, pathname, toggle]);


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