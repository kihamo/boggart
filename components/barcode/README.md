#### OpenVC install

###### Linux

https://gocv.io/getting-started/linux/

```
go get -u -d gocv.io/x/gocv
cd $GOPATH/src/gocv.io/x/gocv
make install
```

###### MacOS

https://gocv.io/getting-started/macos/

```
go get -u -d gocv.io/x/gocv
brew uninstall opencv
brew install opencv
brew install pkgconfig
```

В MacOS терминалу нужно дать доступ к камере



###### Check 

```
go run $GOPATH/src/gocv.io/x/gocv/cmd/version/main.go
```

#### ZBar install

###### Linux

```
sudo apt install zbar-tools
```

###### MacOS

```
brew install zbar
```
