package encryption

// TODO Реализовать шифрование чисел
func Cesar(message string, key byte) string {
	cipherText := ""
	for i := 0; i < len(message); i++ {
		c := message[i]
		if c >= 'a' && c <= 'z' {
			c += key
			if c < 'a' {
				c += 26
			}
		} else if c >= 'A' && c <= 'Z' {
			c += key
			if c < 'A' {
				c += 26
			}
		}
		cipherText += string(c)
	}
	return cipherText
}

func UnCesar(message string, key byte) string {
	cipherText := ""
	for i := 0; i < len(message); i++ {
		c := message[i]
		if c >= 'a' && c <= 'z' {
			c -= key
			if c < 'a' {
				c += 26
			}
		} else if c >= 'A' && c <= 'Z' {
			c -= key
			if c < 'A' {
				c += 26
			}
		}
		cipherText += string(c)
	}
	return cipherText
}

//
