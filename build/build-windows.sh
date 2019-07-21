#!/bin/sh

APP="pullup"

rm -rf $APP


mkdir -p $APP/bin
mkdir -p $APP/czero
mkdir -p $APP/czero/lib

rm -rf ../../../go-czero-import/czero/lib/*
cp -r ../../../go-czero-import/czero/lib_WINDOWS_AMD64/* ../../../go-czero-import/czero/lib/

cp -r ../../../go-czero-import/czero//data $APP/czero/
cp -r ../../../go-czero-import/czero//include $APP/czero/
cp -r ../../../go-czero-import/czero//lib_WINDOWS_AMD64/* $APP/czero/lib/

xgo -ldflags "-H=windowsgui" --deps=https://gmplib.org/download/gmp/gmp-6.0.0a.tar.bz2 --targets=windows/amd64 ../
xgo -ldflags "-H=windowsgui" --deps=https://gmplib.org/download/gmp/gmp-6.0.0a.tar.bz2 --targets=windows/amd64 ./

cp light-wallet-windows-4.0-amd64.exe $APP/bin/
cp build-windows-4.0-amd64.exe $APP/pullup.exe

rm -rf light-wallet-windows-4.0-amd64.exe
rm -rf build-windows-4.0-amd64.exe

find $APP
