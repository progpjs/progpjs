# ProgpJS v2

## Introduction

### What is ProgpJS?

ProgpJS is a fast javascript engine for the Go language, using V8 as backend.

ProgpJS is "fast and fast": fast to execute and fast to develop with, thanks to a code generator handling
the technical stuffs for you. You write a simple Go function, say under what name it must be exposed to javascript
and it's all!

Benchmarks show that it's way faster than Node.js and on par with BunJS and DenoJS. But what is great with ProgpJS
isn't only his speed, but his capacity to easily mix Go code, C++ code and Javascript code, while been very simple to use.
With ProgpJS there isn't technical complexity, nothing, ProgpJS takes in charge all the difficulties for you!

ProgpJS comes with a Node.js compatibility layer for projects needing it. It's in the first stage, but it can be
useful for those needing it. The goal of ProgpJS isn't to be compatible with Node.js, mainly because ProgpJS goal
is to make interact javascript and Go, it's not a Node.js replacement. With Node.js your project is 100% javascript
while with ProgpJS you code your tools in Go (for his high speed) and you use your components through javascript.

### Why will you love ProgpJS ?

* It's very easy to embed ProgpJs in your Go project. Most of ProgpJS functionalities are pluggable, allowing
you to select what you want to add or not. You can add all the builtin function or keep the things minimal and light.

* With ProgpJS the technical parts are handled by the engine through generated code (if you use compiled code)
or reflection (if you use dynamic mode). Its make it very easy and fast to add new functionalities!

* ProgpJS comes with support for Chrome Debugger: you can debug the javascript code through the same debugger
as the one used inside Chrome browser.

* ProgpJS comes with a very fast http stack.
This stack allows to handle request with Go code or javascript code as you wish. Thanks to that you can built
advanced server functionalities. For exemple Go can catch the request, check a cache and send the request
to the javascript engine if the content isn't in the cache yet.

* ProgpJS embeds a very fast Typescript compiler and translate error messages: you can write Typescript code and
use it as-is, without any configuration required. Moreover with ProgpJS error messages points to your typescript code
not the underlying javascript code!

* ProgpJS comes with his own ReactJS implementation, which is optimized for Server Side Rendering.
  It's a works in progress but it's very interesting and it allows you to do thing like
 ```(<MyReactComponent />).toString()``` to render React components, without needing installing/initializing anything.

* ProgpJS allows you to embed your script, resources, and dependencies in the final executable, doing that deploying your code
only require to copy/past a uniq file.

### Linux, Apple, WSL2, X64 and ARM. But not Windows

The current version of ProgpJS isn't supported on Windows but works with WSL2. The reason being that Go don't work
well with the Clang version used to compile V8 on Windows. Previously it was possible to compile v8 with MingGW, but
the recent version can't. It must be possible to create a DLL enclosing V8, but I don't choose to do it mainly
because I will probably add support for another javascript engine which will be fully compatibly with Windows.

ProgJS works well on x64 and ARM64, also Apple Silicon. x86 processors aren't currently supported, since the
v8 engine don't support it. I don't know if some projects require x86 support, if it's the case feel free to contact
me in order to discuss it.

## ProgpJS author

## I'm freelance on complex full stack projects

Hello I'm Johan PIQUET, and I'm the author of ProgpJS. I'm a freelance working mainly on full stack projects 
where a high level of coding competencies are required. My clients are small compagnies but also prestigious compagnies
for whom I carry out project management missions as well as software development missions.

