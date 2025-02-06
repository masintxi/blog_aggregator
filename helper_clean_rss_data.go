package main

import (
	"bytes"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// Extracts and cleans content while removing only the first <div> block after each <!-- BREAK X -->
func extractAndCombineHTML(input string) string {
	// Regular expression to find break markers like <!-- BREAK 1 -->, <!-- BREAK 2 -->, etc.
	breakRegex := regexp.MustCompile(`<!-- BREAK \d+ -->`)

	// Split the input string into sections based on the break markers
	splits := breakRegex.Split(input, -1)

	// Slice to collect cleaned sections
	var sections []string

	// Process each section after the markers
	for i, section := range splits {
		if i == 0 {
			// First section (before any marker), keep as is
			sections = append(sections, cleanHTML(section, false))
			continue
		}

		// Remove only the first <div> block after each marker
		cleanedSection := cleanHTML(section, true)
		sections = append(sections, cleanedSection)
	}

	// Join all sections with a newline
	return strings.Join(sections, "\n")
}

// Cleans HTML while removing only the first <div> block after the break marker
func cleanHTML(input string, removeFirstDiv bool) string {
	doc, err := html.Parse(strings.NewReader(input))
	if err != nil {
		return input // Return original if parsing fails
	}

	var output bytes.Buffer
	var firstDivRemoved = !removeFirstDiv // Flag to track if the first <div> was removed
	seenHeaders := make(map[string]bool)  // Track headers to prevent duplication

	var f func(*html.Node)
	f = func(n *html.Node) {
		// Skip <script> and <style> tags
		if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style") {
			return
		}

		// If it's the first <div> after the break marker, remove it
		if !firstDivRemoved && n.Type == html.ElementNode && n.Data == "div" {
			firstDivRemoved = true
			return // Skip this block
		}

		// Handle headers separately to avoid duplication
		if n.Type == html.ElementNode && isHeaderTag(n.Data) {
			headerText := extractText(n)
			if headerText != "" && !seenHeaders[headerText] {
				output.WriteString("\n" + headerText + "\n") // Add a newline before and after
				seenHeaders[headerText] = true
			}
			return // Avoid re-processing header text inside later
		}

		// Append text content
		if n.Type == html.TextNode {
			text := strings.TrimSpace(n.Data)
			if text != "" {
				output.WriteString(text + " ")
			}
		}

		// Recursively process child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	// Convert to string and clean up spaces/newlines
	cleanedText := output.String()

	// Remove excessive spaces but preserve newlines
	spaceRegex := regexp.MustCompile(`[ \t]+`)
	cleanedText = spaceRegex.ReplaceAllString(cleanedText, " ")

	// Remove extra newlines (\n\n → \n)
	newlineRegex := regexp.MustCompile(`\n{2,}`)
	cleanedText = newlineRegex.ReplaceAllString(cleanedText, "\n")

	// Trim final spaces and newlines
	return strings.TrimSpace(cleanedText)
}

// Checks if the node is a header tag (h1, h2, ..., h6)
func isHeaderTag(tag string) bool {
	return tag == "h1" || tag == "h2" || tag == "h3" || tag == "h4" || tag == "h5" || tag == "h6"
}

// Extracts text from a node, including its children
func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data)
	}
	var result string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result += extractText(c) + " "
	}
	return strings.TrimSpace(result)
}
