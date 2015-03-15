##Conv

This is a golang version encoding converter, just implemented converting from GBK/CP936 to UTF-8.

##Usage

```golang
if file, err := os.Open(path); err != nil {
	log.Fatal(err)
}else {
	defer file.Close()

	var out bytes.Buffer
	GbkToUtf8(file, &out, true)
}
```
