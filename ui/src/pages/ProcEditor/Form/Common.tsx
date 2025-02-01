import { Group, NumberInput, Switch, Text } from "@mantine/core"

export const CommonForm = ({ form }: any) => {
    return (
        <>
            <NumberInput
                label="Retry Count"
                type="retry"
                description="Number of retries for the task"
                placeholder="Enter the retry count"
                defaultValue={0}
                {...form.getInputProps('retry')}
            />

            <NumberInput
                label="Timeout MilSec"
                description="Timeout for the task in milliseconds"
                placeholder="Enter the timeout in milliseconds"
                defaultValue={0}
                {...form.getInputProps('timeoutMilSec')}
            />


            <Group justify="space-between" mt="xs">
                <div>
                    <Text fw={500}>Continue on error </Text>
                    <Text c="dimmed">
                        Continue executing next tasks even if the task fails
                    </Text>
                </div>
                <Switch
                    onLabel="ON" offLabel="OFF"
                    checked={form.values.continueOnError}
                    defaultChecked={false}
                    onChange={(e) => form.setFieldValue("continueOnError", e.target.checked)}
                />
            </Group>
        </>
    )
}