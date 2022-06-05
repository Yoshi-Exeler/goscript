Declaring a variable:
let x: string;
let <name>: <type>

Declaring and initializing a variable:
let v: string = x
let <name>: <type> = <value>

Setting the value of a variable:
v = x
<name> = <value>

Declaring and initializing an array:
let x: []string = []
let <name>: []<type> = <array literal>

count loop:
for let i: int = 0; i < 10; i++ {}
for <declaration>;<condition>;<action> {}

foreach loop:
foreach index, element in list {}
foreach <name1>, <name2> in <name3> {}

function definition:
func myfunc(a: string, b: int) => x: int, y: string {}
func <name>(<p1>: <type>, <p2>: <type>) => <rp1>: <type>, <rp2>: <type> {}

List of primitives:
int8, int16, int32, int64
uint8, uint16, uint32, uint64
string, char
byte
float32, float64
any

Composing Types:
struct MyStruct {
    A: string
    B: int64
    C: byte
}
struct <name> {
    <name>: <type>
    <name>: <type>
    <name>: <type>
}

Modules:
a source file inherits the module name of its parent directory
all files in a directory are joined together internally
cyclic imports are allowed

main.gs
/module1/a.gs
/module1/b.gs
/module1/c.gs
/module2/c.gs
/module2/x.gs
/module3/submodule1/a.gs

the main file must be marked by beginning with
application MyApp

to import a package use the import statement:
import "./module1"
import "./<module name>"
Means import this this module directly using its default name

import "./module1" as mymodule
import "./<module name>" as <module alias>
Means import this this module directly using the name alias mymodule

import "./module1/submodule1/"
import "./<module name>"
Means import this submodule directly using its default name

external remotemodule from "http://someremotething.com/a"
external <alias> from "<source url>"
Means declare the module remotemodule importable, referring to the folder at the url. folder will be git cloned from the git remote

import "ext/module1"
import "<standard module name>"
Means import module1 from the standard library. ext/ modules come from the an external source

import "std/module1"
import "<standard module name>"
Means import module1 from the standard library. std/ modules come from the standard library

Builtin functions:
print(t: string)
prints to stdout