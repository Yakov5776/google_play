package main

import (
   "154.pages.dev/google/play"
   "fmt"
   "net/http"
   "os"
   "time"
   "io"
   "reflect"
   "encoding/base64"
   "crypto/sha1"
   option "154.pages.dev/http"
)

func (f flags) is_token_empty() (bool, error) {
   home, err := os.UserHomeDir()
   if err != nil {
      return true, err
   }
   home += "/google/play/"
   var token play.Refresh_Token
   token.Raw, err = os.ReadFile(home + "token.txt")
   if err != nil {
      return true, err
   }
   if err := token.Unmarshal(); err != nil {
      return true, err
   }
   if token.Values.Get("Token") == "" {
      return true, nil
   }
   return false, nil
}

func (f flags) do_device() error {
   home, err := os.UserHomeDir()
   if err != nil {
      return err
   }
   home += "/google/play/"
   err = os.MkdirAll(home, os.ModePerm)
   name := home + fmt.Sprintf("%v.bin", f.platform)
   var check play.Checkin
   play.Phone.Set(f.platform)
   if err := check.Checkin(play.Phone); err != nil {
      return err
   }
   if err := os.WriteFile(name, check.Raw, 0666); err != nil {
      return err
   }
   fmt.Println("Sleep(9*time.Second)")
   time.Sleep(9*time.Second)
   if err := check.Unmarshal(); err != nil {
      return err
   }
   notoken, err := f.is_token_empty()
   if notoken {
      return nil
   }
   return check.Sync(play.Phone)
}

func (f flags) check_file_sha1(name string, control_hash []byte) (bool, error) {
   _, err := os.Stat(name)
   if err != nil {
      return false, nil
   }
   file, err := os.Open(name)
   if err != nil {
      return false, err
   }
   defer file.Close()
   fmt.Printf("File %s exists, verifying...\n", name)
   hash := sha1.New()
   if _, err := io.Copy(hash, file); err != nil {
      return false, err
   }
   if !reflect.DeepEqual(hash.Sum(nil), control_hash) {
      fmt.Printf("  SHA-1 mismatch, redownloading...\n")
      return false, nil
   }
   fmt.Printf("  SHA-1 OK\n")
   return true, nil
}

func (f flags) download(url, name string, sig string) error {
   h, err := base64.RawURLEncoding.DecodeString(sig)
   if err != nil {
      return err
   }
   //fmt.Printf("SHA-1: %x\n", h)
   matches, err := f.check_file_sha1(name, h)
   if err != nil {
      return err
   }
   if matches {
      return nil
   }
   res, err := http.Get(url)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   file, err := os.Create(name)
   if err != nil {
      return err
   }
   defer file.Close()
   pro := option.Progress_Length(res.ContentLength)
   if _, err := file.ReadFrom(pro.Reader(res)); err != nil {
      return err
   }
   return nil
}

func (f flags) client(a *play.Access_Token, c *play.Checkin) error {
   home, err := os.UserHomeDir()
   if err != nil {
      return err
   }
   home += "/google/play/"
   err = os.MkdirAll(home, os.ModePerm)
   var token play.Refresh_Token
   token.Raw, err = os.ReadFile(home + "token.txt")
   if err != nil {
      return err
   }
   if err := token.Unmarshal(); err != nil {
      return err
   }
   if err := a.Refresh(token); err != nil {
      return err
   }
   c.Raw, err = os.ReadFile(fmt.Sprint(home, f.platform, ".bin"))
   if err != nil {
      return err
   }
   return c.Unmarshal()
}

func (f flags) do_acquire() error {
   var client play.Acquire
   err := f.client(&client.Token, &client.Checkin)
   if err != nil {
      return err
   }
   return client.Acquire(f.app.ID)
}

func (f flags) do_auth() error {
   home, err := os.UserHomeDir()
   if err != nil {
      return err
   }
   home += "/google/play/"
   err = os.MkdirAll(home, os.ModePerm)
   token, err := play.Exchange(f.code)
   if err != nil {
      return err
   }
   return os.WriteFile(home + "token.txt", token.Raw, 0666)
}

func (f flags) do_delivery() error {
   var client play.Delivery
   err := f.client(&client.Token, &client.Checkin)
   if err != nil {
      return err
   }
   client.App = f.app
   if err := client.Delivery(f.single); err != nil {
      return err
   }
   option.Location()
   for _, apk := range client.Config_APKs() {
      if url, ok := apk.URL(); ok {
         if config, ok := apk.Config(); ok {
            if sig, ok := apk.Signature(); ok {
               err := f.download(url, f.app.APK(config), sig)
               if err != nil {
                  return err
               }
            }
         }
      }
   }
   for _, obb := range client.OBB_Files() {
      if url, ok := obb.URL(); ok {
         if role, ok := obb.Role(); ok {
            if vc, ok := obb.Version_Code(); ok {
               if sig, ok := obb.Signature(); ok {
                  err := f.download(url, f.app.OBB(role, vc), sig)
                  if err != nil {
                     return err
                  }
               }
            }
         }
      }
   }
   if url, ok := client.URL(); ok {
      if sig, ok := client.Signature(); ok {
         err := f.download(url, f.app.APK(""), sig)
         if err != nil {
            return err
         }
      }
   }
   return nil
}

func (f flags) do_details() (*play.Details, error) {
   var client play.Details
   err := f.client(&client.Token, &client.Checkin)
   if err != nil {
      return nil, err
   }
   if err := client.Details(f.app.ID, f.single); err != nil {
      return nil, err
   }
   return &client, nil
}

