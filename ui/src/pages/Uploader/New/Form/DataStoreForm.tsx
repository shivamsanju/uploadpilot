import { Stack, Select, TextInput } from "@mantine/core";
import { optionsFilter } from "../../../../utils/filter";
import { useGetStorageConnectors } from "../../../../apis/storage";
import { Connector } from "../../../../types/connector";
import { useMemo } from "react";
import { UseFormReturnType } from "@mantine/form";
import { CreateUploaderForm } from "../../../../types/uploader";
import ErrorCard from "../../../../components/ErrorCard/ErrorCard";


type DataStorePageProps = {
    form: UseFormReturnType<CreateUploaderForm, (values: CreateUploaderForm) => CreateUploaderForm>;
}
const DataStorePage: React.FC<DataStorePageProps> = ({ form }) => {

    const { isPending, error, connectors } = useGetStorageConnectors();

    const selectedConnector = useMemo(() => connectors?.find((connector: Connector) => connector.id === form.values.connectorId), [form.values.connectorId, connectors]);
    return error ? <ErrorCard title={error.name} message={error.message} h="50vh" /> : (
        <Stack gap="md">
            <Select
                withAsterisk
                filter={optionsFilter}
                label="Storage connector"
                placeholder="Select a storage connector"
                data={isPending ? [] : connectors.map((connector: Connector) => ({ ...connector, value: connector.id, label: connector.name }))}
                {...form.getInputProps("connectorId")}
            />

            {selectedConnector?.id && (
                <>
                    <TextInput
                        withAsterisk
                        label="Datastore Name"
                        placeholder="Enter a datastore name"
                        {...form.getInputProps("dataStoreName")}

                    />
                    <TextInput
                        withAsterisk
                        label={selectedConnector.type === "azure" ? "Container Name" : "Bucket Name"}
                        placeholder={`Enter a ${selectedConnector.type === "azure" ? "container" : "bucket"} name`}
                        {...form.getInputProps("bucket")}
                    />
                </>
            )}
        </Stack>
    )
};

export default DataStorePage;
