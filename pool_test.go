package gost

import (
	"fmt"
	"testing"
)

func TestPool(t *testing.T) {
	p := newPool(2, "/home/julien/Applications/phantomjs")
	w, err := p.Get()
	if err != nil {
		t.Fatal(err)
	}

	err = w.GetDriver().Get("http://www.twenga.fr/")
	if err != nil {
		t.Fatal(err)
	}

	var title string
	title, err = w.GetDriver().Title()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("The title of the page is %q", title)
	err = w.GetDriver().Close()
	if err != nil {
		t.Fatal(err)
	}
}
