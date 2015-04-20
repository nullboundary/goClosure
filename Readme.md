# goClosure
A simple utility to concat and minify javascript files via the Google Closure compiler API. 

If you just want a simple binary or go solution for minifying javascript, or if you just don't want or can't install Node.js to use Grunt or Gulp, then goClosure is for you. 

#### Commands:

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

Concat

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

Minify

```	
NAME:
   minify - Minify one js file via Google Closure Compiler API

USAGE:
   command minify [arguments...]

DESCRIPTION:
   goClosure minify <inputfile> <outputfile>


```

All

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
