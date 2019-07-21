#!/bin/sh

APP="SERO-Light.app"
mkdir -p $APP/Contents/{MacOS,Resources}
xgo --deps=https://gmplib.org/download/gmp/gmp-6.0.0a.tar.bz2 --targets=darwin/amd64  ../
xgo --deps=https://gmplib.org/download/gmp/gmp-6.0.0a.tar.bz2 --targets=darwin/amd64  ./
cp light-wallet-darwin-10.6-amd64 $APP/Contents/MacOS/
cp start-darwin-10.6-amd64 $APP/Contents/MacOS/start
cp -r web $APP/Contents/MacOS/
cat > $APP/Contents/Info.plist << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleExecutable</key>
	<string>start</string>
	<key>CFBundleIconFile</key>
	<string>icon.icns</string>
	<key>CFBundleIdentifier</key>
	<string>cash.sero.light</string>
</dict>
</plist>
EOF
cp icons/icon.icns $APP/Contents/Resources/icon.icns

rm -rf light-wallet-darwin-10.6-amd64 start-darwin-10.6-amd64
find $APP
