package value

import "fmt"

type Value float64

func (v Value) String() string {
	return fmt.Sprintf("%.06f", float64(v))
}
