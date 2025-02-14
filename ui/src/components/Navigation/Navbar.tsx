import {
  IconAdjustments,
  IconCircles,
  IconDatabase,
  IconRoute,
  IconGauge,
  IconUsers,
  IconShoppingCartBolt,
} from "@tabler/icons-react";
import { ScrollArea } from "@mantine/core";
import { LinksGroup } from "./NavLinksGroup";
import classes from "./Navbar.module.css";
import { useLocation, useParams } from "react-router-dom";
import WorkspaceSwitcher from "../WorkspaceSwitcher";
import { useMemo } from "react";

const isActive = (pathname: string, item: any) => {
  if (item.label === "Get Started") {
    return pathname === item.link;
  }
  return pathname.includes(item.link);
};
const NavBar = ({
  toggle,
}: {
  toggle: React.Dispatch<React.SetStateAction<boolean>>;
}) => {
  const { pathname } = useLocation();
  const { workspaceId } = useParams();

  const links = useMemo(() => {
    return [
      {
        label: "Get Started",
        icon: IconCircles,
        link: `/workspaces/${workspaceId}`,
      },
      {
        label: "Uploads",
        icon: IconDatabase,
        link: `/workspaces/${workspaceId}/uploads`,
      },
      {
        label: "Processors",
        icon: IconRoute,
        link: `/workspaces/${workspaceId}/processors`,
      },
      {
        label: "Configuration",
        icon: IconAdjustments,
        link: `/workspaces/${workspaceId}/configuration`,
      },
      {
        label: "Users",
        icon: IconUsers,
        link: `/workspaces/${workspaceId}/users`,
      },
      {
        label: "Marketplace",
        icon: IconShoppingCartBolt,
        link: `/workspaces/${workspaceId}/tools`,
      },
      {
        label: "Analytics",
        icon: IconGauge,
        link: `/workspaces/${workspaceId}/analytics`,
      },
    ].map((item) => (
      <LinksGroup
        {...item}
        key={item.label}
        active={isActive(pathname, item)}
        toggle={toggle}
      />
    ));
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
};

export default NavBar;
