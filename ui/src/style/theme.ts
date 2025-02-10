import {
  createTheme,
  TextInput,
  Button,
  Select,
  Textarea,
  TagsInput,
  Badge,
  NumberInput,
  Text,
  virtualColor,
  MultiSelect,
  Input,
  PasswordInput,
  Anchor,
  Burger,
  Paper,
  Switch,
  SegmentedControl,
  Loader,
  ColorInput,
} from "@mantine/core";


type Icolors = readonly [string, string, string, string, string, string, string, string, string, string, ...string[]]

// const magenta: Icolors = [
//   "#ffe9f6",
//   "#ffd1e6",
//   "#faa1c9",
//   "#f66eab",
//   "#f24391",
//   "#f02981",
//   "#f01879",
//   "#d60867",
//   "#c0005c",
//   "#a9004f"
// ]

// const yellowOrange: Icolors= [
//   "#fff8e1",
//   "#ffefcb",
//   "#ffdd9a",
//   "#ffca64",
//   "#ffba38",
//   "#ffb01b",
//   "#ffab09",
//   "#e39500",
//   "#cb8400",
//   "#b07100",
// ]

// const white: Icolors = [
//   "#f5f5f5",
//   "#e7e7e7",
//   "#cdcdcd",
//   "#b2b2b2",
//   "#9a9a9a",
//   "#8b8b8b",
//   "#848484",
//   "#717171",
//   "#656565",
//   "#575757"
// ]
// const blue: Icolors = [
//   "#e2f7ff",
//   "#ceeaff",
//   "#9fd1fb",
//   "#6db8f6",
//   "#43a2f1",
//   "#2894ef",
//   "#128eef",
//   "#007ad6",
//   "#006dc1",
//   "#005eab"
// ]

const lightBlue: Icolors = [
  "#dffbff",
  "#caf2ff",
  "#99e2ff",
  "#64d2ff",
  "#3cc4fe",
  "#23bcfe",
  "#09b8ff",
  "#00a1e4",
  "#008fcd",
  "#007cb6"
]



export const myAppTheme = createTheme({
  primaryColor: "appcolor",
  fontFamily: "Inter",
  headings: {
    fontFamily: "Inter",
  },
  colors: {
    textColor: virtualColor({
      name: "textColor",
      dark: "#F3F5F7",
      light: "#7a7a7b",
    }),
    appcolor: lightBlue,
  },
  defaultRadius: "md",
  components: {
    TextInput: TextInput.extend({
      defaultProps: {
        size: "xs",
        bd: "none",
      },
    }),
    NumberInput: NumberInput.extend({
      defaultProps: {
        size: "xs",
      },
    }),
    MultiSelect: MultiSelect.extend({
      defaultProps: {
        size: "xs",
      },
    }),
    TagsInput: TagsInput.extend({
      defaultProps: {
        size: "xs",
      },
    }),
    Select: Select.extend({
      defaultProps: {
        size: "xs",
      },
    }),
    Textarea: Textarea.extend({
      defaultProps: {
        size: "xs",
      },
    }),
    Input: Input.extend({
      defaultProps: {
        size: "xs",
      },
    }),
    ColorInput: ColorInput.extend({
      defaultProps: {
        size: "xs",
      },
    }),
    PasswordInput: PasswordInput.extend({
      defaultProps: {
        size: "xs",
      },
    }),
    Button: Button.extend({
      defaultProps: {
        size: "xs",
      },
    }),
    Badge: Badge.extend({
      defaultProps: {
        size: "xs",
        variant: "light",
      },
    }),
    Switch: Switch.extend({
      defaultProps: {
        size: "lg",
      },
    }),
    SegmentedControl: SegmentedControl.extend({
      defaultProps: {
        size: "xs",
      },
    }),
    Text: Text.extend({
      defaultProps: {
        size: "xs",
        color: "textColor",
      },
    }),
    Anchor: Anchor.extend({
      defaultProps: {
        size: "xs",
      },
    }),
    Burger: Burger.extend({
      defaultProps: {
        size: "xs",
      },
    }),
    Paper: Paper.extend({
      defaultProps: {
        radius: 0,
        style: {
          borderColor:
            "light-dark(var(--mantine-color-gray-2), var(--mantine-color-dark-8))",
        },
      },
    }),
    Loader: Loader.extend({
      defaultProps: {
        type: "dots",
      },
    }),
  },
});
