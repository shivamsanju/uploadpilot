import { Card, Text, Code, Badge, Group, Divider, ScrollArea, Title } from '@mantine/core';
import { Handle, Position } from 'reactflow';
import '@mantine/code-highlight/styles.css';
import { CodeHighlight } from '@mantine/code-highlight';

const sampleCode = `
func SumOfSquares(nums []int) int {
    sum := 0
    for _, n := range nums {
        sum += n * n
    }
    return sum
}
`

const sampleDocumentation = `
This function iterates over the provided slice of integers, computes the square of each element, and accumulates the total. The documentation string follows Go conventions to explain the purpose and usage.
`
function extractPackageAndFunction(fullString: string) {
    if (typeof fullString !== "string" || !fullString.includes(".")) {
        throw new Error("Input must be a string containing at least one dot ('.')");
    }

    fullString = fullString.replace(/"/g, '');
    // Split the string at the last dot
    const lastDotIndex = fullString.lastIndexOf(".");
    const packageName = fullString.substring(0, lastDotIndex);
    const functionName = fullString.substring(lastDotIndex + 1);

    return { packageName, functionName };
}

const CodeNode = ({ data }: any) => {
    const { packageName, functionName } = extractPackageAndFunction(data.name);

    return (
        <Card shadow="md" radius="md" style={{ width: 500 }} withBorder>
            <Card.Section>
                <Text c="green" variant="light" size="md" style={{ width: '100%', padding: '10px 0', textAlign: 'center' }}>
                    {packageName}
                </Text>
            </Card.Section>

            <Card.Section style={{ padding: '10px', textAlign: 'center' }}>
                <Title w={500} order={4} c="blue">
                    {functionName}
                </Title>
            </Card.Section>

            <Divider my="sm" label="Documentation" labelPosition="center" />

            <Card.Section style={{ padding: '10px', textAlign: 'center' }}>
                <ScrollArea style={{ maxHeight: 100 }}>
                    <Text size="sm" color="dimmed">
                        {data.documentation || sampleDocumentation}
                    </Text>
                </ScrollArea>
            </Card.Section>

            <Divider my="sm" label="Code" labelPosition="center" />

            <Card.Section style={{ padding: '10px' }}>
                <CodeHighlight code={data.code || sampleCode} language="go" />
            </Card.Section>

            <Group p="apart" style={{ marginTop: '10px' }}>
                <Handle type="target" position={Position.Top} id="b" style={{ background: 'blue' }} />
                <Handle type="source" position={Position.Bottom} id="a" style={{ background: 'green' }} />
            </Group>
        </Card>
    );
};

export default CodeNode;
