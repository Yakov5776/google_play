package play

import (
   "154.pages.dev/encoding"
   "154.pages.dev/protobuf"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strconv"
)

type Details struct {
   App Application
   Checkin Checkin
   Token Access_Token
   m protobuf.Message
}

func (d *Details) Details(single bool) error {
   req, err := http.NewRequest("GET", "https://android.clients.google.com", nil)
   if err != nil {
      return err
   }
   req.URL.Path = "/fdfe/details"
   query := url.Values{
      "doc": {d.App.ID},
   }
   if d.App.Version > 0 {
      query.Set("vc", strconv.FormatUint(d.App.Version, 10))
   }
   req.URL.RawQuery = query.Encode()
   authorization(req, d.Token)
   user_agent(req, single)
   if err := x_dfe_device_id(req, d.Checkin); err != nil {
      return err
   }
   if err := x_dfe_userlanguages(req, d.App.Languages); err != nil {
      return err
   }
   res, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return errors.New(res.Status)
   }
   d.m, err = func() (protobuf.Message, error) {
      b, err := io.ReadAll(res.Body)
      if err != nil {
         return nil, err
      }
      return protobuf.Consume(b)
   }()
   if err != nil {
      return err
   }
   d.m.Message(1)
   d.m.Message(2)
   d.m.Message(4)
   return nil
}

// play.google.com/store/apps/details?id=com.google.android.youtube
func (d Details) Downloads() (uint64, bool) {
   d.m.Message(13)
   d.m.Message(1)
   return d.m.Varint(70)
}

func (d Details) Files() []uint64 {
   var files []uint64
   d.m.Message(13)
   d.m.Message(1)
   for _, f := range d.m {
      if f.Number == 17 {
         if m, ok := f.Message(); ok {
            if file, ok := m.Varint(1); ok {
               files = append(files, file)
            }
         }
      }
   }
   return files
}

// play.google.com/store/apps/details?id=com.google.android.youtube
func (d Details) Name() (string, bool) {
   return d.m.String(5)
}

// play.google.com/store/apps/details?id=com.google.android.youtube
func (d Details) Offered_By() (string, bool) {
   return d.m.String(6)
}

// play.google.com/store/apps/details?id=com.google.android.youtube
func (d Details) Price() (float64, bool) {
   d.m.Message(8)
   if v, ok := d.m.Varint(1); ok {
      return float64(v) / 1_000_000, true
   }
   return 0, false
}

// play.google.com/store/apps/details?id=com.google.android.youtube
func (d Details) Price_Currency() (string, bool) {
   d.m.Message(8)
   return d.m.String(2)
}

// play.google.com/store/apps/details?id=com.google.android.youtube
func (d Details) Requires() (string, bool) {
   d.m.Message(13)
   d.m.Message(1)
   d.m.Message(82)
   d.m.Message(1)
   return d.m.String(1)
}

func (d Details) Size() (uint64, bool) {
   d.m.Message(13)
   d.m.Message(1)
   return d.m.Varint(9)
}

func (d Details) String() string {
   var b []byte
   b = append(b, "downloads:"...)
   if v, ok := d.Downloads(); ok {
      b = fmt.Append(b, " ", encoding.Cardinal(v))
   }
   b = append(b, "\nfiles:"...)
   for _, file := range d.Files() {
      if file >= 1 {
         b = append(b, " OBB"...)
      } else {
         b = append(b, " APK"...)
      }
   }
   b = append(b, "\nname:"...)
   if v, ok := d.Name(); ok {
      b = fmt.Append(b, " ", v)
   }
   b = append(b, "\noffered by:"...)
   if v, ok := d.Offered_By(); ok {
      b = fmt.Append(b, " ", v)
   }
   b = append(b, "\nprice:"...)
   if v, ok := d.Price(); ok {
      b = fmt.Append(b, " ", v)
   }
   if v, ok := d.Price_Currency(); ok {
      b = fmt.Append(b, " ", v)
   }
   b = append(b, "\nrequires:"...)
   if v, ok := d.Requires(); ok {
      b = fmt.Append(b, " ", v)
   }
   b = append(b, "\nsize:"...)
   if v, ok := d.Size(); ok {
      b = fmt.Append(b, " ", encoding.Size(v))
   }
   b = append(b, "\nupdated on:"...)
   if v, ok := d.Updated_On(); ok {
      b = fmt.Append(b, " ", v)
   }
   b = append(b, "\nversion code:"...)
   if v, ok := d.Version_Code(); ok {
      b = fmt.Append(b, " ", v)
   }
   b = append(b, "\nversion name:"...)
   if v, ok := d.Version_Name(); ok {
      b = fmt.Append(b, " ", v)
   }
   b = append(b, "\nwhats new:"...)
   if v, ok := d.Whats_New(); ok {
      b = fmt.Append(b, " ", v)
   }
   return string(b)
}

// play.google.com/store/apps/details?id=com.google.android.youtube
func (d Details) Updated_On() (string, bool) {
   d.m.Message(13)
   d.m.Message(1)
   return d.m.String(16)
}

// developer.android.com/guide/topics/manifest/manifest-element
func (d Details) Version_Code() (uint64, bool) {
   d.m.Message(13)
   d.m.Message(1)
   return d.m.Varint(3)
}

// developer.android.com/guide/topics/manifest/manifest-element
func (d Details) Version_Name() (string, bool) {
   d.m.Message(13)
   d.m.Message(1)
   return d.m.String(4)
}

func (d Details) Whats_New() (string, bool) {
   d.m.Message(13)
   d.m.Message(1)
   return d.m.String(15)
}
