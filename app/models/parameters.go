package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"runtime"
	"strconv"
)

type Parameter struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}

type Parameters []Parameter

var (
	TMP_DIRECTORY = "/tmp/"
)

func init() {
	if runtime.GOOS == "windows" {
		TMP_DIRECTORY = "C:\\ProgramData\\Witbe\\storage\\data\\jsonception\\"
		fmt.Println("Running under Windows")
	}
}

func (p *Parameter) Print() {
	fmt.Printf("name : %s\tvalue : %s\ttype : %s\n", p.Name, p.Value, p.Type)
}

func (p *Parameters) CheckTestParamsWithExecParams(others *Parameters) error {
	fmt.Printf("len j : %d\tlen other : %d\n", len(*p), len(*others))
	if len(*p) != len(*others) {
		return errors.New("Number of parameters doesnt match between Test and Script model")
	}
	for i, v := range *p {
		if v.Name != (*others)[i].Name {
			return errors.New("Name parameters differ between Test and Script model")
		}
		if v.Type != (*others)[i].Type {
			return errors.New("Type parameters differ between Test and Script model")
		}
	}
	return nil
}

func (j *Parameters) GiveFileParamsValue(params_name, file_name string) {
	fmt.Printf("params name : %s\t[aras : %s\n", params_name, *j)
	for i, v := range *j {
		if v.Name == params_name {
			(*j)[i].Value = TMP_DIRECTORY + file_name
		}
	}
}

func (j *Parameters) UploadFileFromParameters(m *multipart.Form) error {
	if _, err := os.Stat(TMP_DIRECTORY); os.IsNotExist(err) {
		os.Mkdir(TMP_DIRECTORY, os.ModeDir)
	}

	for fname, _ := range m.File {
		fheaders := m.File[fname]
		for i, _ := range fheaders {
			j.GiveFileParamsValue(fname, fheaders[i].Filename)
			//for each fileheader, get a handle to the actual file
			file, err := fheaders[i].Open()
			defer file.Close() //close the source file handle on function return
			if err != nil {
				return err
			}
			//create destination file making sure the path is writeable.
			dst_path := TMP_DIRECTORY + fheaders[i].Filename
			if _, err := os.Stat(dst_path); os.IsNotExist(err) {
				dst, err := os.Create(dst_path)
				defer dst.Close()                             //close the destination file handle on function return
				defer os.Chmod(dst_path, (os.FileMode)(0644)) //limit access restrictions
				if err != nil {
					return err
				}
				//copy the uploaded file to the destination file
				if _, err := io.Copy(dst, file); err != nil {
					return err
				}
			}
		}
	}

	// set unnamed file default
	for i, v := range *j {
		if v.Type == "file" && v.Value == nil {
			(*j)[i].Value = ""
		}
	}
	fmt.Printf("prams ; %s\n", *j)
	return nil
}

func (j Parameters) Check() error {
	names := make(map[string]int)
	for i, v := range j {
		if v.Name == "" {
			return errors.New("Parameter's name cannot be empty")
		}
		if v.Value == nil {
			continue
		}
		names[v.Name]++
		if names[v.Name] > 1 {
			return errors.New("Parameters must have unique name")
		}
		if v.Value == nil {
			continue
		}
		if value, ok := v.Value.(string); ok {
			if value == "" {
				continue
			}
			if v.Type == "bool" {
				if value != "0" && value != "false" && value != "true" && value != "1" {
					return errors.New("Bool cannot have value : " + value)
				}
				if newVal, err := strconv.ParseBool(value); err != nil {
					return err
				} else {
					fmt.Printf("value : %v\n", newVal)
					if newVal != true && newVal != false {
						return errors.New("Bool value cannot differ from 0 and 1")
					}
					j[i].Value = newVal
				}
			} else if v.Type == "int" {
				if newVal, err := strconv.Atoi(value); err != nil {
					fmt.Printf("wtf : %d\n", newVal)
					return err
				} else {
					j[i].Value = newVal
				}
			}
		}
	}
	return nil
}

func (j Parameters) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *Parameters) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}
