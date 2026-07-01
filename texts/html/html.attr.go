package html

import "github.com/spf13/cast"

type (
	HTMLAttr func(h HTMLWriter)

	HTMLProperty struct {
		Key   string
		Value string
	}
)

func Normal(content string) HTMLAttr { return func(h HTMLWriter) { h.GetBuffer().S(content) } }

func A(name, href string, properties ...HTMLProperty) HTMLAttr {
	return func(h HTMLWriter) {
		h.GetBuffer().S(`<a href="`, href, `"`)
		if len(properties) > 0 {
			for idx := range properties {
				h.GetBuffer().S(" ", properties[idx].Key, `="`, properties[idx].Value, `"`)
			}
		}
		h.GetBuffer().S(">", name, "</a>")
	}
}

func P(content string) HTMLAttr { return func(h HTMLWriter) { h.GetBuffer().S("<p>", content, "</p>") } }

func Br() HTMLAttr { return func(h HTMLWriter) { h.GetBuffer().S("<br />") } }

func Ul(contents ...string) HTMLAttr {
	return func(h HTMLWriter) {
		h.GetBuffer().S("<ul>")
		for idx := range contents {
			h.GetBuffer().S("<li>", contents[idx], "</li>")
		}
		h.GetBuffer().S("</ul>")
	}
}

func Any(tag, content string, properties ...HTMLProperty) HTMLAttr {
	return func(h HTMLWriter) {
		h.GetBuffer().S("<", tag)
		if len(properties) > 0 {
			for idx := range properties {
				h.GetBuffer().S(" ", properties[idx].Key, `="`, properties[idx].Value, `"`)
			}
		}
		if content == "" {
			h.GetBuffer().S(" />")
		} else {
			h.GetBuffer().S(">", content, "</", tag, ">")
		}
	}
}

func Table(options []HTMLAttr, properties ...HTMLProperty) HTMLAttr {
	return func(h HTMLWriter) {
		h.GetBuffer().S("<table")
		if len(properties) > 0 {
			for idx := range properties {
				h.GetBuffer().S(" ", properties[idx].Key, `="`, properties[idx].Value, `"`)
			}
		}
		h.GetBuffer().S(">")
		if len(options) > 0 {
			for idx := range options {
				options[idx](h)
			}
		}
		h.GetBuffer().S("</table>")
	}
}

func Tr(options []HTMLAttr, properties ...HTMLProperty) HTMLAttr {
	return func(h HTMLWriter) {
		h.GetBuffer().S("<tr")
		if len(properties) > 0 {
			for idx := range properties {
				h.GetBuffer().S(" ", properties[idx].Key, `="`, properties[idx].Value, `"`)
			}
		}
		h.GetBuffer().S(">")
		if len(options) > 0 {
			for idx := range options {
				options[idx](h)
			}
		}
		h.GetBuffer().S("</tr>")
	}
}

func Td(content string, properties ...HTMLProperty) HTMLAttr {
	return func(h HTMLWriter) {
		h.GetBuffer().S("<td")
		if len(properties) > 0 {
			for idx := range properties {
				h.GetBuffer().S(" ", properties[idx].Key, `="`, properties[idx].Value, `"`)
			}
		}
		h.GetBuffer().S(">", content, "</td>")
	}
}

func Th(content string, properties ...HTMLProperty) HTMLAttr {
	return func(h HTMLWriter) {
		h.GetBuffer().S("<th")
		if len(properties) > 0 {
			for idx := range properties {
				h.GetBuffer().S(" ", properties[idx].Key, `="`, properties[idx].Value, `"`)
			}
		}
		h.GetBuffer().S(">", content, "</th>")
	}
}

func THead(options []HTMLAttr, properties ...HTMLProperty) HTMLAttr {
	return func(h HTMLWriter) {
		h.GetBuffer().S("<thead")
		if len(properties) > 0 {
			for idx := range properties {
				h.GetBuffer().S(" ", properties[idx].Key, `="`, properties[idx].Value, `"`)
			}
		}
		h.GetBuffer().S(">")
		if len(options) > 0 {
			for idx := range options {
				options[idx](h)
			}
		}
		h.GetBuffer().S("</thead>")
	}
}

func TBody(options []HTMLAttr, properties ...HTMLProperty) HTMLAttr {
	return func(h HTMLWriter) {
		h.GetBuffer().S("<tbody")
		if len(properties) > 0 {
			for idx := range properties {
				h.GetBuffer().S(" ", properties[idx].Key, `="`, properties[idx].Value, `"`)
			}
		}
		h.GetBuffer().S(">")
		if len(options) > 0 {
			for idx := range options {
				options[idx](h)
			}
		}
		h.GetBuffer().S("</tbody>")
	}
}

func H(level int, options []HTMLAttr, properties ...HTMLProperty) HTMLAttr {
	return func(h HTMLWriter) {
		h.GetBuffer().S("<h", cast.ToString(level))
		if len(properties) > 0 {
			for idx := range properties {
				h.GetBuffer().S(" ", properties[idx].Key, `="`, properties[idx].Value, `"`)
			}
		}
		h.GetBuffer().S(">")
		if len(options) > 0 {
			for idx := range options {
				options[idx](h)
			}
		}
		h.GetBuffer().S("</h", cast.ToString(level), ">")
	}
}
