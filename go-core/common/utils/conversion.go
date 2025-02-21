package commonutils

import "encoding/json"

func ConvertDTOToModel[T any, M any](dto *T, model *M) error {
	data, err := json.Marshal(dto)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, model)
}
