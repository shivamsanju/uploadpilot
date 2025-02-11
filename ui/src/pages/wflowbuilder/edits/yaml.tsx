import React, { useState } from "react";
import Editor, { Monaco } from "@monaco-editor/react";
import { parse } from "yaml";
import Ajv, { ErrorObject } from "ajv";
import { schema } from "./schema";
import {
  Box,
  Group,
  LoadingOverlay,
  Text,
  Title,
  useMantineColorScheme,
} from "@mantine/core";
import { IconAlertCircle, IconCircleCheck } from "@tabler/icons-react";
import { DiscardButton } from "../../../components/Buttons/DiscardButton";
import { SaveButton } from "../../../components/Buttons/SaveButton";
import { useUpdateProcessorWorkflowMutation } from "../../../apis/processors";

const YamlEditor: React.FC = () => {
  const [yamlContent, setYamlContent] = useState<string>("name: John\nage: 25");
  const [error, setError] = useState<string | null>(null);
  const { colorScheme } = useMantineColorScheme();

  const { mutateAsync } = useUpdateProcessorWorkflowMutation();
  const handleEditorDidMount = (monaco: Monaco) => {
    monaco.editor.defineTheme("myCustomThemeDark", {
      base: "vs-dark",
      inherit: true,
      rules: [{ token: "comment", fontStyle: "italic" }],
      colors: {
        "editor.background": "#141414",
      },
    });
  };

  const validateYaml = (content: string) => {
    try {
      const parsedData = parse(content) as unknown;
      const ajv = new Ajv();
      const validate = ajv.compile(schema);
      const valid = validate(parsedData);

      if (!valid) {
        setError(
          validate.errors?.map((err: ErrorObject) => err.message).join(", ") ||
            "Invalid YAML"
        );
      } else {
        setError(null);
      }
    } catch (e) {
      setError((e as Error).message);
    }
  };

  const saveYaml = () => {
    if (!error) {
      console.log("Saving YAML:", yamlContent);
    }
  };

  return (
    <Box>
      <LoadingOverlay
        visible={false}
        overlayProps={{ backgroundOpacity: 0 }}
        zIndex={1000}
      />
      <Group justify="space-between" align="center" p="xs">
        <Box>
          <Title order={4} opacity={0.8}>
            Steps
          </Title>
          <Group
            align="center"
            gap={2}
            c={error ? "red" : "dimmed"}
            p={0}
            pt={2}
          >
            {error ? (
              <IconAlertCircle size="12" />
            ) : (
              <IconCircleCheck size="12" />
            )}
            <Text size="xs">{error || "All good"}</Text>
          </Group>
        </Box>
        <Group gap="md">
          <DiscardButton />
          <SaveButton onClick={saveYaml} />
        </Group>
      </Group>

      <Editor
        beforeMount={handleEditorDidMount}
        theme={colorScheme === "dark" ? "myCustomThemeDark" : "vs"}
        language="yaml"
        height="70vh"
        defaultLanguage="yaml"
        value={yamlContent}
        onChange={(value: any) => {
          if (typeof value === "string") {
            setYamlContent(value);
            validateYaml(value);
          }
        }}
        options={{
          minimap: { enabled: false },
          scrollBeyondLastLine: false,
          renderLineHighlight: "none",
          padding: {
            top: 10,
          },
          rulers: [],
        }}
      />
    </Box>
  );
};

export default YamlEditor;
