package main

import "fmt"

type Data struct {
	Title string
}

func (d *Data) String() string {
	return d.Title
}

func main() {
	root := NewNode(Data{"cocoRoot"})
	coco := NewNode(Data{"coco"})
	xixi := NewNode(Data{"xixi"})

	root.AddChild(coco)
	root.AddChild(xixi)

	xixiA := NewNode(Data{"xixiA"})
	xixiB := NewNode(Data{"xixiB"})

	xixi.AddChild(xixiA)
	xixi.AddChild(xixiB)

	fmt.Println(root.IsRoot())

}
