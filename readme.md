# Google Play

Download APK from Google Play or send API requests.

## tool examples

Go [here](//accounts.google.com/embedded/setup/v2/android) in incognito mode and sign in with your Google Account. Then get authorization code (`oauth_token`) cookie from [browser&nbsp;storage][1]. It should be valid for 10 minutes. Then exchange authorization code for refresh token (`aas_et`):

~~~
play -o oauth2_4/0Adeu5B...
~~~

[1]://firefox-source-docs.mozilla.org/devtools-user/storage_inspector

Now create a file containing `X-DFE-Device-ID` (GSF ID):

~~~
play -c -p 2
~~~

2 in this example is platform number. You need to run this command once for each platform you intend to use. Different apps or different app versions may be available for different devices.

Valid platforms are:

- 0 - ARMv6 phone (320x480 160dpi)
- 1 - ARMv7 phone (720x1280 320dpi)
- 2 - ARMv8 phone (1080x1920 480dpi)
- 3 - ARMv8 4K phone (2160x3840 640dpi)
- 4 - x86 phone (720x1280 320dpi)
- 5 - x86-64 phone (1080x1920 480dpi)
- 6 - ARMv7 tablet (1024x600 160dpi)
- 7 - ARMv7 large tablet (1280x800 160dpi)
- 8 - ARMv8 tablet (2560x1600 320dpi)
- 9 - ARMv8 large tablet (2960x1848 240dpi)
- 10 - MIPS device
- 11 - MIPS 64-bit device
- 12 - RISC-V 64-bit device

Get app details:

~~~
> play -a com.google.android.youtube -p 2
creator: Google LLC
file: APK APK APK APK
installation size: 89.03 megabyte
downloads: 14.81 billion
offer: 0 USD
requires: Android 8.0 and up
title: YouTube
upload date: Sep 22, 2023
version: 18.38.37
version code: 1540222400
~~~

Acquire the app if you don't own it already. This only needs to be done once per Google account:

~~~
play -a com.google.android.youtube -p 2 -acquire
~~~

Download APK. You need to specify any valid version code. The latest code is provided by the previous details command. If APK is split, all pieces will be downloaded:

~~~
play -a com.google.android.youtube -p 2 -v 1540222400
~~~
