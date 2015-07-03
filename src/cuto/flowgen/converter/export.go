package converter

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

// Export writes xml to w. Contents of xml is read from d(Definitions object).
func Export(w io.Writer, d *Definitions) error {
	output, err := xml.MarshalIndent(d, "", "    ")
	if err != nil {
		return fmt.Errorf("XML marshal error: %v", err)
	}

	w.Write([]byte(xml.Header))
	w.Write(output)
	return nil
}

// ExportFile writes xml to file.
func ExportFile(filepath string, d *Definitions) error {
	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return fmt.Errorf("Output file open error: %s", err)
	}

	return Export(f, d)
}

// GenerateDefinitions generates BPMN Difinitions object from head element of flow.
func GenerateDefinitions(head Element) *Definitions {
	d := NewDefinitions()

	pre := d.Process.Start.ID
	for current := head; current != nil; current = current.Next() {
		switch current.(type) {
		case *Job:
			pre = d.AppendJob(current.(*Job), pre)
		case *Gateway:
			pre = d.AppendGateway(current.(*Gateway), pre)
		default:
			panic("Unexpected type detected.")
		}
	}

	d.AppendSequenceFlow(NewSequenceFlow(pre, d.Process.End.ID))

	return d
}
