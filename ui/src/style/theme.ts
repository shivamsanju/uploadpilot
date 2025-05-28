import { CodeHighlight } from '@mantine/code-highlight';
import {
  Anchor,
  Badge,
  Burger,
  Button,
  Card,
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
  Tooltip,
  TooltipFloating,
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

const dark: Icolors = [
  '#E2E5E9',
  '#B3BAC5',
  '#6F7D8C',
  '#3E4A5A',
  '#1F2937',
  '#1A222E',
  '#151B25',
  '#10141C',
  '#0C0A15',
  '#08060F',
];

const appcolor: Icolors = [
  '#f2effa',
  '#e0dcef',
  '#bfb5e0',
  '#9c8bd2',
  '#7f68c6',
  '#6c51bf',
  '#6346bd',
  '#5338a7',
  '#493195',
  '#3e2a84',
];

export const myAppTheme = createTheme({
  primaryColor: 'appcolor',
  fontFamily: 'Inter',
  headings: {
    fontFamily: 'Inter',
  },
  colors: {
    textColor: virtualColor({
      name: 'textColor',
      dark: '#E2E5E9',
      light: '#7a7a7b',
    }),
    appcolor: appcolor,
    dark: dark,
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
      },
    }),
    Card: Card.extend({
      defaultProps: {
        radius: 'md',
      },
    }),
    Loader: Loader.extend({
      defaultProps: {
        type: 'dots',
      },
    }),
    CodeHighlight: CodeHighlight.extend({
      defaultProps: {
        bg: 'light-dark(var(--mantine-color-gray-1), var(--mantine-color-dark-7))',
        copiedLabel: 'Copied',
        copyLabel: 'Copy',
      },
    }),
    Tooltip: Tooltip.extend({
      defaultProps: {
        fs: 'xs',
        bg: 'light-dark(var(--mantine-color-dark-8), var(--mantine-color-dark-2))',
      },
    }),
    TooltipFloating: TooltipFloating.extend({
      defaultProps: {
        fs: 'xs',
        bg: 'light-dark(var(--mantine-color-dark-8), var(--mantine-color-dark-2))',
      },
    }),
  },
});

// const white: Icolors = [
//   '#ffffff',
//   '#fafafa',
//   '#e8e8e8',
//   '#d3d3d3',
//   '#bfbfbf',
//   '#b0b0b0',
//   '#a9a9a9',
//   '#989898',
//   '#8c8c8c',
//   '#7e7e7e',
// ];

// const defaultDark: Icolors = [
//   '#C9C9C9',
//   '#b8b8b8',
//   '#828282',
//   '#696969',
//   '#424242',
//   '#3b3b3b',
//   '#2e2e2e',
//   '#242424',
//   '#1f1f1f',
//   '#141414',
// ];
