package gobump

import (
	"encoding/json"
	"archive/zip"
	"errors"
	"bytes"
	"io"
)

func readConfig(zip *zip.ReadCloser)(*ModuleConfig,error){
	var m = &ModuleConfig{}
	
	 for _, f := range zip.File {
            if f.Name == PLUGIN_CONFIG_FILENAME{
				confF, err := f.Open()

				if err != nil {
						return nil, errors.New("Failed to open plugin config: "+err.Error())
				}

				defer confF.Close()

				dec := json.NewDecoder(confF)
				
				if err := dec.Decode(&m); err == io.EOF {
				} else if err != nil {
					return nil, errors.New("Failed to open decode config: "+err.Error())
				}
				return m, nil
			}
    }
    return nil, errors.New("Invalid or incorrectly formed plugin: no configuration file ("+PLUGIN_CONFIG_FILENAME+") detected.")
}


func readFile(filename string, zip *zip.ReadCloser)([]byte,error){
	for _, f := range zip.File {
		if f.Name == filename{
			F, err := f.Open()
			if err != nil {
				return nil, errors.New("Failed to open file for reading: "+err.Error())
			}
			
			buf := new(bytes.Buffer)
			buf.ReadFrom(F)
			return buf.Bytes(), nil

			defer F.Close()
		}
	}
	return nil, errors.New(filename+": 404 file not found")
}
