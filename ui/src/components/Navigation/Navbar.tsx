import { Box, ScrollArea } from '@mantine/core';
import { FC, ReactNode, useMemo } from 'react';
import { useLocation } from 'react-router-dom';
import { useNavbar } from '../../context/NavbarContext';
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
  items,
  footer,
}: {
  items: NavItem[];
  footer?: ReactNode;
}) => {
  const { pathname } = useLocation();
  const { opened } = useNavbar();

  const links = useMemo(() => {
    return items.map(item => (
      <LinksGroup
        {...item}
        key={item.label}
        active={isActive(pathname, item)}
        collapsed={!opened}
      />
    ));
  }, [items, pathname, opened]);

  return (
    <nav className={classes.navbar}>
      <Box mb="lg">
        <MenuButton />
      </Box>
      <ScrollArea className={classes.links}>{links}</ScrollArea>
      <Box px={opened ? 'sm' : 0} mb="xs">
        <UserButton collapsed={!opened} />
      </Box>
    </nav>
  );
};

export default NavBar;
