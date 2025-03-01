import { Box, ScrollArea } from '@mantine/core';
import { FC, ReactNode, useMemo } from 'react';
import { useLocation } from 'react-router-dom';
import { MenuButton } from '../MenuButton/MenuButton';
import UserButton from '../UserMenu';
import classes from './Navbar.module.css';
import { LinksGroup } from './NavLinksGroup';

const isActive = (pathname: string, item: any) => {
  if (item.label === 'Get Started' || item.label === 'Workspaces') {
    return pathname === item.link;
  }
  return pathname.includes(item.link);
};

type NavItem = {
  label: string;
  icon: FC<any>;
  link: string;
  links?: { label: string; link: string }[];
};

const NavBar = ({
  collapsed,
  toggle,
  items,
  footer,
}: {
  collapsed: boolean;
  toggle: React.Dispatch<React.SetStateAction<boolean>>;
  items: NavItem[];
  footer?: ReactNode;
}) => {
  const { pathname } = useLocation();

  const links = useMemo(() => {
    return items.map(item => (
      <LinksGroup
        {...item}
        key={item.label}
        active={isActive(pathname, item)}
        toggle={toggle}
        collapsed={collapsed}
      />
    ));
  }, [items, pathname, toggle, collapsed]);

  return (
    <nav className={classes.navbar}>
      <Box mb="lg">
        <MenuButton toggle={toggle} collapsed={collapsed} />
      </Box>
      <ScrollArea className={classes.links}>{links}</ScrollArea>
      <Box px={collapsed ? 0 : 'sm'} mb="xs">
        <UserButton collapsed={collapsed} />
      </Box>
    </nav>
  );
};

export default NavBar;
