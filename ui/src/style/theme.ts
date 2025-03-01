import {
  Anchor,
  Badge,
  Burger,
  Button,
  ColorInput,
  createTheme,
  Input,
  Loader,
  MultiSelect,
  NumberInput,
  Paper,
  PasswordInput,
  SegmentedControl,
  Select,
  Switch,
  TagsInput,
  Text,
  Textarea,
  TextInput,
  virtualColor,
} from '@mantine/core';
import { DateInput } from '@mantine/dates';

type Icolors = readonly [
  string,
  string,
  string,
  string,
  string,
  string,
  string,
  string,
  string,
  string,
  ...string[],
];

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

const white: Icolors = [
  '#ffffff',
  '#fafafa',
  '#e8e8e8',
  '#d3d3d3',
  '#bfbfbf',
  '#b0b0b0',
  '#a9a9a9',
  '#989898',
  '#8c8c8c',
  '#7e7e7e',
];

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

// const lightBlue: Icolors = [
//   '#e1f8ff',
//   '#cbedff',
//   '#9ad7ff',
//   '#64c1ff',
//   '#3aaefe',
//   '#20a2fe',
//   '#099dff',
//   '#0088e4',
//   '#0079cd',
//   '#0069b6',
// ];

export const myAppTheme = createTheme({
  primaryColor: 'appcolor',
  fontFamily: 'Inter',
  headings: {
    fontFamily: 'Inter',
  },
  colors: {
    textColor: virtualColor({
      name: 'textColor',
      dark: '#F3F5F7',
      light: '#7a7a7b',
    }),
    appcolor: white,
  },
  defaultRadius: 'md',
  components: {
    TextInput: TextInput.extend({
      defaultProps: {
        size: 'xs',
        bd: 'none',
      },
    }),
    NumberInput: NumberInput.extend({
      defaultProps: {
        size: 'xs',
      },
    }),
    MultiSelect: MultiSelect.extend({
      defaultProps: {
        size: 'xs',
      },
    }),
    TagsInput: TagsInput.extend({
      defaultProps: {
        size: 'xs',
      },
    }),
    Select: Select.extend({
      defaultProps: {
        size: 'xs',
      },
    }),
    Textarea: Textarea.extend({
      defaultProps: {
        size: 'xs',
      },
    }),
    Input: Input.extend({
      defaultProps: {
        size: 'xs',
      },
    }),
    ColorInput: ColorInput.extend({
      defaultProps: {
        size: 'xs',
      },
    }),
    PasswordInput: PasswordInput.extend({
      defaultProps: {
        size: 'xs',
      },
    }),
    DateInput: DateInput.extend({
      defaultProps: {
        size: 'xs',
      },
    }),
    Button: Button.extend({
      defaultProps: {
        size: 'xs',
        radius: 'md',
      },
    }),

    Badge: Badge.extend({
      defaultProps: {
        size: 'xs',
        variant: 'light',
        radius: 'xl',
      },
    }),
    Switch: Switch.extend({
      defaultProps: {
        size: 'lg',
      },
    }),
    SegmentedControl: SegmentedControl.extend({
      defaultProps: {
        size: 'xs',
      },
    }),
    Text: Text.extend({
      defaultProps: {
        size: 'xs',
        color: 'textColor',
      },
    }),
    Anchor: Anchor.extend({
      defaultProps: {
        size: 'xs',
      },
    }),
    Burger: Burger.extend({
      defaultProps: {
        size: 'xs',
      },
    }),
    Paper: Paper.extend({
      defaultProps: {
        radius: 'md',
        style: {
          borderColor:
            'light-dark(var(--mantine-color-gray-4), var(--mantine-color-dark-6))',
        },
      },
    }),
    Loader: Loader.extend({
      defaultProps: {
        type: 'dots',
      },
    }),
  },
});
