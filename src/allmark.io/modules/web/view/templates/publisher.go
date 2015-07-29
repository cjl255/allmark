// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

const publisherSnippet = `{{define "publisher-snippet"}}
{{if or .Author.Name .CreationDate}}
<section class="publisher">
{{if and .Author.Name .Author.URL}}

	created by <span class="author" itemprop="author" rel="author">
	<a href="{{ .Author.URL }}" title="{{ .Author.Name }}" target="_blank">
	{{ .Author.Name }}
	</a>
	</span>

{{else if .Author.Name}}

	created by <span class="author" itemprop="author" rel="author">{{ .Author.Name }}</span>

{{end}}
{{if .CreationDate}}

	{{if not .Author.Name}}created{{end}} on <span class="creationdate" itemprop="dateCreated">{{ .CreationDate }}</span>

{{end}}
</section>
{{end}}
{{end}}
`
