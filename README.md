# containers-from-scratch

You need root permissions for this to work.

Also note that the Go code uses some syscall definitions that are only available when building with GOOS=linux.
Basics
credit goes https://yourbasic.org/golang/measure-execution-time/
A slice doesn’t store any data, it just describes a section of an underlying array.

When you change an element of a slice, you modify the corresponding element of its underlying array, and other slices that share the same underlying array will see the change.
A slice can grow and shrink within the bounds of the underlying array.
Slices are indexed in the usual way: s[i] accesses the ith element, starting from zero.
var s []int                   // a nil slice
s1 := []string{"foo", "bar"}
s2 := make([]int, 2)          // same as []int{0, 0}
s3 := make([]int, 2, 4)       // same as new([4]int)[:2]
fmt.Println(len(s3), cap(s3)) // 2 4

a := [...]int{0, 1, 2, 3} // an array
s := a[1:3]               // s == []int{1, 2}        cap(s) == 3
s = a[:2]                 // s == []int{0, 1}        cap(s) == 4
s = a[2:]                 // s == []int{2, 3}        cap(s) == 2
s = a[:]                  // s == []int{0, 1, 2, 3}  cap(s) == 4

s := []int{0, 1, 2, 3, 4} // a slice
s = s[1:4]                // s == []int{1, 2, 3}
s = s[1:2]                // s == []int{2} (index relative to slice)
s = s[:3]                 // s == []int{2, 3, 4} (extend length)

s := []string{"Foo", "Bar"}
for i, v := range s {
    fmt.Println(i, v)
}
The built-in copy function copies elements into a destination slice dst from a source slice src.
func copy(dst, src []Type) int
It returns the number of elements copied, which will be the minimum of len(dst) and len(src). The result does not depend on whether the arguments overlap.

It is legal to copy bytes from a string to a slice of bytes.
copy(dst []byte, src string) int
var s = make([]int, 3)
n := copy(s, []int{0, 1, 2, 3}) // n == 3, s == []int{0, 1, 2}
s := []int{0, 1, 2}
n := copy(s, s[1:]) // n == 2, s == []int{1, 2, 2}
var b = make([]byte, 5)
copy(b, "Hello, world!") // b == []byte("Hello")
The idiomatic way to implement a stack in Go is to use a slice:

to push you use the built-in append function, and
to pop you slice off the top element.
var stack []string

stack = append(stack, "world!") // Push
stack = append(stack, "Hello ")

for len(stack) > 0 {
    n := len(stack) - 1 // Top element
    fmt.Print(stack[n])

    stack = stack[:n] // Pop
}
Watch out for memory leaks
If the stack is permanent and the elements temporary, you may want to remove the top element before popping the stack.

// Pop
stack[n] = "" // Erase element (write zero value)
stack = stack[:n]

3 ways to compare slices

/ Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func Equal(a, b []int) bool {
    if len(a) != len(b) {
        return false
    }
    for i, v := range a {
        if v != b[i] {
            return false
        }
    }
    return true
}

var a []int = nil
var b []int = make([]int, 0)
fmt.Println(reflect.DeepEqual(a, b)) // false

The performance of this function is much worse than for the code above, but it’s useful in test cases where simplicity and correctness are crucial. The semantics, however, are quite complicated.
To remove all elements, set the slice to nil.
a := []string{"A", "B", "C", "D", "E"}
a = nil
fmt.Println(a, len(a), cap(a)) // [] 0 0
This will release the underlying array to the garbage collector (assuming there are no other references).

Note that nil slices and empty slices are very similar:

they look the same when printed,
they have zero length and capacity,
they can be used with the same effect in for loops and append functions.
Keep allocated memory
To keep the underlying array, slice the slice to zero length.
a := []string{"A", "B", "C", "D", "E"}
a = a[:0]
fmt.Println(a, len(a), cap(a)) // [] 0 5
If the slice is extended again, the original data reappears.
fmt.Println(a[:2]) // [A B]
Concatenate slices

a := []int{1, 2}
b := []int{11, 22}
a = append(a, b...) // a == [1 2 11 22]
The ... unpacks b. Without the dots, the code would attempt to append the slice as a whole, which is invalid.

Delete an element from a slice
Fast version (changes order)
a := []string{"A", "B", "C", "D", "E"}
i := 2

// Remove the element at index i from a.
a[i] = a[len(a)-1] // Copy last element to index i.
a[len(a)-1] = ""   // Erase last element (write zero value).
a = a[:len(a)-1]   // Truncate slice.

fmt.Println(a) // [A B E D]

Slow version (maintains order)
a := []string{"A", "B", "C", "D", "E"}
i := 2

// Remove the element at index i from a.
copy(a[i:], a[i+1:]) // Shift a[i+1:] left one index.
a[len(a)-1] = ""     // Erase last element (write zero value).
a = a[:len(a)-1]     // Truncate slice.

fmt.Println(a) // [A B D E]

Find an element in a slice
/ Contains tells whether a contains x.
func Contains(a []string, x string) bool {
    for _, n := range a {
        if x == n {
            return true
        }
    }
    return false
}
/ Find returns the smallest index i at which x == a[i],
// or len(a) if there is no such index.
func Find(a []string, x string) int {
    for i, n := range a {
        if x == n {
            return i
        }
    }
    return len(a)
}

Last item in a slice
a := []string{"A", "B", "C"}
s := a[len(a)-1] // C
To remove it:
a = a[:len(a)-1] // [A B]
Watch out for memory leaks
If the slice is permanent and the element temporary, you may want to remove the reference before removing the element from the slice.

a[len(a)-1] = "" // Erase element (write zero value)
a = a[:len(a)-1] // [A B]
Empty slice vs. nil slice
In practice, nil slices and empty slices can often be treated in the same way:

they have zero length and capacity,
they can be used with the same effect in for loops and append functions,
and they even look the same when printed.
var a []int = nil
fmt.Println(len(a)) // 0
fmt.Println(cap(a)) // 0
fmt.Println(a)      // []
However, if needed, you can tell the difference.
var a []int = nil
var a0 []int = make([]int, 0)

fmt.Println(a == nil)  // true
fmt.Println(a0 == nil) // false

fmt.Printf("%#v\n", a)  // []int(nil)
fmt.Printf("%#v\n", a0) // []int{}
The official Go wiki recommends using nil slices over empty slices.

[…] the nil slice is the preferred style.

Note that there are limited circumstances where a non-nil but zero-length slice is preferred, such as when encoding JSON objects (a nil slice encodes to null, while []string{} encodes to the JSON array []).

When designing interfaces, avoid making a distinction between a nil slice and a non-nil, zero-length slice, as this can lead to subtle programming errors.
Pass a slice to a variadic function
ou can pass a slice s directly to a variadic funtion using the s... notation
func main() {
    primes := []int{2, 3, 5, 7}
    fmt.Println(Sum(primes...)) // 17
}

func Sum(nums ...int) int {
    res := 0
    for _, n := range nums {
        res += n
    }
    return res
}

+++++++++++++++++++++++++++++++
var m map[string]int                // m == nil, len(m) == 0
m1 := make(map[string]float64)      // empty map of string-float64 pairs
m2 := make(map[string]float64, 100) // preallocate room for 100 entries
m3 := map[string]float64{
    "e":  2.71828,
    "pi": 3.1416,
}
fmt.Println(len(m1), len(m2), len(m3)) // 0 0 2
The default zero value of a map is nil. A nil map is equivalent to an empty map except that no elements can be added.
You create a map either by a map literal or a call to the make function, which takes an optional capacity as argument.
The built-in len function retrieves the number of key-value pairs.

Add, find and delete

m := make(map[string]float64)

m["pi"] = 3.1416 // Add a new key-value pair.
fmt.Println(m)   // map[pi:3.1416]

v1 := m["pi"]  // v1 == 3.1416
v2 := m["foo"] // v2 == 0 (zero value)

_, exists := m["pi"] // exists == true
_, exists = m["foo"] // exists == false

if x, ok := m["pi"]; ok { // Prints 3.1416.
    fmt.Println(x)
}

delete(m, "pi") // Delete a key-value pair.
fmt.Println(m)  // map[]
Iteration
m := map[string]float64{
    "e":  2.71828,
    "pi": 3.1416,
}
for key, value := range m { // order not specified
    fmt.Println(key, value)
}
Check if a map contains a key
m := map[string]float64{"pi": 3.1416}

v1 := m["pi"]  // v1 == 3.1416
v2 := m["foo"] // v2 == 0.0 (zero value)

_, exists := m["pi"] // exists == true

if x, ok := m["pi"]; ok {
    fmt.Println(x) // 3.1416
}
Check if a map is empty
if len(m) == 0 {
    // m is empty
}
Count elements in a map
m := map[string]int{
    "key1": 1,
    "key2": 10,
    "key3": 100,
}
fmt.Println(len(m))  // 3
Get slices of keys and values from a map
keys := make([]keyType, 0, len(myMap))
values := make([]valueType, 0, len(myMap))

for k, v := range myMap {
	keys = append(keys, k)
	values = append(values, v)
}

Sort a map by key or value
 map is an unordered collection of key-value pairs.
If you need a stable iteration order, you must maintain a separate data structure.
This example uses a sorted slice of keys to print a map[string]int in key order.
m := map[string]int{"Alice": 23, "Eve": 2, "Bob": 25}

keys := make([]string, 0, len(m))
for k := range m {
	keys = append(keys, k)
}
sort.Strings(keys)

