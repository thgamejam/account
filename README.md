# 项目结构

```
.
├── proto
│   ├── api  // API & Error Proto files & Generated codes
│   │   ├── foo
│   │   │   ├── job
│   │   │   └── service
│   │   └── bar
│   │       └── interface
│   └── conf  // Config Proto Proto files
│       └── interface
└── app  // kratos microservices projects
    ├── foo
    │   ├── job
    │   └── service
    └── bar
        └── interface
```


## Generate other auxiliary files by Makefile
```
# 初始化项目并下载和更新依赖项
make init
# 依赖注入
make wire
# 生成错误文件代码
make error
# 生成配置文件代码
make config
# 生成api文件代码
make api
# 构建
make build
# 生成所有代码
make all
# 显示帮助
make help
```

## Docker
```bash
# build
docker build -t <your-docker-image-name> .

# run
docker run --rm -p 8000:8000 -p 9000:9000 -v </path/to/your/configs>:/data/conf <your-docker-image-name>
```

