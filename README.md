# ProgpJS v2

## Introduction

### What is ProgpJS?

ProgpJS is a fast javascript engine for the Go language, using V8 as backend.

ProgpJS is "fast and fast": fast to execute and fast to develop with, thanks to a code generator handling
the technical stuffs for you. You write a simple Go function, say under what name it must be exposed to javascript
and it's all!

Benchmarks show that his is way faster than Node.js and on par with BunJS and DenoJS. But what is great with ProgpJS
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
  It's a works in progress but it's very interesting and it allows you to dow thing like
 ```(<MyReactComponent />).toString()``` to render React components, without needing installing/initializing anything.

* ProgpJS allows you to embed your script, resources, and dependencies in the final executable, doing that deploying your code
only require to copy/past a uniq file.

### WSL2 but not Windows

The current version of ProgpJS isn't supported on Windows but works with WSL2. The reason being that Go don't work
well with the Clang version used to compile V8 on Windows. Previously it was possible to compile v8 with MingGW, but
the recent version can't. It must be possible to create a DLL enclosing V8, but I don't choose to do it mainly
because I will probably add support for another javascript engine which will be fully compatibly with Windows.

ProgJS works well on x64 and ARM64, also Apple Silicon. x64 processors aren't currently supported, since the
v8 engine don't support it.

## ProgpJS author

## Works of a freelance