for _, k := range keys {
	fmt.Println(k, m[k])
}

Format a string without printing it

s := fmt.Sprintf("Size: %d MB.", 85) // s == "Size: 85 MB."

Type rune: a Unicode code point
The rune type is an alias for int32, and is used to emphasize than an integer represents a code point.
CII defines 128 characters, identified by the code points 0–127. It covers English letters, Latin numbers, and a few other characters.
Unicode, which is a superset of ASCII, defines a codespace of 1,114,112 code points. Unicode version 10.0 covers 139 modern and historic scripts (including runes, but not Klingon) as well as multiple symbol sets.
Strings and UTF-8
Note that a string is a sequence of bytes, not runes.
However, strings often contain Unicode text encoded in UTF-8, which encodes all Unicode code points using one to four bytes. Since Go source code itself is encoded as UTF-8, string literals will automatically get this encoding.
For simple cases where performance is a non-issue, fmt.Sprintf is your friend.
s := fmt.Sprintf("Size: %d MB.", 85) // s == "Size: 85 MB."
Fast concatenation with a string builder
The strings.Builder type is used to efficiently concatenate strings using write methods.

It offers a subset of the bytes.Buffer methods that allows it to safely avoid redundant copying.
The Grow method can be used to preallocate memory when the maximum size of the string is known.
var b strings.Builder
b.Grow(32)
for i, p := range []int{2, 3, 5, 7, 11, 13} {
    fmt.Fprintf(&b, "%d:%d, ", i+1, p)
}
s := b.String()   // no copying
s = s[:b.Len()-2] // no copying (removes trailing ", ")
fmt.Println(s)
1:2, 2:3, 3:5, 4:7, 5:11, 6:13
var buf bytes.Buffer
for i, p := range []int{2, 3, 5, 7, 11, 13} {
    fmt.Fprintf(&buf, "%d:%d, ", i+1, p)
}
buf.Truncate(buf.Len() - 2) // Remove trailing ", "
s := buf.String()           // Copy into a new string
fmt.Println(s)
buf := []byte("Size: ")
buf = strconv.AppendInt(buf, 85, 10)
buf = append(buf, " MB."...)
s := string(buf)
buf := make([]byte, 0, 16)
buf = append(buf, "Size: "...)
buf = strconv.AppendInt(buf, 85, 10)
buf = append(buf, " MB."...)
s := string(buf)
How to split a string into a slice
s := strings.Split("a,b,c", ",")
fmt.Println(s)
// Output: [a b c]
Use the strings.Fields function to split a string into substrings removing white space.
s := strings.Fields(" a \t b \n")
fmt.Println(s)
// Output: [a b]
if "Foo" == "Bar" {
    fmt.Println("Foo and Bar are equal.")
} else {
    fmt.Println("Foo and Bar are not equal.")
}
// Output: Foo and Bar are not equal.
Use <, >, <= or >= to determine lexical order.

