# Google Play

download APK from Google Play or send API requests

## tool examples

[sign in](//accounts.google.com/embedded/setup/v2/android) with your Google
Account. then get authorization code (`oauth_token`) cookie from
[browser&nbsp;storage][1]. should be valid for 10 minutes. then exchange
authorization code for refresh token (`aas_et`):

~~~
play -o oauth2_4/0Adeu5B...
~~~

[1]://firefox-source-docs.mozilla.org/devtools-user/storage_inspector

create a file containing `X-DFE-Device-ID` (GSF ID) for future requests:

~~~
play -d
~~~

get app details:

~~~
> play -a com.google.android.youtube
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

acquire app. only needs to be done once per Google account:

~~~
play -a com.google.android.youtube -acquire
~~~

download APK. you need to specify any valid version code. the latest code is
provided by the previous details command. if APK is split, all pieces will be
downloaded:

~~~
play -a com.google.android.youtube -v 1540222400
~~~
