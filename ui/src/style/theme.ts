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
    appcolor: [
      "#fff8e1",
      "#ffefcb",
      "#ffdd9a",
      "#ffca64",
      "#ffba38",
      "#ffb01b",
      "#ffab09",
      "#e39500",
      "#cb8400",
      "#b07100",
    ],
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
        radius: "sm",
        style: {
          borderColor:
            "light-dark(var(--mantine-color-gray-1), var(--mantine-color-dark-8))",
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