f "Foo" < "Bar" {
    fmt.Println("Foo comes before Bar.")
} else {
    fmt.Println("Foo does not come before Bar.")
}
// Output: Foo does not come before Bar.
Convert string to/from byte slice
When you convert a string to a byte slice, you get a slice that contains the bytes of the string.
fmt.Println([]byte("abc日")
// [97 98 99 230 151 165]
Converting a slice of bytes to a string yields a string whose bytes are the elements of the slice.
b := []byte{'a', 'b', 'c', '\xe6', '\x97', '\xa5'}
s := string(b)
fmt.Println(s)
// Output: abc日
Convert string to/from rune slice
Converting a string to a slice of runes yields a slice whose elements are the Unicode code points of the string.
s := "abc日"
r := []rune(s)
fmt.Printf("%v\n", r) // [97 98 99 26085]
fmt.Printf("%U\n", r) // [U+0061 U+0062 U+0063 U+65E5]
Rune slice to string
Converting a slice of runes to a string yields a string that is the concatenation of the runes converted to UTF-8 encoded strings.

Values outside the range of valid Unicode code points are converted to \uFFFD, the Unicode replacement character �.
r := []rune{'\u0061', '\u0062', '\u0063', '\u65E5', -1}
s := string(r)
fmt.Println(s) // abc日�
Convert int/int64 to string
int to string
Use the strconv.Itoa function to convert an int to a decimal string.
str := strconv.Itoa(123) // str == "123"
nt64 to string
Use strconv.FormatInt to format an integer in a given base.
var n int64 = 32
str := strconv.FormatInt(n, 10) // decimal
fmt.Println(str)                // 32
var n int64 = 32
str := strconv.FormatInt(n, 16) // hexadecimal
fmt.Println(str)                // 20
Use strconv.Atoi to convert/parse a string to an int.
str := "123"
if n, err := strconv.Atoi(str); err == nil {
    fmt.Println(n+1)
} else {
    fmt.Println(str, "is not an integer.")
}
// Output: 124
String to int64
Use strconv.ParseInt to parse a decimal string (base 10) and check if it fits into a 64-bit signed integer.
str := "123"
n, err := strconv.ParseInt(str, 10, 64)
if err == nil {
    fmt.Printf("%d of type %T", n, n)
}
// Output: 123 of type int64
he two numeric arguments represent a base (0, 2 to 36) and a bit size (0 to 64).

If the first argument is 0, the base is implied by the string’s prefix: base 16 for "0x", base 8 for "0", and base 10 otherwise.

The second argument specifies the integer type that the result must fit into. Bit sizes 0, 8, 16, 32, and 64 correspond to int, int8, int16, int32, and int64.
Use the fmt.Sprintf method to convert a floating-point number to a string.
s := fmt.Sprintf("%f", 123.456) // s == "123.456000"
Use the strconv.ParseFloat function to convert a string to a floating-point number with the precision specified by bitSize: 32 for float32, or 64 for float64.

func ParseFloat(s string, bitSize int) (float64, error)
When bitSize is 32, the result still has type float64, but it will be convertible to float32 without changing its value.

f := "3.14159265"
if s, err := strconv.ParseFloat(f, 32); err == nil {
    fmt.Println(s) // 3.1415927410125732
}
if s, err := strconv.ParseFloat(f, 64); err == nil {
    fmt.Println(s) // 3.14159265
}
Raw string literals
Raw string literals, delimited by back quotes, can contain line breaks.
str := `First line
Second line`
fmt.Println(str)
Raw strings literals are interpreted literally and backslashes have no special meaning.

Interpreted string literals
To insert escape characters, use interpreted string literals delimited by double quotes.
str := "\tFirst line\n" +
"Second line"
fmt.Println(str)
Remove duplicate whitespace
space := regexp.MustCompile(`\s+`)
s := space.ReplaceAllString("Hello  \t \n world!", " ")
fmt.Printf("%q", s) // "Hello world!"
\s+ is a regular expression:

the character class \s matches a space, tab, new line, carriage return or form feed,
and + says “one or more of those”.
In other words, the code will replace each whitespace substring with a single space character.
Use the strings.TrimSpace function to remove leading and trailing whitespace as defined by Unicode.
s := strings.TrimSpace("\t Goodbye hair!\n ")
fmt.Printf("%q", s) // "Goodbye hair!"
s := strings.Repeat("da", 2) // "dada"
Reverse a UTF-8 encoded string
This function returns a string with the UTF-8 encoded characters of s in reverse order. Invalid UTF-8 sequences, if any, will be reversed byte by byte.
func ReverseUTF8(s string) string {
    res := make([]byte, len(s))
    prevPos, resPos := 0, len(s)
    for pos := range s {
        resPos -= pos - prevPos
        copy(res[resPos:], s[prevPos:pos])
        prevPos = pos
    }
    copy(res[0:], s[prevPos:])
    return string(res)
}
for _, s := range []string{
	"Ångström",
	"Hello, 世界",
	"\xff\xfe\xfd", // invalid UTF-8
} {
	fmt.Printf("%q\n", ReverseUTF8(s))
}
"mörtsgnÅ"
"界世 ,olleH"
"\xfd\xfe\xff"
_______________________________________
The Format method formats a time.Time.
The time.Parse function parses a date string.
func (t Time) Format(layout string) string
func Parse(layout, value string) (Time, error)
input := "2017-08-31"
layout := "2006-01-02"
t, _ := time.Parse(layout, input)
fmt.Println(t)                       // 2017-08-31 00:00:00 +0000 UTC
fmt.Println(t.Format("02-Jan-2006")) // 31-Aug-2017
Each Time has an associated Location, which is used for display purposes.

The method In returns a time with a specific location. Changing the location in this way changes only the presentation; it does not change the instant in time.

Here is a convenience function that changes the location associated with a time.
// TimeIn returns the time in UTC if the name is "" or "UTC".
// It returns the local time if the name is "Local".
// Otherwise, the name is taken to be a location name in
// the IANA Time Zone database, such as "Africa/Lagos".
func TimeIn(t time.Time, name string) (time.Time, error) {
    loc, err := time.LoadLocation(name)
    if err == nil {
        t = t.In(loc)
    }
    return t, err
}
for _, name := range []string{
	"",
	"Local",
	"Asia/Shanghai",
	"America/Metropolis",
} {
	t, err := TimeIn(time.Now(), name)
	if err == nil {
		fmt.Println(t.Location(), t.Format("15:04"))
	} else {
		fmt.Println(name, "<time unknown>")
	}
}
UTC 19:32
Local 20:32
Asia/Shanghai 03:32
America/Metropolis <time unknown>
Use time.Now and one of time.Unix or time.UnixNano to get a timestamp.
now := time.Now()      // current local time
sec := now.Unix()      // number of seconds since January 1, 1970 UTC
nsec := now.UnixNano() // number of nanoseconds since January 1, 1970 UTC

fmt.Println(now)  // time.Time
fmt.Println(sec)  // int64
fmt.Println(nsec) // int64
2009-11-10 23:00:00 +0000 UTC m=+0.000000000
1257894000
1257894000000000000
Get year, month, day from time
The Date function returns the year, month and day of a time.Time.
year, month, day := time.Now().Date()
fmt.Println(year, month, day)      // For example 2009 November 10
fmt.Println(year, int(month), day) // For example 2009 11 10
type Month int

const (
	January Month = 1 + iota
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)

How to find the day of week

The Weekday function returns returns the day of the week of a time.Time.
weekday := time.Now().Weekday()
fmt.Println(weekday)      // "Tuesday"
fmt.Println(int(weekday)) // "2"
Type Weekday
The time.Weekday type specifies a day of the week (Sunday = 0, …).
type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)
Days between two dates
func main() {
    // The leap year 2000 had 366 days.
    t1 := Date(2000, 1, 1)
    t2 := Date(2001, 1, 1)
    days := t2.Sub(t1).Hours() / 24
    fmt.Println(days) // 366
}
func Date(year, month, day int) time.Time {
    return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
Days in a month
To compute the last day of a month, you can use the fact that time.Date accepts values outside their usual ranges – the values are normalized during the conversion.

To compute the number of days in February, look at the day before March 1.
func main() {
    t := Date(2000, 3, 0) // the day before 2000-03-01
    fmt.Println(t)        // 2000-02-29 00:00:00 +0000 UTC
    fmt.Println(t.Day())  // 29
}

func Date(year, month, day int) time.Time {
    return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
AddDate normalizes its result in the same way. For example, adding one month to October 31 yields December 1, the normalized form of November 31.

t = Date(2000, 10, 31).AddDate(0, 1, 0) // a month after October 31
fmt.Println(t)                          // 2000-12-01 00:00:00 +0000 UTC
Measure execution time

start := time.Now()
// Code to measure
duration := time.Since(start)

// Formatted string, such as "2h3m0.5s" or "4.503μs"
fmt.Println(duration)

// Nanoseconds as int64
fmt.Println(duration.Nanoseconds())
Measure a function call
You can track the execution time of a complete function call with this one-liner, which logs the result to the standard error stream.
func foo() {
    defer duration(track("foo"))
    // Code to measure
}
func track(msg string) (string, time.Time) {
    return msg, time.Now()
}

func duration(msg string, start time.Time) {
    log.Printf("%v: %v\n", msg, time.Since(start))
}

A struct is a typed collection of fields, useful for grouping data into records.
type Student struct {
    Name string
    Age  int
}

var a Student    // a == Student{"", 0}
a.Name = "Alice" // a == Student{"Alice", 0}

var pa *Student   // pa == nil
pa = new(Student) // pa == &Student{"", 0}
pa.Name = "Alice" // pa == &Student{"Alice", 0}

b := Student{ // b == Student{"Bob", 0}
    Name: "Bob",
}

pb := &Student{ // pb == &Student{"Bob", 8}
    Name: "Bob",
    Age:  8,
}

c := Student{"Cecilia", 5} // c == Student{"Cecilia", 5}

Basics
Structs and arrays are copied when used in assignments and passed as arguments to functions. With pointers this can be avoided.

Pointers store addresses of objects. The addresses can be passed around more efficiently than the actual objects.

A pointer has type *T. The keyword new allocates a new object and returns its address.
type Student struct {
    Name string
}

var ps *Student = new(Student) // ps holds the address of the new struct
ps := new(Student)
s := Student{"Alice"} // s holds the actual struct
ps := &s              // ps holds the address of the struct
The & operator can also be used with composite literals. The two lines above can be written as

ps := &Student{"Alice"}
Pointer indirection
For a pointer x, the pointer indirection *x denotes the value which x points to. Pointer indirection is rarely used, since Go can automatically take the address of a variable.
ps := new(Student)
ps.Name = "Alice" // same as (*ps).Name = "Alice"
Pointers as parameters
When using a pointer to modify an object, you’re affecting all code that uses the object.
/ Bob is a function that has no effect.
func Bob(s Student) {
    s.Name = "Bob" // changes only the local copy
}

// Charlie sets pp.Name to "Charlie".
func Charlie(ps *Student) {
    ps.Name = "Charlie"
}

func main() {
    s := Student{"Alice"}

    Bob(s)
    fmt.Println(s) // prints {Alice}

    Charlie(&s)
    fmt.Println(s) // prints {Charlie}
}
Untyped numeric constants with no limits
const a uint = 17
const b = 55
An untyped constant has no limits. When it’s used in a context that requires a type, a type will be inferred and a limit applied.
const big = 10000000000  // Ok, even thought it's too big for an int.
const bigger = big * 100 // Still ok.
var i int = big / 100    // No problem: the new result fits in an int.

// Compile time error: "constant 10000000000 overflows int"
var j int = big
The inferred type is determined by the syntax of the value:

123 gets type int, and
123.4 becomes a float64.
The other possibilities are rune (alias for int32) and complex128.

Enumerations
Go does not have enumerated types. Instead, you can use the special name iota in a single const declaration to get a series of increasing values. When an initialization expression is omitted for a const, it reuses the preceding expressio
const (
    red = iota // red == 0
    blue       // blue == 1
    green      // green == 2
)
Enum with String function
yourbasic.org/golang
Basic solution
A group of constants enumerated with iota might do the job.
const (
    Sunday    int = iota // Sunday == 0
    Monday               // Monday == 1
    Tuesday              // Tuesday == 2
    Wednesday            // …
    Thursday
    Friday
    Saturday
)
The iota keyword represents successive integer constants 0, 1, 2,…
It resets to 0 whenever the word const appears in the source code,
and increments after each const specification.
In the example, we also rely on the fact that expressions are implicitly repeated in a paren­thesized const declaration – this indicates a repetition of the preceding expression and its type.

type Suite int

const (
    Spades Suite = iota
    Hearts
    Diamonds
    Clubs
)
and give it a String function:

func (s Suite) String() string {
	return [...]string{"Spades", "Hearts", "Diamonds", "Clubs"}[s]
}
Here is the new type in action.

var s Suite = Hearts
fmt.Print(s)
switch s {
case Spades:
	fmt.Println(" are best.")
case Hearts:
	fmt.Println(" are second best.")
default:
	fmt.Println(" aren't very good.")
}
Hearts are second best.
Make slices, maps and channels
yourbasic.org/golang
Slices, maps and channels can be created with the built-in make function. The memory is initialized with zero values.

Call	Type	Description
make(T, n)	slice	slice of type T with length n
make(T, n, c)		capacity c
make(T)	map	map of type T
make(T, n)		initial room for approximately n elements
make(T)	channel	unbuffered channel of type T
make(T, n)		buffered channel with buffer size n
s := make([]int, 10, 100)      // slice with len(s) == 10, cap(s) == 100
m := make(map[string]int, 100) // map with initial room for ~100 elements
c := make(chan int, 10)        // channel with a buffer size of 10
Slices, arrays and maps can also be created with composite literals.

s := []string{"f", "o", "o"} // slice with len(s) == 3, cap(s) == 3
a := [...]int{1, 2}          // array with len(a) == 2
m := map[string]float64{     // map with two key-value elements
    "e":  2.71828,
    "pi": 3.1416,
}
Methods for all types
yourbasic.org/golang
Any type declared in a type definition can have methods attached.

A method is a function with a receiver argument.
The receiver appears between the func keyword and the method name.
You can define methods on any type declared in a type definition.
In this example, the Value method is associated with MyType. The method receiver is called p.

type MyType struct {
    n int
}

func (p *MyType) Value() int { return p.n }

func main() {
    pm := new(MyType)
    fmt.Println(pm.Value()) // 0 (zero value)
}
If you convert the value to a different type, the new value will have the methods of the new type, not those of the old type.

type MyInt int

func (m MyInt) Positive() bool { return m > 0 }

func main() {
    var m MyInt = 2
    m = m * m // The operators of the underlying type still apply.

    fmt.Println(m.Positive())         // true
    fmt.Println(MyInt(-1).Positive()) // false

    var n int
    n = int(m) // The conversion is required.
    n = m      // ILLEGAL
}
../main.go:14:4: cannot use m (type MyInt) as type int in assignment
Basics
An interface type consists of a set of method signatures. A variable of interface type can hold any value that implements these methods.

In this example both Temp and *Point implement the MyStringer interface.

type MyStringer interface {
	String() string
}
type Temp int

func (t Temp) String() string {
	return strconv.Itoa(int(t)) + " °C"
}

type Point struct {
	x, y int
}

func (p *Point) String() string {
	return fmt.Sprintf("(%d,%d)", p.x, p.y)
}
Actually, *Temp also implements MyStringer, since the method set of a pointer type *T is the set of all methods with receiver *T or T.

When you call a method on an interface value, the method of its underlying type is executed.

var x MyStringer

x = Temp(24)
fmt.Println(x.String()) // 24 °C

x = &Point{1, 2}
fmt.Println(x.String()) // (1,2)
Structural typing
A type implements an interface by implementing its methods. No explicit declaration is required.

In fact, the Temp, *Temp and *Point types also implement the standard library fmt.Stringer interface. The String method in this interface is used to print values passed as an operand to functions such as fmt.Println.

var x MyStringer

x = Temp(24)
fmt.Println(x) // 24 °C

x = &Point{1, 2}
fmt.Println(x) // (1,2)
The empty interface
The interface type that specifies no methods is known as the empty interface.

interface{}
An empty interface can hold values of any type since every type implements at least zero methods.

var x interface{}

x = 2.4
fmt.Println(x) // 2.4

x = &Point{1, 2}
fmt.Println(x) // (1,2)
The fmt.Println function is a chief example. It takes any number of arguments of any type.

func Println(a ...interface{}) (n int, err error)
Interface values
An interface value consists of a concrete value and a dynamic type: [Value, Type]

In a call to fmt.Printf, you can use %v to print the concrete value and %T to print the dynamic type.

var x MyStringer
fmt.Printf("%v %T\n", x, x) // <nil> <nil>

x = Temp(24)
fmt.Printf("%v %T\n", x, x) // 24 °C main.Temp

x = &Point{1, 2}
fmt.Printf("%v %T\n", x, x) // (1,2) *main.Point

x = (*Point)(nil)
fmt.Printf("%v %T\n", x, x) // <nil> *main.Point
The zero value of an interface type is nil, which is represented as [nil, nil].

Calling a method on a nil interface is a run-time error. However, it’s quite common to write methods that can handle a receiver value [nil, Type], where Type isn’t nil.

You can also use type assertions, type switches and reflection to access the dynamic type of an interface value. Find the type of an object has more details.

Equality
Two interface values are equal

if they have equal concrete values and identical dynamic types,
or if both are nil.
A value t of interface type T and a value x of non-interface type X are equal if

t’s concrete value is equal to x
and t’s dynamic type is identical to X.
var x MyStringer
fmt.Println(x == nil) // true

x = (*Point)(nil)
fmt.Println(x == nil) // false
In the second print statement, the concrete value of x equals nil, but its dynamic type is *Point, which is not nil.

the empty interface in GO corresponds to a void pointer in C and an object reference in Java
it specifies zero methods
interface {}
a variable of empty interface type can hold values of any type since every type implements at least zero methods

var a interface{}

a = 24
fmt.Printf("[%v, %T]\n", a, a) // "[24, int]"

a = &Point{1, 2}
fmt.Printf("[%v, %T]\n", a, a) // "[(1,2), *main.Point]"
The fmt.Println function is a chief example from the standard library. It takes any number of arguments of any type.

func Println(a ...interface{}) (n int, err error)
Naed Return values
In GO return parameters may be named and used as regular variables when the function return these variables
func f() (i int, s string) {
    i = 17
    s = "abc"
    return // same as return i, s
}
func RedFull (r Reader ,buf []byte)(n int,err error) {
    for len(buf) >0 && err==nil {
        var nr int
        nr ,err = r.Read(buf)
        n+=nrbuf = buf[nr:]
    }
    return
    }
}

go has two different error handling mechaniss
most functions return errors
only truely unrecoverable conditons such as out of range index ,produce run time exceptions
go multivalued return makes it easy to return a detailed error mesage alongside the normal return value
By convention such messages have type error a smiple build in interface
type error interface {
    Error() string
}
os.Open functin returns a non-nil errro when it fails to open the file
func Open(name String) (file *File, err error)
f,err:=os.open('filename.txt')
if err !=nil {
    log.Fatal(err)
}
Custom errors
To create a simpsimple string only error you can use errors.New

err := errors.New("Houston we have a problem")
simple errors
// simple string based errors
err1 := errors.New("math: square root of the neagative number")
// with formating
err2 : = = fmt.Errorf(:math : square root of a neagative number %g",x)
custom errors  with data
To define a custom error type you must satisfy the predeclared error interface
type SyntexError struct {
    Line int
    Col int
}
func (e *SyntexError) Error() string {
    return fmt.Sprintf("%d:%d: syntex error ",e.Line,e.Col)

}

type InternalError struct {
    Path string
}
func (e *InternalError) Error() string {
    return fmt.Sprintf("parse %v:internal error",e.Path)
}

if Foo is a function that can return a Syntex Error or an Internal  eror you may handle
the two cases lke that

if err: = Foo(); err ! = nil {

}

In go how to recover from the panic
the built in  recover function can be used to regain  the contorol from the panicking situation
A call to recover stops the unwinding and return the argument passed to the panic
if the goroutine is not panicking returns nil
because the only code that runs while unwinding is inside deferred functions recover is only useful only inside the deferred functions

func main () {
    n:= foo()
    fmt.Println("main received",n )

}
function foo() int {
    defer func() {
        if err:= recover(); err!=nil {
            fmt.Println(err)
        }
    }()
    m:=1
    panic("foo:fail")
    m=2
    return m
}

foo: fail
main received 0
since panic roccured before foo returned a value , n still has its initial zero value

To return  a value during a panic you must use named return value

func main () {
    n:= foo()
    fmt.Println("main received ",n)

}
  func foo() (m int) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			m = 2
		}
	}()
	m = 1
	panic("foo: fail")
	m = 3
	return m
}
foo: fail
main received 2
stack traces are tyically printed to the console when a unexpected occurs they can be useful in debugging
not only do you see where the error happened but also how the program arreived in the place
he stack trace can be read from the bottom up:

