"use client";

import { Box } from "@mantine/core";
import classes from "./Header.module.css";
import { Logo } from "../Logo/Logo";
import UserMenu from "../UserMenu";

interface Props {
  burger?: React.ReactNode;
}

export function AdminHeader({ burger }: Props) {

  return (
    <header className={classes.header}>
      {burger}
      <Logo />
      <Box style={{ flex: 1 }} />

      <UserMenu />
    </header>
  );
}
