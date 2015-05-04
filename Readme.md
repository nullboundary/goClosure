# goClosure
A simple utility to concat and minify javascript files via the Google Closure compiler API. 

If you just want a simple binary or go solution for minifying javascript, or if you just don't want or can't install Node.js to use Grunt or Gulp, then goClosure is for you. 

### go install

```
go install github.com/nullboundary/goClosure
```

## Usage

```	
COMMANDS:
   concat, c	Read an html file and concat the js files in order
   minify, m	Minify one js file via Google Closure Compiler API
   all, a	Both Concat and Minify. Input is same as concat
   help, h	Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --help, -h		show help
   --version, -v	print the version
```

### Concat javascript files into one file

```	
NAME:
   concat - Read an html file and concat the js files in order

USAGE:
   command concat [command options] [arguments...]

DESCRIPTION:
   goClosure concat <input html file> <output js file>

OPTIONS:
   --path, -p "none"	Changes the path root for the js files listed in the html input file. Syntax <oldPath>:<newPath>
   --modify, -m "none"	Changes input html file <script> tags replacing the many old js files with one concated new file

```


Create or use an html file that lists your js files in the order you would like to concat them. 

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Sample</title>
</head>
<body>
  <p>Go Closure</p>
  <script type="text/javascript" src="/assets/js/earlyload.js"></script>
  <!-- External libraries will be ignored -->
  <script type="text/javascript" src="https://ajax.googleapis.com/ajax/libs/jquery/1.11.0/jquery.min.js"></script>
  <!-- -p oldpath:newpath flag can replace the file path to match the local directory -->
  <script type="text/javascript" src="/assets/js/lib.js"></script>
  <script type="text/javascript" src="/assets/js/site.js"></script>
</body>
</html>
```

Example:
```
goClosure concat index.html site.min.js -m /assets/js -p /assets/:www/
```

### Minify one js file via Google Closure Compiler API

```	
NAME:
   minify - Minify one js file via Google Closure Compiler API

USAGE:
   command minify [arguments...]

DESCRIPTION:
   goClosure minify <inputfile> <outputfile>


```
Example:
```
goClosure minify site.js site.min.js
```

### Concat and Minify

```	
NAME:
   all - Both Concat and Minify. Input is same as concat

USAGE:
   command all [command options] [arguments...]

DESCRIPTION:
   goClosure all <input html file> <output js file>

OPTIONS:
   --path, -p "none"	Changes the path root for the js files listed in the html input file. Syntax <oldPath>:<newPath>
   --modify, -m "none"	Changes input html file <script> tags replacing the many old js files with one concated new file

```
Example:
```
goClosure all index.html site.min.js -m /assets/js -p /assets/:www/
```
