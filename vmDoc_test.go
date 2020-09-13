package ottogoquery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

func Test_f(t *testing.T) {
	buf, _ := ioutil.ReadFile("1.html")
	vdoc, err := NewVMDocFromReader(bytes.NewReader(buf))
	if err != nil {
		log.Panicln(err)
	}
	jsCommand := `
		f=docExec("Find","#content > div > div.indexs > h2 > a")
		len = f("Length")
		r = []
		for (i=0;i<len;i++){
			f1=f("clone")
			f1("Eq",i)
			r.push(f1("Text"))
		}
		r = r
		`
	r, err := vdoc.Run(jsCommand)
	if err != nil {
		log.Panicln(err)
	}
	v, _ := r.Export()
	buf2, err := json.Marshal(v)
	log.Println(err)
	fmt.Printf("%s", buf2)

}
