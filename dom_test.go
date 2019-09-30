package dom_test

import (
	"strings"
	"testing"

	"github.com/go-shiori/dom"
	"golang.org/x/net/html"
)

func TestGetElementsByTagName(t *testing.T) {
	htmlSource := `<div>
		<h1></h1>
		<h2></h2><h2></h2>
		<h3></h3><h3></h3><h3></h3>
		<p></p><p></p><p></p><p></p><p></p>
		<div></div><div></div><div></div><div></div><div></div>
		<div><p>Hey it's nested</p></div>
		<div></div>
		<img/><img/><img/><img/><img/><img/><img/><img/>
		<img/><img/><img/><img/>
	</div>`

	doc, err := parseHTMLSource(htmlSource)
	if err != nil {
		t.Errorf("GetElementsByTagName(), failed to parse: %v", err)
	}

	tests := map[string]int{
		"h1":  1,
		"h2":  2,
		"h3":  3,
		"p":   6,
		"div": 7,
		"img": 12,
		"*":   31,
	}

	for tagName, count := range tests {
		t.Run(tagName, func(t *testing.T) {
			if got := len(dom.GetElementsByTagName(doc, tagName)); got != count {
				t.Errorf("GetElementsByTagName() = %v, want %v", got, count)
			}
		})
	}
}

func TestCreateElement(t *testing.T) {
	tests := []struct {
		name     string
		tagName  string
		tagCount int
	}{{
		name:     "3 headings1",
		tagName:  "h1",
		tagCount: 3,
	}, {
		name:     "4 headings2",
		tagName:  "h2",
		tagCount: 4,
	}, {
		name:     "5 headings3",
		tagName:  "h3",
		tagCount: 5,
	}, {
		name:     "10 paragraph",
		tagName:  "p",
		tagCount: 10,
	}, {
		name:     "6 div",
		tagName:  "div",
		tagCount: 6,
	}, {
		name:     "8 image",
		tagName:  "img",
		tagCount: 8,
	}, {
		name:     "22 custom tag",
		tagName:  "custom-tag",
		tagCount: 22,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &html.Node{}
			for i := 0; i < tt.tagCount; i++ {
				doc.AppendChild(dom.CreateElement(tt.tagName))
			}

			if tags := dom.GetElementsByTagName(doc, tt.tagName); len(tags) != tt.tagCount {
				t.Errorf("CreateElement() = %v, want %v", len(tags), tt.tagCount)
			}
		})
	}
}

func TestCreateTextNode(t *testing.T) {
	tests := []string{
		"hello world",
		"this is awesome",
		"all cat is good boy",
		"all dog is good boy as well",
	}

	for _, text := range tests {
		t.Run(text, func(t *testing.T) {
			node := dom.CreateTextNode(text)
			if outerHTML := dom.OuterHTML(node); outerHTML != text {
				t.Errorf("CreateTextNode() = %v, want %v", outerHTML, text)
			}
		})
	}
}

