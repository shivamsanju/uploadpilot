import { themeQuartz, iconSetQuartzLight, colorSchemeDark } from 'ag-grid-community';

export const appTheme = themeQuartz
    .withPart(iconSetQuartzLight)
    .withParams({
        fontFamily: "Poppins",
        browserColorScheme: "dark",
        headerFontSize: 13,
        headerFontWeight: 500,
        fontSize: 12,
        cardShadow: "rgba(0, 0, 0, 0.1) 0px 1px 2px 0px",
    });


export const appDarkTheme = themeQuartz
    .withPart(colorSchemeDark)
    .withPart(iconSetQuartzLight)
    .withParams({
        fontFamily: "Poppins",
        browserColorScheme: "light",
        headerFontSize: 13,
        headerFontWeight: 500,
        fontSize: 12,
        cardShadow: "rgba(0, 0, 0, 0.1) 0px 1px 2px 0px",
    });