package str

type (
	HtmlAttributer interface {
		Register(html *Html)
	}

	AttrHtmlNormal struct{ content string }
	AttrHtmlA      struct {
		name, href string
		properties []HtmlProperty
	}
	AttrHtmlP   struct{ content string }
	AttrHtmlBr  struct{ content string }
	AttrHtmlUl  struct{ contents []string }
	AttrHtmlAny struct {
		tag, content string
		properties   []HtmlProperty
	}
	AttrHtmlTable struct {
		properties []HtmlProperty
		options    []HtmlAttributer
	}
	AttrHtmlTr struct {
		properties []HtmlProperty
		options    []HtmlAttributer
	}
	AttrHtmlTd struct {
		properties []HtmlProperty
		content    string
	}
	AttrHtmlTh struct {
		properties []HtmlProperty
		content    string
	}
	AttrHtmlTHead struct {
		properties []HtmlProperty
		options    []HtmlAttributer
	}
	AttrHtmlTBody struct {
		properties []HtmlProperty
		options    []HtmlAttributer
	}

	HtmlProperty struct {
		Key   string
		Value string
	}
)

func HtmlNormal(content string) AttrHtmlNormal {
	return AttrHtmlNormal{content: content}
}

func (my AttrHtmlNormal) Register(html *Html) {
	html.buffer.S(my.content)
}

func HtmlA(name, href string, properties ...HtmlProperty) AttrHtmlA {
	return AttrHtmlA{name: name, href: href, properties: properties}
}

func (my AttrHtmlA) Register(html *Html) {
	html.buffer.S(`<a href="`, my.href, `"`)
	if len(my.properties) > 0 {
		for idx := range my.properties {
			html.buffer.S(" ", my.properties[idx].Key, `="`, my.properties[idx].Value, `"`)
		}
	}
	html.buffer.S(">", my.name, "</a>")
}

func HtmlP(content string) AttrHtmlP {
	return AttrHtmlP{content: content}
}

func (my AttrHtmlP) Register(html *Html) {
	html.buffer.S("<p>").S(my.content).S("</p>")
}

func HtmlBr() AttrHtmlBr {
	return AttrHtmlBr{content: "<br />"}
}

func (my AttrHtmlBr) Register(html *Html) {
	html.buffer.S(my.content)
}

func HtmlUl(contents ...string) AttrHtmlUl {
	return AttrHtmlUl{contents: contents}
}

func (my AttrHtmlUl) Register(html *Html) {
	if len(my.contents) > 0 {
		html.buffer.S("<ul>")
		for idx := range my.contents {
			html.buffer.S("<li>", my.contents[idx], "</li>")
		}
		html.buffer.S("</ul>")
	}
}

func HtmlAny(tag, content string) AttrHtmlAny {
	return AttrHtmlAny{tag: tag, content: content, properties: []HtmlProperty{}}
}

func (my AttrHtmlAny) Register(html *Html) {
	html.buffer.S("<", my.tag)
	if len(my.properties) > 0 {
		for idx := range my.properties {
			html.buffer.S(" ", my.properties[idx].Key, `="`, my.properties[idx].Value, `"`)
		}
	}
	if my.content == "" {
		html.buffer.S(" />")
	} else {
		html.buffer.S(">", my.content, "</", my.tag, ">")
	}
}

func (my AttrHtmlAny) AppendProperties(properties ...HtmlProperty) AttrHtmlAny {
	my.properties = append(my.properties, properties...)
	return my
}

func HtmlTable(options ...HtmlAttributer) AttrHtmlTable {
	return AttrHtmlTable{options: options, properties: []HtmlProperty{}}
}

func (my AttrHtmlTable) Register(html *Html) {
	html.buffer.S("<table")
	if len(my.properties) > 0 {
		for idx := range my.properties {
			html.buffer.S(" ", my.properties[idx].Key, `="`, my.properties[idx].Value, `"`)
		}
	}
	html.buffer.S(">")
	if len(my.options) > 0 {
		for idx := range my.options {
			my.options[idx].Register(html)
		}
	}
	html.buffer.S("</table>")
}

