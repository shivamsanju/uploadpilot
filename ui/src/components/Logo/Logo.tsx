import { Flex, Text } from "@mantine/core";
import { Link } from "react-router-dom";
import classes from "./Logo.module.css";

interface Props {
  width?: string;
  height?: string;
}

export const Logo: React.FC<Props> = () => {
  return (
    <Flex direction="row" align="center" gap={4}>
      <Link
        to="/"
        style={{ textDecoration: "none" }}
        className={classes.heading}
      >
        <Text fw="bolder" size="xl">
          Upload
          <Text component="span" fw="normal" className={classes.subheading}>
            Pilot
          </Text>
        </Text>
      </Link>
    </Flex>
  );
};
