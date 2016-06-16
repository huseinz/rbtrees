package main

import (
	"fmt"
	"strings"
	"bufio"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type RBTree struct {
	root *RBNode
}

type RBNode struct {
	data                int
	left, right, parent *RBNode
	red                 bool
}

func newtree() *RBTree {
	return &RBTree{}
}

func leftrotate(t *RBNode, x *RBNode) {

	y := x.right
	x.right = y.left

	if y.left != nil {
		y.left.parent = x
	}

	y.parent = x.parent

	if x.parent == nil {
		t = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}

	y.left = x
	x.parent = y
}

func rightrotate(t *RBNode, x *RBNode) {

	y := x.left
	x.left = y.right

	if y.right != nil {
		y.right.parent = x
	}

	y.parent = x.parent

	if x.parent == nil {
		t = y
	} else if x == x.parent.right {
		x.parent.right = y
	} else {
		x.parent.left = y
	}

	y.right = x
	x.parent = y
}

func RBTreeinsert(t *RBTree, data int) {

	x := &RBNode{data, nil, nil, nil, true}
	treeinsert(t, x)

/*
	for x != t.root && x.parent != nil && x.parent.red && x.red {

		parentpt := x.parent
		gparentpt := x.parent.parent

		if gparentpt != nil && parentpt == gparentpt.left{

			unclept := gparentpt.right

			if unclept != nil && unclept.red {
				gparentpt.red = true
				parentpt.red = false
				unclept.red = false
				*x = *gparentpt
			} else {
				if x == parentpt.right {
					leftrotate(t.root, parentpt)
					x = parentpt
					parentpt = x.parent
				}
				rightrotate(t.root, gparentpt)
				tmp := parentpt.red
				parentpt.red = gparentpt.red
				gparentpt.red = tmp
				x = parentpt
			}
		} else if gparentpt != nil{

			unclept := gparentpt.left

			if unclept != nil && unclept.red {
				gparentpt.red = true
				parentpt.red = false
				unclept.red = false
				x = gparentpt
			} else {
				if x == parentpt.left {
					rightrotate(t.root, parentpt)
					x = parentpt
					parentpt = x.parent
				}
				leftrotate(t.root, gparentpt)
				tmp := parentpt.red
				parentpt.red = gparentpt.red
				gparentpt.red = tmp
				x = parentpt
			}
		}
	}
	t.root.red = false
	*/
}

func treeinsert(tree *RBTree, x *RBNode) {

	if tree.root == nil {
		tree.root = x
		return
	}

	t := tree.root

	if x.data == t.data {
		return //perhaps return a value
	}

	if x.data > t.data {
		if t.right == nil {
			t.right = x
			x.parent = t
			return
		}
		treeinsert(&RBTree{t.right}, x)
	} else {
		if t.left == nil {
			t.left = x
			x.parent = t
			return
		}
		treeinsert(&RBTree{t.left}, x)
	}
}

func printinorder(t *RBNode) {

	if t == nil {
		return
	}

	printinorder(t.left)
	fmt.Printf("%v ", t.data)
	if t.red {
		fmt.Println("red")
	} else {
		fmt.Println("black")
	}
	printinorder(t.right)
}

func graph(t *RBNode, fn string) {

	if t == nil {
		fmt.Println("tree is empty")
		return
	}
	f, err := os.Create(fn + ".dot")
	if err != nil {
		fmt.Println("failed to open file\n")
	}

	f.WriteString("digraph G {\n\n")
	f.WriteString("\tnode [shape=circle, style=filled, fontcolor=white]\n")
	f.WriteString("nodesep=0.4;\n")

	helper_graph(t, f)

	f.WriteString("}\n")
	f.Close()

	_, err = exec.LookPath("dot")
	if err != nil {
		fmt.Println("dot could not be found")
		log.Fatal(err)
	}
	err = exec.Command("dot", "-Tpdf", fn+".dot", "-o", fn+".pdf").Run()
}

func helper_graph(root *RBNode, f *os.File) {

	label := strconv.Itoa(root.data)
	if root.red { //easier to read
		f.WriteString("\t" + label + " [fillcolor=red]\n")
	} else {
		f.WriteString("\t" + label + " [fillcolor=black]\n")
	}

	if root.left != nil {
		f.WriteString("\t" + label + " -> " + strconv.Itoa(root.left.data) + "\n")
		helper_graph(root.left, f)
	} else {
		f.WriteString("\t" + label + " -> nill" + label + " [style=invis]\n")
		f.WriteString("\tnill" + label + " [color=white, fillcolor=white]\n")
	}

	if root.right != nil {
		f.WriteString("\t" + label + " -> " + strconv.Itoa(root.right.data) + "\n")
		helper_graph(root.right, f)
	} else {
		f.WriteString("\t" + label + " -> nilr" + label + " [style=invis]\n")
		f.WriteString("\tnilr" + label + " [color=white, fillcolor=white]\n")
	}
}

func randomtree(n int) *RBTree {

	t := newtree()
	rand.Seed(time.Now().UTC().UnixNano())
	elems := rand.Perm(n)

	for _,k := range elems {
		RBTreeinsert(t, k)
	}

	return t
}

func search(t *RBNode, n int) (*RBNode) {

	switch {
	case t == nil:
		return nil
	case n < t.data:
		return search(t.left, n)
	case n > t.data:
		return search(t.right, n)
	default:
		return t
	}
}


func main() {

	t := newtree()

	fmt.Println("type 'h' to see list of commands")
	for true{

		fmt.Print("% ")

		r := bufio.NewReader(os.Stdin)
		cmdstr, err := r.ReadString('\n')
		if err != nil {
			return
		}

		cmd := strings.Fields(cmdstr)

		switch cmd[0] {
		case "r":
			if len(cmd) != 2 {
				fmt.Println("Usage: r <n>")
				continue
			}
			n, nerr := strconv.Atoi(cmd[1])
			if nerr != nil {
				fmt.Println("Usage: r <n>")
				continue
			}
			t = randomtree(n)
		case "s":
			if len(cmd) < 2 {
				fmt.Println("Usage: s <n1 n2 ...>")
				continue
			}
			for _,tok := range cmd[1:]{
				n, nerr := strconv.Atoi(tok)
				if nerr != nil {
					fmt.Println("bad token: " + tok)
					continue
				}
				if search(t.root, n) != nil{
					fmt.Println(tok + " : yes")
				} else {
					fmt.Println(tok + " : no")
				}
			}
		case "p":
			if t.root == nil{
				fmt.Println("tree is empty")
			} else {
				printinorder(t.root)
			}
		case "g":
			if len(cmd) != 2 {
				fmt.Println("Usage: g <name>")
				continue
			}
			graph(t.root, cmd[1])
		case "i":
			if len(cmd) < 2 {
				fmt.Println("Usage: i <n1 n2 ...>")
				continue
			}
			for _,tok := range cmd[1:]{
				n, nerr := strconv.Atoi(tok)
				if nerr != nil {
					fmt.Println("bad token: " + tok)
					continue
				}
				RBTreeinsert(t, n)
			}
		case "h":
			fmt.Println("r <n>\t\t: generate random RB-tree with n nodes with values from 0-n")
			fmt.Println("s <n1, n2, ...>\t: determines if value(s) are present in the tree")
			fmt.Println("i <n1, n2, ...>\t: insert value(s) into RB-tree")
			fmt.Println("d <n1, n2, ...>\t: delete value(s) from RB-tree")
			fmt.Println("g <name>\t: draw RB-tree to PDF and dot files")
			fmt.Println("p\t\t: print tree in-order")
			fmt.Println("h\t\t: print this message")
			fmt.Println("q\t\t: exit program")
		case "q":
			return
		case "lr":
			n, nerr := strconv.Atoi(cmd[1])
			if nerr == nil{
				leftrotate(t.root, search(t.root, n))
			}
		case "rr":
			n, nerr := strconv.Atoi(cmd[1])
			if nerr == nil{
				rightrotate(t.root, search(t.root, n))
			}
		default:
			fmt.Println("unrecognized command: " + cmd[0])
			fmt.Println("type 'h' for list of commands")
		}
	}
}
