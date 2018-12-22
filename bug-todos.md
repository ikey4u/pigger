# bugs

- 标题中的高亮无法渲染, 如
```
### `char **ss` 与 `char *ss[]`
```

- 下面这种代码高亮无法正常渲染
```
                       MISCELLANEOUS COMMANDS

     -<flag>              Toggle a command line option [see OPTIONS below].
     --<name>             Toggle a command line option, by name.
     _<flag>              Display the setting of a command line option.
     __<name>             Display the setting of an option, by name.
     +cmd                 Execute the less cmd each time a new file is examined.
     !command             Execute the shell command with $SHELL.
     |Xcommand            Pipe file between current pos & mark X to shell command.
     v                    Edit the current file with $VISUAL or $EDITOR.
     V                    Print version number of "less".
```

# todo

- 制作安装包
- 不通等级的标题着以不同颜色以示区分
- 两个相邻的 list 项之间可以加入多个空行以使结构清晰, 比如

```
- item one

- item two

- itme three
```

但是渲染后 item 之间不需要留有空行.

- 一个 list 项的内容中可以有空行隔开两段以使结构化清晰, 比如

```
- A list

    Para one

    Para two

    Para three

        The codes

    Para four
```

渲染时段与段之间需要有一个空行.

- 列表项换行添加新方法: 如果一个列表项后空一行, 那么表示要新换行.
比如

```
- The item

    The content of the item
```

等价于

```
- The item:
    The content of the item
```

都表示在 `The item` 之后新起一行写入 item 的内容.

- 博客系统自动删除 .txt 不存在的文章