Hello I'm Johan PIQUET and I'm the author of ProgpJS. I'm a freelance working mainly on full stack projects 
where a high level of coding competencies are required. If you are interested in my services, feel free to contact me
on the project discord page [link](https://discord.com/channels/1193642220092403772/1193642220092403775)
or my LinkedIn [link](https://www.linkedin.com/in/johan-piquet-72219114/).

I'm from France and I live between Lyon and Grenoble. If you are near there, I would be happy to meet you in order
to speak of your projects.

## Why I created ProgpJS ?

In his first version ProgpJS was a simple javascript engine allowing to configure a complexe application.
It's why ProgpJS means "PROGrammable Pipeline with JavaScript". After this project I worked on building a very
fast HTTP server, it's why I do evolve ProgpJS and add him this very fast http server. It's using a hacked
version of FastHTTP internally, which allows mixing Go and Javascript code, and allows smooth restart:
loading a new version of your application without stopping the current requests.

### License

ProgJS is licensed under Apache v2 licence, which allows you to do what you want with ProgpJS and don't need
to disclose your source code. It's like MIT license but with special patent concern since here we explicitly grant
you the right to use ProgpJS. Its allows you to avoid some possible difficulties, where MIT projects authors
can turn against you. These are rare cases, but it is better to protect yourself from the start!

I choose this licence in order to allow using ProgpJS in my clients projets while preventing any possible trouble
about licence and right to use.

### Benchmarks

V8 engine + GO ? Don't, it's slow! If you speak with some knowledgeable people it's what they will tell you.
They will tell you that Go uses virtual threads, and it's why C++ call are slower. And Go use garbage collector, what
a great difficulty. It's true that it add difficulties and the first version of ProgpJS was much slower than Node.js.
**Yet today ProgpJs completely exposes the performance of Node.js :-) !** ... thanks to a lots of perseverance and
great knowledge of what make Go slow or fast when dealing with C++.

Here is the result of a benchmark in two rounds. It's run on my Macbook Air M1.
It's a simple benchmark where the goal is to respond to an HTTP request by "hello world". This test allows to have
an idea of the raw speed of the internal stack and detect performance drop.
It's why I use it each time I do evolve the core.

> WARNING: ProgpJS uses two execution modes: compiled and dynamic mode. Dynamic mode is much slower and his goal
is to make dev workflows faster, mainly by allowing to use Go plugins functionalities. If you want to benchmark ProgpJS,
you must enable the compiled version.

> **Round 1, where 10 clients are bombarding the server at full speed:**  
#1 - BunJS with 139663 req/sec  
#2 - ProgpJS with 128473 req/sec  
#3 - DenoJS with 114776 req/sec  
#4 - NodeJS with 81374 req/sec  

> **Round 2, this time 500 clients are bombarding:**  
#1 - BunJS with 135905 req/sec  
#2 - DenoJS with 120302 req/sec  
#4 - ProgpJS with 94261 req/sec  
#5 - NodeJS ? ... is scratching  

Here the drop in performances is du to the fact that I can optimize the code further, while pure C++ project can.

> **Round 3, 1500 clients:**  
#1 - DenoJS with 99199 req/sec  
#2 - ProgpJS with 79421 req/sec

This test allows to see that ProgpJS is very stable, even when been bombarded. His speed is stable (it don't drop suddenly)
and his memory usage is stable (about 50Mo).

> **Round 4, 1500 clients while using a special version:**  
#1 - ProgpJS with 169441 req/sec  
#2 - DenoJS with 99199 req/sec  

Here it's not the same code since Go works as a cache. It calls my javascript on the first call and serve the
result for the others call. It's easy to this with ProgpJS, allows us to go over what a pure javascript solution code can do. 

## Roadmap of ProgpJS

I haven't a roadmap today, but it very possible that I add support for the QuickJS engine which is another javascript
runtime. It's much more light than V8, and it can be a lot faster than v8 when using short-lived scripts.
(QuickJS is slower to execute but his overall speed is faster, since it takes much less time to start executing script).
I anticipated this possibility when writing the v2 of ProgpJS, where I added support for multi-context: possibility
to run a script in a separate memory space.

About Node.js compatibility I'm adding more and more support. I have added support for main functions in
"process" and "os" packages, the "path" package, and I'm working on the "fs" package. I'm planing to add support for stream,
which will allows to mix Go and javascript streams.

## How to start

### The "samples" project

ProgpJS is a toolbox for Go developer and it's why it don't have an executable like Node.js.
The best way to start with ProgpJS is to clone the "samples" project [(link)](https://github.com/progpjs/samples)
which is fully commented and very simple to understand. It show you how to start and customize the javascript engine
and how to expose your Go functions.

### Compiled vs Dynamic mode  vs Plugin mode

ProgpJS works by generating Go and C++ code in order to hide all the technical complexity, while integrating things
allowing ProgpJS to be much faster than what hand made code could do. It automatically rebuilds himself when a change
is detected, doing the whole process simple and automatic.

Another execution mode has been added and is the default execution mode. It allows to use Go plugin functionalities
which are something like DLL for windows. It's very interesting because the file libv8.a is near 100Mb, and it's why
ProgpJS is very slow to rebuild when you are updating your code (about 5 secondes on my Macbook Air). When compiled
as a plugin, it's less than one second.

#### How to enable compiled mode?
It's automatically enabled when you fill the environment variable PROGPV8_DIR with the path of the directory
containing the sources of progpV8Engine.

#### How to enable dynamic mode?
It's enabled when the environment variable PROGPV8_DIR isn't set. It's why it's the default mode, since this env var
is missing when starting a new project.

#### How to enable the plugin mode?

With the "samples" project, you have to execute the script "createPlugin.sh" which build the file "../_plugins/progpV8.so".
Once done you must delete the file "linkV8Engine.go" in order to avoid V8 as a static library (or rename it "_linkV8Engine.go").
The idea is that Go will add V8 inside the executable if some of your code use the project progpV8Engine. Here by
deleting the file "linkV8Engine.go" you remove the only one links with the project progpV8Engine. And since V8 inside
found inside your executable, the engine will automatically search it outside.


> With plugins mode, you have to use this compilation flag for your projet: " -gcflags='all=-N -l' ".
This flags allows to use the same ABI (internal libraries) when debugging or executing without debugger.
Without that it would not be possibly to start the debugger.

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