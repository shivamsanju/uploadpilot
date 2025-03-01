import { Box, Collapse, Group, Text, Transition } from '@mantine/core';
import { useViewportSize } from '@mantine/hooks';
import { IconChevronRight } from '@tabler/icons-react';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import classes from './NavLinksGroup.module.css';

interface LinksGroupProps {
  icon: React.FC<any>;
  label: string;
  initiallyOpened?: boolean;
  links?: { label: string; link: string }[];
  link?: string;
  active?: boolean;
  isWorkspaceChild?: boolean;
  toggle: React.Dispatch<React.SetStateAction<boolean>>;
  collapsed: boolean;
}

export function LinksGroup({
  icon: Icon,
  label,
  initiallyOpened,
  links,
  link,
  active,
  toggle,
  collapsed,
}: LinksGroupProps) {
  const navigate = useNavigate();
  const { width } = useViewportSize();

  const hasLinks = Array.isArray(links);
  const [opened, setOpened] = useState(initiallyOpened || false);

  const items = (hasLinks ? links : []).map(link => (
    <Text<'a'>
      component="a"
      className={classes.link}
      href={link.link}
      key={link.label}
      onClick={event => {
        if (width < 768) {
          toggle(false);
        }
        event.preventDefault();
        navigate(link.link);
      }}
    >
      {link.label}
    </Text>
  ));

  const handleClick = () => {
    if (hasLinks) {
      setOpened(o => !o);
      return;
    }
    if (link) {
      navigate(link);
    }
    if (width < 768) {
      toggle(true);
    }
  };

  return (
    <>
      <Group
        justify="space-between"
        gap={0}
        onClick={handleClick}
        className={`${active ? classes.active : ''} ${classes.control} ${collapsed ? classes.collapsed : ''}`}
      >
        <Group align="center" gap="md" wrap="nowrap">
          <Icon size={18} />
          <Transition
            mounted={!collapsed}
            transition="fade"
            duration={200}
            timingFunction="ease"
          >
            {styles =>
              collapsed ? (
                <></>
              ) : (
                <Box className={`${classes.label}`} style={{ ...styles }}>
                  {label}
                </Box>
              )
            }
          </Transition>
        </Group>
        {hasLinks && (
          <IconChevronRight
            className={classes.chevron}
            stroke={1.5}
            size={16}
            style={{ transform: opened ? 'rotate(-90deg)' : 'none' }}
          />
        )}
      </Group>
      {hasLinks ? <Collapse in={opened}>{items}</Collapse> : null}
    </>
  );
}
