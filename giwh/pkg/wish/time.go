package wish

import "time"

type Time struct {
	time.Time
}

func (t *Time) String() string {
	return t.Time.Format("2006-01-02 15:04:05")
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(t.Time.Format(`"2006-01-02 15:04:05"`)), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var err error
	t.Time, err = time.Parse(`"2006-01-02 15:04:05"`, string(data))
	return err
}
