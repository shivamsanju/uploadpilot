"use client";

import { Box } from "@mantine/core";
import classes from "./Header.module.css";
import { Logo } from "../Logo/Logo";
import UserMenu from "../UserMenu";
import { useNavigate } from "react-router-dom";

interface Props {
  burger?: React.ReactNode;
}

export function AdminHeader({ burger }: Props) {

  const navigate = useNavigate();
  return (
    <header className={classes.header}>
      {burger}
      <Box onClick={() => navigate("/")} style={{ cursor: "pointer" }}>
        <Logo height="40" width="129.64" />
      </Box>
      <Box style={{ flex: 1 }} />
      <UserMenu />
    </header>
  );
}
