package parsers

import (
	"strings"
)

// DecodeLaTeX converts common LaTeX escape sequences to their UTF-8 equivalents.
// It handles accents like \'a, \"o, \^e, \~n, \c{c} and special characters.
func DecodeLaTeX(input string) string {
	// Common replacements map
	// We use a slice of structs to ensure order if necessary, but a map is usually fine for single chars.
	// However, we want to handle longer sequences first.
	replacements := []struct {
		old string
		new string
	}{
		// Acute accents
		{`{\'a}`, "á"}, {`{\'e}`, "é"}, {`{\'i}`, "í"}, {`{\'o}`, "ó"}, {`{\'u}`, "ú"},
		{`{\'A}`, "Á"}, {`{\'E}`, "É"}, {`{\'I}`, "Í"}, {`{\'O}`, "Ó"}, {`{\'U}`, "Ú"},
		{`{\'\i}`, "í"}, // Special case for i without dot

		// Grave accents
		{`{\` + "`" + `a}`, "à"}, {`{\` + "`" + `e}`, "è"}, {`{\` + "`" + `i}`, "ì"}, {`{\` + "`" + `o}`, "ò"}, {`{\` + "`" + `u}`, "ù"},
		{`{\` + "`" + `A}`, "À"}, {`{\` + "`" + `E}`, "È"}, {`{\` + "`" + `I}`, "Ì"}, {`{\` + "`" + `O}`, "Ò"}, {`{\` + "`" + `U}`, "Ù"},

		// Diaeresis (umlaut)
		{`{\"a}`, "ä"}, {`{\"e}`, "ë"}, {`{\"i}`, "ï"}, {`{\"o}`, "ö"}, {`{\"u}`, "ü"},
		{`{\"A}`, "Ä"}, {`{\"E}`, "Ë"}, {`{\"I}`, "Ï"}, {`{\"O}`, "Ö"}, {`{\"U}`, "Ü"},

		// Tilde
		{`{\~n}`, "ñ"}, {`{\~N}`, "Ñ"},
		{`{\~a}`, "ã"}, {`{\~o}`, "õ"},
		{`{\~A}`, "Ã"}, {`{\~O}`, "Õ"},

		// Circumflex
		{`{\^a}`, "â"}, {`{\^e}`, "ê"}, {`{\^i}`, "î"}, {`{\^o}`, "ô"}, {`{\^u}`, "û"},
		{`{\^A}`, "Â"}, {`{\^E}`, "Ê"}, {`{\^I}`, "Î"}, {`{\^O}`, "Ô"}, {`{\^U}`, "Û"},

		// Cedilla
		{`{\c{c}}`, "ç"}, {`{\c{C}}`, "Ç"},
		{`\c{c}`, "ç"}, {`\c{C}`, "Ç"}, // Sometimes without braces around the whole thing

		// Other symbols
		{`{\ss}`, "ß"},
		{`{\aa}`, "å"}, {`{\AA}`, "Å"},
		{`{\ae}`, "æ"}, {`{\AE}`, "Æ"},
		{`{\o}`, "ø"}, {`{\O}`, "Ø"},
		
		// Simple escapes without outer braces (sometimes found)
		{`\'a`, "á"}, {`\'e`, "é"}, {`\'i`, "í"}, {`\'o`, "ó"}, {`\'u`, "ú"},
		{`\'A`, "Á"}, {`\'E`, "É"}, {`\'I`, "Í"}, {`\'O`, "Ó"}, {`\'U`, "Ú"},
		{`\"a`, "ä"}, {`\"e`, "ë"}, {`\"i`, "ï"}, {`\"o`, "ö"}, {`\"u`, "ü"},
		{`\~n`, "ñ"}, {`\~N`, "Ñ"},
	}

	result := input
	for _, r := range replacements {
		result = strings.ReplaceAll(result, r.old, r.new)
	}
	
	// Strip outer braces if present (e.g. {{Title}} -> {Title} -> Title)
	// BibTeX parser might leave one set of braces if double braces were used.
	if strings.HasPrefix(result, "{") && strings.HasSuffix(result, "}") {
		result = strings.TrimPrefix(result, "{")
		result = strings.TrimSuffix(result, "}")
	}
	
	return result
}
