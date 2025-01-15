import { themeQuartz, iconSetQuartzLight, colorSchemeDark } from 'ag-grid-community';


const commonTheme = {
    accentColor: "#0CA678",
    fontFamily: "Roboto",
    headerFontSize: 13,
    headerFontWeight: 500,
    cardShadow: "rgba(0, 0, 0, 0.1) 0px 1px 2px 0px",
    wrapperBorder: "none",
    wrapperBorderRadius: 10,
    columnBorder: false,
    rowVerticalPaddingScale: 1.3,
    cellHorizontalPaddingScale: 0.7,
    spacing: 6,
    border: 'none',
    backgroundColor: "transparent",
    headerBackgroundColor: "transparent",

}
export const appTheme = themeQuartz
    .withPart(iconSetQuartzLight)
    .withParams({
        ...commonTheme,
        browserColorScheme: "dark",
        rowHoverColor: "rgb(240, 240, 240)",
        borderColor: "rgb(240, 240, 240)",
        foregroundColor: "rgb(46, 46, 46)",
        menuBackgroundColor: "rgb(248, 249, 250)",
    });


export const appDarkTheme = themeQuartz
    .withPart(colorSchemeDark)
    .withPart(iconSetQuartzLight)
    .withParams({
        ...commonTheme,
        browserColorScheme: "dark",
        rowHoverColor: "rgb(46, 46, 46)",
        borderColor: "rgb(46, 46, 46)",
        foregroundColor: "rgb(200, 200, 200)",
        menuBackgroundColor: "#242424",
    });