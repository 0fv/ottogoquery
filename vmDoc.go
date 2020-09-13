package ottogoquery

import (
	"io"

	"github.com/PuerkitoBio/goquery"
	"github.com/robertkrimen/otto"
)

type vmDoc struct {
	doc *goquery.Document
	*otto.Otto
}

func newvmDoc(doc *goquery.Document, vm *otto.Otto) *vmDoc {
	functionInit(vm, doc)
	vmDoc := vmDoc{}
	vmDoc.doc = doc
	vmDoc.Otto = vm
	return &vmDoc
}

//NewVMDocFromReader new otto vm and
func NewVMDocFromReader(r io.Reader) (*vmDoc, error) {
	v := otto.New()
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	return newvmDoc(doc, v), nil
}
