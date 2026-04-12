// implements hexadecimal encoding and decoding
//
// control characters can be two bytes, for example U+0080 to U+009F:
// wikipedia.org/wiki/Unicode_control_characters
package hex

import (
   "encoding/hex"
   "fmt"
   "unicode"
   "unicode/utf8"
)

func binary_rune(r rune, size int) bool {
   if r == utf8.RuneError {
      if size == 1 {
         return true
      }
   }
   if unicode.Is(unicode.C, r) {
      if r != '\n' {
         if r != '\r' {
            return true
         }
      }
   }
   return false
}

// wikipedia.org/wiki/Escape_character
const escape_character = '%'

func Binary(src []byte) bool {
   for len(src) >= 1 {
      r, size := utf8.DecodeRune(src)
      if binary_rune(r, size) {
         return true
      }
      src = src[size:]
   }
   return false
}

func Decode(src []byte) ([]byte, error) {
   var dst []byte
   for len(src) >= 1 {
      b, size, err := decode_byte(src)
      if err != nil {
         return nil, err
      }
      dst = append(dst, b)
      src = src[size:]
   }
   return dst, nil
}

func decode_byte(src []byte) (byte, int, error) {
   if len(src) == 0 {
      return 0, 0, nil
   }
   if src[0] != escape_character {
      return src[0], 1, nil
   }
   if len(src) <= 2 {
      return 0, 0, fmt.Errorf("invalid hex escape %q", src)
   }
   var dst [1]byte
   _, err := hex.Decode(dst[:], src[1:3])
   if err != nil {
      return 0, 0, err
   }
   return dst[0], 3, nil
}

func encode_rune(src []byte) ([]byte, int) {
   r, size := utf8.DecodeRune(src)
   if r != escape_character {
      if !binary_rune(r, size) {
         return src[:size], size
      }
   }
   dst := make([]byte, size*3)
   for i := 0; i < len(dst); i += 3 {
      dst[i] = escape_character
      hex.Encode(dst[i+1:], src[:1])
      src = src[1:]
   }
   return dst, size
}

func Encode(src []byte) []byte {
   var dst []byte
   for len(src) >= 1 {
      b, size := encode_rune(src)
      dst = append(dst, b...)
      src = src[size:]
   }
   return dst
}
