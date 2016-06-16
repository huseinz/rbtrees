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

var niln *RBNode = &RBNode{red : false}

func newtree() *RBTree {
	return &RBTree{}
}

func leftrotate(t *RBNode, x *RBNode) {

	y := x.right
	x.right = y.left

	if y.left != niln {
		y.left.parent = x
	}

	y.parent = x.parent

	if x.parent == niln {
		*t = *y
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

	if y.right != niln {
		y.right.parent = x
	}

	y.parent = x.parent

	if x.parent == niln {
		*t = *y
	} else if x == x.parent.right {
		x.parent.right = y
	} else {
		x.parent.left = y
	}

	y.right = x
	x.parent = y
}

func RBTreeinsert(t *RBTree, data int) {

	z := &RBNode{data, niln, niln, niln, true}
	treeinsert(t, z)

	for z.parent.red {
		if z.parent == z.parent.parent.left{
			unclept := z.parent.parent.right
			if unclept.red {
				z.parent.red = false
				unclept.red = false
				z.parent.parent.red = true
				z = z.parent.parent
			} else {
				if z == z.parent.right{
					z = z.parent
					leftrotate(t.root, z)
				}
				z.parent.red = false
				z.parent.parent.red = true
				rightrotate(t.root, z.parent.parent)
			}
		} else {
			unclept := z.parent.parent.left
			if unclept.red {
				z.parent.red = false
				unclept.red = false
				z.parent.parent.red = true
				z = z.parent.parent
			} else {
				if z == z.parent.left{
					z = z.parent
					rightrotate(t.root, z)
				}
				z.parent.red = false
				z.parent.parent.red = true
				leftrotate(t.root, z.parent.parent)
			}
		}
	}

	t.root.red = false
}

func treeinsert(tree *RBTree, x *RBNode) {

	if tree.root == nil {
		tree.root = x
		tree.root.red = false
		return
	}

	t := tree.root

	if x.data == t.data {
		return //perhaps return a value
	}

	if x.data > t.data {
		if t.right == niln {
			t.right = x
			x.parent = t
			return
		}
		treeinsert(&RBTree{t.right}, x)
	} else {
		if t.left == niln {
			t.left = x
			x.parent = t
			return
		}
		treeinsert(&RBTree{t.left}, x)
	}
}

func printinorder(t *RBNode) {

	if t == niln {
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

	if t == niln {
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

	if root.left != niln {
		f.WriteString("\t" + label + " -> " + strconv.Itoa(root.left.data) + "\n")
		helper_graph(root.left, f)
	} else {
		f.WriteString("\t" + label + " -> nill" + label + " [style=invis]\n")
		f.WriteString("\tnill" + label + " [color=white, fillcolor=white]\n")
	}

	if root.right != niln {
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
	case t == niln:
		return niln
	case n < t.data:
		return search(t.left, n)
	case n > t.data:
		return search(t.right, n)
	default:
		return t
	}
}


func main() {

	niln.left = niln
	niln.right = niln
	niln.parent = niln
	niln.red = false
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
				if search(t.root, n) != niln{
					fmt.Println(tok + " : yes")
				} else {
					fmt.Println(tok + " : no")
				}
			}
		case "p":
			if t.root == niln{
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
