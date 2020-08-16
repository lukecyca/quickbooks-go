package quickbooks

import (
	"fmt"
	"time"
)

const format = "2006-01-02"

// Date represents a Quickbooks date
type Date struct {
	time.Time
}

// UnmarshalJSON removes time from parsed date
func (d *Date) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}

	// Strip time off
	b = b[:10]

	d.Time, err = time.Parse(format, string(b))
	if err != nil {
		return fmt.Errorf("Could not parse Date: %s", err.Error())
	}
	return nil
}

func (d Date) String() string {
	return d.Format(format)
}

func (d *Date) MarshalJSON() ([]byte, error) {
	return []byte("\"" + d.String() + "\""), nil
}
