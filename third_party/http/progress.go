package http

import (
   "154.pages.dev/encoding"
   "fmt"
   "io"
   "net/http"
   "time"
)

func (p Progress) percent() encoding.Percent {
   return encoding.Percent(p.first) / encoding.Percent(p.length)
}

func (p Progress) rate() encoding.Rate {
   return encoding.Rate(p.first) / encoding.Rate(time.Since(p.date).Seconds())
}

func (p Progress) size() encoding.Size {
   return encoding.Size(p.first)
}

// curl.se/docs/manpage.html#--no-progress-meter
type Progress struct {
   // datatracker.ietf.org/doc/html/rfc9110#name-content-range
   first int
   // datatracker.ietf.org/doc/html/rfc9110#name-content-range
   last int64
   // datatracker.ietf.org/doc/html/rfc9110#name-content-range
   length int64
   // datatracker.ietf.org/doc/html/rfc9110#name-content-range
   parts struct {
      // datatracker.ietf.org/doc/html/rfc9110#name-content-range
      last int64
      // datatracker.ietf.org/doc/html/rfc9110#name-content-range
      length int64
   }
   // datatracker.ietf.org/doc/html/rfc9110#name-last-modified
   modified time.Time
   // datatracker.ietf.org/doc/html/rfc9110#name-date
   date time.Time
}

func Progress_Length(length int64) *Progress {
   var p Progress
   p.length = length
   p.modified = time.Now()
   p.date = time.Now()
   return &p
}

func Progress_Parts(length int) *Progress {
   var p Progress
   p.modified = time.Now()
   p.date = time.Now()
   p.parts.length = int64(length)
   return &p
}

// complete-length
//
//   last       parts.last
//  --------   -------------
//   length     parts.length
func (p *Progress) Reader(res *http.Response) io.Reader {
   if p.parts.length >= 1 {
      p.parts.last += 1
      p.last += res.ContentLength
      p.length = p.last * p.parts.length / p.parts.last
   }
   return io.TeeReader(res.Body, p)
}

func (p *Progress) Write(b []byte) (int, error) {
   p.first += len(b)
   if time.Since(p.modified) >= time.Second {
      fmt.Println(p.percent(), " ", p.size(), " ", p.rate())
      p.modified = time.Now()
   }
   return len(b), nil
}
