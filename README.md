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
    - 在保证渲染一致的情况下, 能够生成单 HTML 页面以及导出为 PDF

本项目是本人边学习 hugo 源码而开发的个人项目.


# TODO

- 实现文本到简单 html 的渲染

# 格式

- 代码格式化

    行内格式化使用两个点将待格式化文本包括起来, 如 `..const var a = 10;..`.
    如果想使用一整段代码, 则可以使用 tab 缩进(四个空格), 比如

    ```
    #include <iostream>
    using namspace std;
    int main() {
        cout << "Hello pigger" << endl;
        return 0;
    }
    ```

行内格式化不采用 markdown 的反单引号标记, 因为文本不友好.

- 文章元信息(meta info)

    在文章开头处用 `---` 组成一个节区, 写入相关信息, 可写的信息如下,
    标 `*` 的表示必选

    ```
    ---
    - Title: *
    - Author:
    - Date:
    - Url:

    - Topic:
    - Keyword:
    ---
    ```

- 链接

    使用 `@<somthing>` 表示一个链接, 可以是本地的也可以是远程的,
    如果是本地的, 文件将会被拷贝到静态目录里面, 然后上传到远程服务器上,
    @ 会被渲染成一个链接.

# 架构设计

```
sys/
    etc/   # pigger 系统配置文件
       ├── css/
       │   ├── normalize.css
       │   └── pigger.css
       └── themes/
           └── default.json
    www/      # 生成的 html
        images/ # 图片
        videos/ # 视频
        aritcle-demo.html # 生成的样例文章
        text-demo.html # 生成的样例文章
    tmp/ # 临时文件
index.html # 首页

etc/  # 用户配置文件, 可以覆盖 pigger 系统配置 sys/etc/
    css/
    themes/
home/ # 写作平台
    assets/
        images/
        videos/
    md-demo.md # markdown 文件
    text-demo.txt # 文本文件
    draft/ # 草稿文件, 生成 html 时将会跳过
        drft-demo.md
```

用户只需要保留 `usr/` 和 `home/` 目录即可, 可以方便的实现迁移.