func TestGetAttribute(t *testing.T) {
	tests := []struct {
		name       string
		htmlSource string
		attrName   string
		want       string
	}{{
		name:       "attr id from paragraph",
		htmlSource: `<p id="main-paragraph"></p>`,
		attrName:   "id",
		want:       "main-paragraph",
	}, {
		name:       "attr class from list",
		htmlSource: `<ul class="bullets"></ul>`,
		attrName:   "class",
		want:       "bullets",
	}, {
		name:       "attr style from paragraph",
		htmlSource: `<div style="display: none"></div>`,
		attrName:   "style",
		want:       "display: none",
	}, {
		name:       "attr doesn't exists",
		htmlSource: `<p id="main-paragraph"></p>`,
		attrName:   "class",
		want:       "",
	}, {
		name:       "node has no attributes",
		htmlSource: `<p></p>`,
		attrName:   "id",
		want:       "",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("GetAttribute(), failed to parse: %v", err)
			}

			if got := dom.GetAttribute(node, tt.attrName); got != tt.want {
				t.Errorf("GetAttribute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetAttribute(t *testing.T) {
	tests := []struct {
		name       string
		htmlSource string
		attrName   string
		attrValue  string
		want       string
	}{{
		name:       "set id of paragraph",
		htmlSource: `<p id="main-paragraph"></p>`,
		attrName:   "id",
		attrValue:  "txt-main",
		want:       `<p id="txt-main"></p>`,
	}, {
		name:       "set id from paragraph with several attrs",
		htmlSource: `<p id="main-paragraph" class="title"></p>`,
		attrName:   "id",
		attrValue:  "txt-main",
		want:       `<p id="txt-main" class="title"></p>`,
	}, {
		name:       "set new attr for paragraph",
		htmlSource: `<p></p>`,
		attrName:   "class",
		attrValue:  "title",
		want:       `<p class="title"></p>`,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("SetAttribute(), failed to parse: %v", err)
			}

			dom.SetAttribute(node, tt.attrName, tt.attrValue)
			if outerHTML := dom.OuterHTML(node); outerHTML != tt.want {
				t.Errorf("setAttribute() = %v, want %v", outerHTML, tt.want)
			}
		})
	}
}

func TestRemoveAttribute(t *testing.T) {
	tests := []struct {
		name       string
		htmlSource string
		attrName   string
		want       string
	}{{
		name:       "remove id of paragraph",
		htmlSource: `<p id="main-paragraph"></p>`,
		attrName:   "id",
		want:       `<p></p>`,
	}, {
		name:       "remove id from paragraph with several attrs",
		htmlSource: `<p id="main-paragraph" class="title"></p>`,
		attrName:   "id",
		want:       `<p class="title"></p>`,
	}, {
		name:       "remove inexist attr of paragraph",
		htmlSource: `<p id="main-paragraph"></p>`,
		attrName:   "class",
		want:       `<p id="main-paragraph"></p>`,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("RemoveAttribute(), failed to parse: %v", err)
			}

			dom.RemoveAttribute(node, tt.attrName)
			if outerHTML := dom.OuterHTML(node); outerHTML != tt.want {
				t.Errorf("RemoveAttribute() = %v, want %v", outerHTML, tt.want)
			}
		})
	}
}

