import { useMantineColorScheme } from "@mantine/core";
import { AgGridReact, AgGridReactProps } from "ag-grid-react"
import { appTheme, appDarkTheme } from "../../style/agGridTheme";

export const ThemedAgGridReact = (props: AgGridReactProps) => {
    const { colorScheme } = useMantineColorScheme();
    return (
        <AgGridReact
            className={colorScheme === "dark" ? "ag-dark" : ""}
            {...props}
            theme={colorScheme === "dark" ? appDarkTheme : appTheme}
        />
    )
}