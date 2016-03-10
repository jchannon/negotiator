package negotiator

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

const defaultCSVContentType = "text/csv"

type csvProcessor struct {
	comma       rune
	contentType string
}

// NewCSV creates an output processor that serialises a dataModel in CSV form. With no arguments, the default
// format is comma-separated; you can supply any rune to be used as an alternative separator.
//
// Model values should be one of the following:
//
// * string or []string, or [][]string
//
// * fmt.Stringer or []fmt.Stringer, or [][]fmt.Stringer
//
// * []int or similar (bool, int8, int16, int32, int64, uint8, uint16, uint32, uint63, float32, float64, complex)
//
// * [][]int or similar (bool, int8, int16, int32, int64, uint8, uint16, uint32, uint63, float32, float64, complex)
//
// * struct for some struct in which all the fields are exported and of simple types (as above).
//
// * []struct for some struct in which all the fields are exported and of simple types (as above).
func NewCSV(comma ...rune) ResponseProcessor {
	if len(comma) > 0 {
		return &csvProcessor{comma[0], defaultCSVContentType}
	}
	return &csvProcessor{',', defaultCSVContentType}
}

// Implements ContentTypeSettable for this type.
func (p *csvProcessor) SetContentType(contentType string) ResponseProcessor {
	p.contentType = contentType
	return p
}

func (*csvProcessor) CanProcess(mediaRange string) bool {
	return strings.EqualFold(mediaRange, "text/csv") || strings.EqualFold(mediaRange, "text/*")
}

func (p *csvProcessor) Process(w http.ResponseWriter, dataModel interface{}) error {
	if dataModel == nil {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	w.Header().Set("Content-Type", p.contentType)
	writer := csv.NewWriter(w)
	writer.Comma = p.comma
	return p.flush(writer, p.process(writer, dataModel))
}

func (p *csvProcessor) process(writer *csv.Writer, dataModel interface{}) error {
	switch v := dataModel.(type) {
	case string:
		return writer.Write([]string{v})
	case []string:
		return writer.Write(v)
	case [][]string:
		return writer.WriteAll(v)
	}

	s, ok := dataModel.(fmt.Stringer)
	if ok {
		return writer.Write([]string{s.String()})
	}

	value := reflect.Indirect(reflect.ValueOf(dataModel))

	switch value.Kind() {
	case reflect.Struct:
		if value.NumField() == 0 {
			return nil // nothing to write
		}

		return writeStructFields(writer, value, dataModel)

	case reflect.Array, reflect.Slice:
		if value.Len() == 0 {
			return nil // nothing to write
		}

		v0 := value.Index(0)
		k0 := v0.Kind()

		if reflect.Bool <= k0 && k0 <= reflect.Complex128 {
			return writeArrayOfScalars(writer, value)
		}

		switch k0 {
		//case reflect.Interface:
		//	fmt.Printf("----------- %v\n", v0)

		case reflect.Struct:
			if v0.NumField() == 0 {
				return nil // nothing to write
			}

			_, ok := v0.Interface().(fmt.Stringer)
			if ok {
				return writeArrayOfStringers(writer, value)
			}

			return writeArrayOfStructFields(writer, value, dataModel)

		case reflect.Array, reflect.Slice:
			if v0.Len() == 0 {
				return nil // nothing to write
			}

			v00 := v0.Index(0)
			k00 := v00.Kind()

			if reflect.Bool <= k00 && k00 <= reflect.Complex128 {
				return write2DArrayOfScalars(writer, value)

			} else if k00 == reflect.Struct {
				_, ok := v00.Interface().(fmt.Stringer)
				if ok {
					return write2DArrayOfStringers(writer, value)
				}
			}
		}
	}

	return fmt.Errorf("Unsupported type for CSV: %T", dataModel)
}

func writeArrayOfStructFields(writer *csv.Writer, value reflect.Value, dataModel interface{}) error {
	for j := 0; j < value.Len(); j++ {
		err := writeStructFields(writer, value.Index(j), dataModel)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeStructFields(writer *csv.Writer, str reflect.Value, dataModel interface{}) error {
	sa := make([]string, str.NumField())
	for i := 0; i < str.NumField(); i++ {
		sa[i] = fmt.Sprintf("%v", str.Field(i))
	}
	return writer.Write(sa)
}

func write2DArrayOfStringers(writer *csv.Writer, value reflect.Value) error {
	for j := 0; j < value.Len(); j++ {
		err := writeArrayOfStringers(writer, value.Index(j))
		if err != nil {
			return err
		}
	}
	return nil
}

func writeArrayOfStringers(writer *csv.Writer, value reflect.Value) error {
	sa := make([]string, value.Len())
	for i := 0; i < value.Len(); i++ {
		sa[i] = fmt.Sprintf("%v", value.Index(i).Interface().(fmt.Stringer))
	}
	return writer.Write(sa)
}

func write2DArrayOfScalars(writer *csv.Writer, value reflect.Value) error {
	for j := 0; j < value.Len(); j++ {
		err := writeArrayOfScalars(writer, value.Index(j))
		if err != nil {
			return err
		}
	}
	return nil
}

func writeArrayOfScalars(writer *csv.Writer, vj reflect.Value) error {
	sa := make([]string, vj.Len())
	for i := 0; i < vj.Len(); i++ {
		sa[i] = fmt.Sprintf("%v", vj.Index(i))
	}
	return writer.Write(sa)
}

func (p *csvProcessor) flush(writer *csv.Writer, err error) error {
	if err != nil {
		return err
	}
	writer.Flush()
	return writer.Error()
}
