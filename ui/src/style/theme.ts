import { createTheme, TextInput, Button, Select, Textarea, TagsInput, Badge, Text } from "@mantine/core";

export const theme = createTheme({
    primaryColor: 'teal',
    fontFamily: "Roboto",
    headings: {
        fontFamily: "Roboto",
    },
    defaultRadius: "md",
    components: {
        TextInput: TextInput.extend({
            defaultProps: {
                size: "xs"
            }
        }),
        Button: Button.extend({
            defaultProps: {
                size: "xs"
            }
        }),
        Select: Select.extend({
            defaultProps: {
                size: "xs"
            }
        }),
        Textarea: Textarea.extend({
            defaultProps: {
                size: "xs"
            }
        }),
        TagsInput: TagsInput.extend({
            defaultProps: {
                size: "xs"
            }
        }),
        Badge: Badge.extend({
            defaultProps: {
                size: "xs",
                variant: "light"
            }
        }),
    }
});