func (my AttrHtmlTable) AppendProperties(properties ...HtmlProperty) AttrHtmlTable {
	my.properties = append(my.properties, properties...)
	return my
}

func HtmlTr(options ...HtmlAttributer) AttrHtmlTr {
	return AttrHtmlTr{options: options, properties: []HtmlProperty{}}
}

func (my AttrHtmlTr) Register(html *Html) {
	html.buffer.S("<tr")
	if len(my.properties) > 0 {
		for idx := range my.properties {
			html.buffer.S(" ", my.properties[idx].Key, `="`, my.properties[idx].Value, `"`)
		}
	}
	html.buffer.S(">")
	if len(my.options) > 0 {
		for idx := range my.options {
			my.options[idx].Register(html)
		}
	}
	html.buffer.S("</tr>")
}

func (my AttrHtmlTr) AppendProperties(properties ...HtmlProperty) AttrHtmlTr {
	my.properties = append(my.properties, properties...)
	return my
}

func HtmlTd(content string) AttrHtmlTd {
	return AttrHtmlTd{content: content, properties: []HtmlProperty{}}
}

func (my AttrHtmlTd) Register(html *Html) {
	html.buffer.S("<td")
	if len(my.properties) > 0 {
		for idx := range my.properties {
			html.buffer.S(" ", my.properties[idx].Key, `="`, my.properties[idx].Value, `"`)
		}
	}
	html.buffer.S(">", my.content, "</td>")
}

func (my AttrHtmlTd) AppendProperties(properties ...HtmlProperty) AttrHtmlTd {
	my.properties = append(my.properties, properties...)
	return my
}

func HtmlTh(content string) AttrHtmlTh {
	return AttrHtmlTh{content: content, properties: []HtmlProperty{}}
}

func (my AttrHtmlTh) Register(html *Html) {
	html.buffer.S("<th")
	if len(my.properties) > 0 {
		for idx := range my.properties {
			html.buffer.S(" ", my.properties[idx].Key, `="`, my.properties[idx].Value, `"`)
		}
	}
	html.buffer.S(">", my.content, "</th>")
}

func (my AttrHtmlTh) AppendProperties(properties ...HtmlProperty) AttrHtmlTh {
	my.properties = append(my.properties, properties...)
	return my
}

func HtmlTHead(options ...HtmlAttributer) AttrHtmlTHead {
	return AttrHtmlTHead{options: options, properties: []HtmlProperty{}}
}

func (my AttrHtmlTHead) Register(html *Html) {
	html.buffer.S("<thead")
	if len(my.properties) > 0 {
		for idx := range my.properties {
			html.buffer.S(" ", my.properties[idx].Key, `="`, my.properties[idx].Value, `"`)
		}
	}
	html.buffer.S(">")
	if len(my.options) > 0 {
		for idx := range my.options {
			my.options[idx].Register(html)
		}
	}
	html.buffer.S("</thead>")
}

func (my AttrHtmlTHead) AppendProperties(properties ...HtmlProperty) AttrHtmlTHead {
	my.properties = append(my.properties, properties...)
	return my
}

func HtmlTBody(options ...HtmlAttributer) AttrHtmlTBody {
	return AttrHtmlTBody{options: options, properties: []HtmlProperty{}}
}

func (my AttrHtmlTBody) Register(html *Html) {
	html.buffer.S("<tbody")
	if len(my.properties) > 0 {
		for idx := range my.properties {
			html.buffer.S(" ", my.properties[idx].Key, `="`, my.properties[idx].Value, `"`)
		}
	}
	html.buffer.S(">")
	if len(my.options) > 0 {
		for idx := range my.options {
			my.options[idx].Register(html)
		}
	}
	html.buffer.S("</tbody>")
}

func (my AttrHtmlTBody) AppendProperties(properties ...HtmlProperty) AttrHtmlTBody {
	my.properties = append(my.properties, properties...)
	return my
}
