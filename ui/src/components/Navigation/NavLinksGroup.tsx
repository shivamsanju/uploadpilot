import { useState } from 'react';
import { IconCalendarStats, IconChevronRight } from '@tabler/icons-react';
import { Box, Collapse, Group, Text, ThemeIcon, UnstyledButton } from '@mantine/core';
import classes from './NavLinksGroup.module.css';
import { useNavigate } from 'react-router-dom';

interface LinksGroupProps {
    icon: React.FC<any>;
    label: string;
    initiallyOpened?: boolean;
    links?: { label: string; link: string }[];
    link?: string;
    active?: boolean;
    isWorkspaceChild?: boolean;
}

export function LinksGroup({ icon: Icon, label, initiallyOpened, links, link, active }: LinksGroupProps) {
    const navigate = useNavigate();

    const hasLinks = Array.isArray(links);
    const [opened, setOpened] = useState(initiallyOpened || false);

    const items = (hasLinks ? links : []).map((link) => (
        <Text<'a'>
            component="a"
            className={classes.link}
            href={link.link}
            key={link.label}
            onClick={(event) => {
                event.preventDefault();
                navigate(link.link);
            }}
        >
            {link.label}
        </Text>
    ));

    return (
        <>
            <UnstyledButton onClick={() => hasLinks ? setOpened((o) => !o) : (link ? navigate(link) : "")} className={classes.control}>
                <Group justify="space-between" gap={0}>
                    <Box style={{ display: 'flex', alignItems: 'center' }}>
                        <ThemeIcon variant="light" size={30}>
                            <Icon size={18} />
                        </ThemeIcon>
                        <Box ml="md" className={`${classes.label} ${active ? classes.active : ''}`}>{label}</Box>
                    </Box>
                    {hasLinks && (
                        <IconChevronRight
                            className={classes.chevron}
                            stroke={1.5}
                            size={16}
                            style={{ transform: opened ? 'rotate(-90deg)' : 'none' }}
                        />
                    )}
                </Group>
            </UnstyledButton>
            {hasLinks ? <Collapse in={opened}>{items}</Collapse> : null}
        </>
    );
}

const mockdata = {
    label: 'Releases',
    icon: IconCalendarStats,
    links: [
        { label: 'Upcoming releases', link: '/' },
        { label: 'Previous releases', link: '/' },
        { label: 'Releases schedule', link: '/' },
    ],
};

export function NavbarLinksGroup() {
    return (
        <Box mih={220} p="md">
            <LinksGroup {...mockdata} />
        </Box>
    );
}