import { useMantineColorScheme } from "@mantine/core";
import { AgGridReact, AgGridReactProps } from "ag-grid-react"
import { appTheme, appDarkTheme } from "../../style/agGridTheme";

export const ThemedAgGridReact = (props: AgGridReactProps) => {
    const { colorScheme } = useMantineColorScheme();
    return (
        <AgGridReact {...props} theme={colorScheme === "dark" ? appDarkTheme : appTheme} />
    )
}