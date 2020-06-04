package tool

type Attr struct {
	Name  string
	Value string
}

type AdvAttr struct {
	Name  string
	Value *Ch
}

/**
 * 字符
 */
type Ch struct {
	Value string
	Type  int
}

type RunElem struct {
	Type  string
	Name  string
	Value string
}