If you are searching someone with good relational competencies able to manage a team of young developers
and develop their skills, then feel free to contact me n the project discord page [link](https://discord.com/channels/1193642220092403772/1193642220092403775) or my LinkedIn [link](https://www.linkedin.com/in/johan-piquet-72219114/).
Also as you guess I'm very proficient in Go, C++, Javascript, Node.js and ReactJs.

I'm from France and I live between Lyon and Grenoble. If you are near there, I would be happy
to meet you in order to speak of your projects.

## Why I created ProgpJS ?

In his first version ProgpJS was a simple javascript engine allowing to configure a big application.
It's why ProgpJS means "PROGrammable Pipeline with JavaScript". After this project I worked on a project
requiring a very fast HTTP server, it's why I do evolve ProgpJS and add him this very fast http server.
It's using a hacked version of FastHTTP internally, which was the faster http server for years. This hacked version
allows mixing Go and Javascript code, and allows smooth restart: loading a new version of your application
without stopping the current requests.

### License

ProgJS is licensed under Apache v2 licence, which allows you to do what you want with ProgpJS and don't need
to disclose your source code. It's like MIT license but with special patent concern since here we explicitly grant
you the right to use ProgpJS. Its allows you to avoid some possible difficulties, where MIT projects authors
can turn against you. These are rare cases, but it is better to protect yourself from the start!

I choose this licence in order to allow using ProgpJS in my clients projets while preventing any possible trouble
about licence, ownership and right to use.

### Benchmarks

**Its was very difficult but today ProgpJs completely is near x2 faster than NodeJS !**

V8 engine + GO ? Don't, it's slow! If you speak with some knowledgeable people it's what they will tell you.
They will tell you that Go uses virtual threads, and it's why C++ calls are very slow. And Go uses garbage collector.
And Go encodes his strings and data very differently that what C++ does, adding a great level of difficulties.

Yes it's true, creating a fast javascript engine for Go was very challenging! The first internal version of
ProgpJS was much slower than Node.js. The second internal version was faster than DenoJS and BunJS,
but very difficult to maintains. The third version, which is the first public one, was easy to use and maintains, while been
only a little slower than DenoJS. The current version (v2 public version) has the same speed as the previous version
but add multi-context, which allows to use all core of your server. It can make thing much, much faster if your
application do heaving thing with javascript.

Here is the result of a benchmark in four rounds. It's run on my Macbook Air M1.
It's a simple benchmark where the goal is to respond to an HTTP request by "hello world".
This test allows to have an idea of the raw speed of the internal stack and detect performance drop.
It's why I use it each time I do evolve the core.

> WARNING: ProgpJS uses two execution modes: compiled and dynamic mode. Dynamic mode is slow, his goal
is to make dev workflows faster, mainly by allowing to use Go plugins functionalities. If you want to benchmark ProgpJS,
you must enable the compiled version (see the README of "samples" project).

> **Round 1, where 10 clients are bombarding the server at full speed:**  
#1 - BunJS with 139663 req/sec  
#2 - ProgpJS with 128473 req/sec  
#3 - DenoJS with 114776 req/sec  
#4 - NodeJS with 81374 req/sec  

> **Round 2, this time 500 clients are bombarding:**  
#1 - BunJS with 135905 req/sec  
#2 - DenoJS with 120302 req/sec  
#4 - ProgpJS with 94261 req/sec  
#5 - NodeJS ? ... is scratching (near the 250 clients)

Here the drop in performances is du to the fact that I can't optimize the code further, while pure C++ project can.
This drop is limited, the performance slow down gradually when the number of concurrent requests raise.

> **Round 3, 1500 clients:**  
#1 - DenoJS with 99199 req/sec  
#2 - ProgpJS with 79421 req/sec

This test allows to see that ProgpJS is very stable, even when been bombarded with a lot of connexion.
His speed is stable (it doesn't drop suddenly) and his memory usage is stable (about 80mb here).

> **Round 4, 1500 clients while using a special version:**  
#1 - ProgpJS with 169441 req/sec  
#2 - DenoJS with 99199 req/sec  

Here it's not the same code since Go intercept the http call, check a cache and call the javascript
only when the cache is empty. It's easy to do this with ProgpJS, which allow us to go over what a pure javascript
solution can do. 

## Roadmap of ProgpJS

I haven't a roadmap today, but I'm starting a new project where I will need to use the QuickJS engine
which is another javascript runtime. It's much more light than V8, and it can be a lot faster when using
short-lived scripts. QuickJS is slower to execute, but is a lot faster to start, doing that the overall speed
is better for solutions where starting a lot of small scripts.

About Node.js compatibility I'm adding more and more support, but it takes time, mainly because
I need to do a lot of compatibility tests for each function and track the exact behaviors of Node.js.
The package "path" is fully compatible and tested, while package "process", "os" and "fs" are only partially implemented
and barely tested. My actual goal was mainly to known if thing was missing or if the core of ProgpJS already
implements all the prerequisites.

## How to start

### The "samples" project

ProgpJS is a toolbox for Go developer and it's why it doesn't have an executable like Node.js.
Perhaps when it will be more compatible with Node.js, but it's not the case today, nor his goal.

The best way to start with ProgpJS is to clone the "samples" project [(link)](https://github.com/progpjs/samples)
which is fully commented and very simple to understand. It shows you how to start, how to customize the javascript engine
and how to expose your Go functions.

You can read the file README.md of this project, which explains to you how to create your fist project and use one
of the 3 execution mode: compiled mode, dynamic mode or plugin mode.

**Compiled mode** - ProgpJS works by generating Go and C++ code in order to hide all the technical complexity, while integrating things
allowing ProgpJS to be much faster than what hand made code could do. It automatically rebuilds himself when a change
is detected, doing the whole process simple and automatic.

**Plugin mode** - Another execution mode has been added and is the default execution mode. It allows to use Go plugin functionalities
which are something like DLL for Windows. It's very interesting because the file libv8.a is near 100Mb, and it's why
ProgpJS is very slow to rebuild when you are updating your code (about 5 secondes on my Macbook Air). When compiled
as a plugin, it's less than one second.

**Dynamic mode** - This mode is what ProgJS enavle when he don't find the source code of the project ProgpV8Engine.
It uses reflection to known how too call our Go function, which is slow.

## How to debug?

ProgpJS implements the debugger protocol which allow to debug your javascript code with the debugger of the Chrome browser,
the same used to debug a web page. When using the "samples" project, you can enable the javascript debugger by setting
the environnement variable "PROGP_DEBUG" to 1. 

Once started in debug mode a message is printed in the console, asking you to open the url "chrome://inspect/#devices"
in Chrome browser. Once done you can click on "Open dedicated DevTools for Node" to open the debug window.

> The debugger is limited and can only debug the first script executed with ProgpJS, mainly because this protocol
it has been created for Node.js, which only use one script. There is possible workarounds, but I didn't have time
to work on this hacks.

## Some little samples

### Exposing Go functions to javascript

Here we want to show you add to expose a Go function to javascript. It's very simple.
If you want more samples, you can read the code source of the project *progpjs.samples*.
In the sample here we will expose five Go functions to javascript.

```go
func declareMyJavascriptFunctions(group *progpAPI.FunctionGroup) {
    // The function registry is were all exposed function are declared.
    rg := progpScripts.GetFunctionRegistry()
    
    // This line will help the code generator to know in wich namespace is located our function.
    goNS := rg.UseGoNamespace("github.com/progpjs/samples/v2/modSamples")
    
    // Here it's like a namespace but for the javascript side.
    // We will use the global namespace, in order to automatically add our function
    // in the javascript global namespace (see globalThis), which make our function
    // directly available.
    //
    group := myMod.UseGroupGlobal()
    
    // Now we can declare our functions.
    // The first parameter is the javascript function name, here it will be "testThrowError".
    // The second one is the name of our go function, and the go function himself.
    //
    // Sample from javascript: try { testThrowError("boom") } catch(e) { console.error(e) }
	//
    group.AddFunction("testThrowError", "JsThrowError", JsThrowError)
	
	// Adding a new function only require adding another line.
    // We add "Js" as prefix, it's only a code convention. 
    //
	// Sample from javascript: sayHello({name: "Doe", forename: "John"})
	//
    group.AddFunction("sayHello", "JsSayHello", JsSayHello)
	
	// Here we add an async function.
    group.AddAsyncFunction("myAsyncFunction", "JsSamplesCallAsync", JsSamplesCallAsync)

	// A shared resource is a pointer to a Go object. For exemple the current HTTP call.
    // Here the first function create and send a shared resource to javascript.
    // While the second one get back this resource from JavaScript.
	//
	// Sample javascript usage:
    //  let sr = testReturnSharedResource();
	//  testReceiveSharedResource(sr);
	//
    group.AddFunction("testReturnSharedResource", "JsTestReturnSharedResource", JsTestReturnSharedResource)
    group.AddFunction("testReceiveSharedResource", "JsTestReceiveSharedResource", JsTestReceiveSharedResource)
}

// Here our function is a simple Go function.
// We don't have to known how the engine works, we only write simple Go functions.
//
func JsThrowError(value string) error {
    if value == "boom" { return errors.New("big boom") }
    return nil
}

// We can use struct and struct pointer as input arguments.
// As for return value.
//
func JsSayHello(who Person) string {
    return "Hello " + who.Name + " " + who.Forename    	
}

// Async function name must ends with Async which allows to prevent some errors.
// Here this function do the same thing as JsSayHello, but it's async which mean
// that the thread isn't locked.
//
func JsSamplesCallAsync(who Person, callback progpAPI.JsFunction) {
    // Here it's like writing go fun() { ... }
    // but it's add an handler catch global errors.
    //
    progpAPI.SafeGoRoutine(func() {
        // Add a pause for sample
        progpAPI.PauseMs(100)
        
        // The 2 at the end mean that we return a string
		// as the second parameter. It's required since the callback function
		// generally use the first argument as the error message.
		//
        callback.CallWithString2(JsSayHello(who))
    })
}

func JsTestReceiveSharedResource(sr *progpAPI.SharedResource) {
    println("Received shared resource: ", sr.GetId(), " - Value: ", sr.Value.(string))
}

func JsTestReturnSharedResource() *progpAPI.SharedResource {
	// The second argument is nil here, but can a function called
	// when the shared resource is disposed.
	//
    sr := progpAPI.NewSharedResource("my resource", nil)
    return sr
}

// Here is the struct we use. The json metadata are optional
// and are here to known that 'name" must be converted to "Name".
//
type Person struct {
    Name string `json:"name"`
    Forename string `json:"forename"`
}

```

## Exposing embedded javascript files

It's very simple to embed javascript module withing ProgpJS, mainly because Go add a functionality for that.
In this code we create a file named embed/js/myJsModule.ts at the root of our project.
Once done we only have to add this code in order to bind our javascript module sources:

```go

// Here it's a Go feature allowing to read files embedded in our executable
// (and allowing Go to known that we want to embedded this directory).
//
//go:embed embed/*
var gEmbedFS embed.FS

// Here it's our function declaring our embedded module.
// You must call it when starting your app.
//
func declareMyEmbeddedJavascriptModule() {
    // The provided is what will take the embedded file and return it when our module is requested. 
    // It does anything thing here: say that it's a typescript file.
	//
    provider := progpScripts.ReturnEmbeddedTypescriptModule(gEmbedFS, "embed/js/myJsModule.ts")
	
	// Here we say to the engine that we provide ourselves the module embedded:myJsModule.
    // If you do import("embedded:myJsModule") from javascript then this code will be called.
	//
    progpScripts.AddJavascriptModuleProvider("embedded:myJsModule", provider)
}

```