#!/bin/sh

APP="pullup.app"
rm -rf $APP

mkdir -p $APP/Contents/{MacOS,Resources}
mkdir -p $APP/Contents/MacOS/czero
mkdir -p $APP/Contents/MacOS/czero/lib
mkdir -p $APP/Contents/MacOS/bin

rm -rf ../../../go-czero-import/czero/lib/*
cp -r ../../../go-czero-import/czero/lib_DARWIN_AMD64/* ../../../go-czero-import/czero/lib/

cp -r ../../../go-czero-import/czero/data $APP/Contents/MacOS/czero/
cp -r ../../../go-czero-import/czero/include $APP/Contents/MacOS/czero/
cp -r ../../../go-czero-import/czero/lib_DARWIN_AMD64/* $APP/Contents/MacOS/czero/lib/

xgo --deps=https://gmplib.org/download/gmp/gmp-6.0.0a.tar.bz2 --targets=darwin/amd64  ../
xgo --deps=https://gmplib.org/download/gmp/gmp-6.0.0a.tar.bz2 --targets=darwin/amd64  ./

cp light-wallet-darwin-10.6-amd64 $APP/Contents/MacOS/bin/
cp build-darwin-10.6-amd64 $APP/Contents/MacOS/startup

cat > $APP/Contents/Info.plist << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleExecutable</key>
	<string>startup</string>
	<key>CFBundleIconFile</key>
	<string>icon.icns</string>
	<key>CFBundleIdentifier</key>
	<string>cash.sero.pullup</string>
</dict>
</plist>
EOF
cp icons/icon.icns $APP/Contents/Resources/icon.icns

rm -rf light-wallet-darwin-10.6-amd64
rm -rf build-darwin-10.6-amd64

find $APP
