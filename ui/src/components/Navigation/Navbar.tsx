import { ScrollArea } from "@mantine/core";
import { LinksGroup } from "./NavLinksGroup";
import classes from "./Navbar.module.css";
import { useLocation } from "react-router-dom";
import { FC, ReactNode, useMemo } from "react";

const isActive = (pathname: string, item: any) => {
  if (item.label === "Get Started" || item.label === "Workspaces") {
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
  toggle,
  items,
  footer,
}: {
  toggle: React.Dispatch<React.SetStateAction<boolean>>;
  items: NavItem[];
  footer?: ReactNode;
}) => {
  const { pathname } = useLocation();

  const links = useMemo(() => {
    return items.map((item) => (
      <LinksGroup
        {...item}
        key={item.label}
        active={isActive(pathname, item)}
        toggle={toggle}
      />
    ));
  }, [items, pathname, toggle]);

  return (
    <nav className={classes.navbar}>
      <ScrollArea className={classes.links}>
        <div className={classes.linksInner}>{links}</div>
      </ScrollArea>
      {footer && <div className={classes.footer}>{footer}</div>}
    </nav>
  );
};

export default NavBar;
