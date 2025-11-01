package play

import (
   "154.pages.dev/protobuf"
   "bytes"
   "errors"
   "io"
   "net/http"
)

type AssetDelivery struct {
   App Application
   Checkin Checkin
   Token Access_Token
   m protobuf.Message
}

func (d *AssetDelivery) AssetDelivery() error {
   var m protobuf.Message
   m.Add_String(1, d.App.ID)
   m.Add(2, func(m *protobuf.Message) {
      m.Add_Varint(1, d.App.Version)
      //m.Add_Varint(2, 3)
   })
   //m.Add_Varint(3, 0)
   //m.Add_Varint(4, 0)
   //m.Add_Varint(4, 3)
   //m.Add_Varint(5, 1)
   //m.Add_Varint(5, 2)
   m.Add(6, func(m *protobuf.Message) {
      m.Add_String(1, d.App.AssetModule)
   })
   //m.Add_Varint(8, 0)
   req, err := http.NewRequest(
      "POST",
      "https://android.clients.google.com",
      bytes.NewReader(m.Append(nil)),
   )
   if err != nil {
      return err
   }
   req.URL.Path = "/fdfe/assetModuleDelivery"
   req.Header.Set("Content-Type", "application/x-protobuf")
   authorization(req, d.Token)
   user_agent(req, false)
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
   d.m.Message(151)
   d.m.Message(3)
   return nil
}

// https://developer.android.com/guide/playcore/asset-delivery
func (d AssetDelivery) Asset_Parts() []Asset_Part {
   var files []Asset_Part
   for _, f := range d.m {
      if f.Number == 4 {
         if file, ok := f.Message(); ok {
            files = append(files, Asset_Part{file})
         }
      }
   }
   return files
}

func (d AssetDelivery) Name() (string, bool) {
   return d.m.String(1)
}

func (d AssetDelivery) Version_Code() (uint64, bool) {
   return d.m.Varint(2)
}

func (d AssetDelivery) Asset_Name() (string, bool) {
   d.m.Message(3)
   d.m.Message(1)
   return d.m.String(1)
}

func (d AssetDelivery) Size() (uint64, bool) {
   d.m.Message(3)
   d.m.Message(2)
   return d.m.Varint(1)
}

func (d AssetDelivery) Signature() (string, bool) {
   d.m.Message(3)
   d.m.Message(2)
   return d.m.String(2)
}

func (d AssetDelivery) Parts() []Asset_Part {
   d.m.Message(3)
   d.m.Message(2)

   var files []Asset_Part
   for _, f := range d.m {
      if f.Number == 4 {
         if file, ok := f.Message(); ok {
            files = append(files, Asset_Part{file})
         }
      }
   }
   return files
}

type Asset_Part struct {
   m protobuf.Message
}

func (p Asset_Part) Size() (uint64, bool) {
   return p.m.Varint(1)
}

func (p Asset_Part) Signature() (string, bool) {
   return p.m.String(2)
}

func (p Asset_Part) URL() (string, bool) {
   return p.m.String(3)
}
