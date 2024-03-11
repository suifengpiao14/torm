package torm_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/torm"
)

func TestGetSubTemplateNames(t *testing.T) {
	tplText := `
	{{define "parent"}}
		{{template "childrenLevel1"}}
	{{end}}
	{{define "childrenLevel1"}}
		{{template "childrenLevel2"}}
	{{end}}
	{{define "childrenLevel2"}}
		hello world
	{{end}}
	`

	temp := torm.NewTemplate()
	temp, err := temp.Parse(tplText)
	require.NoError(t, err)

	subTemplateNames, err := torm.GetSubTemplateNames(temp, "parent")
	require.NoError(t, err)
	fmt.Println(subTemplateNames)

}
