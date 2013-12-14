package main

var pkgTemplate = `{{with .PDoc}}
	{{if $.IsMain}}
		{{/* command documentation */}}
		{{comment_md .Doc}}
	{{else}}
		{{/* package documentation */}}
		{{comment_md .Doc}}
		{{example_html $ ""}}

		{{if $.Examples}}
		<div id="pkg-examples">
			<h4>Examples</h4>
			<dl>
			{{range $.Examples}}
			<dd><a class="exampleLink" href="#example_{{.Name}}">{{example_name .Name}}</a></dd>
			{{end}}
			</dl>
		</div>
		{{end}}

{{with .Consts}}
## Constants
{{range .}}
<pre>{{node_html $ .Decl true}}</pre>
{{comment_md .Doc}}
{{end}}
{{end}}
{{with .Vars}}
## Variables
{{range .}}
<pre>{{node_html $ .Decl true}}</pre>
{{comment_md .Doc}}
{{end}}
{{end}}
{{range .Funcs}}
{{/* Name is a string - no need for FSet */}}
{{$name_html := html .Name}}
## func {{$name_html}}
<pre>{{node_html $ .Decl false}}</pre>
{{comment_md .Doc}}
{{example_html $ .Name}}
{{end}}
{{range .Types}}
{{$tname := .Name}}
{{$tname_html := html .Name}}
## type {{$tname_html}}
<pre>{{node_html $ .Decl true}}</pre>
{{comment_md .Doc}}

			{{range .Consts}}
				<pre>{{node_html $ .Decl true}}</pre>
				{{comment_html .Doc}}
			{{end}}

			{{range .Vars}}
				<pre>{{node_html $ .Decl true}}</pre>
				{{comment_html .Doc}}
			{{end}}

			{{example_html $ $tname}}

			{{range .Funcs}}
				{{$name_html := html .Name}}
				<h3 id="{{$name_html}}">func <a href="{{posLink_url $ .Decl}}">{{$name_html}}</a></h3>
				<pre>{{node_html $ .Decl true}}</pre>
				{{comment_html .Doc}}
				{{example_html $ .Name}}
			{{end}}

			{{range .Methods}}
				{{$name_html := html .Name}}
				<h3 id="{{$tname_html}}.{{$name_html}}">func ({{html .Recv}}) <a href="{{posLink_url $ .Decl}}">{{$name_html}}</a></h3>
				<pre>{{node_html $ .Decl true}}</pre>
				{{comment_html .Doc}}
				{{$name := printf "%s_%s" $tname .Name}}
				{{example_html $ $name}}
			{{end}}
		{{end}}
	{{end}}

	{{with $.Notes}}
		{{range $marker, $content := .}}
			<h2 id="pkg-note-{{$marker}}">{{noteTitle $marker | html}}s</h2>
			<ul style="list-style: none; padding: 0;">
			{{range .}}
			<li><a href="{{posLink_url $ .}}">&#x261e;</a> {{html .Body}}</li>
			{{end}}
			</ul>
		{{end}}
	{{end}}
{{end}}

{{with .PAst}}
	<pre>{{node_html $ . false}}</pre>
{{end}}`
