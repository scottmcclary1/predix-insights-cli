package predixinsights

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/fatih/color"
)

const (
	dagTemplateName = "dagfile"
)

var bold = color.New(color.Bold).SprintFunc()

// DAGTemplate struct
type DAGTemplate struct {
	Owner    string
	FlowName string
	Interval string
}

func newFileUploadBuffer(fileName, fileLocation string, fields, values []string) (bytes.Buffer, string, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Add your image file
	file, err := os.Open(fileLocation)
	if err != nil {
		return b, "", err
	}
	defer file.Close()

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return b, "", err
	}

	// Add the other fields
	var formWriter io.Writer
	for i := 0; i < len(fields); i++ {
		formWriter, err = w.CreateFormField(fields[i])
		if err != nil {
			return b, "", err
		}

		_, err = formWriter.Write([]byte(values[i]))
		if err != nil {
			return b, "", err
		}
	}

	formWriter, err = w.CreateFormFile("file", fileName)
	if err != nil {
		return b, "", err
	}
	formWriter.Write(fileContents)
	w.Close()

	return b, w.FormDataContentType(), nil
}

func newFileUploadBufferMultipleFiles(fileDetails []FileDetails) (bytes.Buffer, string, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	var formWriter io.Writer

	for index, fileDetail := range fileDetails {
		// Add your image file
		file, err := os.Open(fileDetail.FileLocation)
		if err != nil {
			return b, "", err
		}
		defer file.Close()

		fileContents, err := ioutil.ReadAll(file)
		if err != nil {
			return b, "", err
		}

		// Add the other fields
		for i := 0; i < len(fileDetail.Fields); i++ {
			formWriter, err = w.CreateFormField(fileDetail.Fields[i])
			if err != nil {
				return b, "", err
			}

			_, err = formWriter.Write([]byte(fileDetail.Values[i]))
			if err != nil {
				return b, "", err
			}
		}

		formWriter, err = w.CreateFormFile(fmt.Sprintf("File%d", index), fileDetail.FileName)
		if err != nil {
			return b, "", err
		}
		formWriter.Write(fileContents)

	}
	w.Close()
	return b, w.FormDataContentType(), nil
}

func newTemplatedUploadBuffer(fileName, fileLocation string, fields, values []string, dt DAGTemplate) (bytes.Buffer, string, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Add the other fields
	var formWriter io.Writer
	for i := 0; i < len(fields); i++ {
		formWriter, err := w.CreateFormField(fields[i])
		if err != nil {
			return b, "", err
		}

		_, err = formWriter.Write([]byte(values[i]))
		if err != nil {
			return b, "", err
		}
	}

	formWriter, err := w.CreateFormFile("file", fileName)
	if err != nil {
		return b, "", err
	}

	fileStr, err := readFile(fileLocation)
	if err != nil {
		return b, "", err
	}

	formWriter, err = writeFromTemplate(dt, fileStr, formWriter)
	if err != nil {
		return b, "", err
	}

	w.Close()

	return b, w.FormDataContentType(), nil
}

func readFile(fileName string) (string, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func writeFromTemplate(dt DAGTemplate, templateContent string, w io.Writer) (io.Writer, error) {
	template := template.New("Template")
	template, err := template.Parse(templateContent)
	if err != nil {
		return nil, err
	}

	err = template.Execute(w, dt)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (ac *Client) dumpRequest(req *http.Request) {
	if ac.Verbose {
		dump, err := httputil.DumpRequestOut(req, true)
		if err == nil {
			fmt.Printf("%s\n%s\n", bold("REQUEST:"), string(dump))
		}
	}
}
func (ac *Client) dumpResponse(res *http.Response) {
	if ac.Verbose {
		dump, err := httputil.DumpResponse(res, true)
		if err == nil {
			fmt.Printf("%s\n%s\n\n", bold("RESPONSE:"), string(dump))
		}
	}
}
