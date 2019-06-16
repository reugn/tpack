# packer

Run Go code as a Unix pipeline command

> [Wiki](https://en.wikipedia.org/wiki/Pipeline_(Unix))  
> In Unix-like computer operating systems, a pipeline is a mechanism for inter-process communication using message passing. A pipeline is a set of processes chained together by their standard streams, so that the output text of each process (stdout) is passed directly as input (stdin) to the next one.

Using shell a lot, processing logs or automating various things and awk/sed are not enough? Utilize packer to write simple go applications that act as a Unix pipeline commands and get all benefits using channels, goroutines, regular expressions and more!

## Usage
Simple etl example (from examples folder)

input.txt
```
abc
+foo
+bar
def
+baz
```
db.json
```json
{
    "foo": "1",
    "bar": "2",
    "baz": "3"
}
```
```go
func main() {
	var db map[string]string
	f, _ := ioutil.ReadFile("db.json")
	json.Unmarshal(f, &db)
	packer.Packer{
		Filter: func(s string) bool {
			return strings.HasPrefix(s, "+")
		},
		Map: func(s string) string {
			s = strings.Replace(s, "+", "", 1)
			return fmt.Sprintf("%s -> %s", s, db[s])
		},
		Reduce: packer.MkString(", "),
	}.Execute()
}
```
```
cat input.txt | go run ./etl.go
```
Result:
```
foo -> 1, bar -> 2, baz -> 3
```

## License
Licensed under the MIT License.