testing.(*T).Run has called testing.tRunner,
which has called bit.TestMax,
which has called bit.(*Set).Max,
which has called panic,
which has called testing.tRunner.func1.
The indented lines show the source file and line number at which the function was called. The hexadecimal numbers refer to parameter values, including values of pointers and internal data structures. Stack Traces in Go has more details.

Print a stack trace
To print the stack trace for the current goroutine, use debug.PrintStack from package runtime/debug.

You can also examine the current stack trace programatically by calling runtime.Stack.

Level of detail
The GOTRACEBACK variable controls the amount of output generated when a Go program fails.

GOTRACEBACK=none omits the goroutine stack traces entirely.
GOTRACEBACK=single (the default) prints a stack trace for the current goroutine, eliding functions internal to the run-time system. The failure prints stack traces for all goroutines if there is no current goroutine or the failure is internal to the run-time.
GOTRACEBACK=all adds stack traces for all user-created goroutines.
GOTRACEBACK=system is like all but adds stack frames for run-time functions and shows goroutines created internally by the run-time.
switch time.Now().Weekday() {
    case time.Saturday:
    fmt.Println(" today is Saturday")
    case time.Sunday:
    fmt.Println("today is Sunday")
    default:
    fmt.Println('today is week day")
    a switch statement runs the first case equal to the condition expresson
    The cases are evaluated from top to bottom stopping when a case succeds
    unlike c and java the case expressions do not need to be constants
    No condition
    hour:= time.Now().Hour()
    switch { // same as switch true
case hour<12:
fmt.Println("good morning")
case hour<17:
fmt.println(" good aftrnoon")
default:
fmt.Println(" good evening")

    }
}

case lists

func whiteSpace(c rune) bool {
    switch c {
    case ' ', '\t', '\n', '\f', '\r':
        return true
    }
    return false
}

fallthrough

switch 2 {
    case 1:
    fmt.Println("1")
    fallthrough
    case 2 :
    fmt.println("2")
    fallthrough
    case 3:
    fmt.Println('3")

}
2
3

A fallthrough statement transfers control to the next case.
It may be used only as the final statement in a clause.

Exit a switch
A break statement terminates execution of the innermost for,switch or select statement
if you need to break out of a surrounding loop not the switch you can put a label on the loop and break of that label.
Loop:
    for _, ch := range "a b\nc" {
        switch ch {
        case ' ': // skip space
            break
        case '\n': // break at newline
            break Loop
        default:
            fmt.Printf("%c\n", ch)
        }
    }



    a
b

Execution order

function Foo( n int  int {
    fmt.Println(n)
    return n
}

func main() {
    switch Foo(2)

    case Foo(1),Foo(2),Foo(3) :

    fmt.Println("first case")
    fallthrough
    case Foo(4) :
    fmt.Println("second case")
}
2
1
2
First case
Second case

First the switch expression is evaluated once.
Then case expressions are evaluated left-to-right and top-to-bottom:
the first one that equals the switch expression triggers execution of the statements of the associated case,
the other cases are skipped.

Three component loop
sum :=0
for i:=1 ; i<5; i++ {
    sum +=1

}
fmt.Println(sum)

1The init statement, i := 1, runs.
2The condition, i < 5, is computed.
3If true, the loop body runs,
otherwise the loop is done.
The post statement, i++, runs.
4Back to step 2.
The scope of i is limited to the loop.
While loop
If you skip the init and post statements, you get a while loop.

n := 1
for n < 5 {
    n *= 2
}
fmt.Println(n) // 8 (1*2*2*2)
The condition, n < 5, is computed.
If true, the loop body runs,
otherwise the loop is done.
Back to step 1.
Infinite loop
If you skip the condition as well, you get an infinite loop.

sum := 0
for {
    sum++ // repeated forever
}
fmt.Println(sum) // never reached
For-each loop
Looping over elements in slices, arrays, maps, channels or strings is often better done with a range loop.
strings:= []string{ "hello" ,"world"}
for i,s :=range strings {
    fmt.Println(i,s)

}
0 hello
1 world
Range statements iterate over slices, arrays, strings, maps or channels.
a := []string{"Foo", "Bar"}
for i, s := range a {
    fmt.Println(i, s)
}

The range expression, a, is evaluated once before beginning the loop.
The iteration values are assigned to the respective iteration variables, i and s, as in an assignment statement.
The second iteration variable is optional.
If a slice or map is nil, the number of iterations is 0.
Strings
For a string, the loop iterates over Unicode code points.
for i, ch := range "日本語" {
    fmt.Printf("%#U starts at byte position %d\n", ch, i)
}
U+65E5 '日' starts at byte position 0
U+672C '本' starts at byte position 3
U+8A9E '語' starts at byte position 6

The index is the first byte of a UTF-8-encoded code point; the second value, of type rune, is the value of the code point.
For an invalid UTF-8 sequence, the second value will be 0xFFFD, and the iteration will advance a single byte.
Maps
The iteration order over maps is not specified and is not guaranteed to be the same from one iteration to the next.
m := map[string]int{
    "one":   1,
    "two":   2,
    "three": 3,
}
for k, v := range m {
    fmt.Println(k, v)
}
two 2
three 3
one 1

If a map entry that has not yet been reached is removed during iteration, this value will not be produced.
If a map entry is created during iteration, that entry may or may not be produced.
Channels
For channels, the iteration values are the successive values sent on the channel until closed.
ch := make(chan int)
go func() {
    ch <- 1
    ch <- 2
    ch <- 3
    close(ch)
}()
for n := range ch {
    fmt.Println(n)

    1
2
3

For a nil channel, the range loop blocks forever.

Do-while loop in 2 different ways
yourbasic.org/golang

do {
    work();
} while (condition);

for ok := true; ok; ok = condition {
    work()
}
for {
    work()
    if !condition {
        break
    }
}

How to use a defer statement
func main() {
    defer fmt.Println("World")
    fmt.Println("Hello")
}
Hello
World

Deferred calls are executed even when the function panics.

func main() {
    defer fmt.Println("World")
    panic("Stop")
    fmt.Println("Hello")
}World
panic: Stop

goroutine 1 [running]:
main.main()
    ../main.go:3 +0xa0

    Deferred function calls are executed in last-in-first-out order, and a deferred function’s arguments are evaluated when the defer statement executes.
    func main() {
    fmt.Println("Hello")
    for i := 1; i <= 3; i++ {
        defer fmt.Println(i)
    }
    fmt.Println("World")
}

Hello
World
3
2
1

Deferred functions may access and modify the surrounding function’s named return parameters.
func foo() (result string) {
    defer func() {
        result = "Change World" // change value at the very last moment
    }()
    return "Hello World"
}

Clean-up
Defer is commonly used to perform clean-up actions, such as closing a file or unlocking a mutex. In this example defer statements are used to ensure that all files are closed before leaving the CopyFile function, whichever way that happens.
func CopyFile(dstName, srcName string) (written int64, err error) {
    src, err := os.Open(srcName)
    if err != nil {
        return
    }
    defer src.Close()

    dst, err := os.Create(dstName)
    if err != nil {
        return
    }
    defer dst.Close()

    return io.Copy(dst, src)
}

Type assertions
Type assertion provides access to an interface 's concrete valueType assertions
A type assertions does not really convert an interface to another data type but it provides an access to an interface
concrete value which is typicaly you wantthe type assertion x.9T) assert that the cocrete value stored in x is of type T
and tha x is not nil
if T is not an interface it assert that the dynamic type of x is identical to Tempif T is an interface it asssers that the dynamic type of x implements Temp

var x interface{} ="foo"
var s string = x.(string)
fmt.Println(s) //"foo"
s,ok :=x.(string)
fmt.println(s,ok)// "foo true
n,ok := x.(int)
fmt.Println(s,ok) // "o false
n = x.(int)        // ILLEGAL
panic: interface conversion: interface {} is string, not int
Type switches
A type switch performs several type assertions in series and runs the first case with a matching type.
var x interface{} = "foo"

switch v := x.(type) {
case nil:
    fmt.Println("x is nil")            // here v has type interface{}
case int:
    fmt.Println("x is", v)             // here v has type int
case bool, string:
    fmt.Println("x is bool or string") // here v has type interface{}
default:
    fmt.Println("type unknown")        // here v has type interface{}
}

x is bool or string

Constructors deconstructed
Go doesn't have explicit constructors. The idiomatic way to set up new data structures is to use proper zero values coupled with factory functions.

Zero value
Try to make the default zero value useful and document its behavior. Sometimes this is all that’s needed.
StopWatch takes advantage of the useful zero values of time.Time, time.Duration and bool.
In turn, users of StopWatch can benefit from its useful zero value.
var clock StopWatch // Ready to use, no initialization needed.
Factory
If the zero value doesn’t suffice, use factory functions named NewFoo or just New.
scanner := bufio.NewScanner(os.Stdin)
err := errors.New("Houston, we have a problem")
Public vs. private
yourbasic.org/golang
PRIVATE KEEP OUT
A package is the smallest unit of private encap­sulation in Go.

All identifiers defined within a package are visible throughout that package.
When importing a package you can access only its exported identifiers.
An identifier is exported if it begins with a capital letter.
Exported and unexported identifiers are used to describe the public interface of a package and to guard against certain programming errors.

Warning: Unexported identifiers is not a security measure and it does not hide or protect any information.
package timer

import "time"

// A StopWatch is a simple clock utility.
// Its zero value is an idle clock with 0 total time.
type StopWatch struct {
    start   time.Time
    total   time.Duration
    running bool
}

// Start turns the clock on.
func (s *StopWatch) Start() {
    if !s.running {
        s.start = time.Now()
        s.running = true
    }
}
The StopWatch and its exported methods can be imported and used in a different package.

package main

import "timer"

func main() {
    clock := new(timer.StopWatch)
    clock.Start()
    if clock.running { // ILLEGAL
        // …
    }
}
Basic guidelines
For a given type, don’t mix value and pointer receivers.
If in doubt, use pointer receivers (they are safe and extendable).
Pointer receivers
You must use pointer receivers

if any method needs to mutate the receiver,
for structs that contain a sync.Mutex or similar synchronizing field (they musn’t be copied).
You probably want to use pointer receivers

for large structs or arrays (it can be more efficient),
in all other cases.
Value receivers
You probably want to use value receivers

for map, func and chan types,
for simple basic types such as int or string,
for small arrays or structs that are value types, with no mutable fields and no pointers.
You may want to use value receivers

for slices with methods that do not reslice or reallocate the slice.
Function types and values
yourbasic.org/golang
Functions in Go are first class citizens.

Function types and function values can be used and passed around just like other values.
type Operator func(x float64) float64

// Map applies op to each element of a.
func Map(op Operator, a []float64) []float64 {
    res := make([]float64, len(a))
    for i, x := range a {
        res[i] = op(x)
    }
    return res
}

func main() {
    op := math.Abs
    a := []float64{1, -2}
    b := Map(op, a)
    fmt.Println(b) // [1 2]

    c := Map(func(x float64) float64 { return 10 * x }, b)
    fmt.Println(c) // [10, 20]
}
Details
A function type describes the set of all functions with the same parameter and result types.

The value of an uninitialized variable of function type is nil.
The parameter names are optional.
The following two function types are identical.

func(x, y int) int
func(int, int) int

Anonymous functions and closures
yourbasic.org/golang
A function literal (or lambda) is a function without a name.
In this example a function literal is passed as the less argument to the sort.Slice function.
func Slice(slice interface{}, less func(i, j int) bool)

people := []string{"Alice", "Bob", "Dave"}
sort.Slice(people, func(i, j int) bool {
    return len(people[i]) < len(people[j])
})
fmt.Println(people)
// Output: [Bob Dave Alice]
You can also use an intermediate variable.

people := []string{"Alice", "Bob", "Dave"}
less := func(i, j int) bool {
    return len(people[i]) < len(people[j])
}
sort.Slice(people, less)
Note that the less function is a closure: it references the people variable, which is declared outside the function.

Closures
Function literals in Go are closures: they may refer to variables defined in an enclosing function. Such variables

are shared between the surrounding function and the function literal,
survive as long as they are accessible.
In this example, the function literal uses the local variable n from the enclosing scope to count the number of times it has been invoked.

// New returns a function Count.
// Count prints the number of times it has been invoked.
func New() (Count func()) {
    n := 0
    return func() {
        n++
        fmt.Println(n)
    }
}

func main() {
    f1, f2 := New(), New()
    f1() // 1
    f2() // 1 (different n)
    f1() // 2
    f2() // 2
}

Read a file line by line
file, err := os.Open("file.txt")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

scanner := bufio.NewScanner(file)
for scanner.Scan() {
    fmt.Println(scanner.Text())
}

if err := scanner.Err(); err != nil {
    log.Fatal(err)
}
Read a whole file into a string (byte slice)
b, err := ioutil.ReadFile("file.txt") // b has type []byte
if err != nil {
    log.Fatal(err)
}
s := string(b)
Read from stdin
Use a bufio.Scanner to read one line at a time from the standard input stream.
scanner := bufio.NewScanner(os.Stdin)
for scanner.Scan() {
	fmt.Println(scanner.Text())
}
if err := scanner.Err(); err != nil {
	log.Println(err)
}
Append text to a file
f, err := os.OpenFile("text.log",
	os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
if err != nil {
	log.Println(err)
}
defer f.Close()
if _, err := f.WriteString("text to append\n"); err != nil {
	log.Println(err)
}
How to write a log message to file
This code appends a log message to the file text.log. It creates the file if it doesn’t already exist.
f, err := os.OpenFile("text.log",
	os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
if err != nil {
	log.Println(err)
}
defer f.Close()

logger := log.New(f, "prefix", log.LstdFlags)
logger.Println("text to append")
logger.Println("more text to append")
log.New creates a new log.Logger that writes to f.
The prefix appears at the beginning of each generated log line.
The flag argument defines which text to prefix to each log entry.
Find the current working directory
yourbasic.org/golang
Use os.Executable to find the path name for the executable that started the current process.
Use filepath.Dir in package path/filepath to extract the path’s directory.
path, err := os.Executable()
if err != nil {
    log.Printf(err)
}
dir := filepath.Dir(path)
fmt.Println(path) // for example /home/user/main
fmt.Println(dir)  // for example /home/user
Warning: There is no guarantee that the path is still pointing to the correct executable. If a symlink was used to start the process, depending on the operating system, the result might be the symlink or the path it pointed to. If a stable result is needed, path/filepath.EvalSymlinks might help.
List files and folders in a directory
yourbasic.org/golang
Use the ioutil.ReadDir function in package io/ioutil. It returns a sorted slice containing elements of type os.FileInfo.

The code in this example prints a sorted list of all file names in the current directory.

files, err := ioutil.ReadDir(".")
if err != nil {
    log.Fatal(err)
}
for _, f := range files {
    fmt.Println(f.Name())
}
Example output:

dev
etc
tmp
usr
Visit files and folders in a directory tree
yourbasic.org/golang
Use the filepath.Walk function in package path/filepath.

It walks a file tree calling a function of type filepath.WalkFunc for each file or directory in the tree, including the root.
The files are walked in lexical order.
Symbolic links are not followed.
The code in this example lists the paths and sizes of all files and directories in the file tree rooted at the current directory.

err := filepath.Walk(".",
    func(path string, info os.FileInfo, err error) error {
    if err != nil {
        return err
    }
    fmt.Println(path, info.Size())
    return nil
})
if err != nil {
    log.Println(err)
}
Create a temporary file or directory
yourbasic.org/golang
File
Use ioutil.TempFile in package io/ioutil to create a temporary file.

file, err := ioutil.TempFile("dir", "prefix")
if err != nil {
    log.Fatal(err)
}
defer os.Remove(file.Name())

fmt.Println(file.Name()) // For example "dir/prefix054003078"
The call to ioutil.TempFile

creates a new file with a name starting with "prefix" in the directory "dir",
opens the file for reading and writing,
and returns the new *os.File.
To put the new file in os.TempDir(), the default directory for temporary files, call ioutil.TempFile with an empty directory string.

Directory
Use ioutil.TempDir in package io/ioutil to create a temporary directory.

dir, err := ioutil.TempDir("dir", "prefix")
if err != nil {
	log.Fatal(err)
}
defer os.RemoveAll(dir)
The call to ioutil.TempDir

creates a new directory with a name starting with "prefix" in the directory "dir"
and returns the path of the new directory.
To put the new directory in os.TempDir(), the default directory for temporary files, call ioutil.TempDir with an empty directory string.
Basics
The io.Writer interface represents an entity to which you can write a stream of bytes.

type Writer interface {
        Write(p []byte) (n int, err error)
}
Write writes up to len(p) bytes from p to the underlying data stream – it returns the number of bytes written and any error encountered that caused the write to stop early.

The standard library provides many Writer implementations, and Writers are accepted as input by many utilities.
Use a built-in writer
Since bytes.Buffer has a Write method you can write directly into the buffer using fmt.Fprintf.

var buf bytes.Buffer
fmt.Fprintf(&buf, "Size: %d MB.", 85)
s := buf.String()) // s == "Size: 85 MB."
Optimize string writes
Some Writers in the standard library have an additional WriteString method. This method can be more efficient than the standard Write method since it writes a string directly without allocating a byte slice.

You can take direct advantage of this optimization by using the io.WriteString() function.

func WriteString(w Writer, s string) (n int, err error)
If w implements a WriteString method, it is invoked directly. Otherwise, w.Write is called exactly once.
Access private fields with reflection
yourbasic.org/golang
With reflection it's possible to read, but not write, unexported fields of a struct defined in another package.

In this example, we access the unexported field len in the List struct in package container/list:
package list

type List struct {
    root Element
    len  int
}
package main

import (
    "container/list"
    "fmt"
    "reflect"
)

func main() {
    l := list.New()
    l.PushFront("foo")
    l.PushFront("bar")

    // Get a reflect.Value fv for the unexported field len.
    fv := reflect.ValueOf(l).Elem().FieldByName("len")
    fmt.Println(fv.Int()) // 2

    // Try to set the value of len.
    fv.Set(reflect.ValueOf(3)) // ILLEGAL
}
The rand.Shuffle function in package math/rand shuffles an input sequence using a given swap function.

a := []int{1, 2, 3, 4, 5, 6, 7, 8}
rand.Seed(time.Now().UnixNano())
rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
[5 8 6 4 3 7 2 1]
Warning: Without the call to rand.Seed you will get the same sequence of pseudo­random numbers each time you run the program.
What’s a seed: Random number generators
yourbasic.org
In reality pseudo­random numbers aren't random at all. They are computed using a fixed determi­nistic algorithm.
The seed is a starting point for a sequence of pseudorandom numbers. If you start from the same seed, you get the very same sequence. This can be quite useful for debugging.

If you want a different sequence of numbers each time, you can use the current time as a seed.
This generator produces a sequence of 97 different numbers, then it starts over again. The seed decides at what number the sequence will start.
/ New returns a pseudorandom number generator Rand with a given seed.
// Every time you call Rand, you get a new "random" number.
func New(seed int) (Rand func() int) {
    current := seed
    return func() int {
        next := (17 * current) % 97
        current = next
        return next
    }
}

func main() {
    rand1 := New(1)
    fmt.Println(rand1(), rand1(), rand1())

    rand2 := New(2)
    fmt.Println(rand2(), rand2(), rand2())
}
17 95 63
34 93 29
The random number generators you’ll find in most programming langauges work just like this, but of course they use a smarter function. Ideally, you want a long sequence with good random properties computed by a function which uses only cheap arithmetic operations. For example, you would typically want to avoid the % modulus operator.
Goroutines are lightweight threads
yourbasic.org/golang

The go statement runs a func­tion in a sepa­rate thread of execu­tion.

You can start a new thread of execution, a goroutine, with the go statement. It runs a function in a different, newly created, goroutine. All goroutines in a single program share the same address space.

go list.Sort() // Run list.Sort in parallel; don’t wait for it.
The following program will print “Hello from main goroutine”. It might also print “Hello from another goroutine”, depending on which of the two goroutines finish first.
func main() {
    go fmt.Println("Hello from another goroutine")
    fmt.Println("Hello from main goroutine")

    // At this point the program execution stops and all
    // active goroutines are killed.
}
The next program will, most likely, print both “Hello from main goroutine” and “Hello from another goroutine”. They may be printed in any order. Yet another possibility is that the second goroutine is extremely slow and doesn’t print its message before the program ends.
func main() {
    go fmt.Println("Hello from another goroutine")
    fmt.Println("Hello from main goroutine")

    time.Sleep(time.Second) // wait for other goroutine to finish
}
Here is a somewhat more realistic example, where we define a function that uses concurrency to postpone an event.
// Publish prints text to stdout after the given time has expired.
// It doesn’t block but returns right away.
func Publish(text string, delay time.Duration) {
    go func() {
        time.Sleep(delay)
        fmt.Println("BREAKING NEWS:", text)
    }() // Note the parentheses. We must call the anonymous function.
}
This is how you might use the Publish function.

func main() {
    Publish("A goroutine starts a new thread.", 5*time.Second)
    fmt.Println("Let’s hope the news will published before I leave.")

    // Wait for the news to be published.
    time.Sleep(10 * time.Second)

    fmt.Println("Ten seconds later: I’m leaving now.")
}
go run publish1.go
Let’s hope the news will published before I leave.
BREAKING NEWS: A goroutine starts a new thread.
Ten seconds later: I’m leaving now
In general it’s not possible to arrange for threads to wait for each other by sleeping. Go’s main method for synchronization is to use channels.

mplementation
Goroutines are lightweight, costing little more than the allocation of stack space. The stacks start small and grow by allocating and freeing heap storage as required.

Internally goroutines act like coroutines that are multiplexed among multiple operating system threads. If one goroutine blocks an OS thread, for example waiting for input, other goroutines in this thread will migrate so that they may continue running.
Channels offer synchronized communication
yourbasic.org/golang
A channel is a mechanism for goroutines to synchronize execution and communicate by passing values.
A new channel value can be made using the built-in function make.
// unbuffered channel of ints
ic := make(chan int)

// buffered channel with room for 10 strings
sc := make(chan string, 10)
To send a value on a channel, use <- as a binary operator. To receive a value on a channel, use it as a unary operator.
ic <- 3   // Send 3 on the channel.
n := <-sc // Receive a string from the channel
The <- operator specifies the channel direction, send or receive. If no direction is given, the channel is bi-directional.
Buffered and unbuffered channels
If the capacity of a channel is zero or absent, the channel is unbuffered and the sender blocks until the receiver has received the value.
If the channel has a buffer, the sender blocks only until the value has been copied to the buffer; if the buffer is full, this means waiting until some receiver has retrieved a value.
Receivers always block until there is data to receive.
Sending or receiving from a nil channel blocks forever.
Closing a channel
The close function records that no more values will be sent on a channel. Note that it is only necessary to close a channel if a receiver is looking for a close.

After calling close, and after any previously sent values have been received, receive operations will return a zero value without blocking.
A multi-valued receive operation additionally returns an indication of whether the channel is closed.
Sending to or closing a closed channel causes a run-time panic. Closing a nil channel also causes a run-time panic.
ch := make(chan string)
go func() {
    ch <- "Hello!"
    close(ch)
}()

fmt.Println(<-ch) // Print "Hello!".
fmt.Println(<-ch) // Print the zero value "" without blocking.
fmt.Println(<-ch) // Once again print "".
v, ok := <-ch     // v is "", ok is false.

// Receive values from ch until closed.
for v := range ch {
    fmt.Println(v) // Will not be executed.
}
n the following example we let the Publish function return a channel, which is used to broadcast a message when the text has been published.
// Publish prints text to stdout after the given time has expired.
// It closes the wait channel when the text has been published.
func Publish(text string, delay time.Duration) (wait <-chan struct{}) {
	ch := make(chan struct{})
	go func() {
		time.Sleep(delay)
		fmt.Println(text)
		close(ch)
	}()
	return ch
}
Note that we use a channel of empty structs to indicate that the channel will only be used for signalling, not for passing data. This is how you might use the function.
wait := Publish("important news", 2 * time.Minute)
// Do some more work.
<-wait // Block until the text has been published.
Select waits on a group of channels
The select statement waits for multiple send or receive opera­tions simul­taneously.
/ blocks until there's data available on ch1 or ch2
select {
case <-ch1:
    fmt.Println("Received from ch1")
case <-ch2:
    fmt.Println("Received from ch2")
}
The statement blocks as a whole until one of the operations becomes unblocked.
If several cases can proceed, a single one of them will be chosen at random.
Send and receive operations on a nil channel block forever. This can be used to disable a channel in a select statement:
ch1 = nil // disables this channel
select {
case <-ch1:
    fmt.Println("Received from ch1") // will not happen
case <-ch2:
    fmt.Println("Received from ch2")
}
Default case
The default case is always able to proceed and runs if all other cases are blocked.


// never blocks
select {
case x := <-ch:
    fmt.Println("Received", x)
default:
    fmt.Println("Nothing available")
}
An infinite random binary sequence
rand := make(chan int)
for {
    select {
    case rand <- 0: // no statement
    case rand <- 1:
    }
    A blocking operation with a timeout
    select {
case news := <-AFP:
    fmt.Println(news)
case <-time.After(time.Minute):
    fmt.Println("Time out: No news in one minute")
}
The function time.After is part of the standard library; it waits for a specified time to elapse and then sends the current time on the returned channel.


A statement that blocks forever
select {}
A select statement blocks until at least one of it’s cases can proceed. With zero cases this will never happen.
Data races explained
yourbasic.org/golang
A data race happens when two goroutines access the same variable concur­rently, and at least one of the accesses is a write.


Data races are quite common and can be very hard to debug.

This function has a data race and it’s behavior is undefined. It may, for example, print the number 1. Try to figure out how that can happen – one possible explanation comes after the code.
func race() {
    wait := make(chan struct{})
    n := 0
    go func() {
        n++ // read, increment, write
        close(wait)
    }()
    n++ // conflicting access
    <-wait
    fmt.Println(n) // Output: <unspecified>
}

How to avoid data races
The only way to avoid data races is to synchronize access to all mutable data that is shared between threads. There are several ways to achieve this. In Go, you would normally use a channel or a lock. (Lower-lever mechanisms are available in the sync and sync/atomic packages.)

The preferred way to handle concurrent data access in Go is to use a channel to pass the actual data from one goroutine to the next. The motto is: “Don’t communicate by sharing memory; share memory by communicating.”

func sharingIsCaring() {
    ch := make(chan int)
    go func() {
        n := 0 // A local variable is only visible to one goroutine.
        n++
        ch <- n // The data leaves one goroutine...
    }()
    n := <-ch // ...and arrives safely in another.
    n++
    fmt.Println(n) // Output: 2
}
In this code the channel does double duty:

it passes the data from one goroutine to another,
and it acts as a point of synchronization.
The sending goroutine will wait for the other goroutine to receive the data and the receiving goroutine will wait for the other goroutine to send the data.

The Go memory model – the conditions under which reads of a variable in one goroutine can be guaranteed to observe values produced by writes to the same variable in a different goroutine – is quite complicated, but as long as you share all mutable data between goroutines through channels you are safe from data races.
How to detect data races
yourbasic.org/golang
Data races can happen easily and are hard to debug. Luckily, the Go runtime is often able to help.
Use -race to enable the built-in data race detector.
$ go test -race [packages]
$ go run -race [packages]
Here’s a program with a data race:
package main
import "fmt"

func main() {
    i := 0
    go func() {
        i++ // write
    }()
    fmt.Println(i) // concurrent read
}
Running this program with the -race options tells us that there’s a race between the write at line 7 and the read at line 9:


$ go run -race main.go
0
==================
WARNING: DATA RACE
Write by goroutine 6:
  main.main.func1()
      /tmp/main.go:7 +0x44

Previous read by main goroutine:
  main.main()
      /tmp/main.go:9 +0x7e

Goroutine 6 (running) created at:
  main.main()
      /tmp/main.go:8 +0x70
==================
Found 1 data race(s)
exit status 66
Details
The data race detector does not perform any static analysis. It checks the memory access in runtime and only for the code paths that are actually executed.

It runs on darwin/amd64, freebsd/amd64, linux/amd64 and windows/amd64.

The overhead varies, but typically there’s a 5-10x increase in memory usage, and 2-20x increase in execution time.
A deadlock happens when a group of goroutines are waiting for each other and none of them is able to proceed.
func main() {
	ch := make(chan int)
	ch <- 1
	fmt.Println(<-ch)
}
The program will get stuck on the channel send operation waiting forever for someone to read the value. Go is able to detect situations like this at runtime. Here is the output from our program:
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan send]:
main.main()
	.../deadlock.go:7 +0x6c
    Debugging tips
A goroutine can get stuck

either because it’s waiting for a channel or
because it is waiting for one of the locks in the sync package.
Common reasons are that

no other goroutine has access to the channel or the lock,
a group of goroutines are waiting for each other and none of them is able to proceed.
Currently Go only detects when the program as a whole freezes, not when a subset of goroutines get stuck.

With channels it’s often easy to figure out what caused a deadlock. Programs that make heavy use of mutexes can, on the other hand, be notoriously difficult to debug.
Waiting for goroutines
A sync.WaitGroup waits for a group of goroutines to finish.

var wg sync.WaitGroup
wg.Add(2)
go func() {
    // Do work.
    wg.Done()
}()
go func() {
    // Do work.
    wg.Done()
}()
wg.Wait()
First the main goroutine calls Add to set the number of goroutines to wait for.
Then two new goroutines run and call Done when finished.
At the same time, Wait is used to block until these two goroutines have finished.

Note: A WaitGroup must not be copied after first use.
Broadcast a signal on a channel
yourbasic.org/golang
When you close a channel, all readers receive a zero value.

n this example the Publish function returns a channel, which is used to broadcast a signal when a message has been published.
// Print text after the given time has expired.
// When done, the wait channel is closed.
func Publish(text string, delay time.Duration) (wait <-chan struct{}) {
    ch := make(chan struct{})
    go func() {
        time.Sleep(delay)
        fmt.Println("BREAKING NEWS:", text)
        close(ch) // Broadcast to all receivers.
    }()
    return ch
}
Notice that we use a channel of empty structs: struct{}. This clearly indicates that the channel will only be used for signalling, not for passing data.

This is how you may use the function.

func main() {
    wait := Publish("Channels let goroutines communicate.", 5*time.Second)
    fmt.Println("Waiting for news...")
    <-wait
    fmt.Println("Time to leave.")
}
Waiting for news...
BREAKING NEWS: Channels let goroutines communicate.
Time to leave.

How to kill a goroutine
One goroutine can't forcibly stop another.

To make a goroutine stoppable, let it listen for a stop signal on a channel.
quit := make(chan struct{})
go func() {
    for {
        select {
        case <-quit:
            return
        default:
            // …
        }
    }
}()
// …
close(quit)
Sometimes it’s convenient to use a single channel for both data and signalling.
// Generator returns a channel that produces the numbers 1, 2, 3,…
// To stop the underlying goroutine, close the channel.
func Generator() chan int {
    ch := make(chan int)
    go func() {
        n := 1
        for {
            select {
            case ch <- n:
                n++
            case <-ch:
                return
            }
        }
    }()
    return ch
}

func main() {
    number := Generator()
    fmt.Println(<-number)
    fmt.Println(<-number)
    close(number)
    // …
}
Timer and Ticker: events in the future
yourbasic.org/golang
Timers and Tickers let you execute code in the future, once or repeatedly.

Timeout (Timer)
time.After waits for a specified duration and then sends the current time on the returned channel:
select {
case news := <-AFP:
	fmt.Println(news)
case <-time.After(time.Hour):
	fmt.Println("No news in an hour.")
}
The underlying time.Timer will not be recovered by the garbage collector until the timer fires. If this is a concern, use time.NewTimer instead and call its Stop method when the timer is no longer needed:
for alive := true; alive; {
	timer := time.NewTimer(time.Hour)
	select {
	case news := <-AFP:
		timer.Stop()
		fmt.Println(news)
	case <-timer.C:
		alive = false
		fmt.Println("No news in an hour. Service aborting.")
	}
}
Repeat (Ticker)
time.Tick returns a channel that delivers clock ticks at even intervals:
go func() {
	for now := range time.Tick(time.Minute) {
		fmt.Println(now, statusUpdate())
	}
}()
The underlying time.Ticker will not be recovered by the garbage collector. If this is a concern, use time.NewTicker instead and call its Stop method when the ticker is no longer needed.

Wait, act and cancel
time.AfterFunc waits for a specified duration and then calls a function in its own goroutine. It returns a time.Timer that can be used to cancel the call:
func Foo() {
    timer = time.AfterFunc(time.Minute, func() {
        log.Println("Foo run for more than a minute.")
    })
    defer timer.Stop()

    // Do heavy work
}

Mutual exclusion lock (mutex)
yourbasic.org/golang
Mutexes let you synchronize data access by explicit locking, without channels.
Sometimes it’s more convenient to synchronize data access by explicit locking instead of using channels. The Go standard library offers a mutual exclusion lock, sync.Mutex, for this purpose.


Use with caution
For this type of locking to be safe, it’s crucial that all accesses to the shared data, both reads and writes, are performed only when a goroutine holds the lock. One mistake by a single goroutine is enough to introduce a data race and break the program.

Because of this you should consider designing a custom data structure with a clean API and make sure that all the synchronization is done internally.

In this example we build a safe and easy-to-use concurrent data structure, AtomicInt, that stores a single integer. Any number of goroutines can safely access this number through the Add and Value methods.
// AtomicInt is a concurrent data structure that holds an int.
// Its zero value is 0.
type AtomicInt struct {
    mu sync.Mutex // A lock than can be held by one goroutine at a time.
    n  int
}

// Add adds n to the AtomicInt as a single atomic operation.
func (a *AtomicInt) Add(n int) {
    a.mu.Lock() // Wait for the lock to be free and then take it.
    a.n += n
    a.mu.Unlock() // Release the lock.
}

// Value returns the value of a.
func (a *AtomicInt) Value() int {
    a.mu.Lock()
    n := a.n
    a.mu.Unlock()
    return n
}

func main() {
    wait := make(chan struct{})
    var n AtomicInt
    go func() {
        n.Add(1) // one access
        close(wait)
    }()
    n.Add(1) // another concurrent access
    <-wait
    fmt.Println(n.Value()) // 2
    Efficient parallel computation
yourbasic.org/golang
cpu
Dividing a large compu­tation into work units for parallel pro­cessing is more of an art than a science.

Here are some rules of thumb.

Each work unit should take about 100μs to 1ms to compute. If the units are too small, the adminis­trative over­head of divi­ding the problem and sched­uling sub-problems might be too large. If the units are too big, the whole computation may have to wait for a single slow work item to finish. This slowdown can happen for many reasons, such as scheduling, interrupts from other processes, and unfortunate memory layout. (Note that the number of work units is independent of the number of CPUs.)
Try to minimize the amount of data sharing. Concurrent writes can be very costly, particularly so if goroutines execute on separate CPUs. Sharing data for reading is often much less of a problem.
Strive for good locality when accessing data. If data can be kept in cache memory, data loading and storing will be dramatically faster. Once again, this is particularly important for writing.
Whatever strategies you are using, don’t forget to benchmark and profile your code.
The following example shows how to divide a costly computation and distribute it on all available CPUs. This is the code we want to optimize.
type Vector []float64

// Convolve computes w = u * v, where w[k] = Σ u[i]*v[j], i + j = k.
// Precondition: len(u) > 0, len(v) > 0.
func Convolve(u, v Vector) Vector {
    n := len(u) + len(v) - 1
    w := make(Vector, n)
    for k := 0; k < n; k++ {
        w[k] = mul(u, v, k)
    }
    return w
}

// mul returns Σ u[i]*v[j], i + j = k.
func mul(u, v Vector, k int) float64 {
    var res float64
    n := min(k+1, len(u))
    j := min(k, len(v)-1)
    for i := k - j; i < n; i, j = i+1, j-1 {
        res += u[i] * v[j]
    }
    return res
}
The idea is simple: identify work units of suitable size and then run each work unit in a separate goroutine. Here is a parallel version of Convolve.

func Convolve(u, v Vector) Vector {
    n := len(u) + len(v) - 1
    w := make(Vector, n)

    // Divide w into work units that take ~100μs-1ms to compute.
    size := max(1, 1000000/n)

    var wg sync.WaitGroup
    for i, j := 0, size; i < n; i, j = j, j+size {
        if j > n {
            j = n
        }
        // These goroutines share memory, but only for reading.
        wg.Add(1)
        go func(i, j int) {
            for k := i; k < j; k++ {
                w[k] = mul(u, v, k)
            }
            wg.Done()
        }(i, j)
    }
    wg.Wait()
    return w
}
When the work units have been defined, it’s often best to leave the scheduling to the runtime and the operating system. However, if needed, you can tell the runtime how many goroutines you want executing code simultaneously.

func init() {
    numcpu := runtime.NumCPU()
    runtime.GOMAXPROCS(numcpu) // Try to use all available CPUs.
}
Dependency injection principle
a top layer containing poem documents
a bogtom layer containing storage entities
a document object needs to access the sercicse of storage object tostore and retrieve the contents .Thus it would seem naturla to add a storage service directly to the document
poet surely need to write poems to a small notebook thus the lead programmer creates the document layer
type Poem struct {
content []byte
storage acmeStorageSrevices.PoemNoteBook
}

func NewPoem() *Poem {
return &Poem {
storage : acmeStorageServices.NewPoemNoteBook(),
}
}

func (p* Poem) Load(title string) {
p.content = p.storage.Load(title)
}
func (p* Poem) Save(title string) {
storage.save(title, p.content)
}
what if poet decides to write a poem on a napkin .The document layer has to be modified .We have created an unwanted dependincy on a particula r
storage type
Abstraction to the rescue
we replace the storage service by the abstraction of that service
type PoemStorage interface {
Load(String) []byte
Save (string, []byte)

}

thus interface describes only a behavior and our poem object can call the interface functions without worringg about the object that implemnts this uinterface
so  now we can define the Poem struct withour any dependency on the storage layer
type Poem struct {
content [] byte
storage PoemStorage

}
we can now assign any type to storage that satisfy this interface

Adding Dependency injection
right now the poem only talks to an empty abstraction .Next step we need a way to connect the real storage to the poem
in the other words we need to inject a dependency on the Poemstorage object into the Poem Layer
we can d this for example through a constructer

 func NePoem (ps *PoemStorage) *Poem {
 return &Poem {
 storage:ps
 }
 }
 when called the constructer recevies an actual PoemStorage object .Yet the returned Poem still just talks to the abstract PoemStorage interface

 finally in main we can wire up all the higher level objects with their low level dependencies
 func main() {
 storage := newNapkin()
 poem:newPoem(sorage)

 }
 package main



 // full code
 import "fmt"
 The “inner ring”
 A Poem contains some poetry and an abstract storage reference.
 type Poem struct {
 	content []byte
 	storage PoemStorage
 }
 PoemStorage is just an interface that defines the behavior of a poem storage. This is all that Poem knows (and needs to know) about storing and retrieving poems. Nothing from the “outer ring” appears here.
 type PoemStorage interface {
 	Type() string        // Return a string describing the storage type.
 	Load(string) []byte  // Load a poem by name.
 	Save(string, []byte) // Save a poem by name.
 }
 NewPoem constructs a Poem object. We use this constructor to inject an object that satisfies the PoemStorage interface.
 func NewPoem(ps PoemStorage) *Poem {
 	return &Poem{
 		content: []byte("I am a poem from a " + ps.Type() + "."),
 		storage: ps,
 	}
 }
 Save simply calls Save on the interface type. The Poem object neither knows nor cares about which actual storage object receives this method call.
 func (p *Poem) Save(name string) {
 	p.storage.Save(name, p.content)
 }
 Load also invokes the injected storage object without knowing it.
 func (p *Poem) Load(name string) {
 	p.content = p.storage.Load(name)
 }
 String makes Poem a Stringer, allowing us to drop it anywhere a string would be expected.
 func (p *Poem) String() string {
 	return string(p.content)
 }
 The “outer ring”
 The notebook
 A Notebook is the classic storage device of a poet.
 type Notebook struct {
 	poems map[string][]byte
 }

 func NewNotebook() *Notebook {
 	return &Notebook{
 		poems: map[string][]byte{},
 	}
 }
 After adding Save and Load, Notebook implicitly satisfies PoemStorage.
 func (n *Notebook) Save(name string, contents []byte) {
 	n.poems[name] = contents
 }

 func (n *Notebook) Load(name string) []byte {
 	return n.poems[name]
 }
 Type returns an informal description of the storage type.
 func (n *Notebook) Type() string {
 	return "Notebook"
 }
 A Napkin is the emergency storage device of a poet. It can store only one poem.
 type Napkin struct {
 	poem []byte
 }

 func NewNapkin() *Napkin {
 	return &Napkin{
 		poem: []byte{},
 	}
 }

 func (n *Napkin) Save(name string, contents []byte) {
 	n.poem = contents
 }

 func (n *Napkin) Load(name string) []byte {
 	return n.poem
 }

 func (n *Napkin) Type() string {
 	return "Napkin"
 }