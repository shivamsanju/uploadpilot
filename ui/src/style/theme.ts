import { createTheme, TextInput, Button, Select, Textarea, TagsInput, Badge } from "@mantine/core";

export const theme = createTheme({
    primaryColor: 'grape',
    fontFamily: "Inter",
    headings: {
        fontFamily: "Inter",
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
