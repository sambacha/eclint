package main

import (
	"fmt"
)

// endOfLines checks the line ending
func endOfLine(eol string, data []byte) error {
	l := len(data)
	switch eol {
	case "lf":
		if l > 0 && data[l-1] != '\n' {
			return fmt.Errorf("line does not end with lf (`\\n`)")
		}
		if l > 1 && data[l-2] == '\r' {
			return fmt.Errorf("line should not end with crlf (`\\r\\n`)")
		}
	case "crlf":
		if l > 0 && data[l-1] != '\n' || (l > 1 && data[l-2] != '\r') {
			return fmt.Errorf("line does not end with crlf (`\\r\\n`)")
		}
	case "cr":
		if l > 0 && data[l-1] != '\r' {
			return fmt.Errorf("line does not end with cr (`\\r`)")
		}
	default:
		return fmt.Errorf("%q is an invalid value for eol, want cr, crlf, or lf", eol)
	}

	return nil
}

// charset checks the line encoding
func charset(charset string, data []byte) error {
	switch charset {
	case "latin1":
		log.V(1).Info("not implemented", "charset", charset)
	case "utf-16be":
		log.V(1).Info("not implemented", "charset", charset)
	case "utf-16le":
		log.V(1).Info("not implemented", "charset", charset)
	case "utf-8":
		log.V(1).Info("not implemented", "charset", charset)
	case "utf-8 bom":
		log.V(1).Info("not implemented", "charset", charset)
	default:
		return fmt.Errorf("%q is an invalid value of charset, want latin1 or some utf variants", charset)
	}

	return nil
}

// indentStyle checks that the line beginnings are either space or tabs
func indentStyle(style string, data []byte) error {
	var c byte
	var x byte
	switch style {
	case "space":
		c = ' '
		x = '\t'
	case "tab":
		c = '\t'
		x = ' '
	default:
		return fmt.Errorf("%q is an invalid value of indent_style, want tab or space", style)
	}

	for i := 0; i < len(data); i++ {
		if data[i] == c {
			continue
		}
		if data[i] == x {
			return fmt.Errorf("pos %d: indentation style mismatch expected %s", i, style)
		}
		break
	}

	return nil
}

// trimTrailingWhitespace
func trimTrailingWhitespace(data []byte) error {
	for i := len(data) - 1; i >= 0; i-- {
		if data[i] == '\r' || data[i] == '\n' {
			continue
		}
		if data[i] == ' ' || data[i] == '\t' {
			return fmt.Errorf("pos %d: looks like a trailing whitespace", i)
		}
		break
	}
	return nil
}
