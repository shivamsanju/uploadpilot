import { Stack, Select, Divider, } from "@mantine/core";
import ImportPolicyForm from "../../../ImportPolicy/New/Form";
import { ImportPolicy } from "../../../../types/importpolicy";
import { useCreateImportPolicyMutation, useGetImportPolicies } from "../../../../apis/importPolicy";
import { optionsFilter } from "../../../../utils/filter";


const ImportPolicyPage = ({ form }: { form: any }) => {
    const { isPending, error, importPolicies } = useGetImportPolicies();
    const { mutateAsync } = useCreateImportPolicyMutation();


    const onSubmit = async (values: ImportPolicy) => {
        mutateAsync(values).then((data) => {
            console.log(data)
            form.setFieldValue("importPolicyId", values.id)
        });
    }




    return (<Stack gap="md">
        <Select
            withAsterisk
            filter={optionsFilter}
            loading={isPending}
            label="Import Policy"
            placeholder="Select an import policy"
            data={isPending ? [] : [{ value: "", label: "Create a new import policy" }, ...importPolicies.map((importPolicy: ImportPolicy) => ({ value: importPolicy.id, label: importPolicy.name }))]}
            {...form.getInputProps('importPolicyId')}
            limit={2}
        />

        {form.values.importPolicyId === "" && (
            <>
                <Divider my="md" label="Create a new import policy" />
                <ImportPolicyForm showCancelButton={false} showSubmitButton onSubmit={onSubmit} />
            </>
        )}
    </Stack>
    )
};

export default ImportPolicyPage;
