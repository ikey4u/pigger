# pigger

用 golang 编写的静态博客生成器

# 项目简介

写该项目有如下几个原因:

- 实践 golang 语言
- 希望写一个适合自己的, 可以自定义的静态博客生成器, 其功能应该有如下几个
    - 基于 markdown 语法并做适当扩展
    - 博文应该是文本友好型, 也就是说, 一篇博文, 在浏览器中和用文本编辑器打开时,
        都能够给人带来舒适的阅读体验, 因此对文本书写格式有要求,
        但是 markdown 语法还是略显复杂
    - 不支持表格, 因为文本不友好
    - 在保证渲染一致的情况下, 能够生成单 HTML 页面以及导出为 PDF

# TODO

- 文本到 HTML 的渲染 (✓)
- 解析文章头(标题, 作者, 时间) (✓)
- 实现列表的渲染 (✓)
- 捆绑静态资源实现代码高亮(css, js 等) (✓)
- 行内代码块渲染(✓)
- 独立代码块(✓)
- 列表内代码块渲染(✓)
- HTML 标记转义(✓)
- UTF-8 字符处理(✓)
- @ 语法自动插入图片(暂时不支持远程 URL 插入图片)功能 (✓)
- @ 语法插入超链接功能(必须以 http 或 ftp 开头) (✓)
    如果链接长度超过 32, 则后续部分以三个省略号来表示.
- 代码块空行功能
- 导出 PDF
- 实时渲染服务

# 格式

- 不支持表格, 因为单独打开文本进行查看的时候不方便

- 列表渲染
    列表的第一行如果如果末尾字符为冒号(忽略右侧空格), 则表示换行,
    也就是说此时列表的第一行作为一个简短的小标题.

- 代码格式化
    行内格式化时使用两个反斜号将待格式化的文本括起来,
    如果想在行内显示一个反斜号, 请务必保证该行只有一个反斜号.

    如果想使用一整段独立的代码, 则可以使用 tab 缩进(四个空格),
    第一行使用 `//:` 来指明高亮的语言, 改行在实际显示时将会被忽略,
    且改行是可选的, 默认的高亮采用的是 C 语言家族. 比如

    ```
    //: c++
    #include <iostream>
    using namspace std;
    int main() {
        cout << "Hello pigger" << endl;
        return 0;
    }
    ```

    列表中也可以使用代码块, 采用 8 个缩进即可.

- 文章元信息(meta info)
    在文章开头处用 `---` 组成一个节区, 写入相关信息, 可写的信息如下,
    标 `*` 的表示必选, 其他选项暂不支持.

    ```
    ---
    - title: *
    - author: *
    - date: *
    ---
    ```

- 链接

    使用 `@<somthing>` 表示一个链接, 可以是本地的也可以是远程的,
    如果是本地的, 文件将会被拷贝到静态目录里面, 然后上传到远程服务器上,
    @ 会被渲染成一个链接.

# 架构设计

```
// 静态资源
etc/   # pigger 系统配置文件
   ├── css/
   │   ├── normalize.css
   │   └── pigger.css
   └── themes/
       └── default.json

// 新建目录
SITE/
    3w/         # 输出目录
        css/
        images/ # 图片
        videos/ # 视频
        aritcle-demo.html # 生成的样例文章
        text-demo.html # 生成的样例文章
        tmp/ # 临时文件
        index.html # 首页
    home/ # 写作平台
        usr/  # 用户配置文件, 可以覆盖或扩展 pigger 系统配置
            css/
            themes/
        assets/
            images/
            videos/
        md-demo.md # markdown 文件
        text-demo.txt # 文本文件
        draft/ # 草稿文件, 生成 html 时将会跳过
            drft-demo.md
```

用户只需要保留 `usr/` 和 `home/` 目录即可, 可以方便的实现迁移.
