/*
 ВНИМАНИЕ!!!!!
ЭТО ТОЛЬКО ПЕРВЫЙ КОМПИЛЯТОР БУДЕТ НАПИСАН НА ГО
ВТОРОЙ БУДЕТ НАПИСАН НА АМЕТИСТЕ!!!!!!
 */

package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("начало компилятора!!\n(а также я люблю го и аметист)\n")
	path := "main.ay"

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("ошибка чтения файла: %v\n", err)
		os.Exit(67)
	}
	code := string(data)
	tokens := lexer(code)

	fmt.Println("\nТокены:")
	for _, t := range tokens {
		fmt.Printf("%-8s | %-8s | %d\n",
			t.token_type, t.name, t.pos)
	}
	//так нужно обязательно, иначе никак
	os.Exit(42)
}

// ============================================
// ЛЕКСЕР
// ============================================

type token struct {
	token_type string
	pos        int
	name       string
}

func lexer(code string) []token {
	var buffer token
	var tokens []token
	in_comment := false
	in_string := false
	p := 0
	label := false

	for i := 0; i < len(code); i++ {
		ch := code[i]
		p++

		if ch == '"' {
			if !in_string {
				in_string = true
				buffer.name += string(ch)
			} else {
				buffer.name += string(ch)
				buffer.token_type = "string"
				buffer.pos = p
				tokens = append(tokens, buffer)
				buffer.name = ""
				in_string = false
			}
			continue
		}

		if in_string {
			buffer.name += string(ch)
			continue
		}
		
		if ch == '*' && i+1 < len(code) && code[i+1] == '/' {
			in_comment = false
			i++
		}
		
		if ch == '/' && i+1 < len(code) && code[i+1] == '*' {
    in_comment = true
    i++
    continue
}
if in_comment {
    if ch == '*' && i+1 < len(code) && code[i+1] == '/' {
        in_comment = false
        i++
    }
    continue
}
		
		if ch == ':' && i+1 < len(code) && code[i+1] == ':' {
    label = true
    i++
    continue
}
if label {
    if ch == ':' && i+1 < len(code) && code[i+1] == ':' {
        label = false
        i++
    	continue
    }
	buffer.name += string(ch)
    continue
}

		// Многосимвольные операторы
		if ch == '-' && i+1 < len(code) && code[i+1] == '>' {
			buffer.token_type = "operator"
			buffer.name = "->"
			buffer.pos = p
			tokens = append(tokens, buffer)
			buffer.name = ""
			i++
			continue
		}

		if ch == '=' && i+1 < len(code) && code[i+1] == '=' {
			buffer.token_type = "operator"
			buffer.name = "=="
			buffer.pos = p
			tokens = append(tokens, buffer)
			buffer.name = ""
			i++
			continue
		}

		if ch == '!' && i+1 < len(code) && code[i+1] == '=' {
			buffer.token_type = "operator"
			buffer.name = "!="
			buffer.pos = p
			tokens = append(tokens, buffer)
			buffer.name = ""
			i++
			continue
		}

		if ch == '<' && i+1 < len(code) && code[i+1] == '=' {
			buffer.token_type = "operator"
			buffer.name = "<="
			buffer.pos = p
			tokens = append(tokens, buffer)
			buffer.name = ""
			i++
			continue
		}

		if ch == '>' && i+1 < len(code) && code[i+1] == '=' {
			buffer.token_type = "operator"
			buffer.name = ">="
			buffer.pos = p
			tokens = append(tokens, buffer)
			buffer.name = ""
			i++
			continue
		}

		if ch == '<' && i+1 < len(code) && code[i+1] == '-' {
			buffer.token_type = "operator"
			buffer.name = "<-"
			buffer.pos = p
			tokens = append(tokens, buffer)
			buffer.name = ""
			i++
			continue
		}

		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') || ch == '_' {
			buffer.name += string(ch)
		} else {
			if buffer.name != "" {
				buffer.token_type = type_indention(buffer)
				buffer.pos = p
				tokens = append(tokens, buffer)
				buffer.name = ""
			}

			if ch != ' ' && ch != '\n' && ch != '\t' && ch != '\r' {
				buffer.token_type = "operator"
				buffer.name = string(ch)
				buffer.pos = p
				tokens = append(tokens, buffer)
				buffer.name = ""
			}
		}
	}

	if buffer.name != "" {
		buffer.token_type = type_indention(buffer)
		buffer.pos = p
		tokens = append(tokens, buffer)
	}
	return tokens
}

func type_indention(buffer token) string {
	switch buffer.name {
	case "carve", "ore", "if", "elif", "else", "gem", "strand", "do", "end", "dig_out", "layer", "asm", "crystal", "import", "mem_alloc", "detonate", "size", "len", "go_to", "explore", "align", "add", "del", "tunnel", "make_tunnel", "t_close", "await", "surface", "true", "false", "nil", "vein", "go", "rem_alloc", "offset", "f_close", "f_carve", "f_open", "f_create", "f_rename", "f_del", "to":
		return "keyword"
	case "int", "bool", "char", "float", "byte", "void":
		return "type"
	default:
		isNumber := true
		for _, c := range buffer.name {
			if c < '0' || c > '9' {
				isNumber = false
				break
			}
		}
		if isNumber {
			return "number"
		}
		return "ident"
	}
}

/*==================================
       ПАРСЕР(однопроходный)
 ==================================*/