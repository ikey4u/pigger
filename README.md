# Pigger

## Note taking && Blog writing in a fucking simple but neat way

![](./docs/logo.png)

## 1. Installation

Download installer for your platform from following links

- windows → [pigmgr_windows_386](https://github.com/ikey4u/pigger/releases/download/v1.0.4/pigmgr_windows_386.exe)
- mac → [pigmgr_darwin_amd64](https://github.com/ikey4u/pigger/releases/download/v1.0.4/pigmgr_darwin_amd64)
- linux → [pigmgr_linux_386](https://github.com/ikey4u/pigger/releases/download/v1.0.4/pigmgr_linux_386)

For linux and mac user, you should make the installer executable, for example

    chmod +x pigmgr_darwin_amd64

Then you can run `pigmgr_xxx` from command line and finish the installation.


## 2. Document format

![](./docs/demodoc.png)
    
Find the PDF in [here](https://raw.githubusercontent.com/ikey4u/pigger/master/docs/demodoc.pdf)

Notice

>You should use a full functional text editor such as `notepad++`,
`vs code editor`, `sublime`, or `vim` and so on.

>But but but not the editor such as windows's notepad! Please!

The text format must be `unix` type (which means that its newline character should be `\n`
but not else something)

## 3. To note taking

A sample note could be found here [http://ahageek.com/writer/posts/msb-and-lsb/index.html](http://ahageek.com/writer/posts/msb-and-lsb/index.html)

You may add a `.txt` suffix to the link  and open the source text file
[http://ahageek.com/writer/posts/msb-and-lsb/index.html.txt](http://ahageek.com/writer/posts/msb-and-lsb/index.html.txt) 

To generate a note (let us assume that the text name is `msb-and-lsb.txt`), you could
run

    pigger msb-and-lsb.txt

Then open `msb-and-lsb/index.html` to see the magic.


## 4. To blog writing

A sample blog is [http://ahageek.com/writer/site.html](http://ahageek.com/writer/site.html)

To generate a static blog(lets assume that the blog root is `writer`), you just need one
command line

    pigger new writer

Then you get a static blog!


In `writer` directory, you write a simple text file(must has a suffix `.txt`) following the
document format, then run

    pigger build

`pigger` will render all your files into beautiful htmls.


Push `writer` into the secondary domain of github, then you are done!

# Support ?


If you like `pigger`, please give `pigger` a star.

![](./docs/givestar.png)

# License

MIT License
