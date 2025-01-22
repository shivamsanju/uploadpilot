import { Flex, Image, useMantineColorScheme } from "@mantine/core";
import { Link } from "react-router-dom";
import DarkLogo from "../../assets/images/full-logo-dark.png"
import LightLogo from "../../assets/images/full-logo.png"

interface Props {
    width?: string;
    height?: string;
    enableOnClick?: boolean;
}

export const Logo2: React.FC<Props> = ({ height, width, enableOnClick }) => {
    const { colorScheme } = useMantineColorScheme();
    return (
        <Flex direction="row" align="center" gap={4}>
            <Link
                onClick={(e) => !enableOnClick && e.preventDefault()}
                to="/"
                style={{ textDecoration: "none" }}
            >
                {colorScheme === "dark" ? <Image src={DarkLogo} alt="logo" h={height} w={width} /> : <Image src={LightLogo} alt="logo" h={height} w={width} />}
            </Link>
        </Flex>
    );
};