func TestHasAttribute(t *testing.T) {
	tests := []struct {
		name       string
		htmlSource string
		attrName   string
		want       bool
	}{{
		name:       "attribute is exist",
		htmlSource: `<p id="main-paragraph"></p>`,
		attrName:   "id",
		want:       true,
	}, {
		name:       "attribute is not exist",
		htmlSource: `<p id="main-paragraph"></p>`,
		attrName:   "class",
		want:       false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("HasAttribute(), failed to parse: %v", err)
			}

			if got := dom.HasAttribute(node, tt.attrName); got != tt.want {
				t.Errorf("HasAttribute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTextContent(t *testing.T) {
	tests := []struct {
		name       string
		htmlSource string
		want       string
	}{{
		name:       "ordinary text node",
		htmlSource: "this is an ordinary text",
		want:       "this is an ordinary text",
	}, {
		name:       "single empty node element",
		htmlSource: "<p></p>",
		want:       "",
	}, {
		name:       "single node with content",
		htmlSource: "<p>Hello all</p>",
		want:       "Hello all",
	}, {
		name:       "single node with content and unnecessary space",
		htmlSource: "<p>Hello all   </p>",
		want:       "Hello all   ",
	}, {
		name:       "nested element",
		htmlSource: "<div><p>Some nested element</p></div>",
		want:       "Some nested element",
	}, {
		name:       "nested element with unnecessary space",
		htmlSource: "<div><p>Some nested element</p>    </div>",
		want:       "Some nested element    ",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("TextContent(), failed to parse: %v", err)
			}

			if got := dom.TextContent(node); got != tt.want {
				t.Errorf("TextContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOuterHTML(t *testing.T) {
	tests := []struct {
		name       string
		htmlSource string
	}{{
		name:       "text node",
		htmlSource: "this is an ordinary text",
	}, {
		name:       "single element",
		htmlSource: "<h1>Hello</h1>",
	}, {
		name:       "nested elements",
		htmlSource: "<div><p>Some nested element</p></div>",
	}, {
		name:       "triple nested elements",
		htmlSource: "<div><p>Some <a>nested</a> element</p></div>",
	}, {
		name:       "mixed nested elements",
		htmlSource: "<div><p>Some <a>nested</a> element</p><p>and more</p></div>",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("dom.OuterHTML(), failed to parse: %v", err)
			}

			if got := dom.OuterHTML(node); got != tt.htmlSource {
				t.Errorf("dom.OuterHTML() = %v, want %v", got, tt.htmlSource)
			}
		})
	}
}

func TestInnerHTML(t *testing.T) {
	tests := []struct {
		name       string
		htmlSource string
		want       string
	}{{
		name:       "text node",
		htmlSource: "this is an ordinary text",
		want:       "",
	}, {
		name:       "single element",
		htmlSource: "<h1>Hello</h1>",
		want:       "Hello",
	}, {
		name:       "nested elements",
		htmlSource: "<div><p>Some nested element</p></div>",
		want:       "<p>Some nested element</p>",
	}, {
		name:       "mixed text and element node",
		htmlSource: "<div><p>Some element</p>with text</div>",
		want:       "<p>Some element</p>with text",
	}, {
		name:       "triple nested elements",
		htmlSource: "<div><p>Some <a>nested</a> element</p></div>",
		want:       "<p>Some <a>nested</a> element</p>",
	}, {
		name:       "mixed nested elements",
		htmlSource: "<div><p>Some <a>nested</a> element</p><p>and more</p></div>",
		want:       "<p>Some <a>nested</a> element</p><p>and more</p>",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("InnerHTML(), failed to parse: %v", err)
			}

			if got := dom.InnerHTML(node); got != tt.want {
				t.Errorf("InnerHTML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestId(t *testing.T) {
	tests := []struct {
		name       string
		htmlSource string
		want       string
	}{{
		name:       "id exists",
		htmlSource: `<p id="main-paragraph"></p>`,
		want:       "main-paragraph",
	}, {
		name:       "id doesn't exist",
		htmlSource: `<p></p>`,
		want:       "",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("ID(), failed to parse: %v", err)
			}

			if got := dom.ID(node); got != tt.want {
				t.Errorf("ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClassName(t *testing.T) {
	tests := []struct {
		name       string
		htmlSource string
		want       string
	}{{
		name:       "class doesn't exist",
		htmlSource: `<p></p>`,
		want:       "",
	}, {
		name:       "class exist",
		htmlSource: `<p class="title"></p>`,
		want:       "title",
	}, {
		name:       "multiple class",
		htmlSource: `<p class="title heading"></p>`,
		want:       "title heading",
	}, {
		name:       "multiple class with unnecessary space",
		htmlSource: `<p class="    title heading    "></p>`,
		want:       "title heading",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("ClassName(), failed to parse: %v", err)
			}

			if got := dom.ClassName(node); got != tt.want {
				t.Errorf("ClassName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChildren(t *testing.T) {
	tests := []struct {
		name       string
		htmlSource string
		want       []string
	}{{
		name:       "has no children",
		htmlSource: "<div></div>",
		want:       []string{},
	}, {
		name:       "has one children",
		htmlSource: "<div><p>Hello</p></div>",
		want:       []string{"<p>Hello</p>"},
	}, {
		name:       "has many children",
		htmlSource: "<div><p>Hello</p><p>I'm</p><p>Happy</p></div>",
		want:       []string{"<p>Hello</p>", "<p>I&#39;m</p>", "<p>Happy</p>"},
	}, {
		name:       "has nested children",
		htmlSource: "<div><p>Hello I'm <span>Happy</span></p></div>",
		want:       []string{"<p>Hello I&#39;m <span>Happy</span></p>"},
	}, {
		name:       "mixed text and element node",
		htmlSource: "<div><p>Hello I'm</p>happy</div>",
		want:       []string{"<p>Hello I&#39;m</p>"},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("Children(), failed to parse: %v", err)
			}

			nodes := dom.Children(node)
			if len(nodes) != len(tt.want) {
				t.Errorf("Children() count = %v, want = %v", len(nodes), len(tt.want))
			}

			for i, child := range nodes {
				wantHTML := tt.want[i]
				childHTML := dom.OuterHTML(child)
				if childHTML != wantHTML {
					t.Errorf("Children() = %v, want = %v", childHTML, wantHTML)
				}
			}
		})
	}
}

func TestChildNodes(t *testing.T) {
	tests := []struct {
		name       string
		htmlSource string
		want       []string
	}{{
		name:       "has no children",
		htmlSource: "<div></div>",
		want:       []string{},
	}, {
		name:       "has one children",
		htmlSource: "<div><p>Hello</p></div>",
		want:       []string{"<p>Hello</p>"},
	}, {
		name:       "has many children",
		htmlSource: "<div><p>Hello</p><p>I'm</p><p>Happy</p></div>",
		want:       []string{"<p>Hello</p>", "<p>I&#39;m</p>", "<p>Happy</p>"},
	}, {
		name:       "has nested children",
		htmlSource: "<div><p>Hello I'm <span>Happy</span></p></div>",
		want:       []string{"<p>Hello I&#39;m <span>Happy</span></p>"},
	}, {
		name:       "mixed text and element node",
		htmlSource: "<div><p>Hello I'm</p>happy</div>",
		want:       []string{"<p>Hello I&#39;m</p>", "happy"},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("ChildNodes(), failed to parse: %v", err)
			}

			nodes := dom.ChildNodes(node)
			if len(nodes) != len(tt.want) {
				t.Errorf("ChildNodes() count = %v, want = %v", len(nodes), len(tt.want))
			}

			for i, child := range nodes {
				wantHTML := tt.want[i]
				childHTML := dom.OuterHTML(child)
				if child.Type == html.TextNode {
					childHTML = dom.TextContent(child)
				}

				if childHTML != wantHTML {
					t.Errorf("ChildNodes() = %v, want = %v", childHTML, wantHTML)
				}
			}
		})
	}
}

func TestFirstElementChild(t *testing.T) {
	tests := []struct {
		name       string
		htmlSource string
		want       string
	}{{
		name:       "has no children",
		htmlSource: "<div></div>",
		want:       "",
	}, {
		name:       "has one children",
		htmlSource: "<div><p>Hey</p></div>",
		want:       "<p>Hey</p>",
	}, {
		name:       "has several children",
		htmlSource: "<div><p>Hey</p><b>bro</b></div>",
		want:       "<p>Hey</p>",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("FirstElementChild(), failed to parse: %v", err)
			}

			firstChild := dom.FirstElementChild(node)
			if got := dom.OuterHTML(firstChild); got != tt.want {
				t.Errorf("FirstElementChild() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNextElementSibling(t *testing.T) {
	tests := []struct {
		name       string
		htmlSource string
		want       string
	}{{
		name:       "has no sibling",
		htmlSource: "<div></div>",
		want:       "",
	}, {
		name:       "has directly element sibling",
		htmlSource: "<div></div><p>Hey</p>",
		want:       "<p>Hey</p>",
	}, {
		name:       "has no element sibling",
		htmlSource: "<div></div>I'm your sibling, you know",
		want:       "",
	}, {
		name:       "has distant element sibling",
		htmlSource: "<div></div>I'm your sibling as well <p>only me matter</p>",
		want:       "<p>only me matter</p>",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("NextElementSibling(), failed to parse: %v", err)
			}

			nextSibling := dom.NextElementSibling(node)
			if got := dom.OuterHTML(nextSibling); got != tt.want {
				t.Errorf("NextElementSibling() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppendChild(t *testing.T) {
	// Child is from inside document
	t.Run("child from existing node", func(t *testing.T) {
		htmlSource := `<div><p>Lonely word</p><span>new friend</span></div>`
		want := `<div><p>Lonely word<span>new friend</span></p></div>`

		doc, err := parseHTMLSource(htmlSource)
		if err != nil {
			t.Errorf("AppendChild(), failed to parse: %v", err)
		}

		p := dom.GetElementsByTagName(doc, "p")[0]
		span := dom.GetElementsByTagName(doc, "span")[0]

		dom.AppendChild(p, span)
		if got := dom.OuterHTML(doc); got != want {
			t.Errorf("AppendChild() = %v, want %v", got, want)
		}
	})

	// Child is new element
	t.Run("child is new element", func(t *testing.T) {
		htmlSource := `<div><p>Lonely word</p><span>new friend</span></div>`
		want := `<div><p>Lonely word<span></span></p><span>new friend</span></div>`

		doc, err := parseHTMLSource(htmlSource)
		if err != nil {
			t.Errorf("AppendChild(), failed to parse: %v", err)
		}

		p := dom.GetElementsByTagName(doc, "p")[0]
		newChild := dom.CreateElement("span")

		dom.AppendChild(p, newChild)
		if got := dom.OuterHTML(doc); got != want {
			t.Errorf("AppendChild() = %v, want %v", got, want)
		}
	})
}

func TestPrependChild(t *testing.T) {
	// Child is from inside document
	t.Run("child from existing node", func(t *testing.T) {
		htmlSource := `<div><p>Lonely word</p><span>new friend</span></div>`
		want := `<div><p><span>new friend</span>Lonely word</p></div>`

		doc, err := parseHTMLSource(htmlSource)
		if err != nil {
			t.Errorf("PrependChild(), failed to parse: %v", err)
		}

		p := dom.GetElementsByTagName(doc, "p")[0]
		span := dom.GetElementsByTagName(doc, "span")[0]

		dom.PrependChild(p, span)
		if got := dom.OuterHTML(doc); got != want {
			t.Errorf("PrependChild() = %v, want %v", got, want)
		}
	})

	// Child is new element
	t.Run("child is new element", func(t *testing.T) {
		htmlSource := `<div><p>Lonely word</p><span>new friend</span></div>`
		want := `<div><p><span></span>Lonely word</p><span>new friend</span></div>`

		doc, err := parseHTMLSource(htmlSource)
		if err != nil {
			t.Errorf("PrependChild(), failed to parse: %v", err)
		}

		p := dom.GetElementsByTagName(doc, "p")[0]
		newChild := dom.CreateElement("span")

		dom.PrependChild(p, newChild)
		if got := dom.OuterHTML(doc); got != want {
			t.Errorf("PrependChild() = %v, want %v", got, want)
		}
	})
}

func TestReplaceChild(t *testing.T) {
	// new child is from existing element
	t.Run("new child from existing element", func(t *testing.T) {
		htmlSource := `<div><p>Lonely word</p><span>new friend</span></div>`
		want := `<div><span>new friend</span></div>`

		doc, err := parseHTMLSource(htmlSource)
		if err != nil {
			t.Errorf("ReplaceNode(), failed to parse: %v", err)
		}

		p := dom.GetElementsByTagName(doc, "p")[0]
		span := dom.GetElementsByTagName(doc, "span")[0]

		dom.ReplaceChild(doc, span, p)
		if got := dom.OuterHTML(doc); got != want {
			t.Errorf("ReplaceNode() = %v, want %v", got, want)
		}
	})

	// new child is new element
	t.Run("new node is new element", func(t *testing.T) {
		htmlSource := `<div><p>Lonely word</p><span>new friend</span></div>`
		want := `<div><span></span><span>new friend</span></div>`

		doc, err := parseHTMLSource(htmlSource)
		if err != nil {
			t.Errorf("ReplaceNode(), failed to parse: %v", err)
		}

		p := dom.GetElementsByTagName(doc, "p")[0]
		newChild := dom.CreateElement("span")

		dom.ReplaceChild(doc, newChild, p)
		if got := dom.OuterHTML(doc); got != want {
			t.Errorf("ReplaceNode() = %v, want %v", got, want)
		}
	})
}

func TestIncludeNode(t *testing.T) {
	htmlSource := `<div>
		<h1></h1><h2></h2><h3></h3>
		<p></p><div></div><img/><img/>
	</div>`

	doc, err := parseHTMLSource(htmlSource)
	if err != nil {
		t.Errorf("IncludeNode(), failed to parse: %v", err)
	}

	allElements := dom.GetElementsByTagName(doc, "*")
	h1 := dom.GetElementsByTagName(doc, "h1")[0]
	h2 := dom.GetElementsByTagName(doc, "h2")[0]
	h3 := dom.GetElementsByTagName(doc, "h3")[0]
	p := dom.GetElementsByTagName(doc, "p")[0]
	div := dom.GetElementsByTagName(doc, "div")[0]
	img := dom.GetElementsByTagName(doc, "img")[0]
	span := dom.CreateElement("span")

	tests := []struct {
		name string
		node *html.Node
		want bool
	}{
		{"h1", h1, true},
		{"h2", h2, true},
		{"h3", h3, true},
		{"p", p, true},
		{"div", div, true},
		{"img", img, true},
		{"span", span, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dom.IncludeNode(allElements, tt.node); got != tt.want {
				t.Errorf("IncludeNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCloneNode(t *testing.T) {
	tests := []struct {
		name       string
		htmlSource string
		want       string
	}{{
		name:       "single div",
		htmlSource: "<div></div>",
	}, {
		name:       "div with one children",
		htmlSource: "<div><p>Hello</p></div>",
	}, {
		name:       "div with many children",
		htmlSource: "<div><p>Hello</p><p>I'm</p><p>Happy</p></div>",
		want:       "<div><p>Hello</p><p>I&#39;m</p><p>Happy</p></div>",
	}, {
		name:       "div with nested children",
		htmlSource: "<div><p>Hello I'm <span>Happy</span></p></div>",
		want:       "<div><p>Hello I&#39;m <span>Happy</span></p></div>",
	}, {
		name:       "div with mixed text and element node",
		htmlSource: "<div><p>Hello I'm</p>happy</div>",
		want:       "<div><p>Hello I&#39;m</p>happy</div>",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := tt.want
			if want == "" {
				want = tt.htmlSource
			}

			node, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("CloneNode(), failed to parse: %v", err)
			}

			clone := dom.CloneNode(node)
			if got := dom.OuterHTML(clone); got != want {
				t.Errorf("CloneNode() = %v, want %v", got, want)
			}
		})
	}
}

func TestGetAllNodesWithTag(t *testing.T) {
	htmlSource := `<div>
		<h1></h1>
		<h2></h2><h2></h2>
		<h3></h3><h3></h3><h3></h3>
		<p></p><p></p><p></p><p></p><p></p>
		<div></div><div></div><div></div><div></div><div></div>
		<div><p>Hey it's nested</p></div>
		<div></div>
		<img/><img/><img/><img/><img/><img/><img/><img/>
		<img/><img/><img/><img/>
	</div>`

	doc, err := parseHTMLSource(htmlSource)
	if err != nil {
		t.Errorf("GetAllNodesWithTag(), failed to parse: %v", err)
	}

	tests := []struct {
		name string
		tags []string
		want int
	}{{
		name: "h1",
		tags: []string{"h1"},
		want: 1,
	}, {
		name: "h1,h2",
		tags: []string{"h1", "h2"},
		want: 3,
	}, {
		name: "h1,h2,h3",
		tags: []string{"h1", "h2", "h3"},
		want: 6,
	}, {
		name: "p",
		tags: []string{"p"},
		want: 6,
	}, {
		name: "p,span",
		tags: []string{"p", "span"},
		want: 6,
	}, {
		name: "div,img",
		tags: []string{"div", "img"},
		want: 19,
	}, {
		name: "span",
		tags: []string{"span"},
		want: 0,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(dom.GetAllNodesWithTag(doc, tt.tags...)); got != tt.want {
				t.Errorf("GetAllNodesWithTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveNodes(t *testing.T) {
	htmlSource := `<div><h1></h1><h1></h1><p></p><img/></div>`

	tests := []struct {
		name   string
		want   string
		filter func(*html.Node) bool
	}{{
		name:   "remove all",
		want:   "<div></div>",
		filter: nil,
	}, {
		name: "remove one tag",
		want: "<div><p></p><img/></div>",
		filter: func(n *html.Node) bool {
			return dom.TagName(n) == "h1"
		},
	}, {
		name: "remove several tags",
		want: "<div><img/></div>",
		filter: func(n *html.Node) bool {
			tag := dom.TagName(n)
			return tag == "h1" || tag == "p"
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := parseHTMLSource(htmlSource)
			if err != nil {
				t.Errorf("RemoveNodes(), failed to parse: %v", err)
			}

			elements := dom.GetElementsByTagName(doc, "*")
			dom.RemoveNodes(elements, tt.filter)

			if got := dom.OuterHTML(doc); got != tt.want {
				t.Errorf("RemoveNodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetTextContent(t *testing.T) {
	textContent := "XXX"
	expectedResult := "<div>" + textContent + "</div>"

	tests := []struct {
		name       string
		htmlSource string
	}{{
		name:       "single div",
		htmlSource: "<div></div>",
	}, {
		name:       "div with one children",
		htmlSource: "<div><p>Hello</p></div>",
	}, {
		name:       "div with many children",
		htmlSource: "<div><p>Hello</p><p>I'm</p><p>Happy</p></div>",
	}, {
		name:       "div with nested children",
		htmlSource: "<div><p>Hello I'm <span>Happy</span></p></div>",
	}, {
		name:       "div with mixed text and element node",
		htmlSource: "<div><p>Hello I'm</p>happy</div>",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("SetTextContent(), failed to parse: %v", err)
			}

			dom.SetTextContent(root, textContent)
			if got := dom.OuterHTML(root); got != expectedResult {
				t.Errorf("SetTextContent() = %v, want %v", got, expectedResult)
			}
		})
	}
}

func parseHTMLSource(htmlSource string) (*html.Node, error) {
	doc, err := html.Parse(strings.NewReader(htmlSource))
	if err != nil {
		return nil, err
	}

	body := dom.GetElementsByTagName(doc, "body")[0]
	return body.FirstChild, nil
}
