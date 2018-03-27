package visitor

import (
	"encoding/json"
	"errors"
	"io"
)

type GetMerchantInfoUnmarshal struct {
}

func (u *GetMerchantInfoUnmarshal) Unmarshal(r io.Reader, v interface{}) error {
	if d, ok := v.(*GetMerchantInfoOutput); ok {
		var items []MerchantInfo
		decoder := json.NewDecoder(r)
		err := decoder.Decode(&items)

		if err != nil {
			return err
		}

		d.Items = items
		return nil
	}

	return errors.New("error data type")
}

type GetAdminUnmarshal struct{}

func (u *GetAdminUnmarshal) Unmarshal(r io.Reader, v interface{}) error {
	if d, ok := v.(*GetAdminOutput); ok {
		var a Admin
		decoder := json.NewDecoder(r)
		err := decoder.Decode(&a)

		if err != nil {
			return err
		}

		d.Admin = a
		return nil
	}

	return errors.New("error data type")
}

type GetMerchantDepartmentUnmarshal struct{}

func (u *GetMerchantDepartmentUnmarshal) Unmarshal(r io.Reader, v interface{}) error {
	if d, ok := v.(*GetMerchantDepartmentOutput); ok {
		var dep MerchantDepartment
		decoder := json.NewDecoder(r)
		err := decoder.Decode(&dep)

		if err != nil {
			return err
		}

		d.Dep = dep
		return nil
	}

	return errors.New("error data type")
}
