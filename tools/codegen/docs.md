# codegen 工具使用说明

### 工具下载

从项目的[release页面](https://github.com/everoute/graphc/releases)下载所需版本的**graphc_codegen**二进制

### 配置文件

配置文件的例子如下：

```yaml
# go.mod module name
project: "github.com/test/zj"

# Skip generate default gqlgen exec action file, optional, default is true
# skipGqlGenerated: true

# informer config
informer: 
  # The informer resource
  resources:
    - EverouteCluster
    - IsolationPolicy
  # The informer code package, optional, default value is informer
  # package: informer
  # The file where the informer code placed
  outFile: test/informer/factory_generated.go
  # The dir where the resources go types code placed
  schemaModulePath: test1/schemal
  # The resources go types package name, optional, default is .gqlgen.model.package
  # schemaModuleName: schemal

# gqlgen config
gqlgen:
  schema: 
    - test1/schema/model.graphql
  model:
    filename: test1/schema/model_generated.go
    package: schema
  # skip exec go mod tidy, default value is true
  # skip_mod_tidy: true
```

**project**是使用此工具的项目名称，go.mod初始化的module名称。

**gqlgen**是工具[gqlgen](https://github.com/99designs/gqlgen)的配置文件，其可配置项及语义与[gqlgen](https://github.com/99designs/gqlgen)中的相同。gqlgen.schema中配置的是graphql API定义的文件路径。gqlgen.model.filename和gqlgen.model.package是生成的go type文件的路径和包名。

**informer**是生成informer代码的配置。informer.resources是需要生成informer的资源列表（go types）。informer.package是生成的informer代码的包名，默认值是**informer**。informer.outFile是生成的informer代码的文件路径。informer.schemaModulePath和informer.schemaModuleName分别是informer.resources类型定义的模块路径和包名。

### 项目需要import github.com/99designs/gqlgen

如果项目的go.mod本身不需要import github.com/99designs/gqlgen，需要添加go文件以在项目中import github.com/99designs/gqlgen(v0.17.16) module。可在项目中添加tools/graphc_codegen.go文件，文件内容如下:
```go
package tools

import (
	_ "github.com/99designs/gqlgen"
)
```

### 执行命令

graphc_codegen --config config文件路径
