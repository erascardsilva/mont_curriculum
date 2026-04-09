package parser

import (
	"regexp"
	"strings"
	"unicode"
)

/*
 * Project: Monte Curriculum
 * Author: Erasmo Cardoso - Software Engineer | Electronics Specialist
 */

// Sanitize aplica regras gramaticalmente estritas para limpar o texto.
func Sanitize(text string) string {
	if text == "" {
		return ""
	}

	// 1. Normalização de Espaços (Preservando quebras de linha únicas)
	text = strings.ReplaceAll(text, "\r", "")
	reMultipleSpaces := regexp.MustCompile(`[ \t]{2,}`)
	text = reMultipleSpaces.ReplaceAllString(text, " ")

	// 2. Pontuação Estrita (Garante espaço após , . ! ? : ;)
	rePunct := regexp.MustCompile(`([,.!?;:])([^\s\d\n])`)
	text = rePunct.ReplaceAllString(text, "$1 $2")
	
	// Caso especial: .Palavra -> . Palavra
	reDotUpper := regexp.MustCompile(`\.([A-Z])`)
	text = reDotUpper.ReplaceAllString(text, ". $1")

	// 3. Dicionário de Separação Agressiva
	keywords := []string{"execução", "integração", "experiência", "atuação", "infraestrutura", "acesso", "autenticação", "análise", "apoio", "desenvolvimento", "manutenção", "instalação", "servidores", "redes", "banco", "dados"}
	for _, kw := range keywords {
		text = strings.ReplaceAll(text, "e"+kw, "e "+kw)
		text = strings.ReplaceAll(text, "com"+kw, "com "+kw)
		text = strings.ReplaceAll(text, "de"+kw, "de "+kw)
		text = strings.ReplaceAll(text, "em"+kw, "em "+kw)
	}

	// 4. Capitalização após pontuação
	text = CapitalizeSentences(text)

	// Limpeza de espaços no início/fim de cada linha
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	text = strings.Join(lines, "\n")

	return strings.TrimSpace(text)
}

// CapitalizeSentences garante que o início de cada frase comece com letra maiúscula.
func CapitalizeSentences(text string) string {
	re := regexp.MustCompile(`(^|[.!?]\s+)([a-z])`)
	return re.ReplaceAllStringFunc(text, func(s string) string {
		runes := []rune(s)
		for i := len(runes) - 1; i >= 0; i-- {
			if unicode.IsLower(runes[i]) {
				runes[i] = unicode.ToUpper(runes[i])
				break
			}
		}
		return string(runes)
	})
}
