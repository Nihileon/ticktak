package sqlc

import "fmt"

func MysqlRealEscapeString(v valType) string {
	sql := fmt.Sprintf("%v", v)
	dest := make([]byte, 0, 2*len(sql))
	var escape byte

	for i := 0; i < len(sql); i++ {
		c := sql[i]
		escape = 0

		switch c {
		case 0:
			escape = '0'
			break
		case '\n':
			escape = 'n'
			break
		case '\r':
			escape = 'r'
			break
		case '\\':
			escape = '\\'
			break
		case '\'':
			escape = '\''
			break
		case '"':
			escape = '"'
			break
		case '\032':
			escape = 'Z'
		}

		if escape != 0 {
			dest = append(dest, '\\', escape)
		} else {
			dest = append(dest, c)
		}
	}

	return string(dest)
}
