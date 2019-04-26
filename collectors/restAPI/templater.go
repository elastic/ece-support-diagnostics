package restAPI

import (
	"bytes"
	"html/template"
	"log"
)

// templateService controls the fields to be templated
func (R *Rest) templateService(Obj interface{}) {
	R.Filename = runTemplate(R.Filename, Obj)
	R.URI = runTemplate(R.URI, Obj)
}

// runTemplate performs the string substitution using the html/template package
func runTemplate(item string, Obj interface{}) string {
	t := template.Must(template.New("testing").Parse(item))
	var tpl bytes.Buffer
	err := t.Execute(&tpl, Obj)
	if err != nil {
		log.Println("executing template:", Obj)
	}
	return tpl.String()
}
