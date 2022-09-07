package template

import (
	"bytes"
	"fmt"
	"math/rand"
	"text/template"
	"time"
)

// ParseRandName renders templates defined with {{.randName}} placeholders.
// This is useful for acceptance tests.
func ParseRandName(rawTemplate string) string {
	t := template.Must(template.New("tpl").Parse(rawTemplate))

	s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)

	var buf bytes.Buffer
	err := t.Execute(&buf, map[string]string{
		"randName": fmt.Sprintf("%d", r1.Intn(10000)),
	})
	if err != nil {
		panic(err)
	}

	return buf.String()
}
