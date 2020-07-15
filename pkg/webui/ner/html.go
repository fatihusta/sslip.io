// Copyright 2020 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ner

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
)

const htmlTemplate = `
<!doctype html>
<html lang="">

<head>
	<meta charset="utf-8">
	<title>spaGO</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<style>
		{{.CSS}}
	</style>
</head>

<body class="flex flex-col">
	<header class="p-2 shadow z-10">
		<a
			class="text-xl font-bold italic"
			href="https://github.com/nlpodyssey/spago"
			target="_blank"
		>
			spa<em class="text-blue">GO</em>
		</a>
	</header>

	<main class="flex-grow overflow-hidden flex">
		<form
			class="bg-gray-200 flex-grow p-2 flex flex-col"
			onsubmit="analyze(); return false;"
		>
			<div class="flex-grow bg-white rounded shadow flex overflow-hidden relative">
				<div
					id="highlightable-text"
					class="absolute inset-0 text-transparent p-2 overflow-auto"
				></div>
				<textarea
					id="text"
					placeholder="Text..."
					class="flex-grow resize-none bg-transparent p-2 z-10 overflow-auto"
					oninput="handleTextareaInput()"
					onscroll="handleTextareaScroll()"
				></textarea>
			</div>
			<div class="mt-2 flex">
				<div class="flex-grow py-2">
					<input id="merge-entities" type="checkbox" checked>
					<label for="merge-entities" class="mr-2">Merge entities</label>

					<input id="filter-not-entities" type="checkbox" checked>
					<label for="filter-not-entities">Filter non-entities</label>
				</div>
				<input
					type="submit"
					id="submit"
					value="Analyze"
					class="rounded shadow cursor-pointer py-2 px-4 bg-blue hover:bg-light-blue"
				>
				<div id="loader" class="hidden"></div>
			</div>
		</form>

		<aside id="entities" class="bg-gray-300 shadow p-2 overflow-auto flex flex-col">
		</aside>
	</main>
	<script>
		{{.Script}}
  	</script>
</body>

</html>
`

var html []byte

func init() {
	t, err := template.New("NER Web UI").Parse(htmlTemplate)
	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		CSS    template.CSS
		Script template.JS
	}{
		CSS:    style,
		Script: script,
	}

	buf := bytes.NewBuffer([]byte{})
	err = t.Execute(buf, data)
	if err != nil {
		log.Fatal(err)
	}
	html = buf.Bytes()
}

// Handler is the server handler function for Named-Entity Recognition web UI
func Handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // that's intended for testing purposes only
	w.Header().Set("Content-Type", "text/html")
	_, err := w.Write(html)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